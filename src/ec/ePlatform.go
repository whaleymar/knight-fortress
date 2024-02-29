package ec

import (
	"sync"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/phys"
)

const (
	SHEETOFFSET_X_GRASS = 0
	SHEETOFFSET_Y_GRASS = 16

	SHEETOFFSET_X_DIRT = 8
	SHEETOFFSET_Y_DIRT = 16
)

func MakeBasicBlock(sheetOffsetX, sheetOffsetY int) Entity {
	entity := Entity{
		0,
		"Platform",
		&ComponentList{},
		mgl32.Vec3{0.0, -0.5, DEPTH_GROUND},
		&sync.RWMutex{},
	}

	entity.Components.Add(&CDrawable{
		gfx.MakeRectVertices(16, 16), // this will get changed on frame one
		gfx.MakeVao(),
		gfx.MakeVbo(),
		[2]float32{1.0, 1.0},
		Sprite{
			[3]int{sheetOffsetX, sheetOffsetY, 0}, // sheet position
			[2]int{8, 8},                          // frame size
			makeStaticAnimationManager(),
		},
		gfx.TEX_MAIN,
		&sync.RWMutex{},
		true,
	})

	entity.Components.Add(&CCollides{
		&phys.AABB{phys.Point{}, phys.Point{0.125, 0.125}}, // TODO hard coded
		phys.RigidBody{phys.RIGIDBODY_STATIC, phys.RBSTATE_STILL},
		true,
	})

	return entity
}
