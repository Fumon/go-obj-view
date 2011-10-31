package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"gl"
	"github.com/jteeuwen/glfw"
	"math"
	"unsafe"
)

const (
	Title = "Cube 01"
	Width = 800
	Height = 800
)

var (
	running bool
	poscolor_buffer gl.Buffer
	cube_ibo gl.Buffer
	program gl.Program
	attrib_loc gl.AttribLocation
	colattrib gl.AttribLocation
	transformattrib gl.UniformLocation

	vert_poscolor = []poscol{ //Three Position followed by three color
		//Front
		{Position: [3]float32{-1.0, -1.0, 1.0}, Color: [3]float32{1.0, 1.0, 0.0}},
		{Position: [3]float32{1.0, -1.0, 1.0}, Color: [3]float32{1.0, 0.0, 0.0}},
		{Position: [3]float32{1.0, 1.0, 1.0}, Color: [3]float32{0.0, 1.0, 0.0}},
		{Position: [3]float32{-1.0, 1.0, 1.0}, Color: [3]float32{0.0, 0.0, 1.0}},
		//Back
		{Position: [3]float32{-1.0, -1.0, -1.0}, Color: [3]float32{1.0, 1.0, 0.0}},
		{Position: [3]float32{1.0, -1.0, -1.0}, Color: [3]float32{1.0, 0.0, 0.0}},
		{Position: [3]float32{1.0, 1.0, -1.0}, Color: [3]float32{0.0, 1.0, 0.0}},
		{Position: [3]float32{-1.0, 1.0, -1.0}, Color: [3]float32{0.0, 0.0, 1.0}},
	}

	vert_index = []uint16{
		//front
		0, 1, 2,
		2, 3, 0,
		//top
		1, 5, 6,
		6, 2, 1,
		//back
		7, 6, 5,
		5, 4, 7,
		//bottom
		4, 0, 3,
		3, 7, 4,
		//left
		4, 5, 1,
		1, 0, 4,
		//right
		3, 2, 6,
		6, 7, 3,
	}
	sizeofposcol int
	offsettocolor int

	//Tick stuff
	//lasttime float64
)

type poscol struct {
	Position [3]float32
	Color [3]float32
}

func main() {
	var err os.Error

	//Init Types
	sizeofposcol = int(unsafe.Sizeof(poscol{}))
	offsettocolor = int(unsafe.Offsetof(poscol{}.Color))


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
	if err = glfw.OpenWindow(Width, Height, 8,8,8,8,0,8, glfw.Windowed);err != nil {
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

	//Default identity transform
	//id := IdMat4()
	//transformattrib.UniformMatrix4fv(1, false, id[:])

	running = true
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		calc_tick()
		draw()
		glfw.SwapBuffers()
		if glfw.Key('Q') == glfw.KeyPress {
			running = false
			break
		}
	}
}

func draw() {
	gl.Enable(gl.DEPTH_TEST)
	//Clear to white
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT|gl.DEPTH_BUFFER_BIT)

	//Use program
	program.Use()
	//Enable both attributes
	attrib_loc.EnableArray()
	defer attrib_loc.DisableArray()
	colattrib.EnableArray()
	defer colattrib.DisableArray()
	//Bind the buffer to the client state array buffer
	poscolor_buffer.Bind(gl.ARRAY_BUFFER)

	//Now define how to read from the buffer
	attrib_loc.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		sizeofposcol, //The 3float appears every 6 floats
		uintptr(0), //0 offset from beginning to start
	)
	colattrib.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		sizeofposcol,
		uintptr(offsettocolor), //Offset 3 floats from beginnning
	)

	//Index Buffer
	cube_ibo.Bind(gl.ELEMENT_ARRAY_BUFFER)
	//

	//Draw
	gl.DrawElementsInternal(gl.TRIANGLES, len(vert_index), gl.UNSIGNED_SHORT, uintptr(0))
}

