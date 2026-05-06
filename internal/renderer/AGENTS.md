# internal/renderer - TUI Rendering

**Package:** `renderer`  
**Domain:** Terminal UI, ASCII art, TrueColor  
**Complexity:** Medium (color gradients, character mapping)

## OVERVIEW
Handles all visual output using tcell. Maps fluid density to ASCII characters and velocity to TrueColor RGB gradients.

## STRUCTURE
```
renderer/
└── renderer.go    # Renderer type and methods
```

## WHERE TO LOOK

| Task | Function | Line |
|------|----------|------|
| Create renderer | `NewRenderer()` | 21 |
| Draw frame | `Render(density, velocity)` | 57 |
| Toggle debug | `SetDebugMode(bool)` | 45 |
| Handle resize | `Resize(w, h)` | 50 |
| Cleanup | `Cleanup()` | 177 |

## RENDERING MODES

**Normal Mode:**
- Density → ASCII character via `densityRamp`
- Velocity → RGB color via `velocityToColor()`

**Debug Mode:**
- Velocity → Arrow character (→↘↓↙←↖↑↗)
- Velocity → RGB color (same as normal)

## CONSTANTS

```go
const densityRamp = " .:-=+*#%@"  // Light to dark
```

## COLOR GRADIENT

**Speed → Color mapping:**
- 0.0-0.2: Blue → Cyan
- 0.2-0.4: Cyan → Green
- 0.4-0.6: Green → Yellow
- 0.6-0.8: Yellow → Orange
- 0.8-1.0: Orange → Red

Max speed normalized to 5.0 (see `velocityToColor()` line 111).

## VELOCITY ARROWS

8-directional arrows based on velocity angle:
- `atan2(vy, vx)` → direction index
- Arrows: `['→', '↘', '↓', '↙', '←', '↖', '↑', '↗']`
- Zero velocity: `·` (middle dot)

## CONVENTIONS

- **Density clamping:** Values clamped to [0, 1] before mapping
- **Index clamping:** Always bounds-check before accessing `densityRamp`
- **Screen clearing:** Called automatically on resize
- **Show() required:** Must call `screen.Show()` after `SetContent()`

## ANTI-PATTERNS

1. **DON'T** call `Render()` without `Show()` - won't display
2. **DON'T** assume density/velocity arrays match screen size - always check bounds
3. **DON'T** modify `densityRamp` - used directly in rendering loop
