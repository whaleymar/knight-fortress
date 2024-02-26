package ec

import (
	// "fmt"

	// "github.com/go-gl/mathgl/mgl32"
	"fmt"

	"github.com/whaleymar/knight-fortress/src/phys"
)

// TODO should just store a position here because it won't always match the entity position origin?
type CCollides struct {
	shape       phys.AABB
	isRigidBody bool
	// bounciness, weight
}

func (comp *CCollides) update(entity *Entity) {}

func (comp *CCollides) getType() ComponentType {
	return CMP_COLLIDES
}

func TryCollideDynamic(entity, other *Entity) {
	// var force phys.Force
	// force = &phys.ImpulseForce{
	// 	mgl32.Vec3{-1, -1, 0}.Normalize(),
	// 	1.0,
	// }
	// (*moveComponent1).velocity = force.Apply((*moveComponent1).velocity)
	colliderEntity := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, entity)
	movableEntity := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, entity)
	nextPosEntity := movableEntity.GetNextPosition(entity)

	colliderOther := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, other)
	movableOther := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, other)
	nextPosOther := movableOther.GetNextPosition(other)

	willCollide := phys.CheckCollision(nextPosEntity, nextPosOther, colliderEntity.shape, colliderOther.shape)
	if !willCollide {
		return
	}
	fmt.Printf("Dynamic Collision between %s and %s\n", entity.name, other.name)

	// TODO handle non-rigid body

}

func TryCollideStaticDynamic(staticEntity, movableEntity *Entity) {
	colliderStatic := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, staticEntity)

	colliderMovable := *GetComponentUnsafe[*CCollides](CMP_COLLIDES, movableEntity)
	moveComponent := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, movableEntity)
	nextPos := moveComponent.GetNextPosition(movableEntity)

	willCollide := phys.CheckCollision(staticEntity.GetPosition(), nextPos, colliderStatic.shape, colliderMovable.shape)
	if !willCollide {
		return
	}
	fmt.Printf("Collision between Static %s and Dynamic %s\n", staticEntity.name, movableEntity.name)

	// determine direction of collision based on dot product of velocity and direction
	// direction := movableEntity.GetPosition().Sub(staticEntity.GetPosition()).Normalize()
	// let's get something basic working
	// acceleration still gets added to speed during movement.Update() so player can force their way into things
	moveComponent.velocity[0] = 0.0
	moveComponent.velocity[1] = 0.0
	moveComponent.accel[0] = 0.0
	moveComponent.accel[1] = 0.0

}
