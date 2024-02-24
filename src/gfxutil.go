package main

const (
	windowTitle = "Gaming"

	windowWidth    = 1280
	windowHeight   = 720
	pixelsPerTexel = 4
	// windowWidth    = 1920
	// windowHeight   = 1080
	// pixelsPerTexel = 6

	// these convert world coordinates to screen coordinates, corrected for aspect ratio
	WORLD_SCALE_RATIO = float32(3.2)
	TEXEL_SCALE_X     = float32(1.0 / (16.0 / WORLD_SCALE_RATIO))
	TEXEL_SCALE_Y     = float32(1.0 / (9.0 / WORLD_SCALE_RATIO))
)
