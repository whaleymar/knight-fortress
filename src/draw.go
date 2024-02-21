package main

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	// "github.com/go-gl/glfw/v3.3/glfw"
)

type cDrawable struct {
	enum        ComponentType
	vertices    []float32
	vao         uint32
	vbo         uint32
	sprite      Sprite
	animManager AnimationManager
}

type Sprite struct {
	img         *image.RGBA
	frameHeight int
	frameWidth  int
	textureIx   uint32
}

type Animation struct {
	fileoffset int
	frameCount int
}

type AnimationManager struct {
	anims     []Animation
	animSpeed float32
	frame     int
	frameTime float32
	animIx    int
}

func (comp *cDrawable) update(entity *Entity) {
	animManager := &comp.animManager
	if animManager.animSpeed == 0.0 { // static image
		return
	}
	animManager.frameTime += DeltaTime.get()
	if animManager.frameTime >= 1/animManager.animSpeed {
		animManager.frameTime = 0.0
		animManager.frame = (animManager.frame + 1) % animManager.getAnimation().frameCount
	}
}

func (comp *cDrawable) getType() ComponentType {
	return CMP_DRAWABLE
}

func (comp *cDrawable) getTexture() (uint32, error) {
	return loadTexture(comp.getFrame(), comp.sprite.textureIx)
}

func (comp *cDrawable) getFrame() *image.RGBA {
	return comp.animManager.makeFrame(comp.sprite)
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

func (animManager *AnimationManager) makeFrame(sprite Sprite) *image.RGBA {
	anim := animManager.getAnimation()

	y0 := sprite.frameHeight * anim.fileoffset
	x0 := sprite.frameWidth * animManager.frame
	rect := image.Rect(x0, y0, x0+sprite.frameWidth, y0+sprite.frameHeight)

	img := sprite.img.SubImage(rect)
	rgba := image.NewRGBA(
		image.Rect(
			0,
			0,
			img.Bounds().Dx(),
			img.Bounds().Dy(),
		),
	)
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)

	// DEBUGGING:
	// saveImage(rgba, fmt.Sprint("tmp/test", glfw.GetTime()))
	// if isTransparent(rgba) {
	// 	fmt.Println("Found transparent image")
	// 	fmt.Println("Index: ", animManager.ix)
	// 	fmt.Println("Frame: ", animManager.frame)
	// 	fmt.Println("anim offset: ", anim.fileoffset)
	// }

	return rgba
}

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
	anim := []Animation{{0, 1}}
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
