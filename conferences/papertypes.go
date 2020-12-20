package conferences

// Paper holds information about a paper submitted to a conference
type Paper struct {
	ID            uint32
	UserID        string
	ConferenceID  uint32
	Title         string
	ElevatorPitch string
	Description   string
	Notes         string
}
