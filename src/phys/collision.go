package phys

import (
	// "fmt"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/whaleymar/knight-fortress/src/math"
)

type CollisionDirection int

const (
	COLLIDE_NONE CollisionDirection = iota
	COLLIDE_HORIZONTAL
	COLLIDE_VERTICAL
	COLLIDE_HV
)

type Collider interface {
	SetPosition(Point)
	GetPosition() Point
	CalculateCenter(mgl32.Vec3) Point
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

func (this *AABB) CalculateCenter(entityPos mgl32.Vec3) Point {
	return Point{this.Half.X + entityPos[0], this.Half.Y + entityPos[1]}
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
