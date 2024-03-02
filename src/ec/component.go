package ec

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
	GetSaveData() interface{}
}

type PODComponentTypeSet interface {
	PODComponent
}

// Plain Old Data
type PODComponent interface {
	getType() ComponentType
}
