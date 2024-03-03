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

const (
	LEVEL_NAME_DEFAULT = "Demo"
	ASSET_DIR_LEVEL    = "assets/level"
)

var _LEVEL_LOCK = &sync.RWMutex{}
var levelPtr *level

func GetCurrentLevel() *level {
	if levelPtr == nil {
		_LEVEL_LOCK.Lock()
		defer _LEVEL_LOCK.Unlock()
		if levelPtr == nil {
			levelPtr = &level{
				LEVEL_NAME_DEFAULT,
				phys.ORIGIN,
				[]uint64{},
			}
		}
	}
	return levelPtr
}

func TryLoadLevel(levelname string) error {
	_LEVEL_LOCK.Lock()
	defer _LEVEL_LOCK.Unlock()

	if levelPtr != nil {
		levelPtr.free()
	}
	newLevel, err := loadLevel(levelname)
	if err != nil {
		return err
	}
	levelPtr = &newLevel
	levelPtr.Reset()
	return nil
}

// TODO needs mutex
type level struct {
	name          string
	startPosition mgl32.Vec3
	entityIDs     []uint64
	// TODO level data path/handle, metadata like the name for quick access
}

func (lvl *level) addChild(uid uint64) {
	lvl.entityIDs = append(lvl.entityIDs, uid)
}

func (lvl *level) free() {
	entityManager := ec.GetEntityManager()
	for _, uid := range lvl.entityIDs {
		entityManager.Remove(uid)
	}
	lvl.entityIDs = nil
}

func (lvl *level) Reset() {
	ec.GetPlayerPtr().SetPosition(lvl.startPosition)
	moveComponent, err := ec.GetComponent[*ec.CMovable](ec.CMP_MOVABLE, ec.GetPlayerPtr())
	if err == nil {
		(*moveComponent).SetVelocity(phys.ORIGIN)
	}
}

func (lvl *level) SaveToFile() error {
	_LEVEL_LOCK.RLock()
	_LEVEL_LOCK.RUnlock()
	type entitydata struct {
		EntityName string
		Position   mgl32.Vec3
	}
	savedata := struct {
		Name          string
		StartPosition mgl32.Vec3
		EntityData    []entitydata
	}{
		Name:          lvl.name,
		StartPosition: lvl.startPosition,
		EntityData:    []entitydata{},
	}

	names := make(map[string]bool)
	entityManager := ec.GetEntityManager()
	for _, uid := range lvl.entityIDs {
		entity, err := entityManager.Get(uid)
		if err != nil {
			continue
		}
		if !names[entity.Name] {
			names[entity.Name] = true
			entity.SaveToFile()
		}
		savedata.EntityData = append(savedata.EntityData, entitydata{entity.Name, entity.GetPosition()})
	}

	return sys.SaveStruct(getLevelPath(lvl.name), savedata)
}

func loadLevel(filename string) (level, error) {
	type entitydata struct {
		EntityName string
		Position   mgl32.Vec3
	}
	savedata := struct {
		Name          string
		StartPosition mgl32.Vec3
		EntityData    []entitydata
	}{}
	err := sys.LoadStruct(getLevelPath(filename), &savedata)
	if err != nil {
		return level{}, err
	}

	newLevel := level{
		savedata.Name,
		savedata.StartPosition,
		[]uint64{},
	}

	entityMgr := ec.GetEntityManager()
	entityCache := map[string]ec.Entity{}
	for _, entityData := range savedata.EntityData {
		entity, ok := entityCache[entityData.EntityName]
		if !ok {
			newEntity, err := ec.LoadEntity(entityData.EntityName)
			if err != nil {
				return level{}, err
			}
			entityCache[entityData.EntityName] = newEntity
			entity = newEntity
		}
		newEntity, err := entity.Copy()
		if err != nil {
			return level{}, err
		}
		newEntity.SetPosition(entityData.Position)
		uid, err := entityMgr.Add(&newEntity)
		if err != nil {
			return level{}, err
		}
		// e2.SetPosition(entityData.Position)
		newLevel.addChild(uid)
	}

	// testing
	// entities := entityMgr.GetEntitiesWithComponent(ec.CMP_ANY)
	// for _, e := range entities {
	// 	fmt.Printf("%p\n", e)
	// 	fmt.Println(e.GetPosition(), "\n")
	// }
	return newLevel, nil
}

func CreateLevelControls() {

	// reset level
	sys.GetControlsManager().Add(sys.ButtonStateMachine{
		Key:   glfw.Key0,
		State: sys.BUTTONSTATE_OFF,
		Callback: func(state sys.ButtonState) {
			GetCurrentLevel().Reset()
		},
		StateTimeLimit:   0.0,
		StateTimeElapsed: 0.0,
		IsAsleep:         true,
	})

	sys.GetControlsManager().Add(sys.ButtonStateMachine{
		Key:   glfw.Key9,
		State: sys.BUTTONSTATE_OFF,
		Callback: func(state sys.ButtonState) {
			currentLevelName := GetCurrentLevel().name
			TryLoadLevel(currentLevelName)
		},
		StateTimeLimit:   0.0,
		StateTimeElapsed: 0.0,
		IsAsleep:         true,
	})
}

func getLevelPath(filename string) string {
	return fmt.Sprintf("%s/%s.yml", ASSET_DIR_LEVEL, filename)
}
