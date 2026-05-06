package renderer

import (
	"github.com/gdamore/tcell/v2"
)

// ASCII density ramp from light to dark (perceptual, ~70 chars)
const densityRamp = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

// maxSpeedSq caps the color scale (max speed ≈ 5.0)
const maxSpeedSq = 25.0

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
				style = tcell.StyleDefault.Foreground(r.velocityToColor(velRow[x], 1.0))
			} else {
				ch = r.densityToChar(d)
				style = tcell.StyleDefault.Foreground(r.velocityToColor(velRow[x], d))
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

// velocityToColor maps velocity and density to RGB color.
// Density controls brightness, velocity controls hue shift.
func (r *Renderer) velocityToColor(v [2]float64, d float64) tcell.Color {
	speedSq := v[0]*v[0] + v[1]*v[1]

	t := speedSq / maxSpeedSq
	if t > 1 {
		t = 1
	}
	if t < 0 {
		t = 0
	}
	if d < 0 {
		d = 0
	} else if d > 1 {
		d = 1
	}

	if d < 0.001 {
		return tcell.ColorDefault
	}

	var vr, vg, vb int32
	if t < 0.5 {
		vr = 0
		vg = int32(t * 2 * 255)
		vb = int32((1 - t*2) * 255)
	} else {
		vr = int32((t - 0.5) * 2 * 255)
		vg = int32((1 - t) * 255 * 2)
		vb = 0
	}

	bright := int32(d * 255)
	rOut := int32(float64(vr)*t + float64(bright)*(1-t))
	gOut := int32(float64(vg)*t + float64(bright)*(1-t))
	bOut := int32(float64(vb)*t + float64(bright)*(1-t))

	if rOut > 255 {
		rOut = 255
	}
	if gOut > 255 {
		gOut = 255
	}
	if bOut > 255 {
		bOut = 255
	}
	return tcell.NewRGBColor(rOut, gOut, bOut)
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
