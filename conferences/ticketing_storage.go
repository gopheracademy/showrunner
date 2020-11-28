package conferences

import (
	"context"
	"database/sql"
	"fmt"

	"encore.dev/storage/sqldb"
	"github.com/lib/pq"
)

// createAttendee creates a new attendee in the database and returns it.
func createAttendee(ctx context.Context, tx *sqldb.Tx, a *Attendee) (*Attendee, error) {
	result := Attendee{}
	sqlStatement := "INSERT INTO attendee (email, coc_accepted) VALUES ($1, $2) RETURNING id, email, coc_accepted"
	sqlArgs := []interface{}{a.Email, a.CoCAccepted}
	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}
	if err := row.Scan(&result.ID, &result.Email, &result.CoCAccepted); err != nil {
		return nil, fmt.Errorf("creating or fetching attendee: %w", err)
	}

	return &result, nil
}

// readAttendeeByEmail returns an attendee for that email if one exists.
func readAttendeeByEmail(ctx context.Context, tx *sqldb.Tx, email string) (*Attendee, error) {
	if email == "" {
		return nil, fmt.Errorf("email is empty")
	}
	return readAttendee(ctx, tx, email, 0)
}

// readAttendeeByID returns an attendee for the given ID if one exists.
func readAttendeeByID(ctx context.Context, tx *sqldb.Tx, id int64) (*Attendee, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is not valid")
	}
	return readAttendee(ctx, tx, "", id)
}

func readAttendee(ctx context.Context, tx *sqldb.Tx, email string, id int64) (*Attendee, error) {
	results := Attendee{}
	sqlStatement := `SELECT id, email, coc_accepted FROM attendee`
	sqlArgs := []interface{}{}
	if email != "" {
		sqlStatement = `SELECT id, email, coc_accepted FROM attendee WHERE email = $1`
		sqlArgs = append(sqlArgs, email)
	}
	if email != "" && id != 0 {
		sqlStatement = `SELECT id, email, coc_accepted FROM attendee WHERE email = $1 AND id = $2`
		sqlArgs = append(sqlArgs, id)
	}
	if email == "" && id != 0 {
		sqlStatement = `SELECT id, email, coc_accepted FROM attendee WHERE id = $1`
		sqlArgs = append(sqlArgs, id)
	}

	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}

	err := row.Scan(&results.ID, &results.Email, &results.CoCAccepted)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading attendee by email: %w", err)
	}

	claims := []SlotClaim{}

	sqlStatement = `SELECT id, ticket_id, redeemed FROM slot_claim 
	WHERE attendee_id = $1`
	sqlArgs = []interface{}{results.ID}
	var rows *sqldb.Rows

	if tx != nil {
		rows, err = sqldb.QueryTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		rows, err = sqldb.Query(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("querying claims for attendee: %w", err)
	}
	for rows.Next() {
		claim := SlotClaim{}
		err := rows.Scan(&claim.ID, &claim.TicketID, &claim.Redeemed)
		if err != nil {
			return nil, fmt.Errorf("scanning slot_claim for attendee: %w", err)
		}
		claims = append(claims, claim)
	}

	results.Claims = claims
	return &results, nil
}

// createConferenceSlot saves a slot in the database.
func createConferenceSlot(ctx context.Context, tx *sqldb.Tx, cslot *ConferenceSlot, conferenceID int64) (*ConferenceSlot, error) {
	var sqlStatement = `INSERT INTO conference_slot (conference_id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING conference_id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public`
	var sqlArgs = []interface{}{
		conferenceID,
		cslot.Name,
		cslot.Description,
		cslot.Cost,
		cslot.Capacity,
		cslot.StartDate,
		cslot.EndDate,
		cslot.PurchaseableFrom,
		cslot.PurchaseableUntil,
		cslot.AvailableToPublic,
	}
	if cslot.DependsOn != 0 {
		sqlStatement = `INSERT INTO conference_slot (conference_id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public, depends_on)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public`
		sqlArgs = append(sqlArgs, cslot.DependsOn)
	}
	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}

	results := ConferenceSlot{}
	err := row.Scan(&results.ID,
		&results.Name,
		&results.Description,
		&results.Cost,
		&results.Capacity,
		&results.StartDate,
		&results.EndDate,
		&results.PurchaseableFrom,
		&results.PurchaseableUntil,
		&results.AvailableToPublic)
	if err != nil {
		return nil, fmt.Errorf("creating new conference slot: %w", err)
	}

	// results.DependsOn = e.DependsOn

	return &results, nil
}

