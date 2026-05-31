package terminal

import (
	"tuigo/terminal/input"
	"tuigo/terminal/resize"
)

type EventType int

const (
	EventUnknown EventType = iota
	EventResize
	EventKey
	EventError
)

type Event struct {
	Type   EventType
	Err    error
	Key    input.Event
	Resize resize.Event
}
