package terminal

import (
	"github.com/NLipatov/tuigo/terminal/input"
	"github.com/NLipatov/tuigo/terminal/resize"
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
