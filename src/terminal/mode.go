package terminal

import "golang.org/x/term"

type Device struct {
	fd           int
	initialState *term.State
}

func NewDevice(fd int) Device {
	return Device{
		fd: fd,
	}
}

func (d *Device) EnableRawMode() error {
	if d.initialState != nil {
		return nil
	}
	oldState, err := term.MakeRaw(d.fd)
	if err != nil {
		return err
	}
	d.initialState = oldState
	return nil
}

func (d *Device) RestoreInitialMode() error {
	if d.initialState == nil {
		return nil
	}
	if err := term.Restore(d.fd, d.initialState); err != nil {
		return err
	}
	d.initialState = nil
	return nil
}

func (d *Device) Size() (width, height int, err error) {
	return term.GetSize(d.fd)
}

func (d *Device) IsTerminal() bool {
	return term.IsTerminal(d.fd)
}
