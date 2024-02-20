package main

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	FILE_SPRITE_PLAYER   = "assets/player.png"
	SPRITE_HEIGHT_PLAYER = 16
	SPRITE_WIDTH_PLAYER  = 16
	ANIM_HOFFSET_PLAYER  = 1
	ANIM_VOFFSET_PLAYER  = 0
	ANIM_HFRAMES_PLAYER  = 2
	ANIM_VFRAMES_PLAYER  = 4
)

type PlayerEntity struct {
	components ComponentManager
	position   mgl32.Vec3
}

// func (entity *PlayerEntity) init() {
//
// }

func (entity *PlayerEntity) getPosition() mgl32.Vec3 {
	return entity.position
}

func (entity *PlayerEntity) setPosition(position mgl32.Vec3) {
	entity.position = position
}

func (entity *PlayerEntity) getManager() ComponentManager {
	return entity.components
}

var lock = &sync.Mutex{} // this is package-wide, maybe move somewhere else (god i miss namespaces) TODO
var playerPtr *PlayerEntity

func getPlayerPtr() *PlayerEntity {
	if playerPtr == nil {
		lock.Lock()
		defer lock.Unlock()
		if playerPtr == nil {
			entity := makePlayerEntity()
			playerPtr = &entity
		}
	}
	return playerPtr
}

func makePlayerEntity() PlayerEntity {
	curVertices := squareVertices
	vao := makeVao(curVertices)
	entity := PlayerEntity{
		&ComponentList{},
		mgl32.Vec3{},
	}

	entity.components.add(&cDrawable{
		CMP_DRAWABLE,
		vao,
		makePlayerSprite(),
		makePlayerAnimationManager(),
	})

	entity.components.add(&cMovable{
		mgl32.Vec3{},
		mgl32.Vec3{},
		0.25,
	})

	return entity
}

func makePlayerSprite() Sprite {
	img, err := loadImage(FILE_SPRITE_PLAYER)
	if err != nil {
		panic(err)
	}

	return Sprite{
		img,
		SPRITE_HEIGHT_PLAYER,
		SPRITE_WIDTH_PLAYER,
	}
}

func makePlayerAnimationManager() AnimationManager {
	// TODO magic numbers
	idleAnim := Animation{ANIM_HOFFSET_PLAYER, 1}
	hAnim := Animation{ANIM_HOFFSET_PLAYER, ANIM_HFRAMES_PLAYER}
	vAnim := Animation{ANIM_VOFFSET_PLAYER, ANIM_VFRAMES_PLAYER}

	mgr := AnimationManager{
		nil,
		4.0, // anim speed (fps)
		0,   // frame
		0.0, // frmae time
		0,   // ix
	}

	mgr.anims = append(mgr.anims, idleAnim)
	mgr.anims = append(mgr.anims, hAnim)
	mgr.anims = append(mgr.anims, vAnim)

	return mgr
}