func calc_tick() {
	//Set the uniform
	//move := float32(math.Sin((glfw.Time() * (math.Pi * 2.0)) / 5.0))
	//angle := float32(glfw.Time() * math.Pi/4.0) //45 degrees a second
	//axis := []float32{1.0, 0.0, 0.0}
	//translate := TranslateMat4([]float32{move, 0.0, 0.0})
	//rotation := AxisAngleRotation(axis, angle)
	//transform := translate.Product(rotation)
	//transformattrib.UniformMatrix4fv(1, false, transform[:])

	transform := TranslateMat4([]float32{0.0, 0.0, -2.0})
	transformattrib.UniformMatrix4fv(1, false, transform[:])
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
	poscolor_buffer.Delete()
}

func init_vbo() (err os.Error) {
	poscolor_buffer = gl.GenBuffer()
	poscolor_buffer.Bind(gl.ARRAY_BUFFER)
	gl.BufferDataCompound(gl.ARRAY_BUFFER, len(vert_poscolor) * sizeofposcol, vert_poscolor, gl.STATIC_DRAW)

	cube_ibo = gl.GenBuffer()
	cube_ibo.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferDataCompound(gl.ELEMENT_ARRAY_BUFFER, len(vert_index) * 2, vert_index, gl.STATIC_DRAW)
	return
}

func init_program() (err os.Error) {
	//Make verrtex shader
	var vs gl.Shader
	if vs, err = loadshader("cube.v.glsl", gl.VERTEX_SHADER); err != nil {
		return
	}
	var fs gl.Shader
	if fs, err = loadshader("cube.f.glsl", gl.FRAGMENT_SHADER); err != nil {
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
	if attrib_loc = program.GetAttribLocation("coord3d"); attrib_loc == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location %v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
	}
	if colattrib = program.GetAttribLocation("v_color"); colattrib == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find color attribute location &v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
	}
	if transformattrib = program.GetUniformLocation("m_transform"); transformattrib == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find transform attribute location %v\n", program.GetInfoLog())
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


//Matrix stuff

//A column major ordered matrix
type mat4 [16]float32

func MakeMat4(a float32) (m *mat4) {
	m = new(mat4)
	for i := range (*m) {
		(*m)[i] = a
	}
	return
}

func IdMat4() (m *mat4) {
	m = new(mat4)
	for i := range m {
		switch i {
		case 0, 5, 10, 15: m[i] = 1.0
		default: m[i] = 0.0
		}
	}
	return
}

func ScaleMat4(s float32) (m *mat4) {
	m = new(mat4)
	for i := range m {
		switch i {
		case 0, 5, 10, 15: m[i] = s
		default: m[i] = 0.0
		}
	}
	return
}

func (m mat4) String() string {
	var result string
	for row := 0; row < 4; row++ {
		result +=
		  fmt.Sprintln(m[row], m[row + 4], m[row + 8], m[row + 12])
	}
	return result
}

//Remember, column major and zero indexed
func (m *mat4) At(row, col int) float32 {
	return m[4*col + row]
}

//a * b in written order
func (a *mat4) Product(b *mat4) (mv *mat4) {
	mv = new(mat4)
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			var sum float32
			for i := 0; i < 4; i++ {
				sum += a.At(row, i) * b.At(i, col)
			}
			mv[col*4 + row] = sum
		}
	}

	return
}

func TranslateMat4(t []float32) (m *mat4) {
	m = IdMat4()
	m[12] = t[0]
	m[13] = t[1]
	m[14] = t[2]
	return
}

func AxisAngleRotation(axis []float32, angle float32) (mv *mat4) {
	mv = IdMat4()
	c := math.Cos(float64(angle))
	s := math.Sin(float64(angle))
	t := 1 - c
	x := float64(axis[0])
	y := float64(axis[1])
	z := float64(axis[2])
	//Exploded
	mv[0] = float32(t * x * x + c)
	mv[1] = float32(t * x * y + z * s)
	mv[2] = float32(t * x * z - y * s)

	mv[4] = float32(t * x * y - z * s)
	mv[5] = float32(t * y * y + c)
	mv[6] = float32(t * y * z + x * s)

	mv[8] = float32(t * x * z + y * s)
	mv[9] = float32(t * y * z - x * s)
	mv[10] = float32(t * z * z + c)

	return
}
