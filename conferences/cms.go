package conferences

import "time"

var _ Contentable = (*Site)(nil)
var _ Component = (*HeroComponent)(nil)

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
	return []Component{
		HeroComponent{},
		HeadlineComponent{},
	}
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
	return []string{"Large", "Medium", "Small", "Thumbnail"}[f]
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
