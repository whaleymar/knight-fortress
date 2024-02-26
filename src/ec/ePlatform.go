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
		mgl32.Vec3{-3.0, 0.0, DEPTH_GROUND},
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
		phys.AABB{0.5, 0.5},
		true,
	})

	return entity
}