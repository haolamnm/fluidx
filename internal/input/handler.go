package input

import (
	"github.com/gdamore/tcell/v2"
)

// EventType represents different input events.
type EventType int

const (
	EventNone EventType = iota
	EventQuit
	EventToggleDebug
	EventResize
	EventMouseMove
	EventMouseClick
	EventKeyArrow
	EventKeySpace
)

// Event represents an input event with associated data.
type Event struct {
	Type   EventType
	X, Y   int
	VX, VY float64
}

// PollEvent blocks on the screen's PollEvent and converts tcell events into
// application-level Events. Returns (Event, false) normally; returns
// (Event{}, true) when the screen has been closed and PollEvent returned nil.
func PollEvent(screen tcell.Screen) (Event, bool) {
	ev := screen.PollEvent()
	if ev == nil {
		return Event{}, true
	}

	switch e := ev.(type) {
	case *tcell.EventKey:
		return handleKey(e), false
	case *tcell.EventMouse:
		return handleMouse(e), false
	case *tcell.EventResize:
		return Event{Type: EventResize}, false
	}
	return Event{}, false
}

var lastMouseX, lastMouseY = -1, -1

func handleKey(e *tcell.EventKey) Event {
	switch e.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return Event{Type: EventQuit}
	case tcell.KeyRune:
		switch e.Rune() {
		case 'q', 'Q':
			return Event{Type: EventQuit}
		case 'd', 'D':
			return Event{Type: EventToggleDebug}
		case ' ':
			return Event{Type: EventKeySpace}
		}
	case tcell.KeyUp:
		return Event{Type: EventKeyArrow, VX: 0, VY: -3}
	case tcell.KeyDown:
		return Event{Type: EventKeyArrow, VX: 0, VY: 3}
	case tcell.KeyLeft:
		return Event{Type: EventKeyArrow, VX: -3, VY: 0}
	case tcell.KeyRight:
		return Event{Type: EventKeyArrow, VX: 3, VY: 0}
	}
	return Event{Type: EventNone}
}

func handleMouse(e *tcell.EventMouse) Event {
	x, y := e.Position()
	buttons := e.Buttons()

	var vx, vy float64
	if lastMouseX >= 0 && lastMouseY >= 0 {
		vx = float64(x-lastMouseX) * 0.8
		vy = float64(y-lastMouseY) * 0.8
	}
	lastMouseX = x
	lastMouseY = y

	if buttons&tcell.Button1 != 0 {
		return Event{Type: EventMouseClick, X: x, Y: y, VX: vx, VY: vy}
	}
	if vx != 0 || vy != 0 {
		return Event{Type: EventMouseMove, X: x, Y: y, VX: vx, VY: vy}
	}
	return Event{Type: EventNone}
}
