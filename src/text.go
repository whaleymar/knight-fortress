package main

import (
	// "os"
	// "path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"

	// "golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

const (
	FONTPATH = "assets/font/font.otf"
)

func loadFont() (*opentype.Font, error) {

	// absPath, _ := filepath.Abs(FONTPATH)
	// b, err := os.ReadFile(absPath)
	// if err != nil {
	// 	return nil, err
	// }
	// return opentype.Parse(b)
	return opentype.Parse(goitalic.TTF)
}

func loadFontFace(font *opentype.Font) (font.Face, error) {
	face, err := opentype.NewFace(font, nil)
	return face, err
}

func pixelCoords(x, y int) fixed.Point26_6 {
	return fixed.P(x, y)
}
