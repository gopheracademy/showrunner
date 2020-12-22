package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// GetAnonPaperParams defines the inputs used by the GetAnonPaper API method
type GetAnonPaperParams struct {
	PaperID uint32
}

// GetAnonPaperResponse defines the output returned by the GetAnonPaper API method
type GetAnonPaperResponse struct {
	AnonPaper AnonPaper
}

// GetAnonPaper retrieves information for a
// specific paper id without identifying
// user info
// encore:api public
func GetAnonPaper(ctx context.Context, params *GetAnonPaperParams) (*GetAnonPaperResponse, error) {

	row := sqldb.QueryRow(
		ctx,
		`
		SELECT id,
		conference_id,
		title,
		elevator_pitch,
		description,
		notes
		FROM paper_submission
		WHERE id = $1
		`, params.PaperID,
	)

	var anonPaper AnonPaper

	err := row.Scan(
		&anonPaper.ID,
		&anonPaper.ConferenceID,
		&anonPaper.Title,
		&anonPaper.ElevatorPitch,
		&anonPaper.Description,
		&anonPaper.Notes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve paper: %w", err)
	}

	return &GetAnonPaperResponse{AnonPaper: anonPaper}, nil

}
