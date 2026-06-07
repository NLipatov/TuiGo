package input

type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeKey
	EventTypeMouse
)

type Event struct {
	Type  EventType
	Key   KeyEvent
	Mouse MouseEvent
}
