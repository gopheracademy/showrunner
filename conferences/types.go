package conferences

import (
	"time"
)

// Conference is a brand like GopherCon
type Conference struct {
	ID     uint32
	Name   string
	Slug   string
	Events []Event
}

// Event is an instance like GopherCon 2020
type Event struct {
	ID        uint32
	Name      string
	Slug      string
	StartDate time.Time
	EndDate   time.Time
	Location  string
	Slots     []EventSlot
}

// EventSlot holds information for any sellable/giftable slot we have in the event for
// a Talk or any other activity that requires admission.
// store: "interface"
type EventSlot struct {
	ID          uint32
	Name        string
	Description string
	Cost        int
	Capacity    int // int should be enough even if we organize glastonbury
	StartDate   time.Time
	EndDate     time.Time
	// DependsOn means that these two Slots need to be acquired together, user must either buy
	// both Slots or pre-own one of the one it depends on.
	DependsOn *EventSlot
	// PurchaseableFrom indicates when this item is on sale, for instance early bird tickets are the first
	// ones to go on sale.
	PurchaseableFrom time.Time
	// PuchaseableUntil indicates when this item stops being on sale, for instance early bird tickets can
	// no loger be purchased N months before event.
	PurchaseableUntil time.Time
	// AvailableToPublic indicates is this is something that will appear on the tickets purchase page (ie, we can
	// issue sponsor tickets and those cannot be bought individually)
	AvailableToPublic bool
}
