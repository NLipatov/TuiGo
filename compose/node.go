package compose

import "github.com/NLipatov/tuigo/canvas"

type Node interface {
	Frame() (canvas.Frame, error)
}
