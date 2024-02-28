package gfx

import (
	"fmt"
	"image"
	// "image/draw"
	"image/png"

	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	STRIDE_SIZE = 5
	FLOAT_SIZE  = 4

	VERTEX_FILE   = "shader/vertex.glsl"
	FRAGMENT_FILE = "shader/fragment.glsl"
)

type ShaderConfig struct {
	vert             uint32
	vertTextureCoord uint32
}

type VAO struct {
	handle uint32
}

type VBO struct {
	handle uint32
}

func InitOpenGL() (uint32, error) {

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

func InitCamera(program uint32) mgl32.Mat4 {
	projection := mgl32.Ortho2D(0, float32(WindowWidth), 0, float32(WindowHeight))
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	return projection
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

func MakeVao() VAO {
	var vao VAO
	gl.GenVertexArrays(1, &vao.handle)
	return vao
}

func (vao VAO) Bind() {
	gl.BindVertexArray(vao.handle)
}

func MakeVbo() VBO {
	var vbo VBO
	gl.GenBuffers(1, &vbo.handle)
	return vbo
}

func (vbo VBO) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handle)
}

func (vbo VBO) Buffer(data []float32) {
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*FLOAT_SIZE, gl.Ptr(data), gl.STATIC_DRAW)
}

func UpdateShaderVars(program uint32) ShaderConfig {
	// this only affects the *current* vao bound to glArrayBuffer

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// vec3 vertices
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, STRIDE_SIZE*FLOAT_SIZE, 0)

	// vec2 texture position
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, false, STRIDE_SIZE*FLOAT_SIZE, 3*FLOAT_SIZE)

	return ShaderConfig{vertAttrib, texCoordAttrib}
}

func loadTexture(rgba *image.RGBA, textureIx uint32) (uint32, error) {

	var texture uint32
	gl.GenTextures(1, &texture)

	// need to bind textures before generating them
	gl.ActiveTexture(gl.TEXTURE0 + textureIx)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
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

func saveImage(img image.Image, name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		log.Printf("failed to encode image: %v", err)
	}
}

func MakeRectVertices(pixelWidth, pixelHeight float32) []float32 {

	return []float32{
		0.0, pixelHeight, 0.0, 0.0, 0.0,
		pixelWidth, pixelHeight, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 1.0,
		pixelWidth, pixelHeight, 0.0, 1.0, 0.0,
		pixelWidth, 0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 1.0,
	}
}

func MakeRectVerticesWithUV(pixelWidth, pixelHeight, xMin, xMax, yMin, yMax float32) []float32 {

	return []float32{
		0.0, pixelHeight, 0.0, xMin, yMin,
		pixelWidth, pixelHeight, 0.0, xMax, yMin,
		0.0, 0.0, 0.0, xMin, yMax,
		pixelWidth, pixelHeight, 0.0, xMax, yMin,
		pixelWidth, 0.0, 0.0, xMax, yMax,
		0.0, 0.0, 0.0, xMin, yMax,
	}
}

var ScreenVertices = MakeRectVertices(WindowWidth, WindowHeight)

var LevelVertices = MakeRectVertices(WindowWidth*16, WindowHeight*16)

var SquareVertices []float32 = MakeRectVertices(64, 64) // placeholder for cDrawable init
