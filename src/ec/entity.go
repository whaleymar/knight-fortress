package ec

import (
	"fmt"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/phys"
)

type Entity struct {
	Uid        uint64
	Name       string
	Components ComponentManager
	Position   mgl32.Vec3
	RWlock     *sync.RWMutex
}

func (entity *Entity) GetPosition() mgl32.Vec3 {
	entity.RWlock.RLock()
	defer entity.RWlock.RUnlock()
	return entity.Position
}

func (entity *Entity) GetBottomLeftPosition() mgl32.Vec3 {
	// if entity is drawable, return position of the bottom left point of its vertex array
	position := entity.GetPosition()
	tmp, err := GetComponent[*CDrawable](CMP_DRAWABLE, entity)
	if err != nil {
		return position
	}
	sizeX, sizeY := (*tmp).getFrameSize() // in texels
	return mgl32.Vec3{
		position[0] - sizeX/gfx.TEXELS_PER_METER/2.0*(*tmp).scale[0],
		position[1] - sizeY/gfx.TEXELS_PER_METER/2.0*(*tmp).scale[1],
		0.0,
	}
}

func (entity *Entity) SetPosition(position mgl32.Vec3) {
	entity.RWlock.Lock()
	defer entity.RWlock.Unlock()
	entity.Position = position
}

func (entity *Entity) GetComponentManager() ComponentManager {
	return entity.Components
}

func (entity *Entity) String() string {
	return string(fmt.Sprint(entity.Name))
}

func (entity *Entity) Equals(other *Entity) bool {
	return entity.Uid == other.Uid
}

func (entity *Entity) GetId() uint64 {
	return entity.Uid
}

func (entity *Entity) Init() error {
	// validate component dependencies
	// should run after all components are added
	cmpCollides, err := GetComponent[*CCollides](CMP_COLLIDES, entity)
	if err == nil {
		_, err := GetComponent[*CMovable](CMP_MOVABLE, entity)
		if err == nil {
			if (*cmpCollides).RigidBody.RBtype == phys.RIGIDBODY_NONE || (*cmpCollides).RigidBody.RBtype == phys.RIGIDBODY_STATIC {
				return fmt.Errorf("Entity with RigidBody enum %d cannot have movable component", (*cmpCollides).RigidBody)
			}
		} else {
			if (*cmpCollides).RigidBody.RBtype == phys.RIGIDBODY_DYNAMIC || (*cmpCollides).RigidBody.RBtype == phys.RIGIDBODY_KINEMATIC {
				return fmt.Errorf("Entity with RigidBody enum %d must have movable component", (*cmpCollides).RigidBody)
			}
		}
		(*cmpCollides).collider.SetPosition(phys.Vec2Point(entity.GetPosition()))
	}
	return nil
}
