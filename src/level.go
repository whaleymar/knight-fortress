package main

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

// TODO
// eventually level data will be in a config/level file and
// each file will be loadable with methods defined here
const (
	SPRITE_HEIGHT_LEVEL = 256
	SPRITE_WIDTH_LEVEL  = 320
)

// I don't like components managing other entities because then I have to recursively search
// all of those entities for their shit
// and almost every entity will be attached to the current level so it makes more sense
// for the main entity manager to do all the work
// and store level data in its own struct without an entity manager
// type cLevelData struct {
// 	playerSpawn mgl32.Vec3
// 	entityMgr   EntityManager
// }
//
// func (comp *cLevelData) update(entity *Entity) {
// 	for _, entity := range comp.entityMgr.getEntitiesWithComponent(CMP_ANY) {
// 		entity.components.update(entity)
// 	}
// }
//
// func (comp *cLevelData) getType() ComponentType {
// 	return CMP_LEVELDATA
// }

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
