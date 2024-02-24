package gfx

import (
	"fmt"
	"image"
	"image/draw"
	"os"
)

const (
	WindowTitle = "Gaming"

	WindowWidth    = 1280
	WindowHeight   = 720
	PixelsPerTexel = 4
	// WindowWidth    = 1920
	// WindowHeight   = 1080
	// PixelsPerTexel = 6

	// these convert world coordinates to screen coordinates, corrected for aspect ratio
	WORLD_SCALE_RATIO = float32(3.2)
	TEXEL_SCALE_X     = float32(1.0 / (16.0 / WORLD_SCALE_RATIO))
	TEXEL_SCALE_Y     = float32(1.0 / (9.0 / WORLD_SCALE_RATIO))
)

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
