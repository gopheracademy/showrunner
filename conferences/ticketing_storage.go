package conferences

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"encore.dev/storage/sqldb"
)

// SQLStorage provides a Postgres Flavored storage backend to store ticketing information.
type SQLStorage struct {
	tx *sqldb.Tx
}

const (
	ticketIDUniqueConstraint = "ticket_id_is_unique"
	tableSlotClaims          = "slot_claim"
)

func query(ctx context.Context, tx *sqldb.Tx, statement string, args ...interface{}) (*sqldb.Rows, error) {
	if tx != nil {
		return sqldb.QueryTx(tx, ctx, statement, args...)
	}
	return sqldb.Query(ctx, statement, args...)
}

func queryRow(ctx context.Context, tx *sqldb.Tx, statement string, args ...interface{}) *sqldb.Row {
	if tx != nil {
		return sqldb.QueryRowTx(tx, ctx, statement, args...)
	}
	return sqldb.QueryRow(ctx, statement, args...)
}

func exec(ctx context.Context, tx *sqldb.Tx, statement string, args ...interface{}) (sql.Result, error) {
	if tx != nil {
		return sqldb.ExecTx(tx, ctx, statement, args...)
	}
	return sqldb.Exec(ctx, statement, args...)
}

// createAttendee creates a new attendee in the database and returns it.
func createAttendee(ctx context.Context, tx *sqldb.Tx, a *Attendee) (*Attendee, error) {
	result := Attendee{}
	row := queryRow(ctx, tx,
		"INSERT INTO attendee (email, coc_accepted) VALUES ($1, $2) RETURNING id, email, coc_accepted",
		a.Email, a.CoCAccepted)

	if err := row.Scan(&result.ID, &result.Email, &result.CoCAccepted); err != nil {
		return nil, fmt.Errorf("creating or fetching attendee: %w", err)
	}

	for _, c := range a.Claims {
		res, err := exec(ctx, tx,
			`INSERT INTO attendee_to_slot_claims (attendee_id, slot_claim_id) 
			VALUES ($1, $2) 
		ON CONFLICT ON CONSTRAINT slot_claim_id_is_unique DO UPDATE SET attendee_id = EXCLUDED.attendee_id`,
			result.ID, c.ID)
		if err != nil {
			return nil, fmt.Errorf("inserting attendee claims: %w", err)
		}

		ra, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("inserting attendee claims: %w", err)
		}
		if ra == 0 {
			return nil, fmt.Errorf("could not create claim")
		}
	}

	result.Claims = a.Claims
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
func (s *SQLStorage) readAttendeeByID(ctx context.Context, tx *sqldb.Tx, id uint64) (*Attendee, error) {
	if id == 0 {
		return nil, fmt.Errorf("id is not valid")
	}
	return readAttendee(ctx, tx, "", id)
}

func readAttendee(ctx context.Context, tx *sqldb.Tx, email string, id uint64) (*Attendee, error) {
	results := Attendee{}
	q := `SELECT id, email, coc_accepted FROM attendee`
	args := []interface{}{}
	if email != "" {
		q = `SELECT id, email, coc_accepted FROM attendee WHERE email = $1`
		args = append(args, email)
	}
	if email != "" && id != 0 {
		q = `SELECT id, email, coc_accepted FROM attendee WHERE email = $1 AND WHERE id = $2`
		args = append(args, id)
	}
	if email == "" && id != 0 {
		q = `SELECT id, email, coc_accepted FROM attendee WHERE id = $1`
		args = append(args, id)
	}
	row := queryRow(ctx, tx, q, args...)
	err := row.Scan(&results.ID, &results.Email, &results.CoCAccepted)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading attendee by email: %w", err)
	}

	claims := []SlotClaim{}

	rows, err := query(ctx, tx,
		`SELECT id, ticket_id, redeemed FROM slot_claim 
	JOIN attendees_to_slot_claims ON slot_claim.id = attendees_to_slot_claims.slot_claim_id
	WHERE attendees_to_slot_claims.id = $1`,
		results.ID)
	if err != nil {
		return nil, fmt.Errorf("querying claims for attendee: %w", err)
	}
	for rows.Next() {
		claim := SlotClaim{}
		err := rows.Scan(&claim.ID, &claim.TicketID, &claim.Redeemed)
		if err != nil {
			return nil, fmt.Errorf("scanning slot_claims for attendee: %w", err)
		}
		claims = append(claims, claim)
	}

	results.Claims = claims
	return &results, nil
}

