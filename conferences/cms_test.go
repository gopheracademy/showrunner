package conferences

import "testing"

var gcdotcom = Domain{
	ID:      1,
	FQDN:    "www.gophercon.com",
	Primary: true,
}

var gc21dotcom = Domain{
	ID:      1,
	FQDN:    "2021.gophercon.com",
	Primary: false,
}
var gcSite = Site{
	ID:      1,
	Name:    "GopherCon",
	Domains: []Domain{gc21dotcom, gcdotcom},
}
var localSite = Site{
	ID:      2,
	Name:    "GopherCon",
	Domains: []Domain{},
}

// AgendaComponent is a component content type
type AgendaComponent struct {
	ID uint32
}

// ComponentName returns the name of the Component
// to satisfy the Component interface
func (h AgendaComponent) ComponentName() string {
	return "AgendaComponent"
}

// ValidFor returns true when the component can be applied
// to a Contentable type, otherwise false
func (h AgendaComponent) ValidFor(c Contentable) bool {
	for _, component := range c.ValidComponents() {
		if component.ComponentName() == h.ComponentName() {
			return true
		}
	}
	return false
}
func TestSite(t *testing.T) {
	// primary domain is set
	if gcSite.PrimaryDomain().ID != gcdotcom.ID {
		t.Errorf("wanted %d, got %d", gcdotcom.ID, gcSite.PrimaryDomain().ID)
	}
	// no domain set
	if localSite.PrimaryDomain().ID != 0 {
		t.Errorf("wanted %d, got %d", 0, localSite.PrimaryDomain().ID)
	}
}

func TestContentable(t *testing.T) {

	hc := HeroComponent{
		ID:       1,
		Headline: "Join us in San Diego",
		Image:    Image{},
		Body:     "We're going to have a great time!",
	}
	if !hc.ValidFor(gcSite) {
		t.Error("HeroComponent not valid for Site")
	}
	gcSite.Components = append(gcSite.Components, hc)
	headline := HeadlineComponent{
		ID:       1,
		Headline: "GopherCon 2021 Agenda",
	}

	if !headline.ValidFor(gcSite) {
		t.Error("HeadlineComponent should not be valid for Site")
	}
	agenda := AgendaComponent{
		ID: 1,
	}

	if agenda.ValidFor(gcSite) {
		t.Error("AgendaComponent should not be valid for Site")
	}
}
