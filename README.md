# tuigo

Blazing fast, zero-allocation TUI renderer for Go.

## Benchmarks

| Workload, frame time | tuigo | vaxis | tcell |
|---|---:|---:|---:|
| 1 dirty cell | **12.26 us, 0 allocs** | 34.41 us, 0 allocs | 68.78 us, 4 allocs |
| 5% dirty cells | **18.94 us, 0 allocs** | 43.68 us, 0 allocs | 89.12 us, 960 allocs |
| Full frame change | **44.21 us, 0 allocs** | 87.33 us, 0 allocs | 440.4 us, 19.20k allocs |

120x40 cells, Go 1.26.3, Apple M4 Pro, `-cpu=1`, `-count=20`.
Renderer workloads; terminal paint is not measured.

Reproduce:

```sh
cd benchmarks
go test -run=^$ -bench=BenchmarkRenderer -benchmem -count=20 -benchtime=1s -cpu=1 ./... > results.txt
go run golang.org/x/perf/cmd/benchstat@v0.0.0-20260610192853-712aea8b4705 results.txt
```

## Features

- `Frame -> Render` model
- zero allocations with reused frame buffers
- keyboard, mouse, resize events
- ANSI diff output
- one runtime dependency: `golang.org/x/term`

## Install

```sh
go get github.com/NLipatov/tuigo/terminal
```

## Demo

```sh
go run ./examples/demo   # rich demo
go run ./examples/hello  # minimal example
```

Press `q`, `Esc`, or `Ctrl+C` to quit.

## Usage

Pre-v1.0: public APIs may change between minor releases.

```go
import (
	"context"
	"os"

	"github.com/NLipatov/tuigo/terminal"
)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()

session, err := terminal.NewSession(ctx, os.Stdin, os.Stdout)
if err != nil {
	return err
}
defer func() {
	_ = session.Close()
}()

events, err := session.Start()
if err != nil {
	return err
}

width, height, err := session.Size()
if err != nil {
	return err
}

frame := buildFrame(width, height)
if err := session.Render(frame); err != nil {
	return err
}

for event := range events {
	switch event.Type {
	case terminal.EventKey:
		// update state from event.Key
	case terminal.EventMouse:
		// update state from event.Mouse
	case terminal.EventResize:
		// rebuild and render for event.Resize.Width/Height
	case terminal.EventError:
		// decide whether to stop
	}
}
```

## Frame Buffers

For zero-allocation rendering, reuse two cell buffers: current and next. Draw
into next, render it, then swap.

```go
current := make([]core.Cell, width*height)
next := make([]core.Cell, width*height)

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
