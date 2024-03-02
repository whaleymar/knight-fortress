package phys

import (
	"fmt"

	"reflect"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/whaleymar/knight-fortress/src/math"
)

type ColliderEnum int

const (
	COLLIDER_AABB ColliderEnum = iota
	COLLIDER_CIRCLE
)

const (
	MIN_DELTA_HALT = 0.005
)

// holds members for all structs that implement Collider (needed to serialize/deserialize)
type SuperCollider struct {
	ColliderType ColliderEnum
	Center       Point
	Half         Point
	Radius       float32
}

func ExtractCollider(super SuperCollider) (Collider, error) {
	switch super.ColliderType {
	case COLLIDER_AABB:
		return &AABB{
			super.Center,
			super.Half,
		}, nil
	default:
		return &AABB{}, fmt.Errorf("Unrecognized Collider Type: %d", super.ColliderType)
	}
}

type Collider interface {
	SetPosition(Point)
	GetPosition() Point
	MakeSuperCollider() SuperCollider
	CheckCollision(*Collider) Hit
	CheckCollisionAABB(*AABB) Hit
}

type Hit struct {
	IsHit    bool
	Position Point
	Delta    Point
	Normal   Point
}

type AABB struct {
	Center Point
	Half   Point
}

func (this *AABB) GetPosition() Point {
	return this.Center
}

func (this *AABB) SetPosition(pos Point) {
	this.Center = pos
}

func (this *AABB) MakeSuperCollider() SuperCollider {
	return SuperCollider{
		ColliderType: COLLIDER_AABB,
		Center:       this.Center,
		Half:         this.Half,
		Radius:       0,
	}
}

func (this *AABB) CheckCollision(other *Collider) Hit {
	switch (*other).(type) {
	case *AABB:
		aabb := (*other).(*AABB)
		return this.CheckCollisionAABB(aabb)
	default:
		panic(fmt.Sprintf("Collision check not implemented for collider with type %v", reflect.TypeOf(other)))

	}
}

func (this *AABB) CheckCollisionAABB(other *AABB) Hit {
	// yoinked from https://noonat.github.io/intersect/
	hit := Hit{}

	dx := other.Center.X - this.Center.X
	px := (other.Half.X + this.Half.X) - mgl32.Abs(dx)
	if px <= 0 {
		return hit
	}

	dy := other.Center.Y - this.Center.Y
	py := (other.Half.Y + this.Half.Y) - mgl32.Abs(dy)
	if py <= 0 {
		return hit
	}
	hit.IsHit = true

	// if math.Abs(px-py) < MIN_DELTA_HALT {
	if px == py {
		sign := math.Sign(dx)
		hit.Delta.X = px * sign
		hit.Normal.X = sign
		hit.Position.X = this.Center.X + this.Half.X*sign

		sign = math.Sign(dy)
		hit.Delta.Y = py * sign
		hit.Normal.Y = sign
		hit.Position.Y = this.Center.Y + this.Half.Y*sign
	} else if px < py {
		sign := math.Sign(dx)
		hit.Delta.X = px * sign
		hit.Normal.X = sign
		hit.Position.X = this.Center.X + this.Half.X*sign
		hit.Position.Y = other.Center.Y
	} else {
		sign := math.Sign(dy)
		hit.Delta.Y = py * sign
		hit.Normal.Y = sign
		hit.Position.X = other.Center.X
		hit.Position.Y = this.Center.Y + this.Half.Y*sign
	}

	return hit
}
