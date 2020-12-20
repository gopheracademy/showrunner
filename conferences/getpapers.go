package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// ListPaperParams defines the inputs used by the ListPapers API method
type ListPapersParams struct {
	ConferenceID uint32
}

// ListPapersResponse defines the output return by hte ListPapers API method
type ListPapersResponse struct {
	Papers []Paper
}

// ListPapers retrieves all the papers submitted for a specific conference
// encore:api public
func ListPapers(ctx context.Context, params *ListPapersParams) (*ListPapersResponse, error) {

	rows, err := sqldb.Query(ctx,
		` SELECT id,
			user_id,
			conference_id,
			title,
			elevator_pitch,
			description,
			notes
			FROM paper_submission
`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all papers: %w", err)
	}

	defer rows.Close()

	var papers []Paper

	for rows.Next() {

		var paper Paper

		err := rows.Scan(
			&paper.ID,
			&paper.UserID,
			&paper.ConferenceID,
			&paper.Title,
			&paper.ElevatorPitch,
			&paper.Description,
			&paper.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		papers = append(papers, paper)
	}

	return &ListPapersResponse{Papers: papers}, nil
}