func readConferenceSlotByID(ctx context.Context, tx *sqldb.Tx, id uint64, loadDeps bool) (*ConferenceSlot, error) {
	results := ConferenceSlot{}
	var row *sqldb.Row
	sqlStatement := `SELECT id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public, COALESCE(depends_on, 0)
	FROM conference_slot 
	WHERE id = $1`
	sqlArgs := []interface{}{id}
	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}

	err := row.Scan(&results.ID,
		&results.Name,
		&results.Description,
		&results.Cost,
		&results.Capacity,
		&results.StartDate,
		&results.EndDate,
		&results.PurchaseableFrom,
		&results.PurchaseableUntil,
		&results.AvailableToPublic,
		&results.DependsOn)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading conference slots by id: %w", err)
	}

	return &results, nil
}

// updateConferenceSlot updates conference slot fields from the passed instance
func updateConferenceSlot(ctx context.Context, tx *sqldb.Tx, cslot *ConferenceSlot, conferenceID int64) error {
	// Regular update
	var sqlStatement = `UPDATE conference_slot 
	SET conference_id = $1, name = $2, description = $3, cost =$4, capacity=$5, 
	start_date = $6, end_date = $7, purchaseable_from = $8, purchaseable_until = $9, 
	available_to_public $10
	WHERE id = $11`
	var args = []interface{}{
		conferenceID,
		cslot.Name,
		cslot.Description,
		cslot.Cost,
		cslot.Capacity,
		cslot.StartDate,
		cslot.EndDate,
		cslot.PurchaseableFrom,
		cslot.PurchaseableUntil,
		cslot.AvailableToPublic,
	}

	// this slot depends on another
	if cslot.DependsOn != 0 {
		sqlStatement = `UPDATE conference_slot
	SET conference_id = $1, name = $2, description = $3, cost =$4, capacity=$5,
	start_date = $6, end_date = $7, purchaseable_from = $8, purchaseable_until = $9,
	available_to_public $10, depends_on = $11
	WHERE id = $12`
		args = append(args, cslot.DependsOn)
	}

	// the id is the last argument of the query, always
	args = append(args, cslot.ID)

	var res sql.Result
	var err error
	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, args...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, args...)
	}

	if err != nil {
		return fmt.Errorf("updating conference slot: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of rows affected by query: %w", err)
	}
	if ra == 0 {
		return fmt.Errorf("no such slot")
	}

	return nil
}

// createSlotClaim saves a slot claim and returns it with the populated ID
func createSlotClaim(ctx context.Context, tx *sqldb.Tx, slotClaim *SlotClaim, attendeeID int64) (*SlotClaim, error) {
	var err error

	// TODO: add a dynamic field that sets the user as waitlist if slot is claimed.
	sqlStatement := `INSERT INTO slot_claim (ticket_id, redeemed, conference_slot_id, attendee_id) VALUES ($1, $2, $3, $4)
		RETURNING id, ticket_id, redeemed`
	sqlArgs := []interface{}{slotClaim.TicketID, slotClaim.Redeemed, slotClaim.ConferenceSlot.ID, attendeeID}

	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}

	results := SlotClaim{}
	err = row.Scan(&results.ID, &results.TicketID, &results.Redeemed)

	if err != nil {
		return nil, fmt.Errorf("saving slot claim: %w", err)
	}

	results.ConferenceSlot = slotClaim.ConferenceSlot
	return &results, nil
}

// updateAttendee saves the passed attendee attributes on top of the existing one.
func updateAttendee(ctx context.Context, tx *sqldb.Tx, attendee *Attendee) (*Attendee, error) {
	sqlStatement := `UPDATE attendee SET email = $1, coc_accepted = $2 WHERE id = $3`
	sqlArgs := []interface{}{attendee.Email, attendee.CoCAccepted, attendee.ID}
	var res sql.Result
	var err error
	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("updating attendee: %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("updating attendee: %w", err)
	}
	if ra == 0 {
		return nil, fmt.Errorf("attendee was not updated")
	}

	claimIDs := make(pq.Int64Array, len(attendee.Claims))
	for i, c := range attendee.Claims {
		claimIDs[i] = c.ID
	}

	sqlStatement = `UPDATE slot_claim SET
	attendee_id = $1
	WHERE id = ANY($2)
	`
	sqlArgs = []interface{}{attendee.ID, claimIDs}

	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("inserting attendee claims: %w", err)
	}
	ra, err = res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to find rows affected by update: %w", err)
	}
	if ra == 0 {
		return nil, fmt.Errorf("no such attedee")
	}

	return attendee, nil
}

const (
	tableClaimPayment                = "claim_payment"
	tableFinancialInstrumentMoney    = "payment_method_money"
	tableMoneyToPayment              = "payment_method_money_to_claim_payment"
	tableFinancialInstrumentDiscount = "payment_method_conference_discount"
	tableDiscountToPayment           = "payment_method_conference_discount_to_claim_payment"
	tableFinancialInstrumentCredit   = "payment_method_credit_note"
	tableCreditToPayment             = "payment_method_credit_note_to_claim_payment"
)

