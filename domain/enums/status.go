package enums

type Status string

const (
	StatusStart Status = "start"
	StatusStop  Status = "stop"
)

func (r Status) String() string {
	return string(r)
}

func (r Status) IsValid() bool {
	switch r {
	case StatusStart,
		StatusStop:
		return true
	default:
		return false
	}
}
