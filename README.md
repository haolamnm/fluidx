# fluidx

A high-performance, real-time fluid dynamics simulator for the terminal. Implements Jos Stam's "Stable Fluids" method for 2D Navier-Stokes equations.

## Features

- **Physics Solver**: Grid-based Navier-Stokes solver with advection, diffusion, and pressure projection
- **ASCII Rendering**: Density mapped to grayscale ASCII ramp (`.:-=+*#%@`)
- **TrueColor Support**: Velocity visualized as color gradients (blue→green→yellow→red)
- **Interactive**: Mouse and keyboard input to stir fluid and inject smoke
- **Debug Mode**: Overlay showing velocity vectors as directional arrows
- **60 FPS Target**: Optimized for smooth real-time simulation

## Installation

```bash
go install github.com/haolamnm/fluidx/cmd/fluidx@latest
```

Or clone and build:

```bash
git clone https://github.com/haolamnm/fluidx.git
cd fluidx
go build ./cmd/fluidx
```

### Prerequisites

- Go 1.22 or later
- A terminal emulator that supports TrueColor (most modern terminals do)

## Usage

```bash
./fluidx
```

### Controls

| Input | Action |
|-------|--------|
| `Mouse Move` | Inject velocity (stir fluid) |
| `Mouse Click/Hold` | Inject density (smoke) |
| `Arrow Keys` | Inject velocity at screen center |
| `Space` | Inject density at screen center |
| `d` | Toggle debug mode (velocity arrows) |
| `q` / `ESC` / `Ctrl+C` | Quit |

## Architecture

```
fluidx/
├── cmd/fluidx/
│   └── main.go          # Application entry point
├── internal/
│   ├── solver/          # Navier-Stokes physics solver
│   │   ├── solver.go    # Main solver API
│   │   └── operations.go # Advection, diffusion, projection
│   ├── renderer/        # TUI rendering engine
│   │   └── renderer.go  # ASCII + TrueColor output
│   └── input/           # Input handling
│       └── handler.go   # Mouse/keyboard events
├── Makefile             # Build automation
└── go.mod
```

## Physics Implementation

The solver implements Jos Stam's "Stable Fluids" method:

1. **Diffusion**: Gauss-Seidel relaxation for viscous spread of velocity and density
2. **Projection**: Poisson solver ensures mass conservation (incompressibility)
3. **Advection**: Semi-Lagrangian backwards particle tracing with bilinear interpolation

### Performance Optimizations

- Double buffering prevents simulation artifacts
- No-slip boundary conditions with negated normal velocity at walls
- Pre-computed grid constants avoid repeated calculations
- Fast `atan2` approximation for velocity arrow direction
- 60 FPS cap with `time.NewTicker` prevents excessive CPU usage

## Development

```bash
# Build
make build

# Test
make test

# Run
make run

# Clean build artifacts
make clean
```

## License

MIT License - See [LICENSE](LICENSE) file