func insertMoneyPayment(ctx context.Context, tx *sqldb.Tx, claimPaymentID uint64, payment *PaymentMethodMoney) (*PaymentMethodMoney, error) {
	if payment.ID != 0 { // SERIAL starts in 1
		return nil, fmt.Errorf("this money payment has already been inserted")
	}
	money := PaymentMethodMoney{}

	sqlStatement := `INSERT INTO payment_method_money (amount, ref) VALUES ($1, $2)
		RETURNING id, amount, ref`
	sqlArgs := []interface{}{payment.Amount, payment.PaymentRef}

	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}
	err := row.Scan(&money.ID, &money.Amount, &money.PaymentRef)
	if err != nil {
		return nil, fmt.Errorf("inserting money payment: %w", err)
	}

	sqlStatement = `INSERT INTO payment_method_money_to_claim_payment (payment_method_money_id, claim_payment_id) VALUES ($1, $2)`
	sqlArgs = []interface{}{money.ID, claimPaymentID}

	var res sql.Result

	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument money to payment: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument money to payment: %w", err)
	}
	if ra == 0 {
		return nil, fmt.Errorf("failed to claim payment")
	}
	return &money, nil
}

func insertDiscountPayment(ctx context.Context, tx *sqldb.Tx, claimPaymentID uint64, payment *PaymentMethodConferenceDiscount) (*PaymentMethodConferenceDiscount, error) {
	if payment.ID != 0 { // SERIAL starts in 1
		return nil, fmt.Errorf("this discount has already been inserted")
	}
	discount := PaymentMethodConferenceDiscount{}

	sqlStatement := `INSERT INTO payment_method_conference_discount (amount, detail) VALUES ($1, $2)
		RETURNING id, amount, detail`
	sqlArgs := []interface{}{payment.Amount, payment.Detail}
	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}
	err := row.Scan(&discount.ID, &discount.Amount, &discount.Detail)
	if err != nil {
		return nil, fmt.Errorf("inserting discount payment: %w", err)
	}

	sqlStatement = `INSERT INTO payment_method_conference_discount_to_claim_payment (payment_method_conference_discount_id, claim_payment_id) VALUES ($1, $2)`
	sqlArgs = []interface{}{discount.ID, claimPaymentID}
	var res sql.Result

	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument discount to payment: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument discount to payment: %w", err)
	}
	if ra == 0 {
		return nil, fmt.Errorf("failed to claim payment")
	}
	return &discount, nil

}

func insertCreditPayment(ctx context.Context, tx *sqldb.Tx, claimPaymentID uint64, payment *PaymentMethodCreditNote) (*PaymentMethodCreditNote, error) {
	if payment.ID != 0 { // SERIAL starts in 1
		return nil, fmt.Errorf("this credit note has already been inserted")
	}
	credit := PaymentMethodCreditNote{}

	sqlStatement := `INSERT INTO payment_method_credit_note (amount, detail) VALUES ($1, $2)
		RETURNING id, amount, detail`
	sqlArgs := []interface{}{payment.Amount, payment.Detail}
	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}
	err := row.Scan(&credit.ID, &credit.Amount, &credit.Detail)
	if err != nil {
		return nil, fmt.Errorf("inserting money payment: %w", err)
	}

	sqlStatement = `INSERT INTO payment_method_credit_note_to_claim_payment (payment_method_credit_note_id, claim_payment_id) VALUES ($1, $2)`
	sqlArgs = []interface{}{credit.ID, claimPaymentID}
	var res sql.Result

	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument credit to payment: %w", err)
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("relating financial instrument credit to payment: %w", err)
	}
	if ra == 0 {
		return nil, fmt.Errorf("failed to claim payment")
	}
	return &credit, nil

}

func createClaimPayment(ctx context.Context, tx *sqldb.Tx, c *ClaimPayment) (*ClaimPayment, error) {
	claimPayments := ClaimPayment{}
	sqlStatement := `INSERT INTO claim_payment (invoice) VALUES ($1)
	RETURNING id, invoice`
	sqlArgs := []interface{}{c.Invoice}

	var row *sqldb.Row

	if tx != nil {
		row = sqldb.QueryRowTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		row = sqldb.QueryRow(ctx, sqlStatement, sqlArgs...)
	}
	err := row.Scan(&claimPayments.ID, &claimPayments.Invoice)
	if err != nil {
		return nil, fmt.Errorf("inserting payment for claims: %w", err)
	}

	processedPayments := make([]FinancialInstrument, len(c.Payment), len(c.Payment))

	for i, cp := range c.Payment {
		switch payment := cp.(type) {
		case *PaymentMethodMoney:
			processedPayments[i], err = insertMoneyPayment(ctx, tx, claimPayments.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting money payment: %w", err)
			}
		case *PaymentMethodConferenceDiscount:
			processedPayments[i], err = insertDiscountPayment(ctx, tx, claimPayments.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting discount payment: %w", err)
			}
		case *PaymentMethodCreditNote:
			processedPayments[i], err = insertCreditPayment(ctx, tx, claimPayments.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting credit payment: %w", err)
			}
		default:
			return nil, fmt.Errorf("not sure how to process payments of type %T", cp)
		}
	}
	claimPayments.ClaimsPaid = c.ClaimsPaid
	claimPayments.Payment = processedPayments
	return &claimPayments, nil
}

