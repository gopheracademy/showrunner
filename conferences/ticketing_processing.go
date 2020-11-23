package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
	uuid "github.com/satori/go.uuid"
)

// claimSlots claims N slots for an attendee.
func claimSlots(ctx context.Context, attendee *Attendee, slots ...ConferenceSlot) ([]SlotClaim, error) {
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	var claims = make([]SlotClaim, 0, len(slots))
	for i := range slots {
		slot := slots[i]
		sc := &SlotClaim{
			ConferenceSlot: &slot,
			TicketID:       uuid.NewV4().String(),
		}
		sc, err = createSlotClaim(ctx, tx, sc)
		if err != nil {
			if atomicErr := sqldb.Rollback(tx); atomicErr != nil {
				err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
			}
			return nil, fmt.Errorf("Claiming a slot: %w", err)
		}
		claims[i] = *sc
	}
	attendee.Claims = append(attendee.Claims, claims...)
	_, err = updateAttendee(ctx, tx, attendee)
	if err != nil {
		if atomicErr := sqldb.Rollback(tx); atomicErr != nil {
			err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
		}
		return nil, fmt.Errorf("Updating claimed slots for attendee: %w", err)
	}
	if err := sqldb.Commit(tx); err != nil {
		return nil, fmt.Errorf("confirming atomic operation: %w", err)
	}
	return claims, nil
}

// payClaims assigns payments and/or credits to a set of claims.
func payClaims(ctx context.Context, attendee *Attendee, claims []SlotClaim,
	payments []FinancialInstrument) (*ClaimPayment, error) {
	ptrClaims := make([]*SlotClaim, len(claims))
	for i := range claims {
		ptrClaims[i] = &claims[i]
	}
	claimPayment := &ClaimPayment{
		ClaimsPayed: ptrClaims,
		Payment:     payments,
	}
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("beginning atomic operation: %w", err)
	}

	claimPayment, err = createClaimPayment(ctx, tx, claimPayment)
	if err != nil {
		if atomicErr := sqldb.Rollback(tx); atomicErr != nil {
			err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
		}
		return nil, fmt.Errorf("paying for claims: %w", err)
	}
	if err := sqldb.Commit(tx); err != nil {
		return nil, fmt.Errorf("confirming atomic operation: %w", err)
	}

	return claimPayment, nil
}

// ErrInvalidCurrency should be returned when paying with the wrong kind of instrument
// for instance covering credit with credit.
type ErrInvalidCurrency struct {
	currencyType AssetType
}

func (e *ErrInvalidCurrency) Error() string {
	return fmt.Sprintf("the debt cannot be covered with %s", e.currencyType)
}

// CoverCredit adds funds to a payment to cover for receivables.
func coverCredit(ctx context.Context,
	existingPayment *ClaimPayment,
	payments []FinancialInstrument) error {
	for i := range payments {
		if payments[i].Type() == ATReceivable {
			return &ErrInvalidCurrency{currencyType: payments[i].Type()}
		}
		existingPayment.Payment = append(existingPayment.Payment, payments[i])
	}
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning atomic operation: %w", err)
	}
	_, err = updateClaimPayment(ctx, tx, existingPayment)
	if err != nil {
		if atomicErr := sqldb.Rollback(tx); atomicErr != nil {
			err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
		}

		return fmt.Errorf("saving new payments %w", err)
	}

	if err := sqldb.Commit(tx); err != nil {
		return fmt.Errorf("confirming atomic operation: %w", err)
	}

	return nil
}

// transferClaims transfer claims from one user the the other, assuming they belong to the first.
func transferClaims(ctx context.Context,
	source, target *Attendee, claims []SlotClaim) (*Attendee, *Attendee, error) {
	var err error
	sourceClaimsMap := map[uint64]bool{}
	for i := range source.Claims {
		sourceClaimsMap[source.Claims[i].ID] = true
	}
	for i := range claims {
		if belongsToSource := sourceClaimsMap[claims[i].ID]; !belongsToSource {
			return nil, nil, fmt.Errorf("%d claim for slot %s does not belong to %s", claims[i].ID, claims[i].ConferenceSlot.Name, source.Email)
		}
	}
	tx, err := sqldb.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("beginning atomic operation: %w", err)
	}
	if source, target, err = changeSlotClaimOwner(ctx, tx, claims, source, target); err != nil {
		if atomicErr := sqldb.Rollback(tx); atomicErr != nil {
			err = fmt.Errorf("%w (also cancelling atomic operation: %v)", err, atomicErr)
		}
		return nil, nil, fmt.Errorf("reowning slot claim: %w", err)
	}
	if err := sqldb.Commit(tx); err != nil {
		return nil, nil, fmt.Errorf("confirming atomic operation: %w", err)
	}

	return source, target, nil
}
