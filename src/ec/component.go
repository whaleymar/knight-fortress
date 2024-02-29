package ec

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

type Component interface {
	update(*Entity)
	getType() ComponentType
	onDelete()
}
