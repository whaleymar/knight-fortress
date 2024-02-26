package phys

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/whaleymar/knight-fortress/src/math"
	"github.com/whaleymar/knight-fortress/src/sys"
)

// TODO interface & have rect and ellipse shapes
type AABB struct {
	// (0,0) is bottom left
	X float32
	Y float32
}

type Force interface {
	Apply(mgl32.Vec3) mgl32.Vec3
}

// this needs to not get overwritten
type ContinuousForce struct {
	Direction mgl32.Vec3
	Magnitude float32
	Time      float32
}

type ImpulseForce struct {
	Direction mgl32.Vec3
	Magnitude float32
}

func (force *ContinuousForce) Apply(curVelocity mgl32.Vec3) mgl32.Vec3 {
	mag := force.Magnitude
	oldTime := force.Time
	force.Time -= sys.DeltaTime.Get()
	if force.Time < 0.0 {
		mag *= oldTime / sys.DeltaTime.Get()
		force.Time = 0.0
	}
	return curVelocity.Add(force.Direction.Mul(mag))
}

func (force *ImpulseForce) Apply(curVelocity mgl32.Vec3) mgl32.Vec3 {
	return curVelocity.Add(force.Direction.Mul(force.Magnitude))
}

func CheckCollision(position1, position2 mgl32.Vec3, shape1, shape2 AABB) bool {
	xOverlap := math.Between(position2[0], position1[0], position2[0]+shape2.X) || math.Between(position1[0], position2[0], position1[0]+shape1.X)
	yOverlap := math.Between(position2[1], position1[1], position2[1]+shape2.Y) || math.Between(position1[1], position2[1], position1[1]+shape1.Y)
	return xOverlap && yOverlap
}