// updateClaimPayment saves the invoice and payments of this claim payment assuming it exists
func updateClaimPayment(ctx context.Context, tx *sqldb.Tx, c *ClaimPayment) (*ClaimPayment, error) {
	sqlStatement := `UPDATE claim_payment SET invoice = $1
	WHERE id = $2`
	sqlArgs := []interface{}{c.Invoice, c.ID}

	var res sql.Result
	var err error

	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, fmt.Errorf("updating payment for claims: %w", err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("updating payment for claims: %w", err)
	}
	if ra == 0 {
		return nil, fmt.Errorf("claim payment was not found")
	}

	processedPayments := make([]FinancialInstrument, len(c.Payment), len(c.Payment))

	for i, cp := range c.Payment {
		switch payment := cp.(type) {
		case *PaymentMethodMoney:
			if payment.ID != 0 {
				processedPayments[i] = payment
				continue
			}
			processedPayments[i], err = insertMoneyPayment(ctx, tx, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting money payment: %w", err)
			}
		case *PaymentMethodConferenceDiscount:
			if payment.ID != 0 {
				processedPayments[i] = payment
				continue
			}
			processedPayments[i], err = insertDiscountPayment(ctx, tx, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting discount payment: %w", err)
			}
		case *PaymentMethodCreditNote:
			if payment.ID != 0 {
				processedPayments[i] = payment
				continue
			}
			processedPayments[i], err = insertCreditPayment(ctx, tx, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting credit payment: %w", err)
			}
		default:
			return nil, fmt.Errorf("not sure how to process payments of type %T", cp)
		}
	}
	newClaim := ClaimPayment{
		ID:         c.ID,
		ClaimsPaid: c.ClaimsPaid,
		Payment:    processedPayments,
		Invoice:    c.Invoice,
	}
	return &newClaim, nil
}

// changeSlotClaimOwner changes the passed claims owner from source to target
func changeSlotClaimOwner(ctx context.Context, tx *sqldb.Tx, slots []SlotClaim, source *Attendee, target *Attendee) (*Attendee, *Attendee, error) {
	if source == nil || target == nil {
		return nil, nil, fmt.Errorf("either source or target is undefined")
	}

	if len(slots) == 0 {
		return nil, nil, fmt.Errorf("no slots to transfer")
	}

	if len(slots) > len(source.Claims) {
		return nil, nil, fmt.Errorf("the passed source lacks those claims")
	}

	claimIDs := make([]int64, 0, len(slots))
	claimIDsIndex := map[int64]bool{}
	for _, slot := range slots {
		if slot.ID == 0 {
			return nil, nil, fmt.Errorf("some slot claims lack IDs, perhaps they have not been saved yet")
		}
		claimIDs = append(claimIDs, slot.ID)
		claimIDsIndex[slot.ID] = true
	}

	sqlStatement := `UPDATE slot_claim SET attendee_id = $1 WHERE attendee_id = $2 AND id = ANY($3)`
	sqlArgs := []interface{}{target.ID, source.ID, pq.Int64Array(claimIDs)}

	var res sql.Result
	var err error

	if tx != nil {
		res, err = sqldb.ExecTx(tx, ctx, sqlStatement, sqlArgs...)
	} else {
		res, err = sqldb.Exec(ctx, sqlStatement, sqlArgs...)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("changing slot claims ownership: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, nil, fmt.Errorf("changing slot claims ownership: %w", err)
	}
	if int64(len(slots)) != affected {
		return nil, nil, fmt.Errorf("got %d claims to change but only changed %d", len(slots), affected)
	}

	newClaims := make([]SlotClaim, 0, len(source.Claims)-len(claimIDs))
	for _, claim := range source.Claims {
		if claimIDsIndex[claim.ID] {
			target.Claims = append(target.Claims, claim)
			continue
		}
		newClaims = append(newClaims, claim)
	}
	source.Claims = newClaims
	return source, target, nil
}
