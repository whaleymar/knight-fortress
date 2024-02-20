package main

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type ComponentType int

const (
	CMP_DRAWABLE ComponentType = iota
	CMP_MOVABLE
)

type ComponentTypeList interface {
	Component
}

func getComponent[T ComponentTypeList](enum ComponentType, entity Entity) (*T, error) {
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

type Entity interface {
	// init() // TODO i have to figure out how generics work to return a T here
	getPosition() mgl32.Vec3
	setPosition(mgl32.Vec3)
	getManager() ComponentManager
}

// idk if i need multiple of these... maybe an event manager
type ComponentManager interface {
	add(Component)
	get(ComponentType) (*Component, error)
	remove(ComponentType) error
	update(Entity)
}

type ComponentList struct {
	components []Component
}

type Component interface {
	update(Entity)
	getType() ComponentType
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

func (components *ComponentList) update(entity Entity) {
	for _, comp := range components.components {
		comp.update(entity)
	}
}
