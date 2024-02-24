package main

// I want a measurement system so that 1 "meter" equals 32 texels
// so a player with a speed of 1 should travel 32 texels per second

const (
	PHYSICS_PLAYER_SPEEDMAX = 0.5
	PHYSICS_FRICTION_COEF   = float32(0.5)
	PHYSICS_MIN_SPEED       = float32(0.0001)
	ACCEL_PLAYER_DEFAULT    = 0.1
)

type TexelUnit float32
