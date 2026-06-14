# tuigo benchmarks

This module keeps benchmark dependencies out of the main tuigo module.

Run from this directory:

```sh
go test -run=^$ -bench=. -benchmem -count=20 -benchtime=1s -cpu=1 ./...
```

For local smoke checks, use a shorter run:

```sh
go test -run=^$ -bench=. -benchmem -benchtime=100ms -cpu=1 ./...
```

For publishable comparisons, save the output and summarize it with `benchstat`:

```sh
go test -run=^$ -bench=. -benchmem -count=20 -benchtime=1s -cpu=1 ./... > results.txt
go run golang.org/x/perf/cmd/benchstat@v0.0.0-20260610192853-712aea8b4705 results.txt
```

For comparisons, keep the full command, Go version, OS, CPU, dependency
versions, and git commit with the results.

## Scope

These benchmarks compare renderer-style workloads over a fixed terminal-sized
grid. They intentionally do not benchmark a full application framework loop.
Treat them as a reproducible harness for investigation and regression tracking,
not as README-ready marketing numbers by themselves.

- `tuigo` uses two preallocated `core.Frame` buffers and renders through
  `terminal/render.Renderer` into a discard writer.
- `tcell` uses `SimulationScreen`, applies changed cells with `Put` for sparse
  changes and `PutStrStyled` for contiguous rows/runs, then flushes with `Show`.
- `vaxis` uses `Window.SetCell` or `Window.Fill`, then flushes with `Render`
  through a fake console.

The tcell driver uses a simulation backend, and the vaxis driver uses a fake
console. The tuigo driver encodes ANSI bytes but writes them to a discard
writer. These benchmarks do not measure terminal syscall or emulator paint
cost.

Bubble Tea is intentionally excluded from the renderer table. It is an
application framework with a model/update/view loop, so comparing it here would
mix framework view construction with low-level renderer flushing.

## Reporting

Keep renderer-level numbers together:

```md
| Workload | tuigo | tcell | vaxis |
|---|---:|---:|---:|
| One cell change |  |  |  |
| Random 5% change |  |  |  |
| Full content change |  |  |  |
```

Use `ns/op`, `B/op`, and `allocs/op` together. For tuigo renderer rows, include
`bytes_out/op` when discussing terminal output volume.

## Workloads

- `no-change/120x40`: no cells change after the initial render.
- `full-content-change/120x40`: every cell changes between two frames.
- `one-cell-change/120x40`: one cell changes.
- `run-40-change/120x40`: one contiguous 40-cell run changes.
- `random-5pct-change/120x40`: 5% of cells change at deterministic positions.
- `full-style-change/120x40`: every cell keeps the same rune and changes style.
- `unicode-content-change/120x40`: every cell changes from ASCII to a multibyte
  Unicode rune.
