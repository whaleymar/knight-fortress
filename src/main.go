package main

import (
	"fmt"
	// "go/build"
	"image"
	"image/draw"
	_ "image/png" // fixes "image: unknown format" error
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 1280
	windowHeight = 720
	windowTitle  = "Title"

	VERTEX_FILE   = "src/vertex.glsl"
	FRAGMENT_FILE = "src/fragment.glsl"

	STRIDE_SIZE = 3 // should be 5 for x,y,z,u,v coord system
	FLOAT_SIZE  = 4
)

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

	gl.UseProgram(program) // i also call in loop

	// set up vertex shader stuff
	// projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	// arg1 seems to affect both Y angle *and* distance from cube
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// positive numbers in center causes cube to move {left, down, right}
	camera := mgl32.LookAtV(mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}) // camera position, center, "up" (seems like rotation)
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// model matrix transforms our vertex coordinates from Model Space (centered about object center) to World Space (centered about global 0,0,0 coordinates)
	// model := mgl32.Ident4()

	// positive numbers causes cube to move {right, up, towards cube}
	model := mgl32.Translate3D(0, 0, -5) // becomes invisible when <=-7
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	// idk what this does but it's texture specific
	// textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	// gl.Uniform1i(textureUniform, 0)

	// TODO what does this do? Maybe names the frag shader output?
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	// texture, err := newTexture("src/square.png")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// Configure the Vertex Array Object (vao)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// allocate Vertex Buffer Object (vbo) and pass pointer to vertex data
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	curVertices := squareVertices
	gl.BufferData(gl.ARRAY_BUFFER, len(curVertices)*FLOAT_SIZE, gl.Ptr(curVertices), gl.STATIC_DRAW)

	// experimenting with circle
	// circleVerts := mgl32.Circle(1.0, 1.0, 100)
	// circleVertsFlat := make([]float32, 2*len(circleVerts))
	// for i := 0; i < len(circleVerts); i++ {
	// 	circleVertsFlat = append(circleVertsFlat, circleVerts[i][0], circleVerts[i][1], 0.0)
	// 	// fmt.Println(circleVertsFlat)
	// 	fmt.Println(circleVerts[i])
	// }
	// fmt.Println(len(circleVertsFlat))
	// fmt.Println(circleVerts[1])
	// fmt.Println(circleVertsFlat[3], circleVertsFlat[4], circleVertsFlat[5])
	// curVertices := circleVertsFlat
	// gl.BufferData(gl.ARRAY_BUFFER, len(squareVertices)*FLOAT_SIZE, gl.Ptr(curVertices), gl.STATIC_DRAW)
	// // gl.BufferData(gl.ARRAY_BUFFER, len(circleVerts)*8, gl.Ptr(circleVerts), gl.STATIC_DRAW)
	// gl.BufferData(gl.ARRAY_BUFFER, len(circleVertsFlat)*4, gl.Ptr(circleVertsFlat), gl.STATIC_DRAW)

	// set "vert" in vertex shader
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	// gl.VertexAttribPointerWithOffset(vertAttrib, 3, gl.FLOAT, false, STRIDE_SIZE*FLOAT_SIZE, 0)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, true, 0, nil)

	// set position vector in vertex shader
	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	// gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 0, nil)
	gl.VertexAttribPointerWithOffset(texCoordAttrib, 2, gl.FLOAT, true, STRIDE_SIZE*FLOAT_SIZE, 0)

	// Configure global settings
	// gl.Enable(gl.DEPTH_TEST)
	// gl.DepthFunc(gl.LESS)
	// gl.ClearColor(1.0, 1.0, 1.0, 1.0) set color of transparent pixel

	angle := 0.0
	previousTime := glfw.GetTime()

	// set millis uniform
	millis := gl.GetUniformLocation(program, gl.Str("millis\x00"))
	// gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	gl.Uniform1f(millis, float32(previousTime))

	for !window.ShouldClose() {
		// if commented, screen is black and nothing renders
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		gl.Uniform1f(millis, float32(previousTime))
		angle += elapsed
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0}) // should be unit vector or objects will be sheared
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 0, 1}) // should be unit vector or objects will be sheared;
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0.55, 0.55, 0.55}) // should be unit vector or objects will be sheared;

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

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

func initOpenGL() (uint32, error) {

	// Initialize Glow
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

// var squareVertices = []float32{
// 	-1.0, -1.0, 0.0, 1.0, 0.0,
// 	1.0, -1.0, 0.0, 0.0, 0.0,
// 	-1.0, 1.0, 0.0, 1.0, 1.0,
// 	1.0, -1.0, 0.0, 0.0, 0.0,
// 	1.0, 1.0, 0.0, 0.0, 1.0,
// 	-1.0, 1.0, 0.0, 1.0, 1.0,
// }

var squareVertices = []float32{
	-1.0, -1.0, 0.0,
	1.0, -1.0, 0.0,
	-1.0, 1.0, 0.0,
	1.0, -1.0, 0.0,
	1.0, 1.0, 0.0,
	-1.0, 1.0, 0.0,
}

var triangleVertices = []float32{
	0, 0.5, 0, // top
	-0.5, -0.5, 0, // left
	0.5, -0.5, 0, // right}
}

// idk why but i need these UV dims or it looks weird
// var triangleVertices = []float32{
// 	0, 0.5, 0, 1.0, 0.0, // top
// 	-0.5, -0.5, 0, 0.0, 0.0, // left
// 	0.5, -0.5, 0, 1.0, 1.0, // right}
// }

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// UV is for UV mapping https://en.wikipedia.org/wiki/UV_mapping
	// a face is 2 triangles, or 6 vertices
	// 6 faces -> 36 vertices total
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}
