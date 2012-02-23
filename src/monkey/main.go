package main

import (
	"errors"
	"flag"
	"fmt"
	
	gl "github.com/chsc/gogl/gl33"
	"github.com/jteeuwen/glfw"
	. "glmatrix"
	"io/ioutil"
	"math"
	"obj"
	"os"
	"unsafe"
)

const (
	Title = "Monkey!!!!!"
)

var (
	Width   = 800
	Height  = 600
	running 		bool
	normalson		bool
	//ibo gl.Buffer //index buffer
	vbo 			Buffer //vertex buffer
	nbo 			Buffer //Normals 
	//nbo gl.Buffer //normals buffer
	//uvbo gl.Buffer //UVs buffer

	program           	Program
	
	attrib_obj_coord  	AttribLocation
	attrib_obj_normal 	AttribLocation
	transformattrib   	UniformLocation
	it3x3attrib       	UniformLocation
	inv_view         		UniformLocation

	lineprogram     	Program
	line_obj_coord 	AttribLocation
	line_transform 	UniformLocation

	monkeymodel       	*filemodel
	monkeynorms       	[]obj.GeomVertex

	//Tick stuff
	model          		*Mat4
	model_pre_tick 	*Mat4
	view           		*Mat4
	projection     		*Mat4
	mvp_tilt       		*Mat4

	//Custom Datatype Variables
	vertnormSize int
	geomvertexsize int
	offsetNorm   int

	filename = flag.String("file", "", "Sets the model to render")
	spinrate = flag.Float64("spin", 4.0, "Sets the spin rate as Pi/x radians per second")
)

//Type to store verticies and normals
type vertnorm struct {
	Vert obj.GeomVertex
	Norm obj.VertexNormal
}

type filemodel struct {
	Geometry    []vertnorm
	VertexCount uint
}

func (f *filemodel) String() string {
	s := fmt.Sprintf("Filemodel with %v Verticies\n", f.VertexCount)
	for i, g := range f.Geometry {
		s = fmt.Sprintf("%vV %v: %v\n", s, i, g)
	}
	return s
}

func resize_event(width, height int) {
	Width = width
	Height = height
	calculate_projection()
	gl.Viewport(0, 0, gl.Sizei(Width), gl.Sizei(Height))
}

func main() {
	flag.Parse()
	var err error


	if *filename == "" {
		flag.PrintDefaults()
		return
	}

	//Init types
	vertnormSize = int(unsafe.Sizeof(vertnorm{}))
	geomvertexsize = int(unsafe.Sizeof(obj.GeomVertex{}))
	offsetNorm = int(unsafe.Offsetof(vertnorm{}.Norm))

	//Init glfw
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in glfw init %v\n", err)
		return
	}
	defer glfw.Terminate()

	//Set a resize handler
	glfw.SetWindowSizeCallback(resize_event)

	//Open window
	if err = glfw.OpenWindow(Width, Height, 8, 8, 8, 8, 32, 23, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "Error in openwindow: %v\n", err)
		return
	}
	defer glfw.CloseWindow()

	glfw.SetSwapInterval(1) //Vsync
	glfw.SetWindowTitle(Title)

	//Init glew
	if errGL := gl.Init(); errGL != nil {
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
	nkeypress := false
	for running && glfw.WindowParam(glfw.Opened) == 1 {
		calc_tick()
		draw()
		glfw.SwapBuffers()
		if glfw.Key('Q') == glfw.KeyPress {
			running = false
			break
		}
		if nkeypress == false && glfw.Key('N') == glfw.KeyPress{
			nkeypress = true
		} else  if nkeypress == true && glfw.Key('N') == glfw.KeyRelease{
			nkeypress = false
			normalson = !normalson
		}
	}
}

func drawmonkey() {
	//Use program
	program.Use()
	//Enable both attributes
	attrib_obj_coord.EnableArray()
	attrib_obj_normal.EnableArray()
	//Bind the buffer to the client state array buffer
	vbo.Bind(gl.ARRAY_BUFFER)

	//Now define how to read from the buffer
	attrib_obj_coord.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		vertnormSize,
		uintptr(0), //0 offset from beginning to start
	)

	attrib_obj_normal.AttribPointerInternal(
		3,
		gl.FLOAT,
		false,
		vertnormSize,
		uintptr(offsetNorm),
	)

	//Index Buffer
	//ibo.Bind(gl.ELEMENT_ARRAY_BUFFER)

	//Draw
	//gl.DrawElementsInternal(gl.TRIANGLES, len(monkeyobj.Geometry.FaceIndicies), gl.UNSIGNED_INT, uintptr(0))

	gl.DrawArrays(gl.TRIANGLES, 0, gl.Sizei(monkeymodel.VertexCount))
}

func drawnormals() {
	lineprogram.Use()
	line_obj_coord.EnableArray()

	nbo.Bind(gl.ARRAY_BUFFER)
	line_obj_coord.AttribPointerInternal(3, gl.FLOAT, false, 0, uintptr(0))
	gl.DrawArrays(gl.LINES, 0, gl.Sizei(len(monkeynorms)))

	line_obj_coord.DisableArray()	
}

func draw() {
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	if(normalson) {drawnormals() }
	drawmonkey()
	//Normals
	
}

func calc_tick() {
	angle := float32(glfw.Time() * math.Pi / *spinrate)
	axis := []float32{0.0, 1.0, 0.0}
	rotation := AxisAngleRotation(axis, angle)

	mtmp := model.Product(rotation.Product(model_pre_tick))
	mt := mtmp.Upper3by3()
	it3x3Model := mt.Inverse().Transpose()

	mvp_tilt = projection.Product(view.Product(mtmp))
	program.Use()
	transformattrib.UniformMatrix4fv(1, false, mvp_tilt[:])
	it3x3attrib.UniformMatrix3fv(1, false, it3x3Model[:])
	inv_view.UniformMatrix4fv(1, false, view.Inverse()[:])
	lineprogram.Use()
	line_transform.UniformMatrix4fv(1, false, mvp_tilt[:])
}

