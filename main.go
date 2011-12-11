package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"gl"
	"github.com/jteeuwen/glfw"
	//"unsafe"
	. "./matrix/_obj/glmatrix"
	"math"
	"../obj_import/_obj/obj"
)

const (
	Title = "Monkey!!!!!"
)

var (
	Width = 800
	Height = 800
	running bool
	//ibo gl.Buffer //index buffer
	vbo gl.Buffer //vertex buffer
	//nbo gl.Buffer //normals buffer
	//uvbo gl.Buffer //UVs buffer

	program gl.Program
	attrib_obj_coord gl.AttribLocation
	attrib_obj_normal gl.AttribLocation
	transformattrib gl.UniformLocation
	monkeymodel *filemodel

	//Tick stuff
	model *Mat4
	model_pre_tick *Mat4
	view *Mat4
	projection *Mat4
	mvp_tilt *Mat4
)

//Type to store verticies and normals
type vertnorm struct {
	Vert obj.GeomVertex
	Norm obj.VertexNormal
}

type filemodel struct {
	Geometry []vertnorm
	VertexCount uint
}

func (f *filemodel) String() string {
	s := fmt.Sprintf("Filemodel with %v Verticies\n", f.VertexCount)
	for i, g := range f.Geometry {
		s = fmt.Sprintf("%vV %v: %v\n", s, i, g)
	}
	return s
}

func resize_event(width,height int) {
	Width = width
	Height = height
	calculate_projection()
	gl.Viewport(0, 0, Width, Height)
}

func main() {
	var err os.Error

	//Init glfw
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in glfw init %v\n", err)
		return
	}
	defer glfw.Terminate()

	//Set a resize handler
	glfw.SetWindowSizeCallback(resize_event)

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
	attrib_obj_coord.EnableArray()
	defer attrib_obj_coord.DisableArray()
	//Bind the buffer to the client state array buffer
	vbo.Bind(gl.ARRAY_BUFFER)

	//Now define how to read from the buffer
	attrib_obj_coord.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		4*3*2,
		uintptr(0), //0 offset from beginning to start
	)

	attrib_obj_normal.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		4*3*2,
		uintptr(4*3),
	)

	//Index Buffer
	//ibo.Bind(gl.ELEMENT_ARRAY_BUFFER)

	//Draw
	//gl.DrawElementsInternal(gl.TRIANGLES, len(monkeyobj.Geometry.FaceIndicies), gl.UNSIGNED_INT, uintptr(0))

	gl.DrawArrays(gl.TRIANGLES, 0, int(monkeymodel.VertexCount))
}

func calc_tick() {
	angle := float32(glfw.Time() * math.Pi/4.0) //45 degrees a second
	axis := []float32{0.0, 1.0, 0.0}
	rotation := AxisAngleRotation(axis, angle)

	mvp_tilt = projection.Product(view.Product(model.Product(rotation.Product(model_pre_tick))))
	transformattrib.UniformMatrix4fv(1, false, mvp_tilt[:])
}

func calculate_projection() {
	projection = StdProjection( float32(math.Pi/4), float32(0.1), float32(10.0), (float32(Width) / float32(Height) ) )
}

func load_model() (err os.Error) {
	var monkeyobj *obj.Object
	//Open and load the monkey
	file, nerr := ioutil.ReadFile("monkey.obj")
	err = nerr
	if err != nil {
		return
	}
	monkeyobj, err = obj.Parse(string(file))
	if err != nil {
		return
	}

	var unique_vert_count uint
	if monkeyobj.IsQuads {
		unique_vert_count = uint(len(monkeyobj.FaceIndicies) * 4)
	} else {
		unique_vert_count = uint(len(monkeyobj.FaceIndicies) * 3)
	}

	//Now make a model
	m := &filemodel{make([]vertnorm, unique_vert_count), unique_vert_count}
	//Load in geometry
	for i, f := range monkeyobj.FaceIndicies {
		m.Geometry[i].Vert = monkeyobj.Verticies[f[0]]
		m.Geometry[i].Norm = monkeyobj.Normals[f[2]]
	}

	monkeymodel = m

	return
}

func init_resources() (err os.Error) {
	//Calculate the cute transform
	model_pre_tick = AxisAngleRotation([]float32{1.0, 0.0, 0.0}, float32(math.Pi/2.0))
	model = TranslateMat4([]float32{0.0, 0.0, -4.0})
	view = ViewLookAt([]float32{0.0, 2.0, 0.0}, []float32{0.0, -2.0, -4.0}, []float32{0.0, 1.0, 0.0})
	calculate_projection()
	mvp_tilt = projection.Product(view.Product(model))

	if err = load_model(); err != nil {
		return
	}

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
}

func init_vbo() (err os.Error) {
	vbo = gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	gl.BufferDataCompound(gl.ARRAY_BUFFER, int(monkeymodel.VertexCount * 3 * 2 * 4), monkeymodel.Geometry, gl.STATIC_DRAW)

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
	if attrib_obj_coord = program.GetAttribLocation("obj_coord"); attrib_obj_coord == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"obj_coord\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
		return
	}
	if attrib_obj_normal = program.GetAttribLocation("obj_normal"); attrib_obj_normal == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"obj_normal\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
		return
	}
	if transformattrib = program.GetUniformLocation("m_transform"); transformattrib == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find transform attribute location \"m_transform\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = os.NewError("Attribute not located")
		return
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
		defer sh.Delete()
		fmt.Fprintf(os.Stderr, "Error in compiling %s\n\t%v\n",
			filename, sh.GetInfoLog())
		err = os.NewError("Compile error")
		return
	}
	return
}
