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

// TODO thread safe entities
type Entity struct {
	uid        uint64
	components ComponentManager
	position   mgl32.Vec3
}

func (entity *Entity) getPosition() mgl32.Vec3 {
	return entity.position
}

func (entity *Entity) setPosition(position mgl32.Vec3) {
	entity.position = position
}

func (entity *Entity) getManager() ComponentManager {
	return entity.components
}

// TODO double check that these *actually* need to be pointers
type EntityManager struct {
	entities []*Entity
	rwlock   sync.RWMutex
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

func (components *ComponentList) add(comp Component) {
	// TODO should sort array by enum for faster lookup and removal
	// should I search for a matching component type and replace? slow but safer
	components.components = append(components.components, comp)
}

func (components *ComponentList) get(enum ComponentType) (*Component, error) {
	for i, comp := range components.components {
		if comp.getType() == enum {
			return &components.components[i], nil
		}
	}
	return nil, fmt.Errorf("Component not found: %d", enum) // TODO after string code gen change this
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
