package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// GetAllParams defines the inputs used by the GetAll API method
type GetAllParams struct {
}

// GetAllResponse defines the output returned by the GetAll API method
type GetAllResponse struct {
	Conferences []Conference
}

// GetAll retrieves all conferences and events
// encore:api public
func GetAll(ctx context.Context, params *GetAllParams) (*GetAllResponse, error) {

	rows, err := sqldb.Query(ctx,
		`SELECT conference.*, event.id, event.name, event.slug, event.start_date, event.end_date, event.location FROM conference LEFT JOIN event ON event.conference_id = conference.id
		`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all conferences: %w", err)
	}

	defer rows.Close()

	idToConference := map[uint32]*Conference{}

	for rows.Next() {
		var conference Conference
		var event Event

		err := rows.Scan(&conference.ID, &conference.Name, &conference.Slug, &event.ID, &event.Name, &event.Slug, &event.StartDate, &event.EndDate, &event.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		if existingConference, ok := idToConference[conference.ID]; ok {
			existingConference.Events = append(existingConference.Events, event)
		} else {
			conference.Events = append(conference.Events, event)
			idToConference[conference.ID] = &conference
		}
	}

	var conferences []Conference

	for _, conference := range idToConference {
		conferences = append(conferences, *conference)
	}

	return &GetAllResponse{
		Conferences: conferences,
	}, nil
}
