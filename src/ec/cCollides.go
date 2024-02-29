package ec

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/phys"
)

var _ = fmt.Println

type CCollides struct {
	collider   phys.Collider
	RigidBody  phys.RigidBody
	IsGrounded bool
	// bounciness, weight, smoothness
}

func (comp *CCollides) update(entity *Entity) {
	if comp.RigidBody.RBtype == phys.RIGIDBODY_DYNAMIC || comp.RigidBody.RBtype == phys.RIGIDBODY_KINEMATIC {
		comp.collider.SetPosition(phys.Vec2Point(entity.GetPosition()))
	}
	if comp.RigidBody.RBtype == phys.RIGIDBODY_DYNAMIC {
		if !comp.IsGrounded {
			cmpMovable := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, entity)
			cmpMovable.accel = phys.PHYSICS_GRAVITY.Apply(cmpMovable.accel)
			if comp.RigidBody.State == phys.RBSTATE_GROUNDED {
				comp.RigidBody.State = phys.RBSTATE_FALLING
			}
		}
	}
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
		collidesMovable.collider.SetPosition(initialPoint)
		return
	}

	// stop movement in that direction and correct position
	var newPoint phys.Point
	if hit.Normal.X != 0 && hit.Normal.Y != 0 {
		moveComponent.velocity[0] = 0.0
		moveComponent.velocity[1] = 0.0
		newPoint = nextPoint.Sub(hit.Delta)
		if hit.Normal.Y == -1 {
			collidesMovable.IsGrounded = true
			collidesMovable.RigidBody.State = phys.RBSTATE_GROUNDED
		}
	} else if hit.Normal.X != 0 {
		moveComponent.velocity[0] = 0.0
		newPoint = phys.Point{nextPoint.X - hit.Delta.X, initialPoint.Y}
	} else {
		moveComponent.velocity[1] = 0.0
		newPoint = phys.Point{initialPoint.X, nextPoint.Y - hit.Delta.Y}
		if hit.Normal.Y == -1 {
			collidesMovable.IsGrounded = true
			collidesMovable.RigidBody.State = phys.RBSTATE_GROUNDED
		}
	}

	collidesMovable.collider.SetPosition(newPoint)
	movableEntity.SetPosition(mgl32.Vec3{newPoint.X, newPoint.Y, movableEntity.GetPosition()[2]})
	// moveComponent.updateAnimation(movableEntity)
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
