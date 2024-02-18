package main

import (
	_ "fmt"
	// "image/png"
	_ "image/png"
	"log"
	// "os"
	_ "os"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth = 1280
	// windowWidth  = 720
	windowHeight = 720
	windowTitle  = "Gaming"

	VERTEX_FILE   = "src/vertex.glsl"
	FRAGMENT_FILE = "src/fragment.glsl"

	STRIDE_SIZE = 5
	FLOAT_SIZE  = 4

	COLOR_CLEAR_R = 0.28627450980392155
	COLOR_CLEAR_G = 0.8705882352941177
	COLOR_CLEAR_B = 0.8509803921568627
	COLOR_CLEAR_A = 1.0
)

var ORIGIN = mgl32.Vec3{0.0, 0.0, 0.0}
var SIZE_STANDARD = mgl32.Vec2{1.0, 1.0}
var ZERO3 = mgl32.Vec3{0.0, 0.0, 0.0}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	program, err := initOpenGL()
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	_, _ = initCamera(program)

	offset := mgl32.Vec3{0.0, 0.0, 0.0}
	offsetUniform := gl.GetUniformLocation(program, gl.Str("offset\x00"))
	gl.Uniform3fv(offsetUniform, 1, &offset[0])

	// curVertices := screenVertices
	curVertices := squareVertices
	// curVertices := smallSquareVertices // TODO move the code that uses this (in main loop) to entity method
	// entity := makePlayerEntity()
	entity := getPlayerPtr()
	initControls(window, entity)

	_ = setShaderVars(program)

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	// texture, err := loadTexture("src/square.png")
	texture, err := entity.getTexture(0)
	if err != nil {
		panic(err)
	}

	// Global settings
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.ClearColor(COLOR_CLEAR_R, COLOR_CLEAR_G, COLOR_CLEAR_B, COLOR_CLEAR_A)

	previousTime := glfw.GetTime()
	elapsed := float32(0.0)

	millis := gl.GetUniformLocation(program, gl.Str("millis\x00"))
	gl.Uniform1f(millis, float32(previousTime))
	// xmod := float32(1.0)
	frame := 0
	for !window.ShouldClose() {
		frame++
		frame = frame % 60
		// fmt.Printf("frame: %d\n", frame)

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		deltaTime := time - previousTime
		elapsed += float32(deltaTime)
		previousTime = time

		gl.Uniform1f(millis, float32(time))
		entity.update(deltaTime)

		// Render
		gl.UseProgram(program)
		gl.Uniform3fv(offsetUniform, 1, &entity.position[0])

		gl.BindVertexArray(entity.vao)

		// bind `texture` to texture uniform at index 0
		texture, err = entity.getTexture(frame)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		nVertices := int32(len(curVertices) / STRIDE_SIZE)
		gl.DrawArrays(gl.TRIANGLES, 0, nVertices)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil) // idk what the last 2 args do
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}
