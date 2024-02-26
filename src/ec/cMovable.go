package ec

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/math"
	"github.com/whaleymar/knight-fortress/src/phys"
	"github.com/whaleymar/knight-fortress/src/sys"
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

type CMovable struct {
	velocity         mgl32.Vec3
	accel            mgl32.Vec3
	speedMax         float32
	isFrictionActive bool
	followTarget     *Entity
}

func (comp *CMovable) update(entity *Entity) {
	if comp.followTarget != nil {
		comp.setFollowVelocity(entity)
	}
	comp.updateKinematics(entity)
	var drawComponent *CDrawable
	if tmp, err := GetComponent[*CDrawable](CMP_DRAWABLE, entity); err != nil {
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

func (comp *CMovable) getType() ComponentType {
	return CMP_MOVABLE
}

func (comp *CMovable) IsMoving() bool {
	return comp.velocity[0] != 0.0 || comp.velocity[1] != 0.0
}

func (comp *CMovable) updateKinematics(entity *Entity) {
	speedMax := comp.speedMax
	velocityMax := float32(speedMax)
	velocityMin := float32(-speedMax)
	zerof := float32(0)

	// if not accelerating, apply friction
	for i := 0; i < 2; i++ {
		if comp.accel[i] != zerof {
			comp.velocity[i] += comp.accel[i]
			comp.velocity[i] = math.Clamp(comp.velocity[i], velocityMin, velocityMax)
		} else if comp.velocity[i] != zerof && comp.isFrictionActive {
			comp.velocity[i] *= (1 - phys.PHYSICS_FRICTION_COEF)
			if mgl32.Abs(comp.velocity[i]) < phys.PHYSICS_MIN_SPEED {
				comp.velocity[i] = zerof
			}
		}
	}

	entity.SetPosition(comp.GetNextPosition(entity))
}

func (comp *CMovable) GetNextPosition(entity *Entity) mgl32.Vec3 {
	return entity.GetPosition().Add(comp.velocity.Mul(sys.DeltaTime.Get()))
}

func (comp *CMovable) setFollowVelocity(entity *Entity) {
	targetPos := comp.followTarget.GetPosition()
	// entity.setPosition(targetPos)
	if tmp, err := GetComponent[*CDrawable](CMP_DRAWABLE, comp.followTarget); err == nil {
		drawComponent := *tmp
		frameSizeX, frameSizeY := drawComponent.getFrameSize()
		targetPos[0] += (frameSizeX * gfx.PixelsPerTexel / gfx.WindowWidth)
		targetPos[1] += (frameSizeY * gfx.PixelsPerTexel / gfx.WindowHeight)
	}
	curPosition := entity.GetPosition()
	if targetPos == curPosition {
		return
	}

	distance := targetPos.Sub(curPosition)
	comp.velocity = mgl32.Vec3{distance[0], distance[1], 0.0}.Mul(FOLLOW_SPEED_MULTIPLIER)
	nextDistance := targetPos.Sub(comp.GetNextPosition(entity))
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
	entity.SetPosition(newPosition)
	comp.velocity = newVelocity
}
