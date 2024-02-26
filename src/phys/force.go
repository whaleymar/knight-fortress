package phys

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/sys"
)

type Force interface {
	Apply(mgl32.Vec3) mgl32.Vec3
}

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
