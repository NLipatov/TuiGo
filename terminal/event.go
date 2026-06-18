package terminal

import (
	"github.com/NLipatov/tuigo/keyboard"
	"github.com/NLipatov/tuigo/mouse"
	"github.com/NLipatov/tuigo/terminal/input"
	"github.com/NLipatov/tuigo/terminal/resize"
)

type EventType int

const (
	EventUnknown EventType = iota
	EventResize
	EventKey
	EventMouse
	EventError
)

type Event struct {
	Type   EventType
	Err    error
	Key    keyboard.KeyEvent
	Mouse  mouse.MouseEvent
	Resize resize.Event
}

func newEventFromInput(in input.Event) (Event, bool) {
	switch in.Type {
	case input.EventTypeKey:
		return Event{
			Type: EventKey,
			Key:  in.Key,
		}, true
	case input.EventTypeMouse:
		return Event{
			Type:  EventMouse,
			Mouse: in.Mouse,
		}, true
	default:
		return Event{}, false
	}
}
