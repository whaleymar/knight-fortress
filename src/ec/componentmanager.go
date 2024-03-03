package ec

import "fmt"

type ComponentManager struct {
	components []Component // TODO this should rarely update and mostly be searched, so should be sorted + have binary search
}

func GetComponent[T ComponentTypeSet](enum ComponentType, entity *Entity) (*T, error) {
	compInterface, err := entity.GetComponentManager().Get(enum)
	if err != nil {
		return nil, fmt.Errorf("No %d component found", enum)
	}

	comp, ok := (*compInterface).(T)
	if !ok {
		return nil, fmt.Errorf("No %d component found", enum)
	}
	return &comp, nil
}

func GetComponentUnsafe[T ComponentTypeSet](enum ComponentType, entity *Entity) *T {
	// use this for components fetched with getEntitiesWithComponent()
	compInterface, _ := entity.GetComponentManager().Get(enum)
	comp := (*compInterface).(T)
	return &comp
}

func GetPODComponent[T PODComponentTypeSet](enum ComponentType, entity *Entity) (*T, error) {
	for _, comp := range entity.Data {
		if comp.getType() == enum {
			typedComponent := (comp).(T)
			return &typedComponent, nil
		}
	}
	return nil, fmt.Errorf("No %d component found", enum)
}

func (components *ComponentManager) Add(comp Component) {
	components.components = append(components.components, comp)
}

func (components *ComponentManager) Get(enum ComponentType) (*Component, error) {
	for i, comp := range components.components {
		if comp.getType() == enum {
			return &components.components[i], nil
		}
	}
	return nil, fmt.Errorf("Component not found: %d", enum)
}

func (components *ComponentManager) Remove(enum ComponentType) error {
	for i, comp := range components.components {
		if comp.getType() == enum {
			comp.onDelete()
			components.components = append(components.components[:i], components.components[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Component not found: %d", enum)
}

func (components *ComponentManager) Update(entity *Entity) {
	for _, comp := range components.components {
		comp.update(entity)
	}
}

func (components *ComponentManager) Clear() {
	for _, comp := range components.components {
		comp.onDelete()
	}
	components.components = nil
}

func (components *ComponentManager) Copy() ([]Component, error) {
	newComponentList := []Component{}
	for _, comp := range components.components {
		newComponent, err := comp.Copy()
		if err != nil {
			return nil, err
		}
		newComponentList = append(newComponentList, newComponent)
	}
	return newComponentList, nil

}
