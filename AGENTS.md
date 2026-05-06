# fluidx - Project Knowledge Base

**Generated:** 2025-04-15
**Language:** Go 1.21
**Type:** Terminal-based fluid dynamics simulator (Jos Stam's Stable Fluids)

## OVERVIEW
Real-time 2D Navier-Stokes solver with ASCII/TrueColor rendering. Implements advection, diffusion, and pressure projection using Gauss-Seidel relaxation.

## STRUCTURE
```
.
├── cmd/fluidx/          # Application entry point
│   └── main.go          # 60 FPS game loop, event handling
├── internal/            # Private implementation
│   ├── solver/          # Navier-Stokes physics (SEE: internal/solver/AGENTS.md)
│   ├── renderer/        # TUI/ASCII output (SEE: internal/renderer/AGENTS.md)
│   └── input/           # Mouse/keyboard events
├── pkg/fluids/          # EMPTY - reserved for public API
├── go.mod               # Module: fluidx, tcell v2 dependency
└── README.md
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Entry point | `cmd/fluidx/main.go` | Main loop, event routing |
| Physics core | `internal/solver/` | Navier-Stokes implementation |
| Rendering | `internal/renderer/` | ASCII ramp, TrueColor gradients |
| Input handling | `internal/input/handler.go` | tcell event abstraction |
| Tests | `internal/solver/solver_test.go` | Only test file in project |

## CODE MAP

| Symbol | Type | File | Role |
|--------|------|------|------|
| `main()` | func | cmd/fluidx/main.go | Entry, error handling |
| `run()` | func | cmd/fluidx/main.go | Main loop with ticker |
| `NewFluidSolver()` | func | internal/solver/solver.go | Solver constructor |
| `FluidSolver.Step()` | method | internal/solver/solver.go | Simulation tick |
| `Renderer.Render()` | method | internal/renderer/renderer.go | Draw frame |
| `Handler.PollEvent()` | method | internal/input/handler.go | Input abstraction |

## CONVENTIONS

**Go Standards:**
- Standard Go project layout (`cmd/`, `internal/`, `pkg/`)
- Co-located tests (`*_test.go`)
- Constructor pattern: `NewType()` functions
- Error wrapping: `fmt.Errorf("...: %w", err)`

**Project-Specific:**
- Grid indexing: `density[y][x]` (row-major)
- Double buffering: `density`/`densityPrev` fields
- Clamp values to [0, 1] before rendering
- 60 FPS target with `time.NewTicker(frameTime)`

## ANTI-PATTERNS (THIS PROJECT)

1. **NO binaries in repo** - `fluidx` binary in root is accidental (not in .gitignore)
2. **NO CI/CD** - No GitHub Actions, Makefile, or automated builds
3. **NO .gitignore** - Binary and `codedb.snapshot` committed
4. **NO tests for renderer/input** - Only solver has tests

## UNIQUE STYLES

**Physics Implementation:**
- Jos Stam's Stable Fluids method
- Gauss-Seidel relaxation (20 iterations default)
- Bilinear interpolation in advection
- No-slip boundary conditions (velocity negated at walls)

**Rendering:**
- ASCII density ramp: `" .:-=+*#%@"`
- TrueColor velocity visualization (blue→green→yellow→red)
- Debug mode: velocity arrows (→↘↓↙←↖↑↗)

## COMMANDS

```bash
# Build
go build ./cmd/fluidx

# Run
./fluidx

# Test
go test ./...

# Install from remote
go install github.com/yourusername/fluidx/cmd/fluidx@latest
```

## CONTROLS

| Input | Action |
|-------|--------|
| Mouse move | Inject velocity (stir fluid) |
| Mouse click/hold | Inject density (smoke) |
| Arrow keys | Inject velocity at center |
| Space | Inject density at center |
| `d` | Toggle debug (velocity arrows) |
| `q`/ESC/Ctrl+C | Quit |

## NOTES

- **Empty pkg/fluids/** - Reserved for future public API types
- **Resize support** - Terminal resize reinitializes solver grid
- **Performance** - 60 FPS cap prevents excessive CPU
- **Dependencies** - Only `tcell/v2` (terminal UI) + `go-colorful` (colors)
