package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"gl"
	"github.com/jteeuwen/glfw"
)

const (
	Title = "Triangle 01"
	triangle_verts = []float32 {
		0.0, 0.8,
		0.8, -0.8,
		-0.8, -0.8
	}
)

var (
	running bool
	triangle_buffer gl.Buffer
	program gl.Program
	attrib_loc gl.AttribLocation
)


func main() {
	var err os.Error

	//Init glfw
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in glfw init %v\n", err)
		return
	}
	defer glfw.Terminate()

	//Open window
	if err = glfw.OpenWindow(640, 480, 8,8,8,8,0,8, gl.Windowed);err != nil {
		fmt.Fprintf(os.Stderr, "Error in openwindow: %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	//Init glew
	if errGL := gl.Init(); errGL != 0 {
		fmt.Fprintf(os.Stderr, "Error in glew init\n")
		return
	}

	//Initialize application resources
	if err = init_resources(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in resources: %v\n", err)
		return
	}
	defer cleanup_resources()

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		
	}
}

func init_resources() (err os.Error) {
	if err = init_vbo(); err != nil {
		return
	}
	if err = init_program(); err != nil {
		return
	}
	return
}

func cleanup_resources() {
	program.Delete()
	triangle_buffer.Delete()
}

func init_vbo() (err os.Error) {
	//Create the buffer
	triangle_buffer = GenBuffer()
	//Bind the buffer to the array buffer register
	triangle_buffer.Bind(gl.ARRAY_BUFFER)
	//Inform openGL that the current buffer should read data from the vertex array.
	BufferData(gl.ARRAY_BUFFER, len(triangle_verts) * 4, triangle_verts, gl.STATIC_DRAW)
	return
}

func init_program() (err os.Error) {
	//Make verrtex shader
	var vs gl.Shader
	if vs, err = loadshader("triangle.v.glsl", gl.VERTEX_SHADER); err != nil {
		return
	}
	var fs gl.Shader
	if fs, err = loadvsader("triangle.f.glsl", gl.FRAGMENT_SHADER); err != nil {
		return
	}

	//Init Program
	program = CreateProgram()
	//Attach shaders to program before linking
	program.AttachShader(vs)
	program.AttachShader(fs)
	//Link program
	program.Link()
	//Check
	if errInt := program.Get(gl.LINK_STATUS); errInt != 0 {
		fmt.Fprintf(os.Stderr, "Failed to link: %v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Failed to link")
		return
	}

	//Find attribute location
	if attrib_loc = program.GetAttribLocation("coord2d"); attrib_loc == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location %v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
	}

	return
}

func loadshader(filename string, shType gl.GLenum)
		(sh gl.Shader, err os.Error) {
	var sourcetext string
	if sourcetext, err = ioutil.ReadFile(filename); err != nil {
		//Error reading file
		return
	}

	sh = gl.CreateShader(shType)
	//Link source to shader
	sh.Source(sourcetext)
	//Compile
	sh.Compile()

	if errint := sh.Get(gl.COMPILE_STATUS); errint == 0 {
		sh.Delete()
		fmt.Fprintf(os.Stderr, "Error in compiling %s\n%v",
			filename, sh.GetInfoLog())
		err = os.NewError("Compile error")
		return
	}
	return
}
