package main

import (
	"fmt"
	"os"
	"time"

	"fluidx/internal/input"
	"fluidx/internal/renderer"
	"fluidx/internal/solver"
)

const (
	targetFPS = 60
	frameTime = time.Second / targetFPS
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	rend, err := renderer.NewRenderer()
	if err != nil {
		return fmt.Errorf("failed to create renderer: %w", err)
	}
	defer rend.Cleanup()

	screen := rend.GetScreen()
	width, height := screen.Size()

	solv := solver.NewFluidSolver(width, height)
	screen.EnableMouse()

	eventChan := make(chan input.Event, 10)
	quitChan := make(chan struct{})

	go func() {
		for {
			event, done := input.PollEvent(screen)
			if done {
				close(quitChan)
				return
			}
			if event.Type != input.EventNone {
				select {
				case eventChan <- event:
				case <-quitChan:
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(frameTime)
	defer ticker.Stop()

	debugMode := false
	running := true

	solv.AddDensity(width/2, height/2, 5.0)

	for running {
		select {
		case <-quitChan:
			running = false

		case ev := <-eventChan:
			switch ev.Type {
			case input.EventQuit:
				running = false

			case input.EventToggleDebug:
				debugMode = !debugMode
				rend.SetDebugMode(debugMode)
				rend.Clear()

			case input.EventResize:
				newWidth, newHeight := screen.Size()
				rend.Resize(newWidth, newHeight)
				solv.Resize(newWidth, newHeight)
				width, height = newWidth, newHeight

			case input.EventMouseMove:
				if ev.X >= 0 && ev.X < width && ev.Y >= 0 && ev.Y < height {
					solv.AddVelocity(ev.X, ev.Y, ev.VX, ev.VY)
				}

			case input.EventMouseClick:
				if ev.X >= 0 && ev.X < width && ev.Y >= 0 && ev.Y < height {
					solv.AddDensity(ev.X, ev.Y, 2.0)
					solv.AddVelocity(ev.X, ev.Y, ev.VX, ev.VY)
				}

			case input.EventKeyArrow:
				centerX, centerY := width/2, height/2
				solv.AddVelocity(centerX, centerY, ev.VX, ev.VY)

			case input.EventKeySpace:
				centerX, centerY := width/2, height/2
				solv.AddDensity(centerX, centerY, 3.0)
			}

		case <-ticker.C:
			solv.Step()
			rend.Render(solv.GetDensity(), solv.GetVelocity())
		}
	}

	return nil
}
