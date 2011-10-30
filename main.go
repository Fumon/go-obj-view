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
)

var (
	running bool
	triangle_buffer gl.Buffer
	color_buffer gl.Buffer
	program gl.Program
	attrib_loc gl.AttribLocation
	colattrib gl.AttribLocation
	triangle_verts = []float32 {
		0.0, 0.8,
		0.8, -0.8,
		-0.8, -0.8,
	}
	vert_colors = []float32 {
		1.0, 0.0, 0.0, //Top Vert Color
		0.0, 1.0, 0.0, //Right Vert Color
		0.0, 0.0, 1.0, //Left Vert Color
	}
)


func main() {
	var err os.Error

	//Init glfw
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in glfw init %v\n", err)
		return
	}
	defer glfw.Terminate()

	//Set some hints for opening the window
	//No Resize
	glfw.OpenWindowHint(glfw.WindowNoResize, 1)

	//Open window
	if err = glfw.OpenWindow(640, 480, 8,8,8,8,0,8, glfw.Windowed);err != nil {
		fmt.Fprintf(os.Stderr, "Error in openwindow: %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1) //Vsync
	glfw.SetWindowTitle(Title)

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
		draw()
		glfw.SwapBuffers()
		if glfw.Key('Q') == glfw.KeyPress {
			running = false
			break
		}
	}
}

func draw() {
	//Clear to white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	//Use program
	program.Use()
	//Bind the attribute to the array register
	attrib_loc.EnableArray()
	defer attrib_loc.DisableArray()
	//Bind the vertex buffer to the client state array buffer
	triangle_buffer.Bind(gl.ARRAY_BUFFER)

	//Size of offset is 0
	offset := uintptr(0)

	//Set up data type of buffer
	attrib_loc.AttribPointerInternal(
		2, //Cardinality of each datum
		gl.FLOAT, //Type
		false, //Do not norm the data
		0, //No stride
		offset, //No offset
	)

	//Bind the colors
	colattrib.EnableArray()
	defer colattrib.DisableArray()
	color_buffer.Bind(gl.ARRAY_BUFFER)
	colattrib.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		0,
		offset,
	)

	//Draw from array.
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
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
	triangle_buffer = gl.GenBuffer()
	//Bind the buffer to the array buffer register
	triangle_buffer.Bind(gl.ARRAY_BUFFER)
	//Inform openGL that the current buffer should read data from the vertex array.
	gl.BufferData(gl.ARRAY_BUFFER, len(triangle_verts) * 4, triangle_verts, gl.STATIC_DRAW)

	//Create the color buffer
	color_buffer = gl.GenBuffer()
	color_buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vert_colors) * 4, vert_colors, gl.STATIC_DRAW)
	return
}

func init_program() (err os.Error) {
	//Make verrtex shader
	var vs gl.Shader
	if vs, err = loadshader("triangle.v.glsl", gl.VERTEX_SHADER); err != nil {
		return
	}
	var fs gl.Shader
	if fs, err = loadshader("triangle.f.glsl", gl.FRAGMENT_SHADER); err != nil {
		return
	}

	//Init Program
	program = gl.CreateProgram()
	//Attach shaders to program before linking
	program.AttachShader(vs)
	program.AttachShader(fs)
	//Link program
	program.Link()
	//Check
	if errInt := program.Get(gl.LINK_STATUS); errInt == 0 {
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
	if colattrib = program.GetAttribLocation("v_color"); colattrib == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find color attribute location &v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
	}

	return
}

func loadshader(filename string, shType gl.GLenum) (sh gl.Shader, err os.Error) {
	sourcetext, nerr := ioutil.ReadFile(filename)
	err = nerr
	if err != nil {
		//Error reading file
		return
	}

	sh = gl.CreateShader(shType)
	//Link source to shader
	sh.Source(string(sourcetext))
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
