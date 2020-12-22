package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// UpdatePaperParams defines the inputs used by the GetPaper API method
type UpdatePaperParams struct {
	Paper *Paper
}

// UpdatePaperResponse defines the output received by the UpdatePaper API method
type UpdatePaperResponse struct {
	Paper Paper
}

// UpdatePaper updates a paper submission for a specific paper id
// encore:api public
func UpdatePaper(ctx context.Context, params *UpdatePaperParams) (*UpdatePaperResponse, error) {

	row := sqldb.QueryRow(ctx,
		`
		UPDATE paper_submission
		SET title = $1,
			elevator_pitch = $2,
			description = $3,
			notes = $4
		WHERE id = $5
		RETURNING id,
			user_id,
			conference_id,
			title,
			elevator_pitch,
			description,
			notes
	`,
		params.Paper.Title,
		params.Paper.ElevatorPitch,
		params.Paper.Description,
		params.Paper.Notes,
		params.Paper.ID,
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
		return nil, fmt.Errorf("failed to update paper submission: %w", err)
	}

	return &UpdatePaperResponse{Paper: paper}, nil
}
