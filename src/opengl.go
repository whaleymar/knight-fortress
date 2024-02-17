package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png" // fixes "image: unknown format" error
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"golang.org/x/image/font"
)

type ShaderConfig struct {
	vert             uint32
	vertTextureCoord uint32
}

func initOpenGL() (uint32, error) {

	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShaderSource, fragmentShaderSource := loadShaders()
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1)) // i have no idea what this does
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func initCamera(program uint32) (mgl32.Mat4, mgl32.Mat4) {
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	return projection, camera
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func loadShaders() (string, string) {
	absPath, _ := filepath.Abs(VERTEX_FILE)
	b, err := os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	vertexShader := string(b) + "\x00"

	absPath, _ = filepath.Abs(FRAGMENT_FILE)
	b, err = os.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	fragmentShader := string(b) + "\x00"

	return vertexShader, fragmentShader
}

func makeVao(points []float32) uint32 {
	// Make a Vertex Array Object (vao)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// allocate Vertex Buffer Object (vbo) and pass pointer to vertex data
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.BufferData(gl.ARRAY_BUFFER, len(points)*FLOAT_SIZE, gl.Ptr(points), gl.STATIC_DRAW)

	return vao
}

func setShaderVars(program uint32) ShaderConfig {
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// vec3 vertices
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, true, 0, nil)

	// vec2 texture position
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, true, STRIDE_SIZE*FLOAT_SIZE, 0)

	return ShaderConfig{vertAttrib, texCoordAttrib}
}

func loadTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

// func loadFontTexture(dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, face font.Face) (uint32, error) {
func loadFontTexture() (uint32, error) {

	const (
		width  = 72
		height = 36
		posX   = 0
		posY   = 0
	)

	gameFont, err := loadFont()
	if err != nil {
		return 0, err
	}
	face, err := loadFontFace(gameFont)
	if err != nil {
		return 0, err
	}

	rgba := image.NewGray(image.Rect(0, 0, width, height))

	// rgba := image.NewRGBA(dr.Bounds())
	// draw.DrawMask(rgba, dr, mask, maskp, nil, maskp, draw.Over)
	// fmt.Println(dr.Bounds(), mask.Bounds(), maskp, advance)

	drawer := font.Drawer{
		Dst:  rgba,
		Src:  image.White,
		Face: face,
		Dot:  pixelCoords(posX, posY),
	}

	drawer.DrawString("Whaley")

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil

}

// TODO
var screenVertices = []float32{
	-4.5, -1.0, 0.0,
	1.0, -1.0, 0.0,
	-1.0, 1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
}

var squareVertices = []float32{
	-1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	-1.0, 1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
}

var smallSquareVertices = []float32{
	-0.5, -0.5, 0.0,
	0.5, -0.5, 0.0,
	-0.5, 0.5, 0.0,
	0.5, -0.5, 0.0,
	0.5, 0.5, 0.0,
	-0.5, 0.5, 0.0,
}

var triangleVertices = []float32{
	0, 0.5, 0, // top
	-0.5, -0.5, 0, // left
	0.5, -0.5, 0, // right}
}
