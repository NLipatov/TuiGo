//go:build unix

package terminal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

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
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	defer signal.Stop(sigCh)
	for {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		case <-sigCh:
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
