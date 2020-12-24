package conferences

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stripe/stripe-go/webhook"

	"encore.dev/beta/auth"
	"encore.dev/storage/sqldb"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

// CheckoutParams is the payload of checkout which contains the information about the
// slot claim purchase about to happen.
type CheckoutParams struct {
	ConferenceID uint64
	VoucherID    string
	SlotsToClaim []uint64
}

// CheckoutResponse returns information about the outcome of the checkout initiation
// process.
type CheckoutResponse struct {
	StripeSessionID       string
	ClaimPayment          *ClaimPayment
	NoStockOfSlot         *ConferenceSlot
	MissingSlotDependency *ConferenceSlot
}

type secrets struct {
	StripeKey      string
	EndpointSecret string
}

// Checkout initiates the checkout process for purchasing slot attendance in a conference.
func Checkout(ctx context.Context, params *CheckoutParams) (*CheckoutResponse, error) {
	usr, ok := auth.Data().(*User)
	if !ok {
		return nil, fmt.Errorf("unable to use %T %v as a user", usr, usr)
	}
	// Check that we have capacity remaining in the slot and that customer is purchasing
	// all dependencies.
	if err := ensureSlotsCanBeClaimed(ctx, params.SlotsToClaim, usr.Email); err != nil {
		noStock := &ErrSlotFull{}
		missingDep := &ErrDependencyUnmet{}
		switch {
		case errors.As(err, &noStock):
			cs, err := readConferenceSlotByID(ctx, nil, noStock.claimFullID, false)
			if err != nil {
				return nil, fmt.Errorf("reading conference slot for slot full error: %w", err)
			}
			return &CheckoutResponse{NoStockOfSlot: cs}, nil
		case errors.As(err, &missingDep):
			cs, err := readConferenceSlotByID(ctx, nil, missingDep.missing, false)
			if err != nil {
				return nil, fmt.Errorf("reading conference slot for missing dependency error: %w", err)
			}
			return &CheckoutResponse{MissingSlotDependency: cs}, nil
		}
		return nil, fmt.Errorf("ensuring slots can be claimed: %w", err)
	}
	// do the unpaid claims.

	slots, err := readConferenceSlotsByIDs(ctx, nil, params.SlotsToClaim)
	if err != nil {
		return nil, fmt.Errorf("reading slots to claim from database: %w", err)
	}
	claims, err := claimSlots(ctx, usr, slots)
	if err != nil {
		return nil, fmt.Errorf("claiming slots: %w", err)
	}

	// Initiate payment
	scrt := secrets{}
	stripe.Key = scrt.StripeKey

	voucherRow := sqldb.QueryRow(ctx, `SELECT
		valid_from,
		valid_to,
		discount_percentage,
		discount_amount_cents,
		discount_percentage_max_amount_cents,
		conference_id
	FROM discount_vouchers 
	WHERE voucher_id = $1 AND
	conference_id = $2 AND
    backing_payment_id IS NULL AND
	transaction_id = ''`, // basically, not used or being used
		params.VoucherID, params.ConferenceID)

	voucherInfo := VoucherInformation{}
	err = voucherRow.Scan(
		&voucherInfo.ValidFrom,
		&voucherInfo.ValidTo,
		&voucherInfo.Percentage,
		&voucherInfo.AmountInCents,
		&voucherInfo.LimitInCents,
		&voucherInfo.ConferenceID,
	)
	if err != nil {
		return nil, fmt.Errorf("reading voucher information: %w", err)
	}

	totalToDiscount := voucherInfo.AmountInCents
	if voucherInfo.Percentage != 0 {
		totalCost := 0
		for _, claim := range claims {
			totalCost += claim.ConferenceSlot.Cost
		}
		totalToDiscountFloating := float64(totalCost) * (float64(voucherInfo.Percentage) / 100)
		if totalToDiscountFloating > float64(voucherInfo.LimitInCents) && voucherInfo.LimitInCents > 0 {
			totalToDiscountFloating = float64(voucherInfo.LimitInCents)
		}
		// yes, there is loss
		totalToDiscount = int64(totalToDiscountFloating)
	}
	remainingDiscount := totalToDiscount

	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(claims))
	usd := string(stripe.CurrencyUSD)
	var claimIDs []int64
	var totalCost int64
	for i := range claims {
		claimIDs = append(claimIDs, claims[i].ID)
		cost := int64(claims[i].ConferenceSlot.Cost)
		if remainingDiscount > 0 {
			// FIXME: Apply discount to payment
			// FIXME: Mark stripe session to voucher for discount
			switch {
			case remainingDiscount >= cost:
				remainingDiscount -= cost
				cost = 0
			case remainingDiscount < cost:
				cost -= remainingDiscount
				remainingDiscount = 0
			}
			totalCost += cost
		}
		lineItems = append(lineItems,
			&stripe.CheckoutSessionLineItemParams{
				Amount:      &cost,
				Currency:    &usd,
				Name:        &claims[i].ConferenceSlot.Name,
				Description: &claims[i].ConferenceSlot.Description,
				Quantity:    stripe.Int64(1),
			},
		)
	}

	if totalCost == 0 {
		p := &PaymentMethodConferenceDiscount{
			Detail:      params.VoucherID,
			AmountCents: totalToDiscount,
		}
		payment, err := payClaims(ctx, usr, claims, []FinancialInstrument{p})
		if err != nil {
			return nil, fmt.Errorf("paying in full with voucher: %w", err)
		}

		_, err = sqldb.Exec(ctx, `UPDATE discount_vouchers 
				SET backing_payment_id = $1 
				WHERE voucher_id = $2 AND
				conference_id = $3`, payment.ID, params.VoucherID, params.ConferenceID)
		if err != nil {
			return nil, fmt.Errorf("setting backing payment ID to voucher: %w", err)
		}
		return &CheckoutResponse{ClaimPayment: payment}, nil
	}

	stripeParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://example.com/success"), // FIXME tricky one
		CancelURL:  stripe.String("https://example.com/cancel"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: lineItems,
		Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	s, err := session.New(stripeParams)
	if err != nil {
		return nil, fmt.Errorf("creating new stripe session to checkout: %w", err)
	}

	// We will use transaction to retrieve vouchers and claims in case of payment or cancellation
	_, err = sqldb.Exec(ctx, `UPDATE discount_vouchers 
				SET transaction_id = $1 
				WHERE voucher_id = $2 AND
				conference_id = $2`, s.ID, params.VoucherID, params.ConferenceID)
	if err != nil {
		return nil, fmt.Errorf("setting session ID to voucher: %w", err)
	}

	_, err = sqldb.Exec(ctx, `UPDATE slot_claim 
				SET transaction_id = $1 
				WHERE id = ANY($2)`, s.ID, claimIDs)
	if err != nil {
		return nil, fmt.Errorf("setting session ID to claims: %w", err)
	}

	return &CheckoutResponse{StripeSessionID: s.ID}, nil
}

