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
		`SELECT conference_slot.id,
		 conference_slot.name,
		 conference_slot.description,
		 conference_slot.cost,
		 conference_slot.capacity,
		 conference_slot.start_date,
		 conference_slot.end_date,
		 conference_slot.purchaseable_from,
		 conference_slot.purchaseable_until,
		 conference_slot.available_to_public,
		 location.id,
		 location.name,
		 location.description,
		 location.address,
		 location.directions,
		 location.google_maps_url,
		 location.capacity,
		 location.venue_id 
		 FROM conference_slot  
		 LEFT JOIN location ON conference_slot.location_id = location.id 
		 WHERE conference_id = $1 
		`, params.ConferenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all conferences: %w", err)
	}

	defer rows.Close()

	conferenceSlots := []ConferenceSlot{}

	for rows.Next() {
		var conferenceSlot ConferenceSlot

		err := rows.Scan(
			&conferenceSlot.ID,
			&conferenceSlot.Name,
			&conferenceSlot.Description,
			&conferenceSlot.Cost,
			&conferenceSlot.Capacity,
			&conferenceSlot.StartDate,
			&conferenceSlot.EndDate,
			&conferenceSlot.PurchaseableFrom,
			&conferenceSlot.PurchaseableUntil,
			&conferenceSlot.AvailableToPublic,
			&conferenceSlot.Location.ID,
			&conferenceSlot.Location.Name,
			&conferenceSlot.Location.Description,
			&conferenceSlot.Location.Address,
			&conferenceSlot.Location.Directions,
			&conferenceSlot.Location.GoogleMapsURL,
			&conferenceSlot.Location.Capacity,
			&conferenceSlot.Location.VenueID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}

		conferenceSlots = append(conferenceSlots, conferenceSlot)

	}

	return &GetConferenceSlotsResponse{
		ConferenceSlots: conferenceSlots,
	}, nil
}
