package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	// "github.com/go-gl/mathgl/mgl32"
)

// type cControllable struct { // this is coupled to cMovable
// 	acceleration mgl32.Vec2
// }

func initControls(window *glfw.Window, playerPointer *PlayerEntity) {
	window.SetKeyCallback(playerControlsCallback)
}

func playerControlsCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Repeat {
		return
	}

	var accel float32
	if action == glfw.Release {
		accel = -0.05 // TODO magic number
	} else {
		accel = 0.05
	}

	playerPointer := getPlayerPtr()
	player := (*PlayerEntity)(playerPointer)

	var moveComponent *cMovable
	if tmp, err := getComponent[*cMovable](CMP_MOVABLE, player); err != nil {
		return
	} else {
		moveComponent = *tmp
	}
	switch key {
	case glfw.KeyW:
		moveComponent.accel[1] += accel
	case glfw.KeyS:
		moveComponent.accel[1] -= accel
	case glfw.KeyA:
		moveComponent.accel[0] -= accel
	case glfw.KeyD:
		moveComponent.accel[0] += accel
	}
}
