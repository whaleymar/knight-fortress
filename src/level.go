package main

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	SPRITE_HEIGHT_LEVEL = 256
	SPRITE_WIDTH_LEVEL  = 320
)

// eventually this will take a level config file which contains all the component data
func makeLevelEntity() Entity {
	entity := Entity{
		0,
		"Level",
		&ComponentList{},
		mgl32.Vec3{0.0, 0.0, DEPTH_BACKGROUND},
		&sync.RWMutex{},
	}

	entity.components.add(&cDrawable{
		// levelVertices,
		screenVertices,
		makeVao(),
		makeVbo(),
		makeLevelSprite(makeStaticAnimationManager()),
		TEX_MAIN,
		&sync.RWMutex{},
		true, // isUvUpdateNeeded
	})

	return entity

}

func makeLevelSprite(animMgr AnimationManager) Sprite {
	return Sprite{
		[3]int{0, 64, 0},
		[2]int{SPRITE_WIDTH_LEVEL, SPRITE_HEIGHT_LEVEL},
		animMgr,
	}

}
