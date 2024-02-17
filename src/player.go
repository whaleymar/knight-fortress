package main

import "github.com/go-gl/mathgl/mgl32"

type DrawableEntity struct {
	position mgl32.Vec3
	size     mgl32.Vec2
	vao      uint32
	velocity mgl32.Vec3
	accel    mgl32.Vec3
}

func makeDrawableEntity(vao uint32) DrawableEntity {
	entity := DrawableEntity{ORIGIN, SIZE_STANDARD, vao, ZERO3, ZERO3} // TODO make size based on vertices
	return entity
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
