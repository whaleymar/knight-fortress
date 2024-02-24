package main

import (
	"cmp"
	"fmt"
	"slices"

	_ "image/png"
	"log"

	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/whaleymar/knight-fortress/src/ec"
	"github.com/whaleymar/knight-fortress/src/gfx"
	"github.com/whaleymar/knight-fortress/src/sys"
)

// unused import error is probably the stupidest thing I've ever seen
var _ = fmt.Println
var _ = cmp.Compare(1, 1)
var _ = slices.Min([]int{1})

const (
	COLOR_CLEAR_R = 0.12
	COLOR_CLEAR_G = 0.13
	COLOR_CLEAR_B = 0.15
	COLOR_CLEAR_A = 1.0
)

var ORIGIN = mgl32.Vec3{0.0, 0.0, 0.0}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	window := initGlfw()
	defer glfw.Terminate()

	program, err := gfx.InitOpenGL()
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	_ = gfx.InitCamera(program)

	// offsets
	tmpOffset := mgl32.Vec3{}
	drawOffsetUniform := gl.GetUniformLocation(program, gl.Str("offset\x00"))
	gl.Uniform3fv(drawOffsetUniform, 1, &tmpOffset[0])

	InitControls(window)
	gfx.InitMainTexture()

	// init entities
	entityManager := ec.GetEntityManager()

	entityManager.Add(ec.GetPlayerPtr())
	entityManager.Add(ec.GetCameraPtr())
	entity := ec.MakeLevelEntity()
	entityManager.Add(&entity)

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
		sys.DeltaTime.Update()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		fpsCh <- 1 / sys.DeltaTime.Get()

		gl.Uniform1f(millis, float32(glfw.GetTime()))

		for _, entity := range entityManager.GetEntitiesWithComponent(ec.CMP_ANY) {
			entity.GetComponentManager().Update(entity)
		}

		// Render
		gl.UseProgram(program) // I don't know why I'm running this every frame but I'm afraid to change it

		texture = gfx.GetTextureManager().GetTextureHandle(gfx.TEX_MAIN)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		// have to sort by depth so things get blended correctly
		drawableEntities := entityManager.GetEntitiesWithComponent(ec.CMP_DRAWABLE)
		slices.SortFunc(drawableEntities, func(e1, e2 *ec.Entity) int {
			return cmp.Compare(e1.GetPosition().Z(), e2.GetPosition().Z())
		})
		for _, entity := range drawableEntities {
			// if entity.name == "Player" {
			// 	fmt.Println(entity.getPosition())
			// }
			screenCoords := ec.GetScreenCoordinates(entity.GetPosition())
			gl.Uniform3fv(drawOffsetUniform, 1, &screenCoords[0])
			drawComponent := *ec.GetComponentUnsafe[*ec.CDrawable](ec.CMP_DRAWABLE, entity)

			nVertices := int32(len(drawComponent.GetVertices()) / gfx.STRIDE_SIZE)
			drawComponent.GetVao().Bind()
			drawComponent.GetVbo().Bind()
			drawComponent.GetVbo().Buffer(drawComponent.GetVertices())

			gfx.UpdateShaderVars(program)
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

	window, err := glfw.CreateWindow(gfx.WindowWidth, gfx.WindowHeight, gfx.WindowTitle, nil, nil) // idk what the last 2 args do
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	return window
}

func InitControls(window *glfw.Window) {
	window.SetKeyCallback(ec.PlayerControlsCallback)
}

func updateFPS(fpsCh <-chan float32) {
	for fps := range fpsCh {
		// Move cursor up and to the beginning of the line
		// fmt.Print("\033[F\033[K")
		// fmt.Printf("FPS: %f\n", fps)
		_ = fps
	}
}
