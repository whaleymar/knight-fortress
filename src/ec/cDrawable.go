package ec

import (
	"fmt"
	"sync"

	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/sys"
)

// define depth order for sorting
const (
	DEPTH_BACKGROUND = iota
	DEPTH_NPC
	DEPTH_GROUND
	DEPTH_PLAYER
)

var SCALE_NORMAL = [2]float32{1.0, 1.0}

type CDrawable struct {
	vertices         []float32
	vao              gfx.VAO
	vbo              gfx.VBO
	scale            [2]float32
	sprite           Sprite
	textureIx        gfx.TextureSlot
	rwmutex          *sync.RWMutex
	isUvUpdateNeeded bool
}

type Sprite struct { // spritesheet
	SheetPosition [3]int // stores texture array position + x,y position in texture atlas
	FrameSize     [2]int
	AnimManager   AnimationManager
}

type AnimationManager struct {
	Anims     []Animation
	AnimSpeed float32
	frame     int
	frameTime float32
	animIx    AnimationIndex
}

type Animation struct {
	TextureOffset [2]int // relative to Sprite.sheetPosition
	FrameCount    int
}

func (comp *CDrawable) update(entity *Entity) {
	animManager := &comp.sprite.AnimManager
	if animManager.AnimSpeed > 0.0 {
		// check if should update animation frame
		animManager.frameTime += sys.DeltaTime.Get()
		if animManager.frameTime >= 1/animManager.AnimSpeed {
			animManager.frameTime = 0.0
			animManager.frame = (animManager.frame + 1) % animManager.getAnimation().FrameCount
			comp.isUvUpdateNeeded = true
		}
	}

	// update UV
	if comp.isUvUpdateNeeded {
		var xMin, xMax, yMin, yMax float32
		sheetSizeX, sheetSizeY := gfx.GetTextureManager().GetTextureSize(comp.textureIx, 0) // TODO hard coded array Index

		pixelOffset := float32(comp.sprite.SheetPosition[0] + comp.sprite.FrameSize[0]*(comp.sprite.AnimManager.getAnimation().TextureOffset[0]+comp.sprite.AnimManager.frame))
		xMin = pixelOffset / sheetSizeX
		xMax = (pixelOffset + float32(comp.sprite.FrameSize[0])) / sheetSizeX

		pixelOffset = float32(comp.sprite.SheetPosition[1] + comp.sprite.FrameSize[1]*comp.sprite.AnimManager.getAnimation().TextureOffset[1])
		yMin = pixelOffset / sheetSizeY
		yMax = (pixelOffset + float32(comp.sprite.FrameSize[1])) / sheetSizeY

		comp.vertices = gfx.MakeRectVerticesWithUV(comp.getPixelSize(false), comp.getPixelSize(true), xMin, xMax, yMin, yMax)
	}
	comp.isUvUpdateNeeded = false
}

func (comp *CDrawable) getType() ComponentType {
	return CMP_DRAWABLE
}

func (comp *CDrawable) onDelete() {
	comp.vao.Free()
	comp.vbo.Free()
}

func (comp *CDrawable) Copy() (Component, error) {
	cHolder := comp.GetSaveData()
	newComp, err := LoadComponentDrawable(cHolder.YamlData)
	return &newComp, err
}

func (comp *CDrawable) GetSaveData() componentHolder {
	data, err := sys.StructToYaml(struct {
		Scale     [2]float32
		PixelData Sprite
	}{
		Scale:     comp.scale,
		PixelData: comp.sprite,
	})
	if err != nil {
		panic(err)
	}
	return makeComponentHolder(comp.getType(), data)
}

func LoadComponentDrawable(componentData string) (CDrawable, error) {
	data := struct {
		Scale     [2]float32
		PixelData Sprite
	}{}
	err := sys.YamlToStruct(componentData, &data)
	if err != nil {
		return CDrawable{}, fmt.Errorf("Couldn't load draw data from %s", componentData)
	}

	return CDrawable{
		vertices:         gfx.SquareVertices,
		vao:              gfx.MakeVao(),
		vbo:              gfx.MakeVbo(),
		scale:            data.Scale,
		sprite:           data.PixelData,
		textureIx:        gfx.TEX_MAIN,
		rwmutex:          &sync.RWMutex{},
		isUvUpdateNeeded: true,
	}, nil
}

func (comp *CDrawable) getAnimation() Animation {
	comp.rwmutex.RLock()
	defer comp.rwmutex.RUnlock()

	return comp.sprite.AnimManager.getAnimation()
}

func (comp *CDrawable) setAnimation(animIx AnimationIndex) {
	comp.rwmutex.Lock()
	defer comp.rwmutex.Unlock()

	if comp.sprite.AnimManager.setAnimation(animIx) {
		comp.isUvUpdateNeeded = true
	}
}

func (comp *CDrawable) getFrameSize() (float32, float32) {
	frameSize := comp.sprite.FrameSize
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

func (comp *CDrawable) getPixelSize(vertical bool) float32 {
	ix := 0
	if vertical {
		ix = 1
	}
	return float32(comp.sprite.FrameSize[ix]*gfx.PixelsPerTexel) * comp.scale[ix]
}

func (animManager *AnimationManager) getAnimation() Animation {
	return animManager.Anims[animManager.animIx]
}

func (animManager *AnimationManager) setAnimation(animIx AnimationIndex) bool {
	if animIx == animManager.animIx {
		return false
	} else if int(animIx) >= len(animManager.Anims) {
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
