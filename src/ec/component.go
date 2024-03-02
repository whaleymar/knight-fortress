package ec

import "fmt"

type ComponentType int

const (
	CMP_ANY ComponentType = iota
	CMP_DRAWABLE
	CMP_COLLIDES
	CMP_MOVABLE
	CMP_SERIALIZE
)

type ComponentTypeSet interface {
	Component
}

type Component interface {
	update(*Entity)
	getType() ComponentType
	onDelete()
	GetSaveData() componentHolder
}

type PODComponentTypeSet interface {
	PODComponent
}

// Plain Old Data
type PODComponent interface {
	getType() ComponentType
}

// for saving/loading
type componentHolder struct {
	CType    ComponentType
	CompData string
}

func makeComponentHolder(cType ComponentType, comp string) componentHolder {
	return componentHolder{
		CType:    cType,
		CompData: comp,
	}
}

func loadComponent(compHolder componentHolder) (*Component, error) {
	var comp Component
	var err error
	switch compHolder.CType {
	case CMP_MOVABLE:
		component, tmpErr := LoadComponentMovable(compHolder.CompData)
		comp = &component
		err = tmpErr
		break
	case CMP_COLLIDES:
		component, tmpErr := LoadComponentCollider(compHolder.CompData)
		comp = &component
		err = tmpErr
		break
	case CMP_DRAWABLE:
		component, tmpErr := LoadComponentDrawable(compHolder.CompData)
		comp = &component
		err = tmpErr
		break
	default:
		return nil, fmt.Errorf("Unrecognized component type: %d", compHolder.CType)
	}
	if err != nil {
		return nil, err
	}

	return &comp, nil
}
