package phys

// I want a measurement system so that 1 "meter" equals 32 texels
// so an entity with a speed of 1 should travel 32 texels per second

const (
	PHYSICS_PLAYER_SPEEDMAX = float32(2.0)
	PHYSICS_FRICTION_COEF   = float32(0.5)
	PHYSICS_MIN_SPEED       = float32(0.0001)
	ACCEL_PLAYER_DEFAULT    = float32(0.5)
)

type Point struct {
	X float32
	Y float32
}
