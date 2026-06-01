# tuigo

Frame-based terminal UI runtime for Go.

tuigo is a lightweight, low-level runtime for applications that want to own
their UI architecture. It provides terminal setup, input and resize events, and
buffered diff rendering. The application owns state, layout, widgets, and frame
construction.

tuigo is pre-v1.0. Public APIs may change between minor releases.

## Why tuigo

- no framework lifecycle or Elm-style architecture;
- explicit `Frame -> Render` model;
- direct event consumption;
- application-owned state, layout, widgets, and redraw policy;
- performance-oriented buffered diff rendering;
- few dependencies.

## Demo

From the repository root:

```sh
go run ./examples/hello
```

Press `q`, `Esc`, or `Ctrl+C` to quit.

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

Applications should reuse frame cell buffers between redraws. Allocate a
`[]core.Cell` when the terminal size changes, mutate it to build the next frame,
wrap it with `core.NewFrame`, then call `Render`.

Terminal resize should be the only reason to allocate or replace the backing
cell buffer. Normal redraws should mutate cells in the existing buffer.

`Render` is synchronous. After it returns, the buffer can be reused for the next
frame.
