package solver

// diffuse performs Gauss-Seidel relaxation for diffusion
// Reads from 'prev', writes to 'curr'
func (s *FluidSolver) diffuse(prev, curr [][]float64, diff float64) {
	// Copy prev to curr as starting point
	for y := 0; y < s.height; y++ {
		copy(curr[y], prev[y])
	}

	// Gauss-Seidel relaxation using pre-computed constant
	a := diff * s.dtWidthHeight
	for k := 0; k < s.iterations; k++ {
		for y := 1; y < s.height-1; y++ {
			for x := 1; x < s.width-1; x++ {
				curr[y][x] = (prev[y][x] + a*(curr[y][x-1]+curr[y][x+1]+curr[y-1][x]+curr[y+1][x])) / (1 + 4*a)
			}
		}
		s.setBoundary(curr)
	}
}

// diffuseVelocity diffuses velocity components directly without copying to scalar fields
func (s *FluidSolver) diffuseVelocity(diff float64) {
	a := diff * s.dtWidthHeight

	// Store original velocity in prev fields for Gauss-Seidel
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.vxFieldPrev[y][x] = s.velocity[y][x][0]
			s.vyFieldPrev[y][x] = s.velocity[y][x][1]
		}
	}

	// Diffuse vx component
	for k := 0; k < s.iterations; k++ {
		for y := 1; y < s.height-1; y++ {
			for x := 1; x < s.width-1; x++ {
				s.velocity[y][x][0] = (s.vxFieldPrev[y][x] + a*(s.velocity[y][x-1][0]+s.velocity[y][x+1][0]+s.velocity[y-1][x][0]+s.velocity[y+1][x][0])) / (1 + 4*a)
			}
		}
		s.setBoundaryVel(s.velocity)
	}

	// Diffuse vy component
	for k := 0; k < s.iterations; k++ {
		for y := 1; y < s.height-1; y++ {
			for x := 1; x < s.width-1; x++ {
				s.velocity[y][x][1] = (s.vyFieldPrev[y][x] + a*(s.velocity[y][x-1][1]+s.velocity[y][x+1][1]+s.velocity[y-1][x][1]+s.velocity[y+1][x][1])) / (1 + 4*a)
			}
		}
		s.setBoundaryVel(s.velocity)
	}
}

// advect moves quantities along the velocity field (backwards trace)
func (s *FluidSolver) advect(d, dPrev [][]float64, vel [][][2]float64) {
	dt0 := s.dtWidth

	for y := 1; y < s.height-1; y++ {
		for x := 1; x < s.width-1; x++ {
			// Trace particle backwards
			x0 := float64(x) - dt0*vel[y][x][0]
			y0 := float64(y) - dt0*vel[y][x][1]

			// Clamp to boundaries
			if x0 < 0.5 {
				x0 = 0.5
			}
			if x0 > float64(s.width)-1.5 {
				x0 = float64(s.width) - 1.5
			}
			if y0 < 0.5 {
				y0 = 0.5
			}
			if y0 > float64(s.height)-1.5 {
				y0 = float64(s.height) - 1.5
			}

			// Bilinear interpolation
			i0 := int(x0)
			i1 := i0 + 1
			j0 := int(y0)
			j1 := j0 + 1

			s1 := x0 - float64(i0)
			s0 := 1 - s1
			t1 := y0 - float64(j0)
			t0 := 1 - t1

			d[y][x] = s0*(t0*dPrev[j0][i0]+t1*dPrev[j1][i0]) +
				s1*(t0*dPrev[j0][i1]+t1*dPrev[j1][i1])
		}
	}
	s.setBoundary(d)
}

