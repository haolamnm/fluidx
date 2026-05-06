# internal/solver - Physics Engine

**Package:** `solver`  
**Domain:** Navier-Stokes fluid simulation  
**Complexity:** High (numerical methods, double buffering)

## OVERVIEW
Implements Jos Stam's "Stable Fluids" method for 2D incompressible flow. Uses grid-based discretization with Gauss-Seidel relaxation for diffusion and pressure projection.

## STRUCTURE
```
solver/
├── solver.go       # FluidSolver type, public API
├── operations.go   # Internal: diffuse, advect, project
└── solver_test.go  # Unit tests (only tests in project)
```

## WHERE TO LOOK

| Task | Function | File |
|------|----------|------|
| Create solver | `NewFluidSolver(width, height)` | solver.go:27 |
| Simulation step | `FluidSolver.Step()` | solver.go:90 |
| Add density | `AddDensity(x, y, amount)` | solver.go:73 |
| Add velocity | `AddVelocity(x, y, vx, vy)` | solver.go:81 |
| Resize grid | `Resize(width, height)` | solver.go:66 |
| Get state | `GetDensity()`, `GetVelocity()` | solver.go:145 |

## CORE ALGORITHMS

**Step() Pipeline:**
1. Velocity diffuse (viscosity) → Gauss-Seidel
2. Velocity project (incompressibility)
3. Velocity advect (self-advection)
4. Velocity project again
5. Density diffuse
6. Density advect

**Key Parameters:**
- `dt = 0.1` - Time step
- `diffusion = 0.001` - Density diffusion rate
- `viscosity = 0.001` - Velocity viscosity
- `iterations = 20` - Gauss-Seidel iterations

## INTERNAL DETAILS

**Grid Layout:**
```go
density  [][]float64        // [y][x] - row major
velocity [][][2]float64     // [y][x][0=vx, 1=vy]
```

**Double Buffering:**
- `density` / `densityPrev` - Prevents artifacts during diffusion
- `velocity` / `velocityPrev` - Same for velocity
- Swapped implicitly via parameter passing

**Boundary Conditions:**
- Scalar fields: `setBoundary()` - Continuity (copy neighbor)
- Velocity: `setBoundaryVel()` - No-slip (negate normal component)

## CONVENTIONS

- **Grid indexing:** Always `density[y][x]` (not [x][y])
- **Bounds checking:** All public methods validate coordinates
- **Clamping:** Values clamped to valid ranges internally
- **Allocation:** `allocateFields()` creates all slices

## ANTI-PATTERNS

1. **DON'T** modify fields directly - use `AddDensity()`/`AddVelocity()`
2. **DON'T** assume grid size persists - call `Resize()` on terminal resize
3. **DON'T** ignore return values from `GetDensity()` - may be reallocated

## TESTING

```bash
go test -v ./internal/solver/
```

**Test Functions:**
- `TestFluidSolverBasic` - Injection and stepping
- `TestFluidSolverDensityPropagation` - Diffusion over 10 steps

**Note:** No benchmarks yet. Consider adding for `Step()` performance.
