package ec

import (
	"fmt"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type ComponentType int

const (
	CMP_ANY ComponentType = iota
	CMP_DRAWABLE
	CMP_COLLIDES
	CMP_MOVABLE
)

type ComponentTypeList interface {
	Component
}

type Entity struct {
	uid        uint64
	name       string
	components ComponentManager
	position   mgl32.Vec3
	rwlock     *sync.RWMutex
}

func (entity *Entity) GetPosition() mgl32.Vec3 {
	entity.rwlock.RLock()
	defer entity.rwlock.RUnlock()
	return entity.position
}

func (entity *Entity) SetPosition(position mgl32.Vec3) {
	entity.rwlock.Lock()
	defer entity.rwlock.Unlock()
	entity.position = position
}

func (entity *Entity) GetComponentManager() ComponentManager {
	return entity.components
}

func (entity *Entity) String() string {
	return string(fmt.Sprint(entity.name))
}

func (entity *Entity) Equals(other *Entity) bool {
	return entity.uid == other.uid
}

func (entity *Entity) GetId() uint64 {
	return entity.uid
}

type EntityManager struct {
	entities []*Entity
	rwlock   sync.RWMutex
	nextId   uint64
}

var _ENTITYMGR_LOCK = &sync.Mutex{}
var entityManagerPtr *EntityManager

func GetEntityManager() *EntityManager {
	if entityManagerPtr == nil {
		_ENTITYMGR_LOCK.Lock()
		defer _ENTITYMGR_LOCK.Unlock()
		if entityManagerPtr == nil {
			eMgr := EntityManager{}
			entityManagerPtr = &eMgr
		}
	}
	return entityManagerPtr
}

func (eMgr *EntityManager) Add(entity *Entity) {
	// enforce uniqueness?
	eMgr.rwlock.Lock()
	defer eMgr.rwlock.Unlock()

	entity.uid = eMgr.nextId
	eMgr.nextId++
	eMgr.entities = append(eMgr.entities, entity)
}

func (eMgr *EntityManager) Remove(uid uint64) {
	eMgr.rwlock.Lock()
	defer eMgr.rwlock.Unlock()

	for i, entity := range eMgr.entities {
		if uid == entity.uid {
			eMgr.entities = append(eMgr.entities[:i], eMgr.entities[i+1:]...)
			return
		}
	}
}

func (eMgr *EntityManager) Get(uid uint64) (*Entity, error) {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()

	for _, entity := range eMgr.entities {
		if uid == entity.uid {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("Entity with ID %d not found", uid)
}

func (eMgr *EntityManager) Len() int {
	return len(eMgr.entities)
}

func (eMgr *EntityManager) GetEntitiesWithComponent(enum ComponentType) []*Entity {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()

	if enum == CMP_ANY {
		return eMgr.entities
	}
	entities := make([]*Entity, 0, len(eMgr.entities))
	for _, entity := range eMgr.entities {
		_, err := entity.GetComponentManager().Get(enum)
		if err != nil {
			continue
		}
		entities = append(entities, entity)
	}
	return entities
}

func (eMgr *EntityManager) GetEntitiesWithManyComponents(enums ...ComponentType) []*Entity {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()
	// for _, entity := range eMgr.entities {
	// 	fmt.Println(entity.name)
	// }
	// fmt.Println("\n")
	entities := make([]*Entity, 0, len(eMgr.entities))

	for _, entity := range eMgr.entities {
		isValid := true
		for _, enum := range enums {
			_, err := entity.GetComponentManager().Get(enum)
			if err != nil {
				isValid = false
				break
			}
		}
		if isValid {
			entities = append(entities, entity)
		}
	}
	// for _, entity := range entities {
	// 	fmt.Println(entity.name)
	// }
	return entities
}

// idk if i need multiple of these... maybe an event manager
type ComponentManager interface {
	Add(Component)
	Get(ComponentType) (*Component, error)
	Remove(ComponentType) error
	Update(*Entity)
}

type ComponentList struct {
	components []Component // TODO this should rarely update and mostly be searched, so should be sorted + have binary search
}

type Component interface {
	update(*Entity)
	getType() ComponentType
}

func GetComponent[T ComponentTypeList](enum ComponentType, entity *Entity) (*T, error) {
	compMgr := entity.GetComponentManager()
	compInterface, err := compMgr.Get(enum)
	if err != nil {
		return nil, fmt.Errorf("No %d component found", enum)
	}

	comp, ok := (*compInterface).(T)
	if !ok {
		return nil, fmt.Errorf("No %d component found", enum)
	}
	return &comp, nil
}

func GetComponentUnsafe[T ComponentTypeList](enum ComponentType, entity *Entity) *T {
	// use this for components fetched with getEntitiesWithComponent()
	compMgr := entity.GetComponentManager()
	compInterface, err := compMgr.Get(enum)
	_ = err
	comp, ok := (*compInterface).(T)
	_ = ok
	return &comp
}

func (components *ComponentList) Add(comp Component) {
	components.components = append(components.components, comp)
}

func (components *ComponentList) Get(enum ComponentType) (*Component, error) {
	for i, comp := range components.components {
		if comp.getType() == enum {
			return &components.components[i], nil
		}
	}
	return nil, fmt.Errorf("Component not found: %d", enum)
}

func (components *ComponentList) Remove(enum ComponentType) error {
	for i, comp := range components.components {
		if comp.getType() == enum {
			components.components = append(components.components[:i], components.components[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Component not found: %d", enum)
}

func (components *ComponentList) Update(entity *Entity) {
	for _, comp := range components.components {
		comp.update(entity)
	}
}
