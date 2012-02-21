// Some wrappers for gl stuff
package main

// #cgo darwin LDFLAGS: -framework OpenGL -lGLEW -lGL
// #cgo windows LDFLAGS: -lglew32 -lopengl32
// #cgo linux LDFLAGS: -lGLEW -lGL
//
// #include <stdlib.h>
//
// #ifdef __APPLE__
// # include <OpenGL/glew.h>
// #else
// # include <GL/glew.h>
// #endif
//
// #undef GLEW_GET_FUN
// #define GLEW_GET_FUN(x) (*x)
import "C"
import gl "github.com/chsc/gogl/gl33"
import "unsafe"

type (
	Object  gl.Uint
	Program Object
	Shader  Object
)

type (
	AttribLocation gl.Uint
	UniformLocation gl.Uint
)

func glString(s string) *gl.Char { return (*gl.Char)(C.CString(s)) }

func CreateShader(type_ gl.Enum) Shader {
	return Shader(gl.CreateShader(type_))
}

func (shader Shader) Delete() { gl.DeleteShader(shader) }

func (shader Shader) GetInfoLog() string {
	var len gl.Int
	gl.GetShaderiv(gl.Uint(shader), gl.Enum(INFO_LOG_LENGTH), &len)

	log := C.malloc(C.size_t(len + 1))
	gl.GetShaderInfoLog(gl.Uint(shader), gl.Sizei(len), nil, (*gl.Char)(log))

	defer C.free(log)

	return C.GoString((*C.char)(log))
}

func (shader Shader) GetSource() string {
	var len gl.Int
	gl.GetShaderiv(gl.Uint(shader), gl.Enum(SHADER_SOURCE_LENGTH), &len)

	log := C.malloc(C.size_t(len + 1))
	gl.GetShaderSource(gl.Uint(shader), gl.Sizei(len), nil, (*gl.Char)(log))

	defer C.free(log)

	return C.GoString((*C.char)(log))
}

func (shader Shader) Source(source string) {

	csource := glString(source)
	defer freeString(csource)

	var one gl.Int = gl.Int(len(source))

	gl.ShaderSource(gl.Uint(shader), 1, &csource, &one)
}

func (shader Shader) Compile() { gl.CompileShader(gl.Uint(shader)) }

func (shader Shader) Get(param gl.Enum) int {
	var rv gl.Int

	gl.GetShaderiv(gl.Uint(shader), gl.Enum(param), &rv)
	return int(rv)
}

// Program

func CreateProgram() Program { return Program(gl.CreateProgram()) }

func (program Program) Delete() { gl.DeleteProgram(gl.Uint(program)) }

func (program Program) AttachShader(shader Shader) {
	gl.AttachShader(gl.Uint(program), gl.Uint(shader))
}

func (program Program) GetAttachedShaders() []Object {
	var len gl.Int
	gl.GetProgramiv(gl.Uint(program), gl.Enum(ACTIVE_UNIFORM_MAX_LENGTH), &len)

	objects := make([]Object, len)
	gl.GetAttachedShaders(gl.Uint(program), gl.Sizei(len), nil, *((**gl.Uint)(unsafe.Pointer(&objects))))
	return objects
}

func (program Program) DetachShader(shader Shader) {
	gl.DetachShader(gl.Uint(program), gl.Uint(shader))
}

func (program Program) TransformFeedbackVaryings(names []string, buffer_mode gl.Enum) {
	if len(names) == 0 {
		gl.TransformFeedbackVaryings(gl.Uint(program), 0, (**gl.Char)(nil), gl.Enum(buffer_mode))
	} else {
		gl_names := make([]*gl.Char, len(names))

		for i := range names {
			gl_names[i] = glString(names[i])
		}

		gl.TransformFeedbackVaryings(gl.Uint(program), gl.Sizei(len(gl_names)), &gl_names[0], gl.Enum(buffer_mode))

		for _, s := range gl_names {
			freeString(s)
		}
	}
}

func (program Program) Link() { gl.LinkProgram(gl.Uint(program)) }

func (program Program) Validate() { gl.ValidateProgram(gl.Uint(program)) }

func (program Program) Use() { gl.UseProgram(gl.Uint(program)) }

func (program Program) GetInfoLog() string {

	var len gl.Int
	gl.GetProgramiv(gl.Uint(program), gl.Enum(INFO_LOG_LENGTH), &len)

	log := C.malloc(C.size_t(len + 1))
	gl.GetProgramInfoLog(gl.Uint(program), gl.Sizei(len), nil, (*gl.Char)(log))

	defer C.free(log)

	return C.GoString((*C.char)(log))

}

func (program Program) Get(param gl.Enum) int {
	var rv gl.Int

	gl.GetProgramiv(gl.Uint(program), gl.Enum(param), &rv)
	return int(rv)
}

func (program Program) GetUniformiv(location UniformLocation, values []int) {
	// no range check
	gl.GetUniformiv(gl.Uint(program), gl.Int(location), (*gl.Int)(unsafe.Pointer(&(values[0]))))
}

