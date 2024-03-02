package ec

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	// "github.com/whaleymar/knight-fortress/src/math"
	"github.com/whaleymar/knight-fortress/src/phys"
	// "github.com/whaleymar/knight-fortress/src/sys"
)

var _ = fmt.Println

type CCollides struct {
	collider   phys.Collider
	RigidBody  phys.RigidBody
	IsGrounded bool // this can only be turned on by a collision and is turned off for moving entities at the end of each frame
	// bounciness, weight, smoothness
}

func (comp *CCollides) update(entity *Entity) {
	if comp.RigidBody.RBtype == phys.RIGIDBODY_DYNAMIC || comp.RigidBody.RBtype == phys.RIGIDBODY_KINEMATIC {
		comp.collider.SetPosition(phys.Vec2Point(entity.GetPosition()))
	}
	if comp.RigidBody.RBtype == phys.RIGIDBODY_DYNAMIC {
		cmpMovable := *GetComponentUnsafe[*CMovable](CMP_MOVABLE, entity)
		cmpMovable.accel = phys.PHYSICS_GRAVITY.Apply(cmpMovable.accel)
		if comp.RigidBody.State == phys.RBSTATE_GROUNDED {
			comp.RigidBody.State = phys.RBSTATE_FALLING
		}
	}
}

func (comp *CCollides) getType() ComponentType {
	return CMP_COLLIDES
}

func (comp *CCollides) onDelete() {}

func (comp *CCollides) GetSaveData() componentHolder {
	return makeComponentHolder(comp.getType(), struct {
		Collider  phys.SuperCollider
		RigidBody phys.RigidBody
	}{
		Collider:  comp.collider.MakeSuperCollider(),
		RigidBody: comp.RigidBody,
	})
}

func LoadComponentCollider(componentData interface{}) (CCollides, error) {
	// colliderdata := struct {
	// 	Collider  phys.SuperCollider
	// 	RigidBody phys.RigidBody
	// }{}

	// err := sys.LoadStruct(path, &colliderdata)
	// if err != nil {
	// 	return CCollides{}, fmt.Errorf("Couldn't load collision data from %s", path)
	// }

	type savedata struct {
		Collider  phys.SuperCollider
		RigidBody phys.RigidBody
	}

	colliderdata, ok := componentData.(savedata)
	if !ok {
		return CCollides{}, fmt.Errorf("Couldn't cast data to SuperCollider")
	}

	collider, err := phys.ExtractCollider(colliderdata.Collider)
	if err != nil {
		return CCollides{}, fmt.Errorf("Error loading collider struct from supercollider due to %s", err)
	}

	return CCollides{
		collider,
		colliderdata.RigidBody,
		true,
	}, nil
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
	// TODO get rid of first condition + elses
	if hit.Normal.X != 0 && hit.Normal.Y != 0 {
		fmt.Println("corner collision")
		// zero out the slower axis
		if moveComponent.velocity[0] >= moveComponent.velocity[1] {
			moveComponent.velocity[1] = 0.0
			newPoint = phys.Point{initialPoint.X, nextPoint.Y - hit.Delta.Y}
		} else {
			moveComponent.velocity[0] = 0.0
			newPoint = phys.Point{nextPoint.X - hit.Delta.X, initialPoint.Y}
		}

		// original code:
		// moveComponent.velocity[0] = 0.0
		// moveComponent.velocity[1] = 0.0
		// newPoint = nextPoint.Sub(hit.Delta)
		if hit.Normal.Y == -1 {
			collidesMovable.IsGrounded = true
			collidesMovable.RigidBody.State = phys.RBSTATE_GROUNDED
		}
	} else if hit.Normal.X != 0 {
		fmt.Println("hit x", hit)
		moveComponent.velocity[0] = 0.0
		newPoint = phys.Point{nextPoint.X - hit.Delta.X, initialPoint.Y}
		// if math.Abs(hit.Delta.X) > phys.MIN_DELTA_HALT {
		// 	fmt.Println("hit x", hit)
		// 	moveComponent.velocity[0] = 0.0
		// 	newPoint = phys.Point{nextPoint.X - hit.Delta.X, initialPoint.Y}
		// } else {
		// 	return
		// }
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
