package main

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	SPRITE_HEIGHT_PLAYER      = 16
	SPRITE_WIDTH_PLAYER       = 16
	ANIM_OFFSET_PLAYER_HRIGHT = 1
	ANIM_OFFSET_PLAYER_HLEFT  = 3
	ANIM_OFFSET_PLAYER_VDOWN  = 0
	ANIM_OFFSET_PLAYER_VUP    = 2
	ANIM_FRAMES_PLAYER_H      = 2
	ANIM_FRAMES_PLAYER_V      = 4
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
	}

	entity.components.add(&cDrawable{
		CMP_DRAWABLE,
		squareVertices,
		// scaleDepth(squareVertices, 0.1),
		makeVao(),
		makeVbo(),
		makePlayerSprite(makePlayerAnimationManager()),
		TEX_MAIN,
	})

	entity.components.add(&cMovable{
		mgl32.Vec3{},
		mgl32.Vec3{},
		0.25, // speedMax
	})

	return entity
}

func makePlayerSprite(animMgr AnimationManager) Sprite {

	return Sprite{
		mgl32.Vec3{0, 0, 0}, // TODO hard coded
		[2]int{SPRITE_WIDTH_PLAYER, SPRITE_HEIGHT_PLAYER},
		animMgr,
	}
}

func makePlayerAnimationManager() AnimationManager {
	// TODO magic numbers
	idleAnim := makeAnimation(1, ANIM_OFFSET_PLAYER_VDOWN, 1)
	hAnimLeft := makeAnimation(0, ANIM_OFFSET_PLAYER_HLEFT, ANIM_FRAMES_PLAYER_H)
	hAnimRight := makeAnimation(0, ANIM_OFFSET_PLAYER_HRIGHT, ANIM_FRAMES_PLAYER_H)
	vAnimUp := makeAnimation(0, ANIM_OFFSET_PLAYER_VUP, ANIM_FRAMES_PLAYER_V)
	vAnimDown := makeAnimation(0, ANIM_OFFSET_PLAYER_VDOWN, ANIM_FRAMES_PLAYER_V)

	mgr := AnimationManager{
		nil,
		4.0, // anim speed (fps)
		0,   // frame
		0.0, // frame time
		0,   // ix
	}

	mgr.anims = append(mgr.anims, idleAnim)
	mgr.anims = append(mgr.anims, hAnimLeft)
	mgr.anims = append(mgr.anims, hAnimRight)
	mgr.anims = append(mgr.anims, vAnimUp)
	mgr.anims = append(mgr.anims, vAnimDown)

	return mgr
}

func makeAnimation(offsetX, offsetY, frameCount int) Animation {
	return Animation{
		[2]int{offsetX, offsetY},
		frameCount,
	}
}
