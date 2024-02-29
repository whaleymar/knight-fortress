package phys

import "github.com/go-gl/mathgl/mgl32"

// This is a measurement system where 1 "meter" equals 32 texels
// so an entity with a speed of 1 should travel 32 texels per second

const (
	PHYSICS_PLAYER_SPEEDMAX = float32(4.0)
	PHYSICS_FRICTION_COEF   = float32(0.5)
	PHYSICS_MIN_SPEED       = float32(0.0001)
	ACCEL_PLAYER_DEFAULT    = float32(0.33)
)

var PHYSICS_GRAVITY Force = &ImpulseForce{mgl32.Vec3{0, -1, 0}, 0.2}
var ORIGIN mgl32.Vec3 = mgl32.Vec3{0.0, 0.0, 0.0}

type Point struct {
	X float32
	Y float32
}

func (p Point) Add(p2 Point) Point {
	return Point{p.X + p2.X, p.Y + p2.Y}
}

func (p Point) Sub(p2 Point) Point {
	return Point{p.X - p2.X, p.Y - p2.Y}
}

func Vec2Point(vec mgl32.Vec3) Point {
	return Point{vec[0], vec[1]}
}
