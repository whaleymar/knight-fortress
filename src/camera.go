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
	// idk how to do this in the shader or if i even should
	return worldCoords.Sub(getCameraPtr().getPosition())
}
