package conferences

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/gofrs/uuid"
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
	DependsOn uint32
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
	ConferenceID      uint32
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

// ClaimPayment represents a payment for N claims
type ClaimPayment struct {
	ID uint64
	// ClaimsPaid would be what in a bill one see as detail.
	ClaimsPaid []*SlotClaim
	Payment    []FinancialInstrument
	Invoice    string // let us fill this once we know how to invoice
}

// TotalDue returns the total cost to cover by this payment.
func (c *ClaimPayment) TotalDue() int64 {
	var totalDue int = 0
	for _, sc := range c.ClaimsPaid {
		totalDue = totalDue + sc.ConferenceSlot.Cost
	}
	return int64(totalDue)
}

// Fulfilled returns true if the payment of this invoice has been covered with either
// money or credit
func (c *ClaimPayment) Fulfilled() bool {
	totalDue := c.TotalDue()
	f, _ := paymentBalanced(totalDue, c.Payment...)
	return f
}

// Paid returns true if the payment of this invoice has been fully paid.
func (c *ClaimPayment) Paid() bool {
	totalDue := c.TotalDue()
	f, _ := paymentFulfilled(totalDue, c.Payment...)
	b, _ := debtBalanced(c.Payment...)
	return f && b
}

// SlotClaim represents one occupancy of one slot.
type SlotClaim struct {
	ID             int64
	ConferenceSlot *ConferenceSlot
	// TicketID should only be valid when combined with the correct Attendee ID/Email
	TicketID uuid.UUID
	// Redeemed represents whether this has been used (ie the Attendee enrolled in front desk
	// or into the online conf system) until this is not true, transfer/refund might be possible.
	Redeemed bool
}

// Attendee is a person attending one or more Slots of the Conference.
type Attendee struct {
	ID    int64
	Email string
	// CoCAccepted, claims cannot be used without this.
	CoCAccepted bool
	Claims      []SlotClaim
}

// Finance Section

// PaymentMethodMoney represents a payment in cash.
type PaymentMethodMoney struct {
	ID          uint64
	PaymentRef  string // stripe payment ID/Log?
	AmountCents int64  // Money is handled in cents as it is done by our payment processor (stripe)
}

// Total implements FinancialInstrument
func (p *PaymentMethodMoney) Total() int64 {
	return p.AmountCents
}

// Type implements FinancialInstrument
func (p *PaymentMethodMoney) Type() AssetType {
	return ATCash
}

var _ FinancialInstrument = &PaymentMethodMoney{}

// PaymentMethodConferenceDiscount represents a discount issued by the event.
type PaymentMethodConferenceDiscount struct {
	ID uint64
	// Detail describes what kind of discount was issued (ie 100% sponsor, 30% grant)
	Detail      string
	AmountCents int64 // Money is handled in cents as it is done by our payment processor (stripe)
}

// Total implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Total() int64 {
	return p.AmountCents
}

// Type implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Type() AssetType {
	return ATDiscount
}

var _ FinancialInstrument = &PaymentMethodConferenceDiscount{}

// PaymentMethodCreditNote represents credit extended to defer payment.
type PaymentMethodCreditNote struct {
	ID          uint64
	Detail      string
	AmountCents int64 // Money is handled in cents as it is done by our payment processor (stripe)
}

// Total implements FinancialInstrument
func (p *PaymentMethodCreditNote) Total() int64 {
	return p.AmountCents
}

// Type implements FinancialInstrument
func (p *PaymentMethodCreditNote) Type() AssetType {
	return ATReceivable
}

var _ FinancialInstrument = &PaymentMethodCreditNote{}

// AssetType is a type of accounting asset.
type AssetType string

const (
	// ATCash in this context means it is money, like a stripe payment
	ATCash AssetType = "cash"
	// ATReceivable in this context means it is a promise of payment
	ATReceivable AssetType = "receivable"
	// ATDiscount in this context means an issued discount (represented as a fixed amount for
	// accounting's sake)
	ATDiscount AssetType = "discount"
)

// FinancialInstrument represents any kind of instrument used to cover a debt.
type FinancialInstrument interface {
	// Total is the total amount fulfilled by this instrument
	Total() int64
	// Type is the type of asset represented
	Type() AssetType
}

// paymentBalanced returns true or false depending on balancing status and missing
// payment amount if any.
func paymentBalanced(amount int64, payments ...FinancialInstrument) (bool, int64) {
	var receivables int64 = 0
	var received int64 = 0
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received += p.Total()
		case ATReceivable:
			receivables += p.Total()
		}
	}
	missing := amount - received - receivables
	return missing <= 0, missing
}

// paymentFulfilled returns true if the passed amount is covered in full.
func paymentFulfilled(amount int64, payments ...FinancialInstrument) (bool, int64) {
	var received int64 = 0
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received += p.Total()
		}
	}
	missing := amount - received
	return missing <= 0, missing
}

// debtBalanced returns true if all credit notes or similar instruments have been covered or an
// amount if not.
func debtBalanced(payments ...FinancialInstrument) (bool, int64) {
	var receivables int64 = 0
	var received int64 = 0
	for _, p := range payments {
		switch p.Type() {
		case ATCash, ATDiscount:
			received += p.Total()
		case ATReceivable:
			receivables += p.Total()
		}
	}
	missing := receivables - received
	return missing <= 0, missing
}

// SponsorshipLevel defines the type that encapsulates the different sponsorship levels
type SponsorshipLevel int

// Scan converts from database to Go value
func (s *SponsorshipLevel) Scan(src interface{}) error {
	if src == nil {
		*s = SponsorshipLevelNone
		return nil
	}
	if iv, err := driver.String.ConvertValue(src); err == nil {
		if v, ok := iv.([]byte); ok {
			sv := string(v)
			switch sv {
			case "none":
				*s = SponsorshipLevelNone
				return nil
			case "diamond":
				*s = SponsorshipLevelDiamond
				return nil
			case "platinum":
				*s = SponsorshipLevelPlatinum
				return nil
			case "gold":
				*s = SponsorshipLevelGold
				return nil
			case "silver":
				*s = SponsorshipLevelSilver
				return nil
			case "bronze":
				*s = SponsorshipLevelBronze
				return nil
			default:
				*s = SponsorshipLevelNone
				return nil
			}
		}
	}
	// otherwise, return an error
	return errors.New("failed to scan SponsorshipLevel")
}

// Value - Implementation of valuer for database/sql
func (s SponsorshipLevel) Value() (driver.Value, error) {
	// value needs to be a base driver.Value type
	// such as bool.
	return s.String(), nil
}

// These are the valid sponsorship levels
// Not all of these may be used each year, for
// example 2020 had no Diamonds
const (
	SponsorshipLevelNone = iota
	SponsorshipLevelDiamond
	SponsorshipLevelPlatinum
	SponsorshipLevelGold
	SponsorshipLevelSilver
	SponsorshipLevelBronze
	SponsorshipLevelOther
)

func (s SponsorshipLevel) String() string {
	return []string{"none", "diamond", "platinum", "gold", "silver", "bronze", "other"}[s]
}

// Sponsor defines a conference sponsor, such as Google
type Sponsor struct {
	ID               uint32
	Name             string
	Address          string
	Website          string
	SponsorshipLevel SponsorshipLevel
	Contacts         []SponsorContactInformation
	ConferenceID     uint32
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

var contactRoleMappings = []string{"marketing", "logistics", "technical", "other", "sole_contact"}

func (c ContactRole) String() string {
	return contactRoleMappings[c]
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
