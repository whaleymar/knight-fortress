package main

import (
	"fmt"
	"image"
	"image/draw"
	// _ "image/png" // fixes "image: unknown format" error
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// windowWidth  = 1280
	windowWidth  = 720
	windowHeight = 720
	windowTitle  = "Title"

	VERTEX_FILE   = "src/vertex.glsl"
	FRAGMENT_FILE = "src/fragment.glsl"

	STRIDE_SIZE = 3 // should be 5 for x,y,z,u,v coord system
	FLOAT_SIZE  = 4

	COLOR_CLEAR_R = 0.28627450980392155
	COLOR_CLEAR_G = 0.8705882352941177
	COLOR_CLEAR_B = 0.8509803921568627
	COLOR_CLEAR_A = 1.0
)

var ORIGIN = mgl32.Vec3{0.0, 0.0, 0.0}
var SIZE_STANDARD = mgl32.Vec2{1.0, 1.0}
var ZERO3 = mgl32.Vec3{0.0, 0.0, 0.0}

type ShaderConfig struct {
	vert             uint32
	vertTextureCoord uint32
}

type DrawableEntity struct {
	position mgl32.Vec3
	size     mgl32.Vec2
	vao      uint32
	velocity mgl32.Vec3
	accel    mgl32.Vec3
}

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

	model := mgl32.Translate3D(0, 0, -2.5) // becomes invisible when <=-7 (probably because of Far camera param)
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	offset := mgl32.Vec3{0.0, 0.0, 0.0}
	offsetUniform := gl.GetUniformLocation(program, gl.Str("offset\x00"))
	gl.Uniform3fv(offsetUniform, 1, &offset[0])

	// curVertices := screenVertices
	// curVertices := squareVertices
	curVertices := smallSquareVertices
	vao := makeVao(curVertices)
	entity := makeDrawableEntity(vao)
	initControls(window, &entity)

	_ = setShaderVars(program)

	// Global settings
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)
	gl.ClearColor(COLOR_CLEAR_R, COLOR_CLEAR_G, COLOR_CLEAR_B, COLOR_CLEAR_A)

	previousTime := glfw.GetTime()
	elapsed := float32(0.0)

	millis := gl.GetUniformLocation(program, gl.Str("millis\x00"))
	gl.Uniform1f(millis, float32(previousTime))
	// xmod := float32(1.0)
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		deltaTime := time - previousTime
		elapsed += float32(deltaTime)
		previousTime = time

		gl.Uniform1f(millis, float32(time))
		entity = entity.update()

		// model = mgl32.Translate3D(elapsed/100.0, elapsed/100.0, elapsed/100.0)
		// trans := mgl32.Translate3D(0.01, 0.01, 0.0)
		// model = model.Mul4(trans)
		// model = mgl32.HomogRotate3D(float32(elapsed), mgl32.Vec3{0.55, 0.55, 0.55}) // should be unit vector or objects will be sheared;

		// offset = offset.Add(mgl32.Vec3{float32(deltaTime), float32(deltaTime), 0.0}.Mul(xmod))
		// fmt.Println(offset)
		// entity.position.Add(offset)
		// if offset[0] > 2 || offset[0] < -2 {
		// 	xmod *= -1
		// }

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		// gl.Uniform3fv(offsetUniform, 1, &offset[0])
		gl.Uniform3fv(offsetUniform, 1, &entity.position[0])

		gl.BindVertexArray(entity.vao)

		// apply texture to vertices
		// gl.ActiveTexture(gl.TEXTURE0)
		// gl.BindTexture(gl.TEXTURE_2D, texture)

		nVertices := int32(len(curVertices) / 3)
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

func initControls(window *glfw.Window, playerPointer *DrawableEntity) {
	window.SetKeyCallback(playerControlsCallback)
	window.SetUserPointer(unsafe.Pointer(playerPointer))
}

func playerControlsCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// fmt.Println(key, scancode, action, mods)
	if action == glfw.Repeat {
		return
	}

	var accel float32
	if action == glfw.Release {
		accel = -0.05
	} else {
		accel = 0.05
	}

	playerPointer := window.GetUserPointer()
	player := (*DrawableEntity)(playerPointer)
	switch key {
	case glfw.KeyW:
		player.accel[1] += accel
	case glfw.KeyS:
		player.accel[1] -= accel
	case glfw.KeyA:
		player.accel[0] -= accel
	case glfw.KeyD:
		player.accel[0] += accel
	}
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

func makeDrawableEntity(vao uint32) DrawableEntity {
	entity := DrawableEntity{ORIGIN, SIZE_STANDARD, vao, ZERO3, ZERO3} // TODO make size based on vertices
	return entity
}

func (entity DrawableEntity) update() DrawableEntity {
	// this is stinky garbage TODO

	speedMax := float32(0.1)
	speedMin := float32(-0.1)
	zero := float32(0)
	cutoff := float32(0.005)
	friction := float32(0.5)

	// X
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
