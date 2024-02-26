package ec

import (
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/sys"
	"sync"
)

// define depth order for sorting
const (
	DEPTH_BACKGROUND = iota
	DEPTH_NPC
	DEPTH_GROUND
	DEPTH_PLAYER
)

type CDrawable struct {
	vertices         []float32
	vao              gfx.VAO
	vbo              gfx.VBO
	sprite           Sprite
	textureIx        gfx.TextureSlot
	rwlock           *sync.RWMutex
	isUvUpdateNeeded bool
}

type Sprite struct { // spritesheet
	sheetPosition [3]int // stores texture array position + x,y position in texture atlas
	frameSize     [2]int
	animManager   AnimationManager
}

type AnimationManager struct {
	anims     []Animation
	animSpeed float32
	frame     int
	frameTime float32
	animIx    AnimationIndex
}

type Animation struct {
	textureOffset [2]int // relative to Sprite.sheetPosition
	frameCount    int
}

func (comp *CDrawable) update(entity *Entity) {
	animManager := &comp.sprite.animManager
	if animManager.animSpeed > 0.0 {
		// check if should update animation frame
		animManager.frameTime += sys.DeltaTime.Get()
		if animManager.frameTime >= 1/animManager.animSpeed {
			animManager.frameTime = 0.0
			animManager.frame = (animManager.frame + 1) % animManager.getAnimation().frameCount
			comp.isUvUpdateNeeded = true
		}
	}

	// update UV
	if comp.isUvUpdateNeeded {
		var xMin, xMax, yMin, yMax float32
		sheetSizeX, sheetSizeY := gfx.GetTextureManager().GetTextureSize(comp.textureIx, 0) // TODO hard coded array Index

		pixelOffset := float32(comp.sprite.sheetPosition[0] + comp.sprite.frameSize[0]*(comp.sprite.animManager.getAnimation().textureOffset[0]+comp.sprite.animManager.frame))
		xMin = pixelOffset / sheetSizeX
		xMax = (pixelOffset + float32(comp.sprite.frameSize[0])) / sheetSizeX

		pixelOffset = float32(comp.sprite.sheetPosition[1] + comp.sprite.frameSize[1]*comp.sprite.animManager.getAnimation().textureOffset[1])
		yMin = pixelOffset / sheetSizeY
		yMax = (pixelOffset + float32(comp.sprite.frameSize[1])) / sheetSizeY

		comp.vertices = gfx.MakeSquareVerticesWithUV(comp.sprite.frameSize[0]*gfx.PixelsPerTexel, comp.sprite.frameSize[1]*gfx.PixelsPerTexel, xMin, xMax, yMin, yMax)
	}
	comp.isUvUpdateNeeded = false
}

func (comp *CDrawable) getType() ComponentType {
	return CMP_DRAWABLE
}

func (comp *CDrawable) getAnimation() Animation {
	comp.rwlock.RLock()
	defer comp.rwlock.RUnlock()

	return comp.sprite.animManager.getAnimation()
}

func (comp *CDrawable) setAnimation(animIx AnimationIndex) {
	comp.rwlock.Lock()
	defer comp.rwlock.Unlock()

	if comp.sprite.animManager.setAnimation(animIx) {
		comp.isUvUpdateNeeded = true
	}
}

func (comp *CDrawable) getFrameSize() (float32, float32) {
	frameSize := comp.sprite.frameSize
	return float32(frameSize[0]), float32(frameSize[1])
}

func (comp *CDrawable) GetVao() gfx.VAO {
	return comp.vao
}

func (comp *CDrawable) GetVbo() gfx.VBO {
	return comp.vbo
}

func (comp *CDrawable) GetVertices() []float32 {
	return comp.vertices
}

func (animManager *AnimationManager) getAnimation() Animation {
	return animManager.anims[animManager.animIx]
}

func (animManager *AnimationManager) setAnimation(animIx AnimationIndex) bool {
	if animIx == animManager.animIx {
		return false
	} else if int(animIx) >= len(animManager.anims) {
		return false
	}
	animManager.animIx = animIx
	animManager.frame = 0
	animManager.frameTime = 0.0
	return true
}

func makeStaticAnimationManager() AnimationManager {
	anim := []Animation{
		{
			[2]int{0, 0},
			1,
		},
	}
	return AnimationManager{
		anim,
		0.0,
		0,
		0.0,
		0,
	}
}
