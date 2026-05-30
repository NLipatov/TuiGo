//go:build unix

package resize

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type Listener struct {
	ctx    context.Context
	out    chan<- Event
	device Device
}

func NewListener(ctx context.Context, ch chan Event, device Device) Listener {
	return Listener{
		ctx:    ctx,
		out:    ch,
		device: device,
	}
}

func (r *Listener) Listen() error {
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
			event := Event{
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
