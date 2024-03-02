package level

import (
	"fmt"
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/whaleymar/knight-fortress/src/ec"
	"github.com/whaleymar/knight-fortress/src/phys"
	"github.com/whaleymar/knight-fortress/src/sys"
)

var _LEVEL_LOCK = &sync.Mutex{}
var levelPtr *level

func GetCurrentLevel() *level {
	if levelPtr == nil {
		_LEVEL_LOCK.Lock()
		defer _LEVEL_LOCK.Unlock()
		if levelPtr == nil {
			levelPtr = &level{}
		}
	}
	return levelPtr
}

type level struct {
	startPosition mgl32.Vec3
	entityIDs     []uint64
	// TODO level data path/handle, metadata like the name for quick access
}

func (lvl *level) addChild(uid uint64) {
	lvl.entityIDs = append(lvl.entityIDs, uid)
}

func (lvl *level) Load() {
	entityManager := ec.GetEntityManager()

	// nSquares := 256
	// squares := make([]ec.Entity, 256)
	// for i := 0; i < nSquares; i++ {
	// 	squares[i] = ec.MakeBasicBlock(ec.SHEETOFFSET_X_GRASS, ec.SHEETOFFSET_Y_GRASS)
	// 	squares[i].SetPosition(mgl32.Vec3{(float32(i) - 128) * 0.25, -0.5, ec.DEPTH_GROUND})
	//
	// 	uid, err := entityManager.Add(&squares[i])
	// 	if err == nil {
	// 		lvl.addChild(uid)
	// 	}
	// }
	//
	// squares[0].SaveToFile()

	square, err := ec.LoadEntity("Block")
	if err == nil {
		uid, err := entityManager.Add(&square)
		if err == nil {
			lvl.addChild(uid)
		}
	} else {
		fmt.Println(err)
	}

	fmt.Println("loading level")
	ec.GetPlayerPtr().SetPosition(lvl.startPosition)
	moveComponent, err := ec.GetComponent[*ec.CMovable](ec.CMP_MOVABLE, ec.GetPlayerPtr())
	if err == nil {
		(*moveComponent).SetVelocity(phys.ORIGIN)
	}
}

func (lvl *level) Reset() {

	entityManager := ec.GetEntityManager()
	for _, uid := range lvl.entityIDs {
		entityManager.Remove(uid)
	}
	lvl.entityIDs = nil

	lvl.Load()
}

func CreateLevelControls() {

	// reset level
	sys.GetControlsManager().Add(sys.ButtonStateMachine{
		glfw.Key0,
		sys.BUTTONSTATE_OFF,
		func(state sys.ButtonState) {
			GetCurrentLevel().Reset()
		},
		0.0,
		0.0,
		false,
	})
}
