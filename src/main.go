package main

import (
	"fmt"
	// "time"

	// "image/png"
	_ "image/png"
	"log"

	// "os"
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

// var SIZE_STANDARD = mgl32.Vec2{1.0, 1.0}
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

	// player position uniform
	offset := mgl32.Vec3{0.0, 0.0, 0.0}
	offsetUniform := gl.GetUniformLocation(program, gl.Str("offset\x00"))
	gl.Uniform3fv(offsetUniform, 1, &offset[0])

	initControls(window)

	_ = setShaderVars(program)

	var texture uint32
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	// init player and entity manager
	entityManager := getEntityManager()
	entityManager.add(*getPlayerPtr())
	// drawComponent, err := getComponent[*cDrawable](CMP_DRAWABLE, player)
	// if err != nil {
	// 	panic("Player is not drawable")
	// }
	// curVertices := (*drawComponent).vertices
	// texture, err := (*drawComponent).getTexture()
	// if err != nil {
	// 	panic(err)
	// }

	// OpenGL settings
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(COLOR_CLEAR_R, COLOR_CLEAR_G, COLOR_CLEAR_B, COLOR_CLEAR_A)

	// Time
	millis := gl.GetUniformLocation(program, gl.Str("millis\x00"))
	gl.Uniform1f(millis, float32(glfw.GetTime()))
	fpsCh := make(chan float32)
	go updateFPS(fpsCh)

	for !window.ShouldClose() {
		DeltaTime.update()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		fpsCh <- 1 / DeltaTime.get()

		gl.Uniform1f(millis, float32(glfw.GetTime()))
		for _, entity := range entityManager.getEntitiesWithComponent(CMP_ANY) {
			// fmt.Println("updating: ", i)
			entity.components.update(entity)
		}

		// Render
		gl.UseProgram(program) // I don't know why I'm running this every frame but I'm afraid to change it

		for _, entity := range entityManager.getEntitiesWithComponent(CMP_DRAWABLE) {
			// fmt.Println("drawing: ", i)
			gl.Uniform3fv(offsetUniform, 1, &entity.position[0])

			// bind `texture` to texture uniform at index 0
			prevTexture := texture // this might fail on frame 0
			var drawComponent *cDrawable
			if tmp, err := getComponent[*cDrawable](CMP_DRAWABLE, entity); err == nil {
				// CURRENT BUG:
				// this seems like an openGL issue
				// drawComponent state seems fine

				// fmt.Println("rendering")
				drawComponent = *tmp
				// fmt.Println(drawComponent.animManager.frame)
				// saveImage(drawComponent.getFrame(), fmt.Sprintf("tmp/test%f", glfw.GetTime()))
				texture, err = drawComponent.getTexture()
				if err != nil {
					texture = prevTexture
				}
				gl.ActiveTexture(gl.TEXTURE0)
				gl.BindTexture(gl.TEXTURE_2D, texture)

				nVertices := int32(len(drawComponent.vertices) / STRIDE_SIZE)
				gl.DrawArrays(gl.TRIANGLES, 0, nVertices)
			} else {
				fmt.Println(err)
			}

		}

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

func updateFPS(fpsCh <-chan float32) {
	for fps := range fpsCh {
		// Move cursor up and to the beginning of the line
		// fmt.Print("\033[F\033[K")
		// fmt.Printf("FPS: %f\n", fps)
		_ = fps
	}
}
