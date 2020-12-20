package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// GetPaperParams defines the inputs used by the GetPaper API method
type GetPaperParams struct {
	PaperID uint32
}

// GetPaperResponse defines the output returned by the GetPaper API method
type GetPaperResponse struct {
	Paper Paper
}

// GetPaper retrieves information for a specific paper id
// encore:api public
func GetPaper(ctx context.Context, params *GetPaperParams) (*GetPaperResponse, error) {

	row := sqldb.QueryRow(
		ctx,
		`
		SELECT id,
		user_id,
		conference_id,
		title,
		elevator_pitch,
		description,
		notes
		FROM paper_submission
		WHERE id = $1
		`, params.PaperID,
	)

	var paper Paper

	err := row.Scan(
		&paper.ID,
		&paper.UserID,
		&paper.ConferenceID,
		&paper.Title,
		&paper.ElevatorPitch,
		&paper.Description,
		&paper.Notes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve paper: %w", err)
	}

	return &GetPaperResponse{Paper: paper}, nil

}
