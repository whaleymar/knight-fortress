package main

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
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
	PHYSICS_PLAYER_SPEEDMAX   = 0.5
)

var _PLAYER_LOCK = &sync.Mutex{}
var playerPtr *Entity

func getPlayerPtr() *Entity {
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
		&ComponentList{},
		mgl32.Vec3{0.0, 0.0, DEPTH_PLAYER},
		&sync.RWMutex{},
	}

	entity.components.add(&cDrawable{
		CMP_DRAWABLE,
		squareVertices,
		makeVao(),
		makeVbo(),
		makePlayerSprite(makePlayerAnimationManager()),
		TEX_MAIN,
		&sync.RWMutex{},
		true, // isUvUpdateNeeded
	})

	entity.components.add(&cMovable{
		mgl32.Vec3{},
		mgl32.Vec3{},
		PHYSICS_PLAYER_SPEEDMAX,
		true, // frictionActive
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
