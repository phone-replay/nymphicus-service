package enum

type SessionStatus int

const (
	InProgress SessionStatus = iota
	Complete
	Error
)

func (os SessionStatus) String() string {
	return [...]string{"InProgress", "Complete", "Error"}[os]
}
