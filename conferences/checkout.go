package conferences

import (
	"context"
	"errors"
	"fmt"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/checkout/session"
)

// CheckoutParams is the payload of checkout which contains the information about the
// slot claim purchase about to happen.
type CheckoutParams struct {
	SlotsToClaim []uint64
}

// CheckoutResponse returns information about the outcome of the checkout initiation
// process.
type CheckoutResponse struct {
	StripeSessionID       string
	NoStockOfSlot         *ConferenceSlot
	MissingSlotDependency *ConferenceSlot
}

// Checkout initiates the checkout process for purchasing slot attendance in a conference.
func Checkout(ctx context.Context, params *CheckoutParams) (*CheckoutResponse, error) {
	// Check that we have capacity remaining in the slot and that customer is puchasing
	// all dependencies.
	attendeeEmail := "" // FIXME: Find this from user in request.
	if err := ensureSlotsCanBeClaimed(ctx, params.SlotsToClaim, attendeeEmail); err != nil {
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
	attendee, err := readAttendeeByEmail(ctx, nil, attendeeEmail)
	if err != nil {
		return nil, fmt.Errorf("reading attendee to claim slots")
	}
	slots, err := readConferenceSlotsByIDs(ctx, nil, params.SlotsToClaim)
	if err != nil {
		return nil, fmt.Errorf("reading slots to claim from database: %w", err)
	}
	claims, err := claimSlots(ctx, attendee, slots)
	if err != nil {
		return nil, fmt.Errorf("claiming slots: %w", err)
	}

	// Initiate payment
	stripe.Key = "sk_test_4eC39HqLyjWDarjtT1zdp7dc"

	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(claims))
	usd := string(stripe.CurrencyUSD)
	for i := range claims {
		cost := int64(claims[i].ConferenceSlot.Cost)
		//FIXME apply voucher
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

	stripeParams := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String("https://example.com/success"), // FIXME triky one
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

	return &CheckoutResponse{StripeSessionID: s.ID}, nil
}
