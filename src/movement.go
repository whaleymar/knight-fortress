package main

import (
	// "fmt"

	// "fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// movement based animations here: fall, jump, crouch
type AnimationIndex int

const (
	ANIM_IDLE AnimationIndex = iota
	ANIM_MOVE_HLEFT
	ANIM_MOVE_HRIGHT
	ANIM_MOVE_VUP
	ANIM_MOVE_VDOWN
)

const (
	PHYSICS_FRICTION_COEF = float32(0.5)
	PHYSICS_MIN_SPEED     = float32(0.000001)
)

type cMovable struct {
	velocity         mgl32.Vec3
	accel            mgl32.Vec3
	speedMax         float32
	isFrictionActive bool
	followTarget     *Entity
}

func (comp *cMovable) update(entity *Entity) {
	if comp.followTarget != nil {
		comp.setFollowVelocity(entity)
	}
	comp.updateKinematics(entity)
	var drawComponent *cDrawable
	if tmp, err := getComponent[*cDrawable](CMP_DRAWABLE, entity); err != nil {
		return
	} else {
		drawComponent = *tmp
	}

	// update animation based on velocity
	// vertical animation takes priority (eventually falling would be a thing)
	var animType AnimationIndex
	if comp.velocity[1] > 0 {
		animType = ANIM_MOVE_VUP
	} else if comp.velocity[1] < 0 {
		animType = ANIM_MOVE_VDOWN
	} else if comp.velocity[0] > 0 {
		animType = ANIM_MOVE_HRIGHT
	} else if comp.velocity[0] < 0 {
		animType = ANIM_MOVE_HLEFT
	} else {
		animType = ANIM_IDLE
	}
	drawComponent.setAnimation(animType)
}

func (comp *cMovable) getType() ComponentType {
	return CMP_MOVABLE
}

func (comp *cMovable) updateKinematics(entity *Entity) {
	speedMax := comp.speedMax
	velocityMax := float32(speedMax)
	velocityMin := float32(-speedMax)
	zerof := float32(0)

	// if player is not accelerating, apply friction
	for i := 0; i < 2; i++ {
		if comp.accel[i] != zerof {
			comp.velocity[i] += comp.accel[i]
			comp.velocity[i] = clamp(comp.velocity[i], velocityMin, velocityMax)
		} else if comp.velocity[i] != zerof && comp.isFrictionActive {
			comp.velocity[i] *= PHYSICS_FRICTION_COEF
			if mgl32.Abs(comp.velocity[i]) < PHYSICS_MIN_SPEED {
				comp.velocity[i] = zerof
			}
		}
	}

	entity.setPosition(
		entity.getPosition().Add(
			comp.velocity.Mul(
				DeltaTime.get())))
}

func (comp *cMovable) setFollowVelocity(entity *Entity) {
	targetPos := comp.followTarget.getPosition()
	if tmp, err := getComponent[*cDrawable](CMP_DRAWABLE, comp.followTarget); err == nil {
		drawComponent := *tmp
		frameSizeX, frameSizeY := drawComponent.getFrameSize()
		targetPos[0] += (frameSizeX * pixelsPerTexel / windowWidth)
		targetPos[1] += (frameSizeY * pixelsPerTexel / windowHeight)
	}
	distance := targetPos.Sub(entity.getPosition())
	comp.velocity = mgl32.Vec3{distance[0], distance[1], 0.0}.Mul(5)
	// entity.position = targetPos // for testing TODO
}