func (program Program) GetUniformfv(location UniformLocation, values []float32) {
	// no range check
	gl.GetUniformfv(gl.Uint(program), gl.Int(location), (*gl.Float)(unsafe.Pointer(&(values[0]))))
}

func (program Program) GetUniformLocation(name string) UniformLocation {

	cname := glString(name)
	defer freeString(cname)

	return UniformLocation(gl.GetUniformLocation(gl.Uint(program), cname))
}

func (program Program) GetAttribLocation(name string) AttribLocation {

	cname := glString(name)
	defer freeString(cname)

	return AttribLocation(gl.GetAttribLocation(gl.Uint(program), cname))
}

func (program Program) BindAttribLocation(index AttribLocation, name string) {

	cname := glString(name)
	defer freeString(cname)

	gl.BindAttribLocation(gl.Uint(program), gl.Uint(index), cname)

} 

// Attribute location

func (indx AttribLocation) Attrib1f(x float32) {
	gl.VertexAttrib1f(gl.Uint(indx), gl.Float(x))
}

func (indx AttribLocation) Attrib1fv(values []float32) {
	//no range check
	gl.VertexAttrib1fv(gl.Uint(indx), (*gl.Float)(unsafe.Pointer(&values[0])))
}

func (indx AttribLocation) Attrib2f(x float32, y float32) {
	gl.VertexAttrib2f(gl.Uint(indx), gl.Float(x), gl.Float(y))
}

func (indx AttribLocation) Attrib2fv(values []float32) {
	//no range check
	gl.VertexAttrib2fv(gl.Uint(indx), (*gl.Float)(unsafe.Pointer(&values[0])))
}

func (indx AttribLocation) Attrib3f(x float32, y float32, z float32) {
	gl.VertexAttrib3f(gl.Uint(indx), gl.Float(x), gl.Float(y), gl.Float(z))
}

func (indx AttribLocation) Attrib3fv(values []float32) {
	//no range check
	gl.VertexAttrib3fv(gl.Uint(indx), (*gl.Float)(unsafe.Pointer(&values[0])))
}

func (indx AttribLocation) Attrib4f(x float32, y float32, z float32, w float32) {
	gl.VertexAttrib4f(gl.Uint(indx), gl.Float(x), gl.Float(y), gl.Float(z), gl.Float(w))
}

func (indx AttribLocation) Attrib4fv(values []float32) {
	//no range check
	gl.VertexAttrib4fv(gl.Uint(indx), (*gl.Float)(unsafe.Pointer(&values[0])))
}

func (indx AttribLocation) AttribPointer(size uint, normalized bool, stride int, pointer interface{}) {
	t, p := GetGLenumType(pointer)
	gl.VertexAttribPointer(gl.Uint(indx), gl.Int(size), gl.Enum(t), glBool(normalized), gl.Sizei(stride), p)
}

func (indx AttribLocation) EnableArray() {
	gl.EnableVertexAttribArray(gl.Uint(indx))
}

func (indx AttribLocation) DisableArray() {
	gl.DisableVertexAttribArray(gl.Uint(indx))
}


// Vertex Arrays
type VertexArray Object

func GenVertexArray () VertexArray {
	var a gl.Uint
	gl.GenVertexArrays(1, &a)
	return VertexArray(a)
}

func GenVertexArrays (arrays []VertexArray) {
	gl.GenVertexArrays(gl.Sizei(len(arrays)), (*gl.Uint)(&arrays[0]))
}

func (array VertexArray) Delete () {
	gl.DeleteVertexArrays(1, (*gl.Uint)(&array))
}

func DeleteVertexArrays (arrays []VertexArray) {
	gl.DeleteVertexArrays(gl.Sizei(len(arrays)), (*gl.Uint)(&arrays[0]))
}

func (array VertexArray) Bind () {
	gl.BindVertexArray(gl.Uint(array))
}

// UniformLocation
//TODO

func (location UniformLocation) Uniform1f(x float32) {
	gl.Uniform1f(gl.Int(location), gl.Float(x))
}

func (location UniformLocation) Uniform2f(x float32, y float32) {
	gl.Uniform2f(gl.Int(location), gl.Float(x), gl.Float(y))
}

func (location UniformLocation) Uniform3f(x float32, y float32, z float32) {
	gl.Uniform3f(gl.Int(location), gl.Float(x), gl.Float(y), gl.Float(z))
}

func (location UniformLocation) Uniform1fv(v []float32) {
	panic("unimplemented")
	//	gl.Uniform1fv(gl.Int(location), (*C.float)(&v[0]));
}

func (location UniformLocation) Uniform1i(x int) {
	gl.Uniform1i(gl.Int(location), gl.Int(x))
}

func (location UniformLocation) Uniform1iv(v []int) {
	panic("unimplemented")
	//	gl.Uniform1iv(gl.Int(location), (*C.int)(&v[0]));
}

