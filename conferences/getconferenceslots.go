package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// GetConferenceSlotsParams defines the inputs used by the GetConferenceSlots API method
type GetConferenceSlotsParams struct {
	ConferenceID int
}

// GetConferenceSlotsResponse defines the output returned by the GetConferenceSlots API method
type GetConferenceSlotsResponse struct {
	ConferenceSlots []ConferenceSlot
}

// GetConferenceSlots retrieves all event slots for a specific event id
// encore:api public
func GetConferenceSlots(ctx context.Context, params *GetConferenceSlotsParams) (*GetConferenceSlotsResponse, error) {

	rows, err := sqldb.Query(ctx,
		`SELECT id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public FROM conference_slot WHERE conference_id = $1
		`, params.ConferenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all conferences: %w", err)
	}

	defer rows.Close()

	conferenceSlots := []ConferenceSlot{}

	for rows.Next() {
		var conferenceSlot ConferenceSlot

		err := rows.Scan(&conferenceSlot.ID, &conferenceSlot.Name, &conferenceSlot.Description, &conferenceSlot.Cost, &conferenceSlot.Capacity, &conferenceSlot.StartDate, &conferenceSlot.EndDate, &conferenceSlot.PurchaseableFrom, &conferenceSlot.PurchaseableUntil, &conferenceSlot.AvailableToPublic)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		conferenceSlots = append(conferenceSlots, conferenceSlot)

	}

	return &GetConferenceSlotsResponse{
		ConferenceSlots: conferenceSlots,
	}, nil
}
