package solver

import (
	"testing"
)

func TestFluidSolverBasic(t *testing.T) {
	s := NewFluidSolver(10, 10)

	s.AddDensity(5, 5, 1.0)
	if s.density[5][5] != 1.0 {
		t.Errorf("Expected density at (5,5) to be 1.0, got %f", s.density[5][5])
	}

	s.AddVelocity(5, 5, 2.0, 3.0)
	if s.velocity[5][5][0] != 2.0 || s.velocity[5][5][1] != 3.0 {
		t.Errorf("Expected velocity at (5,5) to be (2.0, 3.0), got (%f, %f)",
			s.velocity[5][5][0], s.velocity[5][5][1])
	}

	s.Step()
}

func TestFluidSolverDensityPropagation(t *testing.T) {
	s := NewFluidSolver(20, 20)

	s.AddDensity(10, 10, 10.0)

	for y := 8; y < 12; y++ {
		for x := 8; x < 12; x++ {
			s.AddVelocity(x, y, 1.0, 0.5)
		}
	}

	for i := 0; i < 10; i++ {
		s.Step()
	}

	centerDensity := s.density[10][10]
	maxDensity := 0.0

	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			d := s.density[y][x]
			if d > maxDensity {
				maxDensity = d
			}
		}
	}

	if centerDensity == 10.0 {
		t.Error("Density didn't diffuse from center")
	}
	if maxDensity > 10.0 {
		t.Errorf("Density grew beyond initial max: %f", maxDensity)
	}
}

func TestFluidSolverVelocityAdvection(t *testing.T) {
	s := NewFluidSolver(20, 20)

	s.AddVelocity(10, 10, 5.0, 0)

	for i := 0; i < 5; i++ {
		s.Step()
	}

	center := s.velocity[10][10][0]
	if center > 4.5 {
		t.Errorf("Center velocity didn't advect away: vx at (10,10) = %f (expected < 4.5)", center)
	}

	rightOf := s.velocity[10][12][0]
	if rightOf < 0.05 {
		t.Errorf("Velocity didn't advect rightward: vx at (12,10) = %f (expected >= 0.05)", rightOf)
	}
}

func TestFluidSolverVelocityDirection(t *testing.T) {
	s := NewFluidSolver(60, 60)

	blobSize := 5
	cx, cy := 30, 30
	for y := cy - blobSize/2; y <= cy+blobSize/2; y++ {
		for x := cx - blobSize/2; x <= cx+blobSize/2; x++ {
			s.AddVelocity(x, y, 0, -3.0)
		}
	}

	for i := 0; i < 8; i++ {
		s.Step()
	}

	argMaxY := cy
	maxVy := 0.0
	for y := 0; y < 60; y++ {
		for x := 0; x < 60; x++ {
			vy := s.velocity[y][x][1]
			if vy < maxVy {
				maxVy = vy
				argMaxY = y
			}
		}
	}

	if argMaxY >= cy {
		t.Errorf("Velocity peak didn't move upward: argmax y=%d, injection zone center y=%d", argMaxY, cy)
	}

	totalVyAbove := 0.0
	for y := 0; y < cy-blobSize; y++ {
		for x := 0; x < 60; x++ {
			if s.velocity[y][x][1] < 0 {
				totalVyAbove += -s.velocity[y][x][1]
			}
		}
	}
	if totalVyAbove < 0.01 {
		t.Errorf("No upward velocity in cells above injection zone: totalVyAbove=%f", totalVyAbove)
	}
}

func TestFluidSolverDensityDecay(t *testing.T) {
	s := NewFluidSolver(40, 40)

	s.AddDensity(18, 18, 5.0)
	s.AddDensity(20, 20, 10.0)
	s.AddDensity(22, 22, 5.0)

	initialSum := sumDensity(s)

	const steps = 80
	for i := 0; i < steps; i++ {
		s.Step()
	}

	decayedSum := sumDensity(s)

	if decayedSum > initialSum*1.03 {
		t.Errorf("Density grew instead of decaying: initial=%f, final=%f", initialSum, decayedSum)
	}
	if decayedSum > initialSum*0.90 {
		t.Errorf("Density decayed too slowly: ratio=%f (expected < 0.90 after %d steps)", decayedSum/initialSum, steps)
	}
	if decayedSum < initialSum*0.70 {
		t.Errorf("Density decayed too fast: ratio=%f (expected >= 0.70 after %d steps)", decayedSum/initialSum, steps)
	}
}

func TestFluidSolverFieldSanity(t *testing.T) {
	s := NewFluidSolver(30, 30)

	for x := 5; x <= 25; x++ {
		for y := 5; y <= 25; y++ {
			s.AddDensity(x, y, 10.0)
			s.AddVelocity(x, y, 2.0, -1.0)
		}
	}

	for i := 0; i < 20; i++ {
		s.Step()
	}

	hasNonzeroDensity := false
	hasNonzeroDX := false
	hasPositiveDY := false

	for y := 0; y < 30; y++ {
		for x := 0; x < 30; x++ {
			if s.density[y][x] > 0.001 {
				hasNonzeroDensity = true
			}
			if s.velocity[y][x][0] > 0.001 {
				hasNonzeroDX = true
			}
			if s.velocity[y][x][1] > 0.001 {
				hasPositiveDY = true
			}
		}
	}

	if !hasNonzeroDensity {
		t.Error("All density vanished prematurely")
	}
	if !hasNonzeroDX {
		t.Error("All x-velocity vanished — possible dead field")
	}
	if !hasPositiveDY {
		t.Error("No upward velocity component — boundary handling may be absorbing all motion")
	}
}

func TestFluidSolverReset(t *testing.T) {
	s := NewFluidSolver(10, 10)

	s.AddDensity(5, 5, 1.0)
	s.AddVelocity(5, 5, 2.0, 3.0)

	s.Reset()

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if s.density[y][x] != 0 {
				t.Errorf("density[%d][%d] = %f after Reset", y, x, s.density[y][x])
			}
			if s.velocity[y][x][0] != 0 {
				t.Errorf("velocity[%d][%d][0] = %f after Reset", y, x, s.velocity[y][x][0])
			}
			if s.velocity[y][x][1] != 0 {
				t.Errorf("velocity[%d][%d][1] = %f after Reset", y, x, s.velocity[y][x][1])
			}
		}
	}
}

func sumDensity(s *FluidSolver) float64 {
	sum := 0.0
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			sum += s.density[y][x]
		}
	}
	return sum
}
