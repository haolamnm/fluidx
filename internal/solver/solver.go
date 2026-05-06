package solver

const decayRate = 0.002 // density dissipation per step (~6s half-life at 60 FPS)

// FluidSolver implements Jos Stam's Stable Fluids method for 2D Navier-Stokes
// Grid layout: density[y][x], velocity[y][x] = [vx, vy]
type FluidSolver struct {
	width, height int
	dt            float64
	diffusion     float64
	viscosity     float64
	iterations    int

	// Pre-computed constants for performance
	dtWidth   float64 // dt * width
	dtHeight  float64 // dt * height
	invWidth  float64 // 1.0 / width
	invHeight float64 // 1.0 / height

	// Current fields
	density  [][]float64
	velocity [][][2]float64

	// Previous fields (double buffering)
	densityPrev [][]float64

	// Temporary scalar fields for velocity operations
	vxField, vyField         [][]float64
	vxFieldPrev, vyFieldPrev [][]float64
	pressure                 [][]float64

	// Pre-allocated temporary field for divergence calculation
	divField [][]float64
}

// NewFluidSolver creates a solver with given grid dimensions
func NewFluidSolver(width, height int) *FluidSolver {
	s := &FluidSolver{
		width:      width,
		height:     height,
		dt:         0.1,
		diffusion:  0.001,
		viscosity:  0.001,
		iterations: 20,
	}

	s.updateConstants()
	s.allocateFields()
	return s
}

func (s *FluidSolver) updateConstants() {
	s.dtWidth = s.dt * float64(s.width)
	s.dtHeight = s.dt * float64(s.height)
	s.invWidth = 1.0 / float64(s.width)
	s.invHeight = 1.0 / float64(s.height)
}

func (s *FluidSolver) allocateFields() {
	s.density = make([][]float64, s.height)
	s.densityPrev = make([][]float64, s.height)
	s.velocity = make([][][2]float64, s.height)
	s.vxField = make([][]float64, s.height)
	s.vyField = make([][]float64, s.height)
	s.vxFieldPrev = make([][]float64, s.height)
	s.vyFieldPrev = make([][]float64, s.height)
	s.pressure = make([][]float64, s.height)
	s.divField = make([][]float64, s.height)

	for y := 0; y < s.height; y++ {
		s.density[y] = make([]float64, s.width)
		s.densityPrev[y] = make([]float64, s.width)
		s.velocity[y] = make([][2]float64, s.width)
		s.vxField[y] = make([]float64, s.width)
		s.vyField[y] = make([]float64, s.width)
		s.vxFieldPrev[y] = make([]float64, s.width)
		s.vyFieldPrev[y] = make([]float64, s.width)
		s.pressure[y] = make([]float64, s.width)
		s.divField[y] = make([]float64, s.width)
	}
}

// Resize reinitializes the solver with new dimensions
func (s *FluidSolver) Resize(width, height int) {
	s.width = width
	s.height = height
	s.updateConstants()
	s.allocateFields()
}

// AddDensity injects density at position (x, y)
func (s *FluidSolver) AddDensity(x, y int, amount float64) {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return
	}
	s.density[y][x] += amount
}

// AddVelocity injects velocity at position (x, y)
func (s *FluidSolver) AddVelocity(x, y int, vx, vy float64) {
	if x < 0 || x >= s.width || y < 0 || y >= s.height {
		return
	}
	s.velocity[y][x][0] += vx
	s.velocity[y][x][1] += vy
}

// Step advances the simulation by one time step
func (s *FluidSolver) Step() {
	// Velocity step
	s.diffuseVelocity(s.viscosity)
	s.project(s.velocity, s.pressure)
	s.velocityToScalar(s.velocity, s.vxFieldPrev, s.vyFieldPrev)
	s.advect(s.vxField, s.vxFieldPrev, s.velocity)
	s.advect(s.vyField, s.vyFieldPrev, s.velocity)
	s.scalarToVelocity(s.vxField, s.vyField, s.velocity)
	s.project(s.velocity, s.pressure)

	// Density step
	s.diffuse(s.density, s.densityPrev, s.diffusion)
	s.advect(s.density, s.densityPrev, s.velocity)
	s.decayDensity()
}

// velocityToScalar extracts velocity components to scalar fields
func (s *FluidSolver) velocityToScalar(vel [][][2]float64, vx, vy [][]float64) {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			vx[y][x] = vel[y][x][0]
			vy[y][x] = vel[y][x][1]
		}
	}
}

// scalarToVelocity combines scalar fields back to velocity
func (s *FluidSolver) scalarToVelocity(vx, vy [][]float64, vel [][][2]float64) {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			vel[y][x][0] = vx[y][x]
			vel[y][x][1] = vy[y][x]
		}
	}
}

// Reset clears all fields to zero
func (s *FluidSolver) Reset() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.density[y][x] = 0
			s.velocity[y][x][0] = 0
			s.velocity[y][x][1] = 0
		}
	}
}

// decayDensity reduces density by decayRate to prevent saturation
func (s *FluidSolver) decayDensity() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.density[y][x] *= (1 - decayRate)
		}
	}
}

// GetDensity returns the density field
func (s *FluidSolver) GetDensity() [][]float64 {
	return s.density
}

// GetVelocity returns the velocity field
func (s *FluidSolver) GetVelocity() [][][2]float64 {
	return s.velocity
}

// GetDimensions returns grid width and height
func (s *FluidSolver) GetDimensions() (int, int) {
	return s.width, s.height
}
