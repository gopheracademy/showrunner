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
	Events []Event
}

// GetAll retrieves all conferences and events
// encore:api public
func GetAll(ctx context.Context, params *GetAllParams) (*GetAllResponse, error) {

	rows, err := sqldb.Query(ctx,
		`SELECT event.*, conference.id, conference.name, conference.slug, conference.start_date, conference.end_date, conference.location FROM event LEFT JOIN conference ON conference.event_id = event.id
		`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all conferences: %w", err)
	}

	defer rows.Close()

	idToEvent := map[uint32]*Event{}

	for rows.Next() {
		var event Event
		var conference Conference

		conference.Slots = []ConferenceSlot{}

		err := rows.Scan(&event.ID, &event.Name, &event.Slug, &conference.ID, &conference.Name, &conference.Slug, &conference.StartDate, &conference.EndDate, &conference.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		if existingEvent, ok := idToEvent[event.ID]; ok {
			existingEvent.Conferences = append(existingEvent.Conferences, conference)
		} else {
			event.Conferences = append(event.Conferences, conference)
			idToEvent[event.ID] = &event
		}
	}

	events := []Event{}

	for _, event := range idToEvent {
		events = append(events, *event)
	}

	return &GetAllResponse{
		Events: events,
	}, nil
}
