package main

import (
	"fmt"
	"sync"

	"github.com/go-gl/mathgl/mgl32"
)

type ComponentType int

const (
	CMP_ANY ComponentType = iota
	CMP_DRAWABLE
	CMP_MOVABLE
)

type ComponentTypeList interface {
	Component
}

type Entity struct {
	uid        uint64
	components ComponentManager
	position   mgl32.Vec3
	rwlock     *sync.RWMutex
}

func (entity *Entity) getPosition() mgl32.Vec3 {
	entity.rwlock.RLock()
	defer entity.rwlock.RUnlock()
	return entity.position
}

func (entity *Entity) setPosition(position mgl32.Vec3) {
	entity.rwlock.Lock()
	defer entity.rwlock.Unlock()
	entity.position = position
}

func (entity *Entity) getManager() ComponentManager {
	return entity.components
}

func (entity *Entity) String() string {
	return string(fmt.Sprint(entity.uid))
}

type EntityManager struct {
	entities []*Entity
	rwlock   sync.RWMutex
	nextId   uint64
}

var _ENTITYMGR_LOCK = &sync.Mutex{}
var entityManagerPtr *EntityManager

func getEntityManager() *EntityManager {
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

func (eMgr *EntityManager) add(entity Entity) {
	// enforce uniqueness?
	eMgr.rwlock.Lock()
	defer eMgr.rwlock.Unlock()

	entity.uid = eMgr.nextId
	eMgr.nextId++
	eMgr.entities = append(eMgr.entities, &entity)
}

func (eMgr *EntityManager) remove(uid uint64) {
	eMgr.rwlock.Lock()
	defer eMgr.rwlock.Unlock()

	for i, entity := range eMgr.entities {
		if uid == entity.uid {
			eMgr.entities = append(eMgr.entities[:i], eMgr.entities[i+1:]...)
			return
		}
	}
}

func (eMgr *EntityManager) get(uid uint64) (*Entity, error) {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()

	for _, entity := range eMgr.entities {
		if uid == entity.uid {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("Entity with ID %d not found", uid)
}

func (eMgr *EntityManager) getEntitiesWithComponent(enum ComponentType) []*Entity {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()

	if enum == CMP_ANY {
		return eMgr.entities
	}
	entities := make([]*Entity, 0, len(eMgr.entities))
	for _, entity := range eMgr.entities {
		_, err := entity.getManager().get(enum)
		if err != nil {
			continue
		}
		entities = append(entities, entity)
	}
	return entities
}

// idk if i need multiple of these... maybe an event manager
type ComponentManager interface {
	add(Component)
	get(ComponentType) (*Component, error)
	remove(ComponentType) error
	update(*Entity)
}

type ComponentList struct {
	components []Component
}

type Component interface {
	update(*Entity)
	getType() ComponentType
}

func getComponent[T ComponentTypeList](enum ComponentType, entity *Entity) (*T, error) {
	compMgr := entity.getManager()
	compInterface, err := compMgr.get(enum)
	if err != nil {
		return nil, fmt.Errorf("No %d component found", enum)
	}

	comp, ok := (*compInterface).(T)
	if !ok {
		return nil, fmt.Errorf("No %d component found", enum)
	}
	return &comp, nil
}

func getComponentUnsafe[T ComponentTypeList](enum ComponentType, entity *Entity) *T {
	// use this for components fetched with getEntitiesWithComponent()
	compMgr := entity.getManager()
	compInterface, err := compMgr.get(enum)
	_ = err
	comp, ok := (*compInterface).(T)
	_ = ok
	return &comp
}

func (components *ComponentList) add(comp Component) {
	components.components = append(components.components, comp)
}

func (components *ComponentList) get(enum ComponentType) (*Component, error) {
	for i, comp := range components.components {
		if comp.getType() == enum {
			return &components.components[i], nil
		}
	}
	return nil, fmt.Errorf("Component not found: %d", enum)
}

func (components *ComponentList) remove(enum ComponentType) error {
	for i, comp := range components.components {
		if comp.getType() == enum {
			components.components = append(components.components[:i], components.components[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Component not found: %d", enum)
}

func (components *ComponentList) update(entity *Entity) {
	for _, comp := range components.components {
		comp.update(entity)
	}
}