func calculate_projection() {
	projection = StdProjection(float32(math.Pi/4), float32(0.1), float32(10.0), (float32(Width) / float32(Height)))
}

func load_model() (err error) {
	var monkeyobj *obj.Object
	//Open and load the monkey
	file, nerr := ioutil.ReadFile(*filename)
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
		unique_vert_count = uint(len(monkeyobj.FaceIndicies))
	} else {
		unique_vert_count = uint(len(monkeyobj.FaceIndicies))
	}

	//Now make a model
	m := &filemodel{make([]vertnorm, unique_vert_count), unique_vert_count}
	//Load in geometry
	for i, f := range monkeyobj.FaceIndicies {
		m.Geometry[i].Vert = monkeyobj.Verticies[f[0]]
		m.Geometry[i].Norm = monkeyobj.Normals[f[2]]
	}

	monkeymodel = m

	monkeynorms = make([]obj.GeomVertex, unique_vert_count * 2)
	g := 0
	for _, v := range monkeymodel.Geometry {
		monkeynorms[g] = v.Vert
		monkeynorms[g+1] = obj.GeomVertex{v.Vert[0] + v.Norm[0] * 0.2, v.Vert[1] + v.Norm[1] * 0.2, v.Vert[2] + v.Norm[2] * 0.2}
		g += 2
	}

	return
}

func init_resources() (err error) {
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

func init_vbo() (err error) {
	vbo = GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	BufferDataCompound(gl.ARRAY_BUFFER, int(monkeymodel.VertexCount)*vertnormSize, monkeymodel.Geometry, gl.STATIC_DRAW)


	nbo = GenBuffer()
	nbo.Bind(gl.ARRAY_BUFFER)
	BufferDataCompound(gl.ARRAY_BUFFER, int(len(monkeynorms)) * geomvertexsize, monkeynorms, gl.STATIC_DRAW)

	return
}

func init_program() (err error) {
	//Make verrtex shader
	var vs Shader
	if vs, err = loadshader("cube.v.glsl", gl.VERTEX_SHADER); err != nil {
		return
	}
	var fs Shader
	if fs, err = loadshader("cube.f.glsl", gl.FRAGMENT_SHADER); err != nil {
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
	if errInt := program.Get(gl.LINK_STATUS); errInt == 0 {
		fmt.Fprintf(os.Stderr, "Failed to link: %v\n", program)
		program.Delete()
		err = errors.New("Failed to link")
		return
	}

	//Find attribute location
	if attrib_obj_coord = program.GetAttribLocation("obj_coord"); attrib_obj_coord == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"obj_coord\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = errors.New("Attribute not located")
		return
	}
	if attrib_obj_normal = program.GetAttribLocation("obj_normal"); attrib_obj_normal == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"obj_normal\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = errors.New("Attribute not located")
		return
	}
	if transformattrib = program.GetUniformLocation("m_transform"); transformattrib == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"m_transform\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = errors.New("Attribute not located")
		return
	}
	if it3x3attrib = program.GetUniformLocation("m_3x3itModel"); it3x3attrib == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"m_3x3itModel\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = errors.New("Attribute not located")
		return
	}
	if inv_view = program.GetUniformLocation("m_inv_view"); inv_view == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"m_inv_view\"\n\t%v\n", program.GetInfoLog())
		program.Delete()
		err = errors.New("Attribute not located")
		return
	}

	//Line program
	//Make verrtex shader
	var linevs Shader
	if linevs, err = loadshader("line.v.glsl", gl.VERTEX_SHADER); err != nil {
		return
	}
	var linefs Shader
	if linefs, err = loadshader("line.f.glsl", gl.FRAGMENT_SHADER); err != nil {
		return
	}
	lineprogram = CreateProgram()
	lineprogram.AttachShader(linevs)
	lineprogram.AttachShader(linefs)
	lineprogram.Link()

	if errInt := lineprogram.Get(gl.LINK_STATUS); errInt == 0 {
		fmt.Fprintf(os.Stderr, "Failed to link: %v\n", lineprogram)
		lineprogram.Delete()
		err = errors.New("Failed to link")
		return
	}

	if line_obj_coord = lineprogram.GetAttribLocation("obj_coord"); line_obj_coord == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"obj_coord\"\n\t%v\n", lineprogram.GetInfoLog())
		lineprogram.Delete()
		err = errors.New("Attribute not located")
		return
	}
	if line_transform = lineprogram.GetUniformLocation("m_transform"); line_transform == -1 {
		fmt.Fprintf(os.Stderr, "Failed to find attribute location \"m_transform\"\n\t%v\n", lineprogram.GetInfoLog())
		lineprogram.Delete()
		err = errors.New("Attribute not located")
		return
	}

	return
}

func loadshader(filename string, shType gl.Enum) (sh Shader, err error) {
	sourcetext, nerr := ioutil.ReadFile(filename)
	err = nerr
	if err != nil {
		//Error reading file
		return
	}

	sh = CreateShader(shType)
	//Link source to shader
	sh.Source(string(sourcetext))
	//Compile
	sh.Compile()

	if errint := sh.Get(gl.COMPILE_STATUS); errint == 0 {
		defer sh.Delete()
		fmt.Fprintf(os.Stderr, "Error in compiling %s\n\t%v\n",
			filename, sh.GetInfoLog())
		err = errors.New("Compile error")
		return
	}
	return
}
