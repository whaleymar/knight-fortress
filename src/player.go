package main

import (
	"fmt"
	"image"
	"image/draw"
	// "image/png"
	// "log"
	"os"
	"sync"

	// "unsafe"

	// "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	FILE_SPRITE_PLAYER   = "assets/player.png"
	SPRITE_HEIGHT_PLAYER = 16
	SPRITE_WIDTH_PLAYER  = 16
	ANIM_HOFFSET_PLAYER  = 1
	ANIM_VOFFSET_PLAYER  = 0
	ANIM_HFRAMES_PLAYER  = 3
	ANIM_VFRAMES_PLAYER  = 5
)

type DrawableEntity struct {
	position mgl32.Vec3
	size     mgl32.Vec2
	vao      uint32
	velocity mgl32.Vec3
	accel    mgl32.Vec3
	sprite   Sprite // can't mix go ptrs with C pointers (when passed to player control GLFW function)
	// sprite unsafe.Pointer
	frame int
}

type Sprite struct {
	pixels *image.RGBA
	height int
	width  int
	hAnim  Animation
	vAnim  Animation
}

type Animation struct {
	fileoffset int
	frames     int
}

var lock = &sync.Mutex{}
var playerPtr *DrawableEntity

func getPlayerPtr() *DrawableEntity {
	if playerPtr == nil {
		lock.Lock()
		defer lock.Unlock()
		if playerPtr == nil {
			entity := makePlayerEntity()
			playerPtr = &entity
		}
	}
	return playerPtr
}

func makeDrawableEntity(vao uint32, sprite Sprite) DrawableEntity {
	// func makeDrawableEntity(vao, spriteIx uint32) DrawableEntity {
	entity := DrawableEntity{
		ORIGIN,
		SIZE_STANDARD,
		vao,
		ZERO3,
		ZERO3,
		sprite,
		// gl.Ptr(&sprite), // can't do unsafe ptr to struct?
		0,
	} // TODO make size based on vertices
	return entity
}

func makePlayerEntity() DrawableEntity {
	curVertices := squareVertices
	vao := makeVao(curVertices)
	entity := makeDrawableEntity(vao, makePlayerSprite())
	return entity
}

func makePlayerSprite() Sprite {
	pixels, err := loadImage(FILE_SPRITE_PLAYER)
	if err != nil {
		panic(err)
	}

	hAnim := Animation{ANIM_HOFFSET_PLAYER, ANIM_HFRAMES_PLAYER}
	vAnim := Animation{ANIM_VOFFSET_PLAYER, ANIM_VFRAMES_PLAYER}
	return Sprite{
		pixels,
		SPRITE_HEIGHT_PLAYER,
		SPRITE_WIDTH_PLAYER,
		hAnim,
		vAnim,
	}
}

func (entity DrawableEntity) update() DrawableEntity {
	// this is stinky garbage TODO
	// magic numbers TODO

	speedMax := float32(0.1)
	speedMin := float32(-0.1)
	zero := float32(0)
	cutoff := float32(0.005)
	friction := float32(0.5)

	for i := 0; i < 2; i++ {
		if entity.accel[i] != zero {
			entity.velocity[i] += entity.accel[i]
			if entity.velocity[i] > speedMax {
				entity.velocity[i] = speedMax
			} else if entity.velocity[i] < speedMin {
				entity.velocity[i] = speedMin
			}
		} else if entity.velocity[i] != zero {
			entity.velocity[i] *= friction
			if (entity.velocity[i] > zero && entity.velocity[i] < cutoff) || (entity.velocity[i] < zero && entity.velocity[i] > -cutoff) {
				entity.velocity[i] = zero
			}
		}
	}

	entity.position = entity.position.Add(entity.velocity)

	// fmt.Println(entity.accel)
	// fmt.Println(entity.velocity)
	// fmt.Println(entity.position)
	// fmt.Println("")
	return entity
}

func (sprite Sprite) getFrame(animEnum, frame int) image.Image {
	var anim Animation
	switch animEnum {
	case 0:
		anim = sprite.hAnim
	case 1:
		anim = sprite.vAnim
	default:
		anim = sprite.hAnim
	}
	_ = anim

	y0 := sprite.height * anim.fileoffset
	x0 := sprite.width * frame % anim.frames
	rect := image.Rect(x0, y0, x0+sprite.width, y0+sprite.height)
	return sprite.pixels.SubImage(rect)
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
