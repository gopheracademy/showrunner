package conferences

import (
	"context"
	"testing"
	"time"

	"encore.dev/storage/sqldb"
)

const validFromDateTimestamp = 1445444940 // October 21 2015 4:28 pm

func TestCreateDiscountVoucher(t *testing.T) {
	ctx := context.Background()
	_, err := CreateDiscountVoucher(ctx,
		&CreateDiscountVoucherParams{
			VoucherInformation: &VoucherInformation{
				Percentage:    10,
				LimitInCents:  1000,
				AmountInCents: 10,
				ValidFrom:     time.Unix(validFromDateTimestamp, 0),
				ValidTo:       time.Unix(validFromDateTimestamp, 0).Add(24 * time.Hour),
				ConferenceID:  1,
			},
		})
	if err == nil {
		t.Fatal("expected an error due to amount and percentage specified")
	}

	validFromExpected := time.Unix(validFromDateTimestamp, 0)
	validToExpected := time.Unix(validFromDateTimestamp, 0).Add(24 * time.Hour)
	for _, tcase := range []VoucherInformation{
		{
			Percentage:   10,
			LimitInCents: 1000,
			ValidFrom:    validFromExpected,
			ValidTo:      validToExpected,
			ConferenceID: 1,
		},
		{
			AmountInCents: 2000,
			ValidFrom:     validFromExpected,
			ValidTo:       validToExpected,
			ConferenceID:  1,
		},
	} {
		result, err := CreateDiscountVoucher(ctx,
			&CreateDiscountVoucherParams{
				VoucherInformation: &tcase,
			})
		if err != nil {
			t.Fatalf("expected no error but got: %v", err)
		}
		row := sqldb.QueryRow(ctx,
			`SELECT discount_percentage, 
	discount_percentage_max_amount_cents, 
	discount_amount_cents, 
	valid_from, valid_to, 
	conference_id, 
	spent
	FROM discount_vouchers
	WHERE voucher_id = $1`, result.VoucherID)
		var spent bool
		var validFrom, validTo time.Time
		var discountAmountCents, discountMaxAmountCents int64
		var conferenceID, discountPercentage int
		if err := row.Scan(&discountPercentage,
			&discountMaxAmountCents,
			&discountAmountCents, &validFrom, &validTo, &conferenceID, &spent,
		); err != nil {
			t.Fatalf("could not scan discount back: %v", err)
		}
		if spent {
			t.Fatal("voucher is spent, should not be")
		}
		if discountAmountCents != tcase.AmountInCents {
			t.Fatalf("expected discount amount to be %d cents it is %d", tcase.AmountInCents, discountAmountCents)
		}
		if discountPercentage != tcase.Percentage {
			t.Fatalf("expected discount to be %d%% it is %d%%", tcase.Percentage, discountPercentage)
		}
		if discountMaxAmountCents != tcase.LimitInCents {
			t.Fatalf("discount should be max %d cents it is %d", tcase.LimitInCents, discountMaxAmountCents)
		}
		if validFrom.Unix() != tcase.ValidFrom.Unix() {
			t.Fatalf("voucher should be valid until %v but it is until %v", tcase.ValidFrom, validFrom)
		}
		if validTo.Unix() != tcase.ValidTo.Unix() {
			t.Fatalf("voucher should be valid until %v but it is until %v", tcase.ValidTo, validTo)
		}
		if conferenceID != tcase.ConferenceID {
			t.Fatalf("expected conference ID %d got %d", conferenceID, tcase.ConferenceID)
		}
	}
}
