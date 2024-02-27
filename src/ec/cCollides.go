package ec

import (
	// "fmt"

	// "github.com/go-gl/mathgl/mgl32"
	// "fmt"

	// "github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/phys"
)

type CCollides struct {
	collider  phys.Collider
	rigidBody phys.RigidBodyType
	// bounciness, weight
}

func (comp *CCollides) update(entity *Entity) {
	comp.collider.SetPosition(phys.Vec2Point(entity.GetPosition()))
}

func (comp *CCollides) getType() ComponentType {
	return CMP_COLLIDES
}

func TryCollideStaticDynamic(staticEntity, movableEntity *Entity) {
	collidesStatic := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, staticEntity)

	collidesMovable := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, movableEntity)
	moveComponent := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, movableEntity)

	// check next position
	initialPoint := collidesMovable.collider.GetPosition()
	nextPoint := phys.Vec2Point(moveComponent.GetNextPosition(movableEntity))
	collidesMovable.collider.SetPosition(nextPoint)

	hit := collidesMovable.collider.CheckCollision(&collidesStatic.collider)

	if !hit.IsHit {
		// reset position
		collidesMovable.collider.SetPosition(initialPoint)
		return
	}

	// stop movement in that direction and correct position
	var newPoint phys.Point
	if hit.Normal.X != 0 && hit.Normal.Y != 0 {
		moveComponent.velocity[0] = 0.0
		moveComponent.velocity[1] = 0.0
		newPoint = nextPoint.Sub(hit.Delta)
	} else if hit.Normal.X != 0 {
		moveComponent.velocity[0] = 0.0
		newPoint = phys.Point{nextPoint.X - hit.Delta.X, initialPoint.Y}
	} else {
		moveComponent.velocity[1] = 0.0
		newPoint = phys.Point{initialPoint.X, nextPoint.Y - hit.Delta.Y}
	}

	collidesMovable.collider.SetPosition(newPoint)
	movableEntity.SetPosition(mgl32.Vec3{newPoint.X, newPoint.Y, movableEntity.GetPosition()[2]})
}

func TryCollideDynamic(entity, other *Entity) {
	// var force phys.Force
	// force = &phys.ImpulseForce{
	// 	mgl32.Vec3{-1, -1, 0}.Normalize(),
	// 	1.0,
	// }
	// (*moveComponent1).velocity = force.Apply((*moveComponent1).velocity)
	// collidesEntity := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, entity)
	// movableEntity := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, entity)
	// nextPosEntity := movableEntity.GetNextPosition(entity)
	//
	// collidesOther := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, other)
	// movableOther := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, other)
	// nextPosOther := movableOther.GetNextPosition(other)
	//
	// willCollide, _ := phys.CheckCollision(nextPosEntity, nextPosOther, colliderEntity.shape, colliderOther.shape)
	// if !willCollide {
	// 	return
	// }
	// fmt.Printf("Dynamic Collision between %s and %s\n", entity.name, other.name)

	// TODO handle non-rigid body

}
