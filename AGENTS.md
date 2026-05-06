# AGENTS.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

**These guidelines are working if:** fewer unnecessary changes in diffs, fewer rewrites due to overcomplication, and clarifying questions come before implementation rather than after mistakes.

# Project: fluidx

Real-time 2D fluid dynamics simulator for the terminal (Jos Stam's Stable Fluids).

**Stack:** Go 1.22, tcell/v2 (terminal UI), go-colorful (colors)

## Structure

```
cmd/fluidx/main.go       # Entry, 60 FPS game loop, event routing
internal/solver/         # Navier-Stokes: advection, diffusion, pressure projection
  solver.go              # FluidSolver struct, public API, grid: density[y][x], velocity[y][x][0|1]
  operations.go          # diffuse(), advect(), project(), setBoundary(), setBoundaryVel()
  solver_test.go         # Tests for basic ops and density propagation
internal/renderer/       # ASCII ramp + TrueColor velocity visualization
  renderer.go            # Render(), densityToChar(), velocityToColor(), fastAtan2()
internal/input/          # tcell event → Event struct conversion
  handler.go             # PollEvent(screen) returns (Event, done)
pkg/fluids/              # Empty — reserved for public API
```

## Conventions

- Row-major grid: `field[y][x]`
- Double buffering: `density`/`densityPrev`, `velocity`/`velocityPrev`
- Clamp density [0,1] before rendering
- 60 FPS via `time.NewTicker(frameTime)`
- Gauss-Seidel relaxation, 20 iterations
- No-slip boundaries: normal velocity negated at walls
- Bilinear interpolation in advection
- Pre-computed constants: `dtWidth`, `dtWidthHeight`, `invWidth`
- Only solver has tests (renderer/input untested)
