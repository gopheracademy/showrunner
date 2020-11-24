package conferences

// Component is a structured content piece for
// a web page
type Component interface {
	ComponentName() string
	ValidFor(Contentable) bool
}

// HeroComponent is a component content type
type HeroComponent struct {
	ID       uint32
	Headline string
	Image    Image
	Body     string
}

// ComponentName returns the name of the Component
// to satisfy the Component interface
func (h HeroComponent) ComponentName() string {
	return "HeroComponent"
}

// ValidFor returns true when the component can be applied
// to a Contentable type, otherwise false
func (h HeroComponent) ValidFor(c Contentable) bool {
	for _, component := range c.ValidComponents() {
		if component.ComponentName() == h.ComponentName() {
			return true
		}
	}
	return false
}

// HeadlineComponent is a component content type
type HeadlineComponent struct {
	ID       uint32
	Headline string
}

// ComponentName returns the name of the Component
// to satisfy the Component interface
func (h HeadlineComponent) ComponentName() string {
	return "HeadlineComponent"
}

// ValidFor returns true when the component can be applied
// to a Contentable type, otherwise false
func (h HeadlineComponent) ValidFor(c Contentable) bool {
	for _, component := range c.ValidComponents() {
		if component.ComponentName() == h.ComponentName() {
			return true
		}
	}
	return false
}
