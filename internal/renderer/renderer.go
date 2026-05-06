package renderer

import (
	"github.com/gdamore/tcell/v2"
)

// ASCII density ramp from light to dark
const densityRamp = " .:-=+*#%@"

// Renderer handles fluid visualization using tcell
type Renderer struct {
	screen    tcell.Screen
	width     int
	height    int
	debugMode bool
}

// NewRenderer creates a new renderer
func NewRenderer() (*Renderer, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()

	width, height := screen.Size()

	return &Renderer{
		screen:    screen,
		width:     width,
		height:    height,
		debugMode: false,
	}, nil
}

// SetDebugMode toggles debug overlay
func (r *Renderer) SetDebugMode(enabled bool) {
	r.debugMode = enabled
}

// Resize updates renderer dimensions
func (r *Renderer) Resize(width, height int) {
	r.width = width
	r.height = height
	r.screen.Clear()
}

// Render draws the fluid state
func (r *Renderer) Render(density [][]float64, velocity [][][2]float64) {
	height := r.height
	width := r.width
	if len(density) < height {
		height = len(density)
	}

	for y := 0; y < height; y++ {
		row := density[y]
		velRow := velocity[y]
		rowWidth := width
		if len(row) < rowWidth {
			rowWidth = len(row)
		}

		for x := 0; x < rowWidth; x++ {
			d := row[x]

			if d < 0 {
				d = 0
			} else if d > 1 {
				d = 1
			}

			var ch rune
			var style tcell.Style

			if r.debugMode {
				ch = r.velocityToArrow(velRow[x])
				style = tcell.StyleDefault.Foreground(r.velocityToColor(velRow[x]))
			} else {
				ch = r.densityToChar(d)
				style = tcell.StyleDefault.Foreground(r.velocityToColor(velRow[x]))
			}

			r.screen.SetContent(x, y, ch, nil, style)
		}
	}
	r.screen.Show()
}

// densityToChar maps density (0-1) to ASCII character
func (r *Renderer) densityToChar(d float64) rune {
	idx := int(d * float64(len(densityRamp)-1))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(densityRamp) {
		idx = len(densityRamp) - 1
	}
	return rune(densityRamp[idx])
}

// velocityToColor maps velocity to RGB color using squared speed (avoids Sqrt)
func (r *Renderer) velocityToColor(v [2]float64) tcell.Color {
	speedSq := v[0]*v[0] + v[1]*v[1]
	maxSpeedSq := 25.0

	t := speedSq / maxSpeedSq
	if t > 1 {
		t = 1
	}

	var red, green, blue int32

	if t < 0.5 {
		red = 0
		green = int32(t * 2 * 255)
		blue = int32((1 - t*2) * 255)
	} else {
		red = int32((t - 0.5) * 2 * 255)
		green = int32((1 - t) * 255 * 2)
		blue = 0
	}

	return tcell.NewRGBColor(red, green, blue)
}

// velocityToArrow converts velocity vector to arrow character
func (r *Renderer) velocityToArrow(v [2]float64) rune {
	vx, vy := v[0], v[1]
	if vx == 0 && vy == 0 {
		return '·'
	}

	angle := fastAtan2(vy, vx)
	direction := int((angle+0.3927)/0.7854) & 7

	arrows := []rune{'→', '↘', '↓', '↙', '←', '↖', '↑', '↗'}
	return arrows[direction]
}

// fastAtan2 returns approximate angle in radians
func fastAtan2(y, x float64) float64 {
	if x == 0 {
		if y > 0 {
			return 1.5708
		}
		return -1.5708
	}

	absY := y
	if absY < 0 {
		absY = -absY
	}

	angle := 0.0
	if x > 0 {
		z := absY / x
		angle = z / (1 + 0.28*z*z)
		if y < 0 {
			angle = -angle
		}
	} else {
		z := absY / -x
		angle = 3.14159 - z/(1+0.28*z*z)
		if y < 0 {
			angle = -3.14159 + z/(1+0.28*z*z)
		}
	}

	return angle
}

// GetScreen returns the underlying tcell screen
func (r *Renderer) GetScreen() tcell.Screen {
	return r.screen
}

// Cleanup shuts down the renderer
func (r *Renderer) Cleanup() {
	r.screen.Fini()
}

// Clear clears the screen
func (r *Renderer) Clear() {
	r.screen.Clear()
}
