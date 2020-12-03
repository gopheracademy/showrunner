package conferences

import (
	"context"
	"fmt"
	"log"

	"encore.dev/beta/auth"
	"encore.dev/storage/sqldb"
)

// GetConferenceSponsorsParams defines the inputs used by the GetConferenceSponsors API method
type GetConferenceSponsorsParams struct {
	ConferenceID uint32
}

// GetConferenceSponsorsResponse defines the output returned by the GetConferenceSponsors API method
type GetConferenceSponsorsResponse struct {
	Sponsors []Sponsor
}

// GetConferenceSponsors retrieves the sponsors for a specific conference
// encore:api auth
func GetConferenceSponsors(ctx context.Context, params *GetConferenceSponsorsParams) (*GetConferenceSponsorsResponse, error) {
	usr := auth.Data().(*Data)
	log.Println(usr.UserID)

	rows, err := sqldb.Query(ctx,
		`SELECT sponsor.id,
		 sponsor.name,
		 sponsor.sponsorship_level
		 FROM sponsor 
		 WHERE sponsor.conference_id = $1
		`, params.ConferenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve current conference: %w", err)
	}

	defer rows.Close()

	sponsors := []Sponsor{}

	for rows.Next() {
		var sponsor Sponsor

		err := rows.Scan(
			&sponsor.ID,
			&sponsor.Name,
			&sponsor.SponsorshipLevel,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}
		sponsors = append(sponsors, sponsor)

	}

	return &GetConferenceSponsorsResponse{
		Sponsors: sponsors,
	}, nil
}
