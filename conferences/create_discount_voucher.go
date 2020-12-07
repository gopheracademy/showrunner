package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
	"github.com/gofrs/uuid"
)

// CreateDiscountVoucherParams defines payload accepted by the CreateDiscountVoucher endpoint.
type CreateDiscountVoucherParams struct {
	VoucherInformation *VoucherInformation
}

// CreateDiscountVoucherResponse defines response by the CreateDiscountVoucher endpoint.
type CreateDiscountVoucherResponse struct {
	VoucherID string
}

// CreateDiscountVoucher allows creation of a discount voucher for a given conference.
func CreateDiscountVoucher(ctx context.Context, params *CreateDiscountVoucherParams) (*CreateDiscountVoucherResponse, error) {
	if params.VoucherInformation == nil {
		return nil, fmt.Errorf("voucher information is required")
	}

	if params.VoucherInformation.Percentage > 0 && params.VoucherInformation.AmountInCents > 0 {
		return nil, fmt.Errorf("discount percentage and amount are mutually exclussive")
	}

	voucherID, err := uuid.DefaultGenerator.NewV4()
	if err != nil {
		return nil, fmt.Errorf("generating new voucher ID: %w", err)
	}
	if params.VoucherInformation.Percentage > 0 {
		_, err = sqldb.Exec(ctx,
			`INSERT INTO discount_vouchers 
	(voucher_id, 
		valid_from, 
		valid_to, 
		discount_percentage, 
		discount_percentage_max_amount_cents, 
		conference_id) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
			voucherID,
			params.VoucherInformation.ValidFrom,
			params.VoucherInformation.ValidTo,
			params.VoucherInformation.Percentage,
			params.VoucherInformation.LimitInCents,
			params.VoucherInformation.ConferenceID)
		if err != nil {
			return nil, fmt.Errorf("inserting percentage discount voucher into db: %w", err)
		}
		return &CreateDiscountVoucherResponse{
			VoucherID: voucherID.String(),
		}, nil
	}
	_, err = sqldb.Exec(ctx,
		`INSERT INTO discount_vouchers 
(voucher_id, 
	valid_from, 
	valid_to, 
	discount_percentage, 
	conference_id) 
VALUES ($1, $2, $3, $4, $5)`,
		voucherID,
		params.VoucherInformation.ValidFrom,
		params.VoucherInformation.ValidTo,
		params.VoucherInformation.Percentage,
		params.VoucherInformation.ConferenceID)
	if err != nil {
		return nil, fmt.Errorf("inserting fixed amount voucher into db: %w", err)
	}

	return &CreateDiscountVoucherResponse{
		VoucherID: voucherID.String(),
	}, nil
}
