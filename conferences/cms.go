package conferences

import "time"

// Site is the parent type of all content
type Site struct {
	ID         uint32
	Name       string
	Domains    []Domain
	Events     []Event
	Components []Component
}

// TypeID returns the TypeID of the site to
// satisfy the Contentable interface
func (s Site) TypeID() uint32 {
	return s.ID
}

// TypeName returns the TypeName of the site to
// satisfy the Contentable interface
func (s Site) TypeName() string {
	return "Site"
}

// ValidComponents returns a list of components
// that may be added to the Site to satisfy
// the Contentable interface
func (s Site) ValidComponents() []Component {
	return []Component{}
}

// PrimaryDomain returns the preferred domain
// name for a site
func (s Site) PrimaryDomain() Domain {
	for _, d := range s.Domains {
		if d.Primary {
			return d
		}
	}
	// no domains, must be localhost!
	return Domain{
		ID:      0,
		FQDN:    "127.0.0.1",
		Primary: false,
	}
}

// Domain represents an internet address
type Domain struct {
	ID      uint32
	FQDN    string
	Primary bool
}

// Image represents an image to be
// referenced in content
type Image struct {
	ID        uint32
	Name      string
	AltText   string
	Caption   string
	Width     int
	Height    int
	Ext       string
	URL       string
	Hash      string
	Mime      string  // "image/png"
	Size      float32 // size in bytes
	Formats   []ImageFormat
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ImageFormat stores the details about alternate
// Image sizes
type ImageFormat struct {
	ID     uint32
	Format Format
	Ext    string
	URL    string
	Hash   string
	Mime   string  // "image/png"
	Size   float32 // size in bytes
	Width  int
	Height int
}

// Format is a description of an ImageFormat
type Format int

// Image Formats
const (
	Large Format = iota
	Medium
	Small
	Thumbnail
)

func (f Format) String() string {
	return [...]string{"Large", "Medium", "Small", "Thumbnail"}[f]
}

// Contentable represents a type that may
// have content attached
type Contentable interface {
	TypeID() uint32
	// Name of type "Speaker"
	TypeName() string
	// List of valid components for this type
	ValidComponents() []Component
}

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
