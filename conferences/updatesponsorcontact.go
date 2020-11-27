package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// UpdateSponsorContactParams defines the inputs used by the UpdateSponsorContactParams API method
type UpdateSponsorContactParams struct {
	SponsorContactInformation *SponsorContactInformation
}

// UpdateSponsorContactResponse defines the output returned by the UpdateSponsorContactResponse API method
type UpdateSponsorContactResponse struct {
}

// UpdateSponsorContact retrieves all conferences and events
// encore:api public
func UpdateSponsorContact(ctx context.Context, params *UpdateSponsorContactParams) (*UpdateSponsorContactResponse, error) {

	if params.SponsorContactInformation == nil {
		return nil, fmt.Errorf("SponsorContactInformation is required")
	}

	if int(params.SponsorContactInformation.Role) > len(contactRoleMappings)-1 || params.SponsorContactInformation.Role < 0 {
		return nil, fmt.Errorf("invalid role provided")
	}

	result, err := sqldb.Exec(ctx, `
	UPDATE sponsor_contact_information
	SET name = $1,
		role = $2,
		email = $3,
		phone = $4
	WHERE id = $5`, params.SponsorContactInformation.Name, params.SponsorContactInformation.Role.String(), params.SponsorContactInformation.Email, params.SponsorContactInformation.Phone, params.SponsorContactInformation.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update sponsor contact information: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no such contact found")
	}

	return &UpdateSponsorContactResponse{}, nil
}
