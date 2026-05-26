//go:build windows

package terminal

import (
	"context"
	"time"
)

const pollInterval = time.Millisecond * 150

type ResizeListener struct {
	ctx    context.Context
	out    chan<- ResizeEvent
	device Device
}

func NewResizeListener(ctx context.Context, ch chan ResizeEvent, device Device) ResizeListener {
	return ResizeListener{
		ctx:    ctx,
		out:    ch,
		device: device,
	}
}

func (r *ResizeListener) Listen() error {
	prevWidth, prevHeight, err := r.device.Size()
	if err != nil {
		return err
	}
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	for {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		case <-pollTicker.C:
			width, height, err := r.device.Size()
			if err != nil {
				return err
			}
			if width == prevWidth && height == prevHeight {
				continue
			}
			prevWidth, prevHeight = width, height
			event := ResizeEvent{
				Width:  width,
				Height: height,
			}
			select {
			case <-r.ctx.Done():
				return r.ctx.Err()
			case r.out <- event:
			}
		}
	}
}
