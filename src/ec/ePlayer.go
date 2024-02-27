package ec

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/phys"
)

const (
	SPRITE_LOC_PLAYER_X       = 0
	SPRITE_LOC_PLAYER_Y       = 0
	SPRITE_LOC_PLAYER_Z       = 0
	SPRITE_HEIGHT_PLAYER      = 16
	SPRITE_WIDTH_PLAYER       = 16
	ANIM_OFFSET_PLAYER_HRIGHT = 1
	ANIM_OFFSET_PLAYER_HLEFT  = 3
	ANIM_OFFSET_PLAYER_VDOWN  = 0
	ANIM_OFFSET_PLAYER_VUP    = 2
	ANIM_FRAMES_PLAYER_H      = 2
	ANIM_FRAMES_PLAYER_V      = 4
	ANIM_FPS_DEFAULT          = 4.0
)

var _PLAYER_LOCK = &sync.Mutex{}
var playerPtr *Entity

func GetPlayerPtr() *Entity {
	if playerPtr == nil {
		_PLAYER_LOCK.Lock()
		defer _PLAYER_LOCK.Unlock()
		if playerPtr == nil {
			entity := makePlayerEntity()
			playerPtr = &entity
		}
	}
	return playerPtr
}

func makePlayerEntity() Entity {
	entity := Entity{
		0,
		"Player",
		&ComponentList{},
		mgl32.Vec3{0.0, 0.0, DEPTH_PLAYER},
		&sync.RWMutex{},
	}

	entity.components.Add(&CDrawable{
		gfx.SquareVertices,
		gfx.MakeVao(),
		gfx.MakeVbo(),
		SCALE_NORMAL,
		makePlayerSprite(makePlayerAnimationManager()),
		gfx.TEX_MAIN,
		&sync.RWMutex{},
		true, // isUvUpdateNeeded
	})

	entity.components.Add(&CMovable{
		mgl32.Vec3{},
		mgl32.Vec3{},
		phys.PHYSICS_PLAYER_SPEEDMAX,
		true, // frictionActive
		nil,
		true,
	})

	entity.components.Add(&CCollides{
		&phys.AABB{phys.Point{}, phys.Point{0.25, 0.25}}, // TODO hard coded size -- would be nice to have a method that converts sprites to AABB
		phys.RIGIDBODY_DYNAMIC,
	})

	return entity
}

func makePlayerSprite(animMgr AnimationManager) Sprite {

	return Sprite{
		[3]int{SPRITE_LOC_PLAYER_X, SPRITE_LOC_PLAYER_Y, SPRITE_LOC_PLAYER_Z},
		[2]int{SPRITE_WIDTH_PLAYER, SPRITE_HEIGHT_PLAYER},
		animMgr,
	}
}

func makePlayerAnimationManager() AnimationManager {
	idleAnim := makeAnimation(1, ANIM_OFFSET_PLAYER_VDOWN, 1)
	hAnimLeft := makeAnimation(0, ANIM_OFFSET_PLAYER_HLEFT, ANIM_FRAMES_PLAYER_H)
	hAnimRight := makeAnimation(0, ANIM_OFFSET_PLAYER_HRIGHT, ANIM_FRAMES_PLAYER_H)
	vAnimUp := makeAnimation(0, ANIM_OFFSET_PLAYER_VUP, ANIM_FRAMES_PLAYER_V)
	vAnimDown := makeAnimation(0, ANIM_OFFSET_PLAYER_VDOWN, ANIM_FRAMES_PLAYER_V)

	mgr := AnimationManager{
		[]Animation{
			idleAnim,
			hAnimLeft,
			hAnimRight,
			vAnimUp,
			vAnimDown,
		},
		ANIM_FPS_DEFAULT,
		0,   // frame
		0.0, // frame time
		0,   // ix
	}

	return mgr
}

func makeAnimation(offsetX, offsetY, frameCount int) Animation {
	return Animation{
		[2]int{offsetX, offsetY},
		frameCount,
	}
}

var _CONTROLLER_LOCK = &sync.RWMutex{}
var controllerAcceleration mgl32.Vec3

func GetControllerAccel() mgl32.Vec3 {
	_CONTROLLER_LOCK.RLock()
	defer _CONTROLLER_LOCK.RUnlock()
	return controllerAcceleration
}

func PlayerControlsCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Repeat {
		return
	}

	var multiplier float32
	if action == glfw.Release {
		// sets accel in that direction to zero
		multiplier = -1.0
	} else {
		multiplier = 1.0
	}

	_CONTROLLER_LOCK.Lock()
	defer _CONTROLLER_LOCK.Unlock()
	switch key {
	case glfw.KeyW:
		controllerAcceleration[1] += phys.ACCEL_PLAYER_DEFAULT * multiplier
	case glfw.KeyS:
		controllerAcceleration[1] += -phys.ACCEL_PLAYER_DEFAULT * multiplier
	case glfw.KeyA:
		controllerAcceleration[0] += -phys.ACCEL_PLAYER_DEFAULT * multiplier
	case glfw.KeyD:
		controllerAcceleration[0] += phys.ACCEL_PLAYER_DEFAULT * multiplier
	}
}
