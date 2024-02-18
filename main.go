package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// GLFW event handling must run on the main thread.
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		log.Fatalln("Failed to initialize GLFW:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	windowWidth := 800
	windowHeight := 600
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Your Game Title", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// set up VAO (vertex array object)
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	cube_vertices := []float32{
		-1.0, -1.0, -1.0, // triangle 1 : begin
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0, // triangle 1 : end
		1.0, 1.0, -1.0, // triangle 2 : begin
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0, // triangle 2 : end
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
	}

	// identify the vertex buffer
	var vertexbuffer uint32
	// generate 1 buffer, put the resulting identifier in vertexbuffer
	gl.GenBuffers(1, &vertexbuffer)
	// the following commands will talk about our 'vertexbuffer' buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	// give our vertices to OpenGL
	gl.BufferData(gl.ARRAY_BUFFER, len(cube_vertices)*4, gl.Ptr(cube_vertices), gl.STATIC_DRAW)

	cube_colors := []float32{
		0.583, 0.771, 0.014,
		0.609, 0.115, 0.436,
		0.327, 0.483, 0.844,
		0.822, 0.569, 0.201,
		0.435, 0.602, 0.223,
		0.310, 0.747, 0.185,
		0.597, 0.770, 0.761,
		0.559, 0.436, 0.730,
		0.359, 0.583, 0.152,
		0.483, 0.596, 0.789,
		0.559, 0.861, 0.639,
		0.195, 0.548, 0.859,
		0.014, 0.184, 0.576,
		0.771, 0.328, 0.970,
		0.406, 0.615, 0.116,
		0.676, 0.977, 0.133,
		0.971, 0.572, 0.833,
		0.140, 0.616, 0.489,
		0.997, 0.513, 0.064,
		0.945, 0.719, 0.592,
		0.543, 0.021, 0.978,
		0.279, 0.317, 0.505,
		0.167, 0.620, 0.077,
		0.347, 0.857, 0.137,
		0.055, 0.953, 0.042,
		0.714, 0.505, 0.345,
		0.783, 0.290, 0.734,
		0.722, 0.645, 0.174,
		0.302, 0.455, 0.848,
		0.225, 0.587, 0.040,
		0.517, 0.713, 0.338,
		0.053, 0.959, 0.120,
		0.393, 0.621, 0.362,
		0.673, 0.211, 0.457,
		0.820, 0.883, 0.371,
		0.982, 0.099, 0.879,
	}
	// color buffer
	var colorbuffer uint32
	gl.GenBuffers(1, &colorbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colorbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, len(cube_colors)*4, gl.Ptr(cube_colors), gl.STATIC_DRAW)

	// Load shaders
	programId, err := loadShaders("shaders/shader_vert.glsl", "shaders/shader_frag.glsl")
	if err != nil {
		panic(err)
	}

	/// calculate the MVP matrix
	// Projection matrix
	projection := mgl32.Perspective(mgl32.DegToRad(80.0), float32(windowWidth)/float32(windowHeight), 0.1, 100.0)

	// Camera matrix
	view := mgl32.LookAtV(mgl32.Vec3{4, 3, -3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	// Model matrix
	scale := mgl32.Scale3D(2.0, 2.0, 2.0)
	rotate := mgl32.HomogRotate3D(mgl32.DegToRad(0.0), mgl32.Vec3{0, 0, 1})
	translate := mgl32.Translate3D(0.0, 0.0, 0.0)
	model := translate.Mul4(rotate).Mul4(scale)

	gl.ClearColor(0.0, 0.0, 0.4, 0.0)
	// prevent rendering the back side of the cube
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	glfw.GetCurrentContext().SetInputMode(glfw.StickyKeysMode, glfw.True)
	// Main loop
	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(programId)

		gl.UniformMatrix4fv(gl.GetUniformLocation(programId, gl.Str("model\x00")), 1, false, &model[0])
    gl.UniformMatrix4fv(gl.GetUniformLocation(programId, gl.Str("view\x00")), 1, false, &view[0])
    gl.UniformMatrix4fv(gl.GetUniformLocation(programId, gl.Str("projection\x00")), 1, false, &projection[0])
		gl.Uniform1f(gl.GetUniformLocation(programId, gl.Str("angle\x00")), float32(glfw.GetTime()))

		// 1st attribute buffer : vertices
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
		gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false, 0, 0)

		// 2nd attribute buffer : colors
		gl.EnableVertexAttribArray(1)
		gl.BindBuffer(gl.ARRAY_BUFFER, colorbuffer)
		gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 0, 0)

		// do the thing!!
		gl.DrawArrays(gl.TRIANGLES, 0, 12*3)
		gl.DisableVertexAttribArray(0)
		gl.DisableVertexAttribArray(1)

		// Update window
		window.SwapBuffers()

		// Poll for window events.
		glfw.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
	}
}

func loadShaders(vertexFilePath string, fragmentFilePath string) (uint32, error) {
	vertexShaderId, err := compileShader(vertexFilePath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fragmentShaderId, err := compileShader(fragmentFilePath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	fmt.Println("Linking shaders")
	programId := gl.CreateProgram()
	gl.AttachShader(programId, vertexShaderId)
	gl.AttachShader(programId, fragmentShaderId)
	gl.LinkProgram(programId)

	var status int32
	gl.GetProgramiv(programId, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programId, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link shaders: %v", log)
	}
	gl.DetachShader(programId, vertexShaderId)
	gl.DetachShader(programId, fragmentShaderId)
	gl.DeleteShader(vertexShaderId)
	gl.DeleteShader(fragmentShaderId)

	return programId, nil
}

func compileShader(path string, shaderType uint32) (uint32, error) {
	fmt.Println("Compiling shader:", path)
	shader := gl.CreateShader(shaderType)
	shaderCode, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	source, free := gl.Strs(string(shaderCode))
	gl.ShaderSource(shader, 1, source, nil)
	defer free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile %v: %v", path, log)
	}

	return shader, nil
}