// createConferenceSlot saves a slot in the database.
func createConferenceSlot(ctx context.Context, tx *sqldb.Tx, cslot *ConferenceSlot, conferenceID int64) (*ConferenceSlot, error) {
	var sqlSentence = `INSERT INTO conference_slot (conference_id, name, description, cost, capacity, start_date, end_date, purchaseable_form, purchaseable_until, available_to_public)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING conference_id, name, description, cost, capacity, start_date, end_date, purchaseable_form, purchaseable_until, available_to_public`
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
	/*if e.DependsOn != nil {
		sqlSentence = `INSERT INTO conference_slot (conference_id, name, description, cost, capacity, start_date, end_date, purchaseable_form, purchaseable_until, available_to_public, depends_on_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, name, description, cost, capacity, start_date, end_date, purchaseable_form, purchaseable_until, available_to_public`
		args = append(args, e.DependsOn.ID)
	}*/
	row := queryRow(ctx, tx, sqlSentence, args...)
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
	row := queryRow(ctx, tx,
		`SELECT (id, name, description, cost, capacity, start_date, end_date, purchaseable_form, purchaseable_until, available_to_public, conference_id)
	FROM conference_slot 
	WHERE id = $1`, id)

	var conferenceID uint64

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
		&conferenceID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading conference slots by id: %w", err)
	}

	/*
			if dependsOnID != 0 && loadDeps {
				results.DependsOn, err = readConferenceSlotByID(ctx, tx, dependsOnID, false)
				if err != nil {
					return nil, fmt.Errorf("loading dependency: %w", err)
				}
			}


		conference := Conference{}
		row = queryRow(ctx, tx,
			`SELECT id, name, slug, start_date, end_date, location FROM conference WHERE id = $1`, conferenceID)
		err = row.Scan(&conference.ID, &conference.Name, &conference.Slug, &conference.StartDate, &conference.EndDate, &conference.Location)
		if err != nil {
			return nil, fmt.Errorf("reading conference by id: %w", err)
		}

		results.Conference = &conference
	*/
	return &results, nil
}

