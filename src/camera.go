package main

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

var _CAMERA_LOCK = &sync.Mutex{}
var cameraPtr *Entity

func getCameraPtr() *Entity {
	if cameraPtr == nil {
		_CAMERA_LOCK.Lock()
		defer _CAMERA_LOCK.Unlock()
		if cameraPtr == nil {
			entity := makeCameraEntity()
			cameraPtr = &entity
		}
	}
	return cameraPtr
}

func makeCameraEntity() Entity {
	entity := Entity{
		0,
		"Camera",
		&ComponentList{},
		mgl32.Vec3{},
		&sync.RWMutex{},
	}

	entity.components.add(&cMovable{
		mgl32.Vec3{},
		mgl32.Vec3{},
		mgl32.InfPos,
		false,
		getPlayerPtr(),
	})

	return entity
}

func getScreenCoordinates(worldCoords mgl32.Vec3) mgl32.Vec3 {
	// get position relative to camera
	// convert from meters (1 meter == 32 texels) to pixels
	centered := worldCoords.Sub(getCameraPtr().getPosition())
	return mgl32.Vec3{centered[0] * TEXEL_SCALE_X, centered[1] * TEXEL_SCALE_Y, centered[2]}
}
