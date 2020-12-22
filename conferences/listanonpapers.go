package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// ListAnonPapersParams defines the inputs used by the ListAnonPapers API method
type ListAnonPapersParams struct {
	ConferenceID uint32
}

// ListAnonPapersResponse defines the output returned by the ListAnonPapers API method
type ListAnonPapersResponse struct {
	AnonPapers []AnonPaper
}

// ListAnonPapers retrieves all the papers
// submitted for a specific conference without
// user identification information
// encore:api public
func ListAnonPapers(ctx context.Context, params *ListAnonPapersParams) (*ListAnonPapersResponse, error) {

	rows, err := sqldb.Query(ctx,
		` SELECT id,
			conference_id,
			title,
			elevator_pitch,
			description,
			notes
			FROM paper_submission
			ORDER BY id
`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all papers: %w", err)
	}

	defer rows.Close()

	var anonPapers []AnonPaper

	for rows.Next() {

		var anonPaper AnonPaper

		err := rows.Scan(
			&anonPaper.ID,
			&anonPaper.ConferenceID,
			&anonPaper.Title,
			&anonPaper.ElevatorPitch,
			&anonPaper.Description,
			&anonPaper.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		anonPapers = append(anonPapers, anonPaper)
	}

	return &ListAnonPapersResponse{AnonPapers: anonPapers}, nil
}
