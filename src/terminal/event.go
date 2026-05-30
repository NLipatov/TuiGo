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
)

type Event struct {
	Type   EventType
	Key    input.Event
	Resize resize.Event
}
