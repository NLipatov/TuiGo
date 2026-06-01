# tuigo

Frame-based terminal UI runtime for Go.

tuigo is pre-v1.0. Public APIs may change between minor releases.

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
