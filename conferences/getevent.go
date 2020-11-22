package conferences

import (
	"context"
	"fmt"

	"encore.dev/storage/sqldb"
)

// GetEventSlotsParams defines the inputs used by the GetEventSlots API method
type GetEventSlotsParams struct {
	EventID int
}

// GetEventSlotsResponse defines the output returned by the GetEventSlots API method
type GetEventSlotsResponse struct {
	EventSlots []EventSlot
}

// GetEventSlots retrieves all event slots for a specific event id
// encore:api public
func GetEventSlots(ctx context.Context, params *GetEventSlotsParams) (*GetEventSlotsResponse, error) {

	rows, err := sqldb.Query(ctx,
		`SELECT id, name, description, cost, capacity, start_date, end_date, purchaseable_from, purchaseable_until, available_to_public FROM event_slot WHERE event_id = $1
		`, params.EventID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all conferences: %w", err)
	}

	defer rows.Close()

	var eventSlots []EventSlot

	for rows.Next() {
		var eventSlot EventSlot

		err := rows.Scan(&eventSlot.ID, &eventSlot.Name, &eventSlot.Description, &eventSlot.Cost, &eventSlot.Capacity, &eventSlot.StartDate, &eventSlot.EndDate, &eventSlot.PurchaseableFrom, &eventSlot.PurchaseableUntil, &eventSlot.AvailableToPublic)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		eventSlots = append(eventSlots, eventSlot)

	}

	return &GetEventSlotsResponse{
		EventSlots: eventSlots,
	}, nil
}
