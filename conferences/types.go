package conferences

import (
	"time"
)

// Event is a brand like GopherCon
type Event struct {
	ID          uint32
	Name        string
	Slug        string
	Conferences []Conference
}

// Conference is an instance like GopherCon 2020
type Conference struct {
	ID        uint32
	Name      string
	Slug      string
	StartDate time.Time
	EndDate   time.Time
	Venue     Venue
	Slots     []ConferenceSlot
}

// ConferenceSlot holds information for any sellable/giftable slot we have in the event for
// a Talk or any other activity that requires admission.
// store: "interface"
type ConferenceSlot struct {
	ID          uint32
	Name        string
	Description string
	Cost        int
	Capacity    int // int should be enough even if we organize glastonbury
	StartDate   time.Time
	EndDate     time.Time
	// DependsOn means that these two Slots need to be acquired together, user must either buy
	// both Slots or pre-own one of the one it depends on.
	// DependsOn *ConferenceSlot // Currently removed as it broke encore
	// PurchaseableFrom indicates when this item is on sale, for instance early bird tickets are the first
	// ones to go on sale.
	PurchaseableFrom time.Time
	// PuchaseableUntil indicates when this item stops being on sale, for instance early bird tickets can
	// no loger be purchased N months before event.
	PurchaseableUntil time.Time
	// AvailableToPublic indicates is this is something that will appear on the tickets purchase page (ie, we can
	// issue sponsor tickets and those cannot be bought individually)
	AvailableToPublic bool
	Location          Location
}

// Venue defines a venue that hosts a conference, such as DisneyWorld
type Venue struct {
	ID            uint32
	Name          string
	Description   string
	Address       string
	Directions    string
	GoogleMapsURL string
	Capacity      int
}

// Location defines a location for a venue, such as a room or event space
type Location struct {
	ID            uint32
	Name          string
	Description   string
	Address       string
	Directions    string
	GoogleMapsURL string
	Capacity      int
	VenueID       uint32
}

// SponsorshipLevel defines the type that encapsulates the different sponsorship levels
type SponsorshipLevel int

// These are the valid sponsorship levels
const (
	SponsorshipLevelPlatinum SponsorshipLevel = iota
	SponsorshipLevelGold
	SponsorshipLevelSilver
	SponsorshipLevelBronze
)

func (s SponsorshipLevel) String() string {
	return []string{"platinum", "gold", "silver", "bronze"}[s]
}

// Sponsor defines a conference sponsor, such as Google
type Sponsor struct {
	ID               uint32
	Name             string
	Address          string
	Website          string
	SponsorshipLevel SponsorshipLevel
	Contacts         []SponsorContactInformation
}

// ContactRole defines the type that encapsulates the different contact roles
type ContactRole int

// These are the valid contact roles
const (
	ContactRoleMarketing ContactRole = iota
	ContactRoleLogistics
	ContactRoleTechnical
	ContactRoleOther
	ContactRoleSoleContact
)

func (c ContactRole) String() string {
	return []string{"marketing", "logistics", "technical", "other", "sole_contact"}[c]
}

// SponsorContactInformation defines a contact
//and their information for a sponsor
type SponsorContactInformation struct {
	ID    uint32
	Name  string
	Role  ContactRole
	Email string
	Phone string
}
