package ec

import (
	"fmt"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/phys"
	"github.com/whaleymar/knight-fortress/src/sys"
)

const (
	ASSET_DIR_ENTITY = "assets/entity"
)

type Entity struct {
	Uid        uint64 // unique at runtime
	Name       string
	Components *ComponentManager
	Data       []PODComponent
	Position   mgl32.Vec3
	*sync.RWMutex
}

func MakeEntity(name string) Entity {
	return Entity{
		0,
		name,
		&ComponentManager{},
		[]PODComponent{},
		mgl32.Vec3{},
		&sync.RWMutex{},
	}
}

func (entity *Entity) SaveToFile() error {
	saveComp, err := GetPODComponent[*CSerialize](CMP_SERIALIZE, entity)
	if err != nil {
		return fmt.Errorf("Can't save %s entity because it does not have a serialize component", entity.Name)
	}
	var componentlist []componentHolder
	for _, component := range entity.Components.components {
		componentlist = append(componentlist, component.GetSaveData())
	}

	return sys.SaveStruct(getEntityPath((*saveComp).FileName), componentlist)
}

func LoadEntity(filename string) (Entity, error) {
	// TODO change to struct
	var data []componentHolder
	err := sys.LoadStruct(getEntityPath(filename), &data)
	if err != nil {
		return Entity{}, err
	}
	entity := MakeEntity(filename)
	for _, componentHolder := range data {
		component, err := loadComponent(componentHolder)
		if err != nil {
			fmt.Println(err)
			continue
		}
		entity.GetComponentManager().Add(*component)
	}
	return entity, nil
}

func (entity *Entity) GetPosition() mgl32.Vec3 {
	entity.RLock()
	defer entity.RUnlock()
	return entity.Position
}

func (entity *Entity) GetDrawPosition() mgl32.Vec3 {
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
	entity.Lock()
	defer entity.Unlock()
	entity.Position = position
}

func (entity *Entity) GetComponentManager() *ComponentManager {
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

func (entity *Entity) Copy() (Entity, error) {
	newComponentList, err := entity.Components.Copy()
	if err != nil {
		return Entity{}, err
	}
	return Entity{
		entity.Uid,
		entity.Name,
		&ComponentManager{newComponentList},
		entity.Data,
		entity.Position,
		&sync.RWMutex{},
	}, nil

}

func getEntityPath(filename string) string {
	return fmt.Sprintf("%s/%s.yml", ASSET_DIR_ENTITY, filename)
}
