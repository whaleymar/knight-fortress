package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	FILE_SPRITE_LEVEL = "assets/newbarktown.png"
)

// eventually this will take a level config file which contains all the component data
func makeLevelEntity() Entity {
	vertices := screenVertices
	vao, vbo := makeVao(vertices)
	entity := Entity{
		1, // TODO hard coded id
		&ComponentList{},
		mgl32.Vec3{},
	}

	entity.components.add(&cDrawable{
		CMP_DRAWABLE,
		vertices,
		vao,
		vbo,
		makeLevelSprite(),
		makeStaticAnimationManager(),
	})

	return entity

}

func makeLevelSprite() Sprite {
	img, err := loadImage(FILE_SPRITE_LEVEL)
	if err != nil {
		panic(err)
	}
	return Sprite{
		img,
		img.Bounds().Max.Y,
		img.Bounds().Max.X,
		1, // TODO hard coded
	}

}