//encore:api public raw
func SuccessWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// do we have logging?
		//fmt.Errorf("reading request body: %w", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key
	// You can find your endpoint's secret in your webhook settings
	endpointSecret := secrets{}.EndpointSecret
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), endpointSecret)

	if err != nil {
		//fmt.Errorf("verifying webhook signature: %w", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Handle the checkout.session.completed event
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			// fmt.Errorf("parsing webhook JSON: %w", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// FIXME: Load user
		usr := &User{}

		rows, err := sqldb.Query(req.Context(),
			`SELECT id, ticket_id, redeemed FROM slot_claim WHERE transaction_id = $1`,
			session.ID)
		if err != nil {
			// fmt.Errorf("reading claims for payment: %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var claims []SlotClaim
		for rows.Next() {

			claim := SlotClaim{}
			err := rows.Scan(&claim.ID, &claim.TicketID, &claim.Redeemed)
			if err != nil {
				// fmt.Errorf("scanning slot_claim for attendee: %w", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			claims = append(claims, claim)
		}

		voucherRow := sqldb.QueryRow(req.Context(), `SELECT
		valid_from,
		valid_to,
		discount_percentage,
		discount_amount_cents,
		discount_percentage_max_amount_cents,
		conference_id
	FROM discount_vouchers 
	WHERE transaction_id = $1`,
			session.ID)

		voucherInfo := VoucherInformation{}
		err = voucherRow.Scan(
			&voucherInfo.ValidFrom,
			&voucherInfo.ValidTo,
			&voucherInfo.Percentage,
			&voucherInfo.AmountInCents,
			&voucherInfo.LimitInCents,
			&voucherInfo.ConferenceID,
		)
		if err != nil {
			// fmt.Errorf("reading voucher information: %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		totalToDiscount := voucherInfo.AmountInCents
		if voucherInfo.Percentage != 0 {
			totalCost := 0
			for _, claim := range claims {
				totalCost += claim.ConferenceSlot.Cost
			}
			totalToDiscountFloating := float64(totalCost) * (float64(voucherInfo.Percentage) / 100)
			if totalToDiscountFloating > float64(voucherInfo.LimitInCents) && voucherInfo.LimitInCents > 0 {
				totalToDiscountFloating = float64(voucherInfo.LimitInCents)
			}
			// yes, there is loss
			totalToDiscount = int64(totalToDiscountFloating)
		}

		paymentMethods := []FinancialInstrument{&PaymentMethodMoney{
			AmountCents: session.PaymentIntent.Amount,
			PaymentRef:  session.ID,
		}}

		if totalToDiscount > 0 {
			paymentMethods = append(paymentMethods,
				&PaymentMethodConferenceDiscount{
					AmountCents: totalToDiscount,
					Detail:      session.ID,
				})
		}

		payment, err := payClaims(req.Context(), usr, claims, paymentMethods)
		if err != nil {
			// fmt.Errorf("registering claims payment: %w", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payment) // FIXME: do something with this err?
	}

	w.WriteHeader(http.StatusOK)
}