// project makes the velocity field mass-conserving (incompressible)
func (s *FluidSolver) project(vel [][][2]float64, p [][]float64) {
	// Use pre-allocated divergence field (zero it first)
	div := s.divField
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			div[y][x] = 0
		}
	}

	h := s.invWidth

	for y := 1; y < s.height-1; y++ {
		for x := 1; x < s.width-1; x++ {
			div[y][x] = -0.5 * h * (vel[y][x+1][0] - vel[y][x-1][0] + vel[y+1][x][1] - vel[y-1][x][1])
			p[y][x] = 0
		}
	}
	s.setBoundary(div)
	s.setBoundary(p)

	// Gauss-Seidel relaxation for pressure
	for k := 0; k < s.iterations; k++ {
		for y := 1; y < s.height-1; y++ {
			for x := 1; x < s.width-1; x++ {
				p[y][x] = (div[y][x] + p[y][x-1] + p[y][x+1] + p[y-1][x] + p[y+1][x]) / 4
			}
		}
		s.setBoundary(p)
	}

	// Subtract pressure gradient from velocity
	for y := 1; y < s.height-1; y++ {
		for x := 1; x < s.width-1; x++ {
			vel[y][x][0] -= 0.5 * (p[y][x+1] - p[y][x-1]) / h
			vel[y][x][1] -= 0.5 * (p[y+1][x] - p[y-1][x]) / h
		}
	}
	s.setBoundaryVel(vel)
}

// setBoundary handles boundary conditions (continuity)
func (s *FluidSolver) setBoundary(x [][]float64) {
	// Horizontal boundaries
	for i := 1; i < s.width-1; i++ {
		x[0][i] = x[1][i]
		x[s.height-1][i] = x[s.height-2][i]
	}
	// Vertical boundaries
	for j := 1; j < s.height-1; j++ {
		x[j][0] = x[j][1]
		x[j][s.width-1] = x[j][s.width-2]
	}
	// Corners
	x[0][0] = 0.5 * (x[1][0] + x[0][1])
	x[0][s.width-1] = 0.5 * (x[1][s.width-1] + x[0][s.width-2])
	x[s.height-1][0] = 0.5 * (x[s.height-2][0] + x[s.height-1][1])
	x[s.height-1][s.width-1] = 0.5 * (x[s.height-2][s.width-1] + x[s.height-1][s.width-2])
}

// setBoundaryVel handles velocity boundary conditions (no-slip)
func (s *FluidSolver) setBoundaryVel(vel [][][2]float64) {
	// Horizontal boundaries - negate normal component
	for i := 1; i < s.width-1; i++ {
		vel[0][i][1] = -vel[1][i][1]
		vel[s.height-1][i][1] = -vel[s.height-2][i][1]
		vel[0][i][0] = vel[1][i][0]
		vel[s.height-1][i][0] = vel[s.height-2][i][0]
	}
	// Vertical boundaries
	for j := 1; j < s.height-1; j++ {
		vel[j][0][0] = -vel[j][1][0]
		vel[j][s.width-1][0] = -vel[j][s.width-2][0]
		vel[j][0][1] = vel[j][1][1]
		vel[j][s.width-1][1] = vel[j][s.width-2][1]
	}
	// Corners
	vel[0][0][0] = 0.5 * (vel[1][0][0] + vel[0][1][0])
	vel[0][0][1] = 0.5 * (vel[1][0][1] + vel[0][1][1])
	vel[0][s.width-1][0] = 0.5 * (vel[1][s.width-1][0] + vel[0][s.width-2][0])
	vel[0][s.width-1][1] = 0.5 * (vel[1][s.width-1][1] + vel[0][s.width-2][1])
	vel[s.height-1][0][0] = 0.5 * (vel[s.height-2][0][0] + vel[s.height-1][1][0])
	vel[s.height-1][0][1] = 0.5 * (vel[s.height-2][0][1] + vel[s.height-1][1][1])
	vel[s.height-1][s.width-1][0] = 0.5 * (vel[s.height-2][s.width-1][0] + vel[s.height-1][s.width-2][0])
	vel[s.height-1][s.width-1][1] = 0.5 * (vel[s.height-2][s.width-1][1] + vel[s.height-1][s.width-2][1])
}
