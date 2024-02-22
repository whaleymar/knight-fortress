package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	// "github.com/go-gl/glfw/v3.3/glfw"
)

// define depth order for sorting
const (
	DEPTH_BACKGROUND = iota
	DEPTH_NPC
	DEPTH_PLAYER
)

type cDrawable struct {
	enum      ComponentType
	vertices  []float32
	vao       VAO
	vbo       VBO
	sprite    Sprite
	textureIx TextureSlot
}

type Sprite struct { // spritesheet
	sheetPosition mgl32.Vec3 // stores texture array position + x,y position in texture atlas
	frameSize     [2]int
	animManager   AnimationManager
}

type AnimationManager struct {
	anims     []Animation
	animSpeed float32
	frame     int
	frameTime float32
	animIx    int
}

type Animation struct {
	textureOffset [2]int // relative to Sprite.sheetPosition
	frameCount    int
}

func (comp *cDrawable) update(entity *Entity) {
	animManager := &comp.sprite.animManager
	if animManager.animSpeed > 0.0 {
		// check if should update animation frame
		animManager.frameTime += DeltaTime.get()
		if animManager.frameTime >= 1/animManager.animSpeed {
			animManager.frameTime = 0.0
			animManager.frame = (animManager.frame + 1) % animManager.getAnimation().frameCount
		}
	}

	// update UV
	// todo these casts are so ugly
	// TODO methods for these so I don't go insane
	// TODO only do once for static animation
	var xMin, xMax, yMin, yMax float32
	sheetSizeX, sheetSizeY := getTextureManager().getTextureSize(comp.textureIx, 0) // TODO hard coded array Index

	pixelOffset := comp.sprite.sheetPosition.X() + float32(comp.sprite.frameSize[0]*(comp.sprite.animManager.getAnimation().textureOffset[0]+comp.sprite.animManager.frame))
	xMin = pixelOffset / float32(sheetSizeX)
	xMax = (pixelOffset + float32(comp.sprite.frameSize[0])) / float32(sheetSizeX)

	pixelOffset = comp.sprite.sheetPosition.Y() + float32(comp.sprite.frameSize[1]*comp.sprite.animManager.getAnimation().textureOffset[1])
	yMin = pixelOffset / float32(sheetSizeY)
	yMax = (pixelOffset + float32(comp.sprite.frameSize[1])) / float32(sheetSizeY)

	comp.vertices = scaleDepth(makeSquareVerticesWithUV(comp.sprite.frameSize[0]*pixelsPerTexel, comp.sprite.frameSize[1]*pixelsPerTexel, xMin, xMax, yMin, yMax), 0.2)
}

func (comp *cDrawable) getType() ComponentType {
	return CMP_DRAWABLE
}

func (animManager *AnimationManager) getAnimation() Animation {
	return animManager.anims[animManager.animIx]
}

func (animManager *AnimationManager) setAnimation(animIx int) {
	if animIx == animManager.animIx {
		return
	} else if animIx >= len(animManager.anims) {
		return
	}
	animManager.animIx = animIx
	animManager.frame = 0
	animManager.frameTime = 0.0
}

// func (animManager *AnimationManager) makeFrame(sprite Sprite) *image.RGBA {
// 	anim := animManager.getAnimation()
//
// 	y0 := sprite.frameHeight * anim.fileoffset
// 	x0 := sprite.frameWidth * animManager.frame
// 	rect := image.Rect(x0, y0, x0+sprite.frameWidth, y0+sprite.frameHeight)
//
// 	img := sprite.img.SubImage(rect)
// 	rgba := image.NewRGBA(
// 		image.Rect(
// 			0,
// 			0,
// 			img.Bounds().Dx(),
// 			img.Bounds().Dy(),
// 		),
// 	)
// 	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
//
// 	// DEBUGGING:
// 	// saveImage(rgba, fmt.Sprint("tmp/test", glfw.GetTime()))
// 	// if isTransparent(rgba) {
// 	// 	fmt.Println("Found transparent image")
// 	// 	fmt.Println("Index: ", animManager.ix)
// 	// 	fmt.Println("Frame: ", animManager.frame)
// 	// 	fmt.Println("anim offset: ", anim.fileoffset)
// 	// }
//
// 	return rgba
// }

func loadImage(filename string) (*image.RGBA, error) {
	imgFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Image %q not found on disk: %v", filename, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
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

func isTransparent(img *image.RGBA) bool {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, alpha := img.At(x, y).RGBA()
			if alpha != 0 {
				return false
			}
		}
	}

	return true
}
