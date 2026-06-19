# tuigo

Blazing fast TUI renderer and event loop for Go.

## Benchmarks

Frame time and allocations per render.

| Workload | tuigo | vaxis | tcell |
|---|---:|---:|---:|
| 1 dirty cell | **8.14 us, 0 allocs** | 34.77 us, 0 allocs | 68.75 us, 4 allocs |
| 5% dirty cells | **14.31 us, 0 allocs** | 44.17 us, 0 allocs | 88.62 us, 960 allocs |
| Full frame change | **32.06 us, 0 allocs** | 87.56 us, 0 allocs | 440.5 us, 19.20k allocs |

120x40 cells, Go 1.26.4, Apple M4 Pro, `-cpu=1`, `-count=20`.
Renderer workloads; terminal emulator paint is not measured.

Reproduce with Go benchmarks and pinned `benchstat`:

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
- minimal dependencies: `golang.org/x/term` and `github.com/rivo/uniseg`

## Demo

Rich demo:

```sh
go run ./examples/demo
```

Minimal example:

```sh
go run ./examples/hello
```

Press `q`, `Esc`, or `Ctrl+C` to quit.

## Install

```sh
go get github.com/NLipatov/tuigo@latest
```

Pre-v1.0: public APIs may change between minor releases.
