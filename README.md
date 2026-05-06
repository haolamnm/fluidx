# fluidx

Real-time 2D fluid dynamics simulator for the terminal. Implements Jos Stam's "Stable Fluids" method (Navier-Stokes).

## Install

```bash
go install github.com/haolamnm/fluidx/cmd/fluidx@latest
```

Requires Go 1.22+, TrueColor terminal.

## Usage

```bash
./fluidx
```

| Input | Action |
|-------|--------|
| Mouse move | Stir fluid |
| Click/hold | Inject smoke |
| Arrows | Velocity at center |
| Space | Smoke at center |
| `d` | Toggle debug arrows |
| `q` / ESC / Ctrl+C | Quit |

## Build

```bash
make build    # or: go build ./cmd/fluidx
make test     # or: go test ./...
make run
make clean
```

## License

MIT
