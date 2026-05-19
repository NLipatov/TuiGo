package tuigo

type Renderer struct {
	frame, oldFrame Frame
	firstRender     bool
}

func New(frame Frame) Renderer {
	return Renderer{
		frame:       frame,
		firstRender: true,
	}
}

func (r *Renderer) NextFrame(newFrame Frame) error {
	if newFrame.Width() != r.frame.Width() ||
		newFrame.Height() != r.frame.Height() {
		r.frame = newFrame
		r.firstRender = true
		return nil
	}
	r.oldFrame = r.frame
	r.frame = newFrame
	return nil
}

func (r *Renderer) Render() error {
	for y := range r.frame.Height() {
		for x := range r.frame.Width() {
			cell, err := r.frame.CellAt(x, y)
			if err != nil {
				return err
			}
			if r.firstRender {
				//render the cell
			} else {
				oldCell, err := r.oldFrame.CellAt(x, y)
				if err != nil {
					return err
				}
				if cell != oldCell {
					// rerender this cell
				}
			}
		}
	}
	if r.firstRender {
		r.firstRender = false
	}
	r.oldFrame = r.frame
	return nil
}
