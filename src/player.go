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
	ANIM_HFRAMES_PLAYER  = 2
	ANIM_VFRAMES_PLAYER  = 4
)

type DrawableEntity struct {
	position  mgl32.Vec3
	size      mgl32.Vec2
	vao       uint32
	velocity  mgl32.Vec3
	accel     mgl32.Vec3
	sprite    Sprite
	frame     int
	frameTime float64
	animSpeed float64
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
		0,   // frame
		0.0, // frametime
		4.0, // anim frames per second
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

func (entity *DrawableEntity) update(deltaTime float64) { // todo deltatime should prob be singleton
	// this is stinky garbage TODO
	// magic numbers TODO
	// deltatime TODO

	speedMax := 0.01
	velocityMax := float32(speedMax)
	velocityMin := float32(-speedMax)
	zero := float32(0)
	cutoff := float32(0.0005)
	friction := float32(0.5)

	for i := 0; i < 2; i++ {
		if entity.accel[i] != zero {
			entity.velocity[i] += entity.accel[i]
			if entity.velocity[i] > velocityMax {
				entity.velocity[i] = velocityMax
			} else if entity.velocity[i] < velocityMin {
				entity.velocity[i] = velocityMin
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

	entity.frameTime += deltaTime
	if entity.frameTime >= 1/entity.animSpeed {
		entity.frame = (entity.frame + 1) % entity.sprite.vAnim.frames // TODO anim getter based on ix
		entity.frameTime = 0.0
	}
}

func (entity *DrawableEntity) getTexture(frame int) (uint32, error) {
	return loadTextureFromMemory(entity.getSprite())

}

func (entity *DrawableEntity) getSprite() image.Image {
	// return entity.sprite.getFrame(entity.animIx, entity.frame)
	return entity.sprite.getFrame(1, entity.frame) // TODO animIx
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
	x0 := sprite.width * (frame % anim.frames) // dont think i need the modulo TODO
	// fmt.Println(x0, y0)
	rect := image.Rect(x0, y0, x0+sprite.width, y0+sprite.height)
	// fmt.Println(rect)
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
