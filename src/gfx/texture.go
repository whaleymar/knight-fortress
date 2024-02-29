package gfx

import (
	"fmt"
	"image"
	"sync"

	"github.com/go-gl/gl/v4.1-core/gl"
)

const (
	OPENGL_MAX_TEXTURES = 32 // might need compile flag since this varies
	SIZE_TEX_MAIN_X     = 512
	SIZE_TEX_MAIN_Y     = 512
	SIZE_TEX_MAIN_Z     = 1
	FILE_MAIN_SPRITES   = "assets/sprites.png"
)

type TextureSlot uint32

const (
	TEX_MAIN TextureSlot = iota
)

var _TEXTURE_MGR_LOCK = &sync.Mutex{} // i dont think opengl can run outside main thread. if it can, lock down methods
var textureMgrPtr *TextureManager

func InitMainTexture() {
	texMgr := GetTextureManager()
	texArray := makeTextureArray(
		[3]uint32{SIZE_TEX_MAIN_X, SIZE_TEX_MAIN_Y, SIZE_TEX_MAIN_Z},
	)

	img, err := loadImage(FILE_MAIN_SPRITES)
	if err != nil {
		panic(err)
	}
	texArray.addTexture(img)

	err = texMgr.register(texArray, TEX_MAIN)
	if err != nil {
		panic("Couldn't load main texture atlas")
	}
}

func GetTextureManager() *TextureManager {
	if textureMgrPtr == nil {
		_TEXTURE_MGR_LOCK.Lock()
		defer _TEXTURE_MGR_LOCK.Unlock()
		if textureMgrPtr == nil {
			textureMgrPtr = &TextureManager{}
		}
	}
	return textureMgrPtr
}

type TextureManager struct {
	textureHandles [OPENGL_MAX_TEXTURES]uint32
	textureArrays  [OPENGL_MAX_TEXTURES]TextureArray
	allocMask      [OPENGL_MAX_TEXTURES]bool
}

type TextureArray struct {
	sheets []*image.RGBA
	size   [3]uint32
	// tileSize [2]uint32
}

func (texMgr *TextureManager) register(texArray TextureArray, textureIx TextureSlot) error {
	if texMgr.allocMask[textureIx] {
		return fmt.Errorf("Texture index %d is in use and cannot be overwritten", textureIx)
	}

	var texture uint32
	gl.GenTextures(1, &texture)

	// TODO 2D texture array
	// bind textures before generating them
	gl.ActiveTexture(gl.TEXTURE0 + uint32(textureIx))
	gl.BindTexture(gl.TEXTURE_2D, texture)

	// get closest pixel color (for pixelated look, no interpolation)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	// repeat at end of texture
	// https://registry.khronos.org/OpenGL-Refpages/gl4/html/glTexParameter.xhtml
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

	rgba := texArray.sheets[0] // TODO array
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y), // n layers for texture array
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	texMgr.textureHandles[textureIx] = texture
	texMgr.textureArrays[textureIx] = texArray

	texMgr.allocMask[textureIx] = true

	return nil
}

func (texMgr *TextureManager) GetTextureHandle(ix TextureSlot) uint32 {
	return texMgr.textureHandles[ix]
}

func (texMgr *TextureManager) GetTextureSize(slotIx TextureSlot, arrayIx int) (float32, float32) {
	rect := texMgr.textureArrays[slotIx].sheets[arrayIx].Bounds()
	return float32(rect.Max.X), float32(rect.Max.Y)
}

func (texMgr *TextureManager) Free(textureIx TextureSlot) {
	handle := texMgr.GetTextureHandle(textureIx)
	gl.DeleteTextures(1, &handle)
	texMgr.allocMask[textureIx] = false
}

func makeTextureArray(arraySize [3]uint32) TextureArray {
	return TextureArray{
		nil,
		arraySize,
	}
}

func (texArray *TextureArray) addTexture(img *image.RGBA) error {
	var x, y, z int = int(texArray.size[0]), int(texArray.size[1]), int(texArray.size[2])
	if len(texArray.sheets) == z {
		return fmt.Errorf("Texture array at max capacity")
	}
	if img.Bounds().Max.X != x || img.Bounds().Max.Y != y {
		return fmt.Errorf("Image size (%d, %d) did not match TextureArray size of (%d, %d)", img.Bounds().Max.X, img.Bounds().Max.Y, x, y)
	}
	texArray.sheets = append(texArray.sheets, img)
	return nil
}
