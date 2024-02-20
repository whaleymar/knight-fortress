package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

// movement based animations here: fall, jump
// TODO custom type
const (
	ANIM_IDLE = iota
	ANIM_HMOVE
	ANIM_VMOVE
)

type cMovable struct {
	velocity mgl32.Vec3
	accel    mgl32.Vec3
	speedMax float32
}

func (comp *cMovable) update(entity *Entity) {
	// TODO magic numbers
	// also code kinda disgusting

	speedMax := comp.speedMax * DeltaTime.get()
	velocityMax := float32(speedMax)
	velocityMin := float32(-speedMax)
	zero1d := float32(0)
	// zero2d := mgl32.Vec2{0,0}
	cutoffSpeed := float32(speedMax / 20.0)
	frictionCoefficient := float32(0.5)

	for i := 0; i < 2; i++ {
		if comp.accel[i] != zero1d {
			comp.velocity[i] += comp.accel[i]
			if comp.velocity[i] > velocityMax {
				comp.velocity[i] = velocityMax
			} else if comp.velocity[i] < velocityMin {
				comp.velocity[i] = velocityMin
			}
		} else if comp.velocity[i] != zero1d {
			comp.velocity[i] *= frictionCoefficient
			if (comp.velocity[i] > zero1d && comp.velocity[i] < cutoffSpeed) || (comp.velocity[i] < zero1d && comp.velocity[i] > -cutoffSpeed) {
				comp.velocity[i] = zero1d
			}
		}
	}

	entity.setPosition(entity.getPosition().Add(comp.velocity))

	var drawComponent *cDrawable
	if tmp, err := getComponent[*cDrawable](CMP_DRAWABLE, entity); err != nil {
		return
	} else {
		drawComponent = *tmp
	}

	// update animation based on velocity
	var animIx int
	if comp.velocity[1] != 0 {
		animIx = ANIM_VMOVE
	} else if comp.velocity[0] != 0 {
		animIx = ANIM_HMOVE
	} else {
		animIx = ANIM_IDLE
	}
	drawComponent.animManager.setAnimation(animIx)
}

func (comp *cMovable) getType() ComponentType {
	return CMP_MOVABLE
}
