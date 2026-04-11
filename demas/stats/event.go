package stats

type EventType string

const (
	EventBorn    EventType = "BORN"
	EventDead    EventType = "DEAD"
	EventMeeting EventType = "MEETING"
)

type Event struct {
	Score           float64
	SameIsland      bool
	PrimaryIslandID int
	EventType       EventType
}

