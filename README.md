# tuigo

Frame-based terminal UI runtime for Go.

tuigo is a lightweight, low-level runtime for applications that want to own
their UI architecture. It provides terminal setup, input and resize events, and
buffered diff rendering. The application owns state, layout, widgets, and frame
construction.

Note: tuigo is pre-v1.0. Public APIs may change between minor releases.

## Why tuigo

- performance-oriented buffered diff rendering;
- explicit `Frame -> Render` model;
- single dependency: `golang.org/x/term`;
- no framework-managed lifecycle or state;
- direct event consumption;
- application-owned state, layout, widgets, and redraw policy.

## Demo

From the repository root:

```sh
go run ./examples/hello
```

Press `q`, `Esc`, or `Ctrl+C` to quit.

For a richer visual demo:

```sh
go run ./examples/demo
```

## Usage

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

session, err := terminal.NewSession(ctx, os.Stdin, os.Stdout)
events, err := session.Start()
defer session.Close()

width, height, err := session.Size()
frame := buildFrame(width, height)
err = session.Render(frame)

for event := range events {
	switch event.Type {
	case terminal.EventKey:
		// update application state
	case terminal.EventResize:
		// rebuild and render a frame for event.Resize.Width/Height
	case terminal.EventError:
		// decide whether to stop
	}
}
```

`core.Frame` represents a width x height grid backed by a flat slice of
`core.Cell` values. Each cell has a rune, foreground color, and background
color.

## Frame buffers

To improve rendering performance, use two preallocated frame buffers: the
current frame and the next frame. Allocate both `[]core.Cell` buffers when the
terminal size changes, mutate the next buffer, wrap it with `core.NewFrame`, and
call `Render`. After a successful render, swap the current and next buffers. Do
not mutate the buffer from the last successful `Render` call until another
buffer has been rendered successfully.

```go
current, next := make([]core.Cell, width*height), make([]core.Cell, width*height)

draw(next, width, height, state)
frame, err := core.NewFrame(width, height, next)
if err != nil {
	return err
}
if err := session.Render(frame); err != nil {
	return err
}
current, next = next, current
```
