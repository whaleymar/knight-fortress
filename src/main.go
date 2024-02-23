package main

import (
	"cmp"
	"fmt"
	"slices"

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

// unused import error is probably the stupidest thing I've ever seen
var _ = fmt.Println
var _ = cmp.Compare(1, 1)
var _ = slices.Min([]int{1})

const (
	windowWidth = 1280
	// windowWidth  = 720
	windowHeight = 720
	windowTitle  = "Gaming"

	VERTEX_FILE   = "shader/vertex.glsl"
	FRAGMENT_FILE = "shader/fragment.glsl"

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
	initMainTexture()

	// init player and entity manager
	entityManager := getEntityManager()

	entityManager.add(*getPlayerPtr())
	entityManager.add(makeLevelEntity())

	var texture uint32
	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	// OpenGL settings
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(COLOR_CLEAR_R, COLOR_CLEAR_G, COLOR_CLEAR_B, COLOR_CLEAR_A)
	// gl.ClearColor(0.0, 0.0, 0.0, 0.0)

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
			entity.components.update(entity)
		}

		// Render
		gl.UseProgram(program) // I don't know why I'm running this every frame but I'm afraid to change it

		texture = getTextureManager().getTextureHandle(TEX_MAIN)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// have to sort by depth so things get blended correctly
		drawableEntities := entityManager.getEntitiesWithComponent(CMP_DRAWABLE)
		slices.SortFunc(drawableEntities, func(e1, e2 *Entity) int {
			return cmp.Compare(e1.getPosition().Z(), e2.getPosition().Z())
		})
		for _, entity := range drawableEntities {
			gl.Uniform3fv(offsetUniform, 1, &entity.position[0])
			drawComponent := *getComponentUnsafe[*cDrawable](CMP_DRAWABLE, entity)

			nVertices := int32(len(drawComponent.vertices) / STRIDE_SIZE)
			drawComponent.vao.bind()
			drawComponent.vbo.bind()
			drawComponent.vbo.buffer(drawComponent.vertices)

			updateShaderVars(program)
			gl.DrawArrays(gl.TRIANGLES, 0, nVertices)
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
