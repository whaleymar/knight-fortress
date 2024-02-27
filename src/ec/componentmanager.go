package ec

import "fmt"

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
	compInterface, _ := compMgr.Get(enum)
	comp := (*compInterface).(T)
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
