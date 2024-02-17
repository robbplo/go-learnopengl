package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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

	window, err := glfw.CreateWindow(800, 600, "Your Game Title", nil, nil)
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

	vertices := []float32{
		0.0, 1.0, 0.0,
		1.0, -1.0, 0.0,
		-1.0, -1.0, 0.0,
	}

	// identify the vertex buffer
	var vertexbuffer uint32
	// generate 1 buffer, put the resulting identifier in vertexbuffer
	gl.GenBuffers(1, &vertexbuffer)
	// the following commands will talk about our 'vertexbuffer' buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	// give our vertices to OpenGL
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Load shaders
	programId, err := loadShaders("shaders/SimpleVertexShader.glsl", "shaders/SimpleFragmentShader.glsl")
	if err != nil {
		panic(err)
	}

	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	glfw.GetCurrentContext().SetInputMode(glfw.StickyKeysMode, glfw.True)
	for !shouldClose(window) {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(programId)

		// 1st attribute buffer : vertices
		gl.EnableVertexAttribArray(0)
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.DrawArrays(gl.TRIANGLES, 0, 3) // Starting from vertex 0; 3 vertices total -> 1 triangle
		gl.DisableVertexAttribArray(0)

		// Update window
		window.SwapBuffers()

		// Poll for window events.
		glfw.PollEvents()
	}
}

func shouldClose(window *glfw.Window) bool {
	return window.GetKey(glfw.KeyEscape) == glfw.Press || window.ShouldClose()
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
	// Read shader code from file
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
