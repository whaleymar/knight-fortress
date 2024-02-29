package ec

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
)

var _CAMERA_LOCK = &sync.Mutex{}
var cameraPtr *Entity

func GetCameraPtr() *Entity {
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

	entity.Components.Add(&CMovable{
		mgl32.Vec3{},
		mgl32.Vec3{},
		mgl32.InfPos,
		false,
		GetPlayerPtr(),
		false,
	})

	return entity
}

func GetScreenCoordinates(worldCoords mgl32.Vec3) mgl32.Vec3 {
	// get position relative to camera
	// convert from meters (1 meter == 32 texels) to pixels
	centered := worldCoords.Sub(GetCameraPtr().GetPosition())
	return mgl32.Vec3{centered[0] * gfx.TEXEL_SCALE_X, centered[1] * gfx.TEXEL_SCALE_Y, centered[2]}
}
