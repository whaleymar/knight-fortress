package phys

import (
	"fmt"

	"reflect"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/whaleymar/knight-fortress/src/math"
)

type RigidBodyType int

const (
	RIGIDBODY_NONE RigidBodyType = iota
	RIGIDBODY_STATIC
	RIGIDBODY_DYNAMIC
	RIGIDBODY_KINEMATIC
)

type Collider interface {
	SetPosition(Point)
	GetPosition() Point
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

	if px < py {
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
