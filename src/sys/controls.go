package sys

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type ButtonState int
type ButtonCallback func(ButtonState)

const (
	BUTTONSTATE_OFF ButtonState = iota
	BUTTONSTATE_ON
	BUTTONSTATE_HELD
)

var _CONTROLS_MUTEX = &sync.Mutex{}
var controlsPtr *ControlsManager

func GetControlsManager() *ControlsManager {
	if controlsPtr == nil {
		_CONTROLS_MUTEX.Lock()
		defer _CONTROLS_MUTEX.Unlock()
		if controlsPtr == nil {
			controlsPtr = &ControlsManager{}
		}
	}
	return controlsPtr
}

type ControlsManager struct {
	buttons []*ButtonStateMachine
}

func (controlManager *ControlsManager) Add(bsm ButtonStateMachine) {
	_CONTROLS_MUTEX.Lock()
	defer _CONTROLS_MUTEX.Unlock()
	controlManager.buttons = append(controlManager.buttons, &bsm)
}

func (controlManager *ControlsManager) Update() {
	for _, button := range controlManager.buttons {
		button.Update()
	}
}

type ButtonStateMachine struct {
	Key              glfw.Key
	State            ButtonState
	Callback         ButtonCallback
	StateTimeLimit   float32
	StateTimeElapsed float32
	IsAsleep         bool
}

func ControlsCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	for _, button := range GetControlsManager().buttons {
		if key == button.Key {
			button.ProcessInput(action)
		}
	}
}

func (bsm *ButtonStateMachine) ProcessInput(action glfw.Action) {
	switch action {
	case glfw.Press:
		bsm.State = BUTTONSTATE_ON
		bsm.StateTimeElapsed = 0.0
		bsm.IsAsleep = false
		break
	case glfw.Repeat:
		break
	case glfw.Release:
		bsm.State = BUTTONSTATE_OFF
		break
	}
}

func (bsm *ButtonStateMachine) Update() {
	if bsm.IsAsleep {
		return
	}
	if bsm.State == BUTTONSTATE_OFF {
		defer bsm.Sleep()
	}
	if bsm.StateTimeLimit > 0.0 {
		bsm.StateTimeElapsed += DeltaTime.Get()
		if bsm.StateTimeElapsed >= bsm.StateTimeLimit {
			bsm.State = BUTTONSTATE_OFF
			defer bsm.Sleep()
		}
	}
	bsm.Callback(bsm.State)
	bsm.State = BUTTONSTATE_HELD
}

func (bsm *ButtonStateMachine) Sleep() {
	bsm.IsAsleep = true
}
