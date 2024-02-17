package main

import (
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func initControls(window *glfw.Window, playerPointer *DrawableEntity) {
	window.SetKeyCallback(playerControlsCallback)
	window.SetUserPointer(unsafe.Pointer(playerPointer))
}

func playerControlsCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// fmt.Println(key, scancode, action, mods)
	if action == glfw.Repeat {
		return
	}

	var accel float32
	if action == glfw.Release {
		accel = -0.05
	} else {
		accel = 0.05
	}

	playerPointer := window.GetUserPointer()
	player := (*DrawableEntity)(playerPointer)
	switch key {
	case glfw.KeyW:
		player.accel[1] += accel
	case glfw.KeyS:
		player.accel[1] -= accel
	case glfw.KeyA:
		player.accel[0] -= accel
	case glfw.KeyD:
		player.accel[0] += accel
	}
}
