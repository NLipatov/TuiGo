package resize

type Device interface {
	Size() (width, height int, err error)
}
