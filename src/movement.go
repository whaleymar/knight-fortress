package main

import (
	"fmt"

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
	FOLLOW_CLAMP_DISTANCE   = float32(0.0001)
	FOLLOW_SPEED_MIN        = float32(0.1)
	FOLLOW_SPEED_MULTIPLIER = float32(10.0)
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

	// if not accelerating, apply friction
	for i := 0; i < 2; i++ {
		if comp.accel[i] != zerof {
			comp.velocity[i] += comp.accel[i]
			comp.velocity[i] = clamp(comp.velocity[i], velocityMin, velocityMax)
		} else if comp.velocity[i] != zerof && comp.isFrictionActive {
			comp.velocity[i] *= (1 - PHYSICS_FRICTION_COEF)
			if mgl32.Abs(comp.velocity[i]) < PHYSICS_MIN_SPEED {
				comp.velocity[i] = zerof
			}
		}
	}

	entity.setPosition(comp.getStepDistance(entity))
}

func (comp *cMovable) getStepDistance(entity *Entity) mgl32.Vec3 {
	return entity.getPosition().Add(comp.velocity.Mul(DeltaTime.get()))
}

func (comp *cMovable) setFollowVelocity(entity *Entity) {
	targetPos := comp.followTarget.getPosition()
	// entity.setPosition(targetPos)
	if tmp, err := getComponent[*cDrawable](CMP_DRAWABLE, comp.followTarget); err == nil {
		drawComponent := *tmp
		frameSizeX, frameSizeY := drawComponent.getFrameSize()
		targetPos[0] += (frameSizeX * pixelsPerTexel / windowWidth)
		targetPos[1] += (frameSizeY * pixelsPerTexel / windowHeight)
	}
	curPosition := entity.getPosition()
	if targetPos == curPosition {
		return
	}

	distance := targetPos.Sub(curPosition)
	comp.velocity = mgl32.Vec3{distance[0], distance[1], 0.0}.Mul(FOLLOW_SPEED_MULTIPLIER)
	nextDistance := targetPos.Sub(comp.getStepDistance(entity))
	newPosition := curPosition
	newVelocity := comp.velocity

	for i := 0; i < 2; i++ {
		if distance[i] == 0.0 {
			continue
		}
		// if we're really close OR we're going to overshoot the target, set position to target
		shouldClamp := mgl32.Abs(distance[i]) <= FOLLOW_CLAMP_DISTANCE
		shouldClamp = shouldClamp || ((nextDistance[i] >= 0) != (distance[i] >= 0))
		if !shouldClamp {
			continue
		}
		newPosition[i] = targetPos[i]
		newVelocity[i] = 0.0
	}
	entity.setPosition(newPosition)
	comp.velocity = newVelocity
}
