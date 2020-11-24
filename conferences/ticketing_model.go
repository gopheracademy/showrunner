package conferences

import "github.com/gopheracademy/manager/def"

// EventSlot holds information for any sellable/giftable slot we have in the event for
// a Talk or any other activity that requires admission.
type EventSlot struct {
	ID          uint64     `gaum:"field_name:id"`
	Event       *def.Event `gaum:"field_name:event"`
	Name        string     `gaum:"field_name:name"`
	Description string     `gaum:"field_name:description"`
	Cost        int64      `gaum:"field_name:cost"`
	Capacity    int        `gaum:"field_name:capacity"`   // int should be enough even if we organize glastonbury
	StartDate   uint64     `gaum:"field_name:start_date"` // StartDate is Unix timestamp, seconds since Epoch (1/1/1970 UTC)
	EndDate     uint64     `gaum:"field_name:end_date"`   // EndDate is Unix timestamp, seconds since Epoch (1/1/1970 UTC)
	// DependsOn means that these two Slots need to be acquired together, user must either buy
	// both Slots or pre-own one of the one it depends on.
	DependsOn *EventSlot
	// PurchaseableFrom indicates when this item is on sale, for instance early bird tickets are the first
	// ones to go on sale.
	PurchaseableFrom uint64 `gaum:"field_name:purchaseable_from"` // PurchaseableFrom is Unix timestamp, seconds since Epoch (1/1/1970 UTC)
	// PuchaseableUntil indicates when this item stops being on sale, for instance early bird tickets can
	// no loger be purchased N months before event.
	PurchaseableUntil uint64 `gaum:"field_name:purchaseable_until"` // PurchaseableUntil is Unix timestamp, seconds since Epoch (1/1/1970 UTC)
	// AvailableToPublic indicates is this is something that will appear on the tickets purchase page (ie, we can
	// issue sponsor tickets and those cannot be bought individually)
	AvailableToPublic bool `gaum:"field_name:available_to_public"`
}

// ClaimPayment represents a payment for N claims
type ClaimPayment struct {
	ID uint64 `gaum:"field_name:id"`
	// ClaimsPayed would be what in a bill one see as detail.
	ClaimsPayed []*SlotClaim
	Payment     []FinancialInstrument
	Invoice     string `gaum:"field_name:invoice"` // let us fill this once we know how to invoice
}

// TotalDue returns the total cost to cover by this payment.
func (c *ClaimPayment) TotalDue() int64 {
	var totalDue int64 = 0
	for _, sc := range c.ClaimsPayed {
		totalDue = totalDue + sc.EventSlot.Cost
	}
	return totalDue
}

// Fulfilled returns true if the payment of this invoice has been fulfilled
func (c *ClaimPayment) Fulfilled() bool {
	totalDue := c.TotalDue()
	f, _ := PaymentBalanced(totalDue, c.Payment...)
	b, _ := DebtBalanced(c.Payment...)
	return f && b
}

// SlotClaim represents one occupancy of one slot.
type SlotClaim struct {
	ID        uint64 `gaum:"field_name:id"`
	EventSlot *EventSlot
	// TicketID should only be valid when combined with the correct Attendee ID/Email
	TicketID string `gaum:"field_name:ticket_id"` // uuid
	// Redeemed represents whether this has been used (ie the Attendee enrolled in front desk
	// or into the online conf system) until this is not true, transfer/refund might be possible.
	Redeemed bool `gaum:"field_name:redeemed"`
}

// Attendee is a person attending one or more Slots of the Conference.
type Attendee struct {
	ID    uint64 `gaum:"field_name:id"`
	Email string `gaum:"field_name:email"`
	// CoCAccepted, claims cannot be used without this.
	CoCAccepted bool `gaum:"field_name:co_c_accepted"`
	Claims      []SlotClaim
}

// Finance Section

// PaymentMethodMoney represents a payment in cash.
type PaymentMethodMoney struct {
	ID         uint64 `gaum:"field_name:id"`
	PaymentRef string `gaum:"field_name:ref"`    // stripe payment ID/Log?
	Amount     int64  `gaum:"field_name:amount"` // Money is handled in ints to ease use of OTO, do not divide
}

// Total implements FinancialInstrument
func (p *PaymentMethodMoney) Total() int64 {
	return p.Amount
}

// Type implements FinancialInstrument
func (p *PaymentMethodMoney) Type() AssetType {
	return ATCash
}

var _ FinancialInstrument = &PaymentMethodMoney{}

// PaymentMethodConferenceDiscount represents a discount issued by the event.
type PaymentMethodConferenceDiscount struct {
	ID uint64 `gaum:"field_name:id"`
	// Detail describes what kind of discount was issued (ie 100% sponsor, 30% grant)
	Detail string `gaum:"field_name:detail"`
	Amount int64  `gaum:"field_name:amount"` // Money is handled in ints to ease use of OTO, do not divide
}

// Total implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Total() int64 {
	return p.Amount
}

// Type implements FinancialInstrument
func (p *PaymentMethodConferenceDiscount) Type() AssetType {
	return ATDiscount
}

var _ FinancialInstrument = &PaymentMethodConferenceDiscount{}

// PaymentMethodCreditNote represents credit extended to defer payment.
type PaymentMethodCreditNote struct {
	ID     uint64 `gaum:"field_name:id"`
	Detail string `gaum:"field_name:detail"`
	Amount int64  `gaum:"field_name:amount"` // Money is handled in ints to ease use of OTO, do not divide
}

// Total implements FinancialInstrument
func (p *PaymentMethodCreditNote) Total() int64 {
	return p.Amount
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
// oto: "skip"
type FinancialInstrument interface {
	// Total is the total amount fulfilled by this instrument
	Total() int64
	// Type is the type of asset represented
	Type() AssetType
}

// PaymentBalanced returns true or false depending on balancing status and missing
// payment amount if any.
func PaymentBalanced(amount int64, payments ...FinancialInstrument) (bool, int64) {
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

// DebtBalanced returns true if all cretid notes or similar instruments have been covered or an
// amount if not.
func DebtBalanced(payments ...FinancialInstrument) (bool, int64) {
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
