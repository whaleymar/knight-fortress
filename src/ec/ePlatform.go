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
		mgl32.Vec3{-8.0, 0.0, DEPTH_GROUND},
		&sync.RWMutex{},
	}

	entity.components.Add(&CDrawable{
		gfx.MakeSquareVertices(16, 16),
		gfx.MakeVao(),
		gfx.MakeVbo(),
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
		&phys.AABB{phys.Point{}, phys.Point{0.25, 0.25}}, // TODO hard coded
		phys.RIGIDBODY_STATIC,
	})

	return entity
}
