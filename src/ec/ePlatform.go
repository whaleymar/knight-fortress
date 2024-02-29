package ec

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/phys"
)

func MakePlatformBasic() Entity {
	entity := Entity{
		0,
		"Platform",
		&ComponentList{},
		mgl32.Vec3{0.0, -0.5, DEPTH_GROUND},
		&sync.RWMutex{},
	}

	entity.components.Add(&CDrawable{
		gfx.MakeRectVertices(16, 16),
		gfx.MakeVao(),
		gfx.MakeVbo(),
		[2]float32{4.0, 1.0},
		Sprite{
			[3]int{0, 80, 0},
			[2]int{16, 16},
			makeStaticAnimationManager(),
		},
		gfx.TEX_MAIN,
		&sync.RWMutex{},
		true,
	})

	entity.components.Add(&CCollides{
		&phys.AABB{phys.Point{}, phys.Point{1.0, 0.25}}, // TODO hard coded
		phys.RigidBody{phys.RIGIDBODY_STATIC, phys.RBSTATE_STILL},
		true,
	})

	return entity
}
