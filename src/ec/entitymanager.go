package ec

import (
	"fmt"
	"sync"
)

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

func (eMgr *EntityManager) Add(entity *Entity) (uint64, error) {
	// enforce uniqueness?
	eMgr.rwlock.Lock()
	defer eMgr.rwlock.Unlock()

	err := entity.Init()
	if err != nil {
		return 0, fmt.Errorf("Entity %s has error: %s\n", entity.String(), err)
	}

	uid := eMgr.nextId
	entity.Uid = uid
	eMgr.nextId++
	eMgr.entities = append(eMgr.entities, entity)

	return uid, nil
}

func (eMgr *EntityManager) Remove(uid uint64) {
	eMgr.rwlock.Lock()
	defer eMgr.rwlock.Unlock()

	for i, entity := range eMgr.entities {
		if uid == entity.Uid {
			entity.Components.Clear()
			eMgr.entities = append(eMgr.entities[:i], eMgr.entities[i+1:]...)
			return
		}
	}
}

func (eMgr *EntityManager) Get(uid uint64) (*Entity, error) {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()

	for _, entity := range eMgr.entities {
		if uid == entity.Uid {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("Entity with ID %d not found", uid)
}

func (eMgr *EntityManager) GetEntityWithName(name string) (*Entity, error) {
	eMgr.rwlock.RLock()
	defer eMgr.rwlock.RUnlock()

	for _, entity := range eMgr.entities {
		if name == entity.Name {
			return entity, nil
		}
	}
	return nil, fmt.Errorf("Entity with name %s not found", name)
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
	return entities
}