func (location UniformLocation) Uniform2fv(v []float32) {
	panic("unimplemented")
	//	gl.Uniform2fv(gl.Int(location), (*C.float)(&v[0]));
}

func (location UniformLocation) Uniform2i(x int, y int) {
	gl.Uniform2i(gl.Int(location), gl.Int(x), gl.Int(y))
}

func (location UniformLocation) Uniform2iv(v []int32) {
	panic("unimplemented")
	//	gl.Uniform2iv(gl.Int(location), (*C.int)(&v[0]));
}


func (location UniformLocation) Uniform3fv(v []float32) {
	panic("unimplemented")
	//	gl.Uniform3fv(gl.Int(location), (*C.float)(&v[0]));
}

func (location UniformLocation) Uniform3i(x int, y int, z int) {
	gl.Uniform3i(gl.Int(location), gl.Int(x), gl.Int(y), gl.Int(z))
}

func (location UniformLocation) Uniform3iv(v []int32) {
	panic("unimplemented")
	//	gl.Uniform3iv(gl.Int(location), (*C.int)(&v[0]));
}

func (location UniformLocation) Uniform4f(x float32, y float32, z float32, w float32) {
	gl.Uniform4f(gl.Int(location), gl.Float(x), gl.Float(y), gl.Float(z), gl.Float(w))
}

func (location UniformLocation) Uniform4fv(v []float32) {
	panic("unimplemented")
	//	gl.Uniform4fv(gl.Int(location), (*C.float)(&v[0]));
}

func (location UniformLocation) Uniform4i(x int, y int, z int, w int) {
	gl.Uniform4i(gl.Int(location), gl.Int(x), gl.Int(y), gl.Int(z), gl.Int(w))
}

func (location UniformLocation) Uniform4iv(v []int32) {
	panic("unimplemented")
	//	gl.Uniform4iv(gl.Int(location), (*C.int)(&v[0]));
}


 //Buffer Objects

type Buffer Object

// Create single buffer object
func GenBuffer() Buffer {
	var b C.GLuint
	gl.GenBuffers(1, &b)
	return Buffer(b)
}

// Fill slice with new buffers
func GenBuffers(buffers []Buffer) {
	gl.GenBuffers(gl.Sizei(len(buffers)), (*gl.Uint)(&buffers[0]))
}

// Delete buffer object
func (buffer Buffer) Delete() {
	b := gl.Uint(buffer)
	gl.DeleteBuffers(1, &b)
}

// Delete all textures in slice
func DeleteBuffers(buffers []Buffer) {
	gl.DeleteBuffers(gl.Sizei(len(buffers)), (*gl.Uint)(&buffers[0]))
}

// Bind this buffer as target
func (buffer Buffer) Bind(target gl.Enum) {
	gl.BindBuffer(gl.Enum(target), gl.Uint(buffer))
}

// Bind this buffer as index of target
func (buffer Buffer) BindBufferBase(target gl.Enum, index uint) {
	gl.BindBufferBase(gl.Enum(target), gl.Uint(index), gl.Uint(buffer))
}

// Bind this buffer range as index of target
func (buffer Buffer) BindBufferRange(target gl.Enum, index uint, offset int, size uint) {
	gl.BindBufferRange(gl.Enum(target), gl.Uint(index), gl.Uint(buffer), gl.Intptr(offset), gl.Sizeiptr(size))
}

// Creates and initializes a buffer object's data store
func BufferData(target gl.Enum, size int, data interface{}, usage gl.Enum) {
	_, p := GetGLenumType(data)
	gl.BufferData(gl.Enum(target), gl.Sizeiptr(size), p, gl.Enum(usage))
}

//  Update a subset of a buffer object's data store
func BufferSubData(target gl.Enum, offset int, size int, data interface{}) {
	_, p := GetGLenumType(data)
	gl.BufferSubData(gl.Enum(target), gl.Intptr(offset), gl.Sizeiptr(size), p)
}

// Returns a subset of a buffer object's data store
func GetBufferSubData(target gl.Enum, offset int, size int, data interface{}) {
	_, p := GetGLenumType(data)
	gl.GetBufferSubData(gl.Enum(target), gl.Intptr(offset), gl.Sizeiptr(size), p)
}

//  Map a buffer object's data store
func MapBuffer(target gl.Enum, access gl.Enum) {
	gl.MapBuffer(gl.Enum(target), gl.Enum(access))
}

//  Unmap a buffer object's data store
func UnmapBuffer(target gl.Enum) bool {
	return goBool(gl.UnmapBuffer(gl.Enum(target)))
}

// Return buffer pointer
func glGetBufferPointerv(target gl.Enum, pname gl.Enum, params []unsafe.Pointer) {
	gl.GetBufferPointerv(gl.Enum(target), gl.Enum(pname), &params[0])
}

// Return parameters of a buffer object
func GetBufferParameteriv(target gl.Enum, pname gl.Enum, params []int32) {
	gl.GetBufferParameteriv(gl.Enum(target), gl.Enum(pname), (*gl.Int)(&params[0]))
}