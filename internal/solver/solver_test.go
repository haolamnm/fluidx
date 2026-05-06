package solver

import (
	"testing"
)

func TestFluidSolverBasic(t *testing.T) {
	s := NewFluidSolver(10, 10)

	// Test adding density
	s.AddDensity(5, 5, 1.0)
	if s.density[5][5] != 1.0 {
		t.Errorf("Expected density at (5,5) to be 1.0, got %f", s.density[5][5])
	}

	// Test adding velocity
	s.AddVelocity(5, 5, 2.0, 3.0)
	if s.velocity[5][5][0] != 2.0 || s.velocity[5][5][1] != 3.0 {
		t.Errorf("Expected velocity at (5,5) to be (2.0, 3.0), got (%f, %f)",
			s.velocity[5][5][0], s.velocity[5][5][1])
	}

	// Test that Step doesn't panic
	s.Step()

	// After step, density should have changed (diffused/advected)
	// It might be slightly different due to diffusion
	t.Logf("Density at (5,5) after step: %f", s.density[5][5])
}

func TestFluidSolverDensityPropagation(t *testing.T) {
	s := NewFluidSolver(20, 20)

	// Add density at center
	s.AddDensity(10, 10, 10.0)

	// Add some velocity to move it
	for y := 8; y < 12; y++ {
		for x := 8; x < 12; x++ {
			s.AddVelocity(x, y, 1.0, 0.5)
		}
	}

	// Run several steps
	for i := 0; i < 10; i++ {
		s.Step()
	}

	// Check that density has spread (not all at center anymore)
	centerDensity := s.density[10][10]
	totalDensity := 0.0
	maxDensity := 0.0

	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			d := s.density[y][x]
			totalDensity += d
			if d > maxDensity {
				maxDensity = d
			}
		}
	}

	t.Logf("Center density after 10 steps: %f", centerDensity)
	t.Logf("Max density: %f", maxDensity)
	t.Logf("Total density: %f", totalDensity)

	// Density should have spread from center
	if centerDensity == 10.0 {
		t.Error("Density didn't diffuse from center")
	}
}
