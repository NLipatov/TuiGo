package tuigo

type Renderer struct {
	frame Frame
}

func New(frame Frame) Renderer {
	return Renderer{
		frame: frame,
	}
}

func (r *Renderer) NextFrame(frame Frame) error {
	r.frame = frame
	return nil
}

func (r *Renderer) Render() error {
	for y := range r.frame.Height() {
		for x := range r.frame.Width() {
			_, err := r.frame.CellAt(x, y)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