// updateConferenceSlot updates conference slot fields from the passed instance
func updateConferenceSlot(ctx context.Context, tx *sqldb.Tx, cslot *ConferenceSlot, conferenceID int64) error {
	var sqlStatement = `UPDATE conference_slot 
	SET conference_id = $1, name = $2, description = $3, cost =$4, capacity=$5, 
	start_date = $6, end_date = $7, purchaseable_form = $8, purchaseable_until = $9, 
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
	/* if e.DependsOn != nil {
		sqlStatement = `UPDATE conference_slot
	SET conference_id = $1, name = $2, description = $3, cost =$4, capacity=$5,
	start_date = $6, end_date = $7, purchaseable_form = $8, purchaseable_until = $9,
	available_to_public $10, depends_on_id = $11
	WHERE id = $12`
		args = append(args, e.DependsOn.ID)
	} */
	args = append(args, cslot.ID)

	res, err := exec(ctx, tx, sqlStatement, args...)

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
func createSlotClaim(ctx context.Context, tx *sqldb.Tx, slotClaim *SlotClaim) (*SlotClaim, error) {
	var err error

	row := queryRow(ctx, tx,
		`INSERT INTO slot_claim (ticket_id, redeemed, conference_slot_id) VALUES ($1, $2, $3)
		RETURNING id, ticket_id, redeemed`,
		slotClaim.TicketID, slotClaim.Redeemed, slotClaim.ConferenceSlot.ID)

	results := SlotClaim{}
	err = row.Scan(&results.ID, &results.TicketID, &results.Redeemed)

	if err != nil {
		return nil, fmt.Errorf("saving slot claim: %w", err)
	}

	results.ConferenceSlot = slotClaim.ConferenceSlot
	return &results, nil
}

const (
	tableAttendee               = "attendee"
	tableAttendeeSlotClaims     = "attendee_to_slot_claims"
	slotClaimIDUniqueConstraint = "slot_claim_id_is_unique"
)

// updateAttendee saves the passed attendee attributes on top of the existing one.
func updateAttendee(ctx context.Context, tx *sqldb.Tx, attendee *Attendee) (*Attendee, error) {
	res, err := exec(ctx, tx,
		`UPDATE attendee SET email = $1, coc_accepted = $2 WHERE id = $3`,
		attendee.Email, attendee.CoCAccepted, attendee.ID)

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

	for _, c := range attendee.Claims {

		res, err = exec(ctx, tx,
			`INSERT INTO attendee_to_slot_claims (attendee_id, slot_claim_id) 
			VALUES ($1, $2) 
		ON CONFLICT ON CONSTRAINT slot_claim_id_is_unique DO UPDATE SET attendee_id = $3`,
			attendee.ID, c.ID, attendee.ID)
		if err != nil {
			return nil, fmt.Errorf("inserting attendee claims: %w", err)
		}
		ra, err := res.RowsAffected()
		if err != nil {
			return nil, fmt.Errorf("failed to find rows affected by update: %w", err)
		}
		if ra == 0 {
			return nil, fmt.Errorf("no such attedee")
		}
	}

	return attendee, nil
}

const (
	tableClaimPayment                = "claim_payment"
	tableFinancialInstrumentMoney    = "payment_method_money"
	tableMoneyToPayment              = "payment_method_money_to_claim_payment"
	tableFinancialInstrumentDiscount = "payment_method_conference_slotdiscount"
	tableDiscountToPayment           = "payment_method_conference_slotdiscount_to_claim_payment"
	tableFinancialInstrumentCredit   = "payment_method_credit_note"
	tableCreditToPayment             = "payment_method_credit_note_to_claim_payment"
)

func insertMoneyPayment(ctx context.Context, tx *sqldb.Tx, claimPaymentID uint64, payment *PaymentMethodMoney) (*PaymentMethodMoney, error) {
	if payment.ID != 0 { // SERIAL starts in 1
		return payment, nil
	}
	money := PaymentMethodMoney{}

	row := queryRow(ctx, tx,
		`INSERT INTO payment_method_money (amount, ref) VALUES ($1, $2)
		RETURNING id, amount, ref`,
		payment.Amount, payment.PaymentRef)
	err := row.Scan(&money.ID, &money.Amount, &money.PaymentRef)
	if err != nil {
		return nil, fmt.Errorf("inserting money payment: %w", err)
	}

	res, err := exec(ctx, tx,
		`INSERT INTO payment_method_money_to_claim_payment (payment_method_money_id, claim_payment_id) VALUES ($1, $2)`,
		money.ID, claimPaymentID)

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
		return payment, nil
	}
	discount := PaymentMethodConferenceDiscount{}

	row := queryRow(ctx, tx,
		`INSERT INTO payment_method_conference_slotdiscount (amount, detail) VALUES ($1, $2)
		RETURNING id, amount, detail`,
		payment.Amount, payment.Detail)
	err := row.Scan(&discount.ID, &discount.Amount, &discount.Detail)
	if err != nil {
		return nil, fmt.Errorf("inserting discount payment: %w", err)
	}

	res, err := exec(ctx, tx,
		`INSERT INTO payment_method_conference_slotdiscount_to_claim_payment (payment_method_conference_slotdiscount_id, claim_payment_id) VALUES ($1, $2)`,
		discount.ID, claimPaymentID)

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
		return payment, nil
	}
	credit := PaymentMethodCreditNote{}

	row := queryRow(ctx, tx,
		`INSERT INTO payment_method_credit_note (amount, detail) VALUES ($1, $2)
		RETURNING id, amount, detail`,
		payment.Amount, payment.Detail)
	err := row.Scan(&credit.ID, &credit.Amount, &credit.Detail)
	if err != nil {
		return nil, fmt.Errorf("inserting money payment: %w", err)
	}

	res, err := exec(ctx, tx,
		`INSERT INTO payment_method_conference_slotdiscount_to_claim_payment (payment_method_credit_note_to_claim_payment, claim_payment_id) VALUES ($1, $2)`,
		credit.ID, claimPaymentID)

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
	row := queryRow(ctx, tx,
		`INSERT INTO claim_payment (invoice) VALUES ($1)
	RETURNING id, invoice`, c.Invoice)
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
	res, err := exec(ctx, tx,
		`UPDATE claim_payment SET invoice = $1
	WHERE id = $2`, c.Invoice)
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
			processedPayments[i], err = insertMoneyPayment(ctx, tx, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting money payment: %w", err)
			}
		case *PaymentMethodConferenceDiscount:
			processedPayments[i], err = insertDiscountPayment(ctx, tx, c.ID, payment)
			if err != nil {
				return nil, fmt.Errorf("inserting discount payment: %w", err)
			}
		case *PaymentMethodCreditNote:
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

const attendeesToSlotTable = "attendees_to_slot_claims"

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

	claimIDs := make([]interface{}, 0, len(slots))
	claimIDsIndex := map[uint64]bool{}
	idHolders := make([]string, 0, len(slots))
	for i, slot := range slots {
		if slot.ID == 0 {
			return nil, nil, fmt.Errorf("some slot claims lack IDs, perhaps they have not been saved yet")
		}
		claimIDs = append(claimIDs, slot.ID)
		claimIDsIndex[slot.ID] = true
		idHolders = append(idHolders, strconv.Itoa(i+3))
	}
	args := append([]interface{}{target.ID, source.ID}, claimIDs...)
	res, err := exec(ctx, tx,
		fmt.Sprintf(`UPDATE attendees_to_slot_claims SET attendee_id = $1 WHERE attendee_id = $2 AND slot_claim_id IN (%s)`, strings.Join(idHolders, ",")),
		args...)

	if err != nil {
		return nil, nil, fmt.Errorf("changing slot claims ownershio: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return nil, nil, fmt.Errorf("changing slot claims ownershio: %w", err)
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
