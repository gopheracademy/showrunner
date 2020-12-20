package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// AddPaperParams defines the inputs used by the AddPaper API method
type AddPaperParams struct {
	Paper *Paper
}

// AddPaperResponse defines the output returned by the AddPaper API method
type AddPaperResponse struct {
	PaperID uint32
}

// AddPaper inserts a paper into the paper_submissions table
// encore:api public
func AddPaper(ctx context.Context, params *AddPaperParams) (*AddPaperResponse, error) {

	row := sqldb.QueryRow(ctx,
		`INSERT INTO paper_submission (
			user_id,
  		conference_id,
  		title,
  		elevator_pitch,
  		description,
  		notes
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
			) RETURNING id`,
		params.Paper.UserID,
		params.Paper.ConferenceID,
		params.Paper.Title,
		params.Paper.ElevatorPitch,
		params.Paper.Description,
		params.Paper.Notes,
	)

	var paperID uint32
	err := row.Scan(
		&paperID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add paper: %w", err)
	}

	return &AddPaperResponse{PaperID: paperID}, nil
}
