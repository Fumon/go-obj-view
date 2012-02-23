/*
	An importer for OBJ objects. Originally created to allow rendering in OpenGL from go.

	TODO:
	Error reporting with line numbers
	Multiple object importing (Will change interface)
	Quad importing
*/
package obj

import (
	"errors"
	"fmt" //      "io/ioutil"


	"strconv"
	"strings"
)

type Geometry struct {
	Verticies    []GeomVertex
	Normals      []VertexNormal
	FaceIndicies []FaceTriplet
	IsQuads      bool
}

type Object struct {
	Name string
	Geometry
	UVs         []TextureVertex2D
	Initialized bool
	HasNormals  bool
	HasUVs      bool
}

func (o *Object) String() string {
	tmp := fmt.Sprintf("Object of %v verticies, %v normals, %v face triplets\nVerticies:\n", len(o.Verticies), len(o.Normals), len(o.FaceIndicies))
	for i, d := range o.Verticies {
		tmp = fmt.Sprintf("%v%v %v: %v\n", tmp, "V", i, d)
	}
	for i, d := range o.Normals {
		tmp = fmt.Sprintf("%v%v %v: %v\n", tmp, "N", i, d)
	}
	for i, d := range o.FaceIndicies {
		tmp = fmt.Sprintf("%v%v %v: %v\n", tmp, "F", i, d)
	}
	return tmp
}

func initObject() *Object {
	return &Object{"", Geometry{make([]GeomVertex, 0), make([]VertexNormal, 0), make([]FaceTriplet, 0), false}, make([]TextureVertex2D, 0), false, false, false}
}

func Parse(s string) (obj *Object, err error) {
	err = nil

	//Split into lines
	lines := strings.Split(s, "\n")
	//Create an Object to put data into.
	obj = initObject()
	vertsStarted := false
	facesStarted := false

	for _, e := range lines {
		//Split up into components
		f := strings.Fields(e)
		if len(f) < 2 {
			continue
		}
		//Take actions
		switch f[0] {
		case "o":
			if vertsStarted == true {
				//We already have an object definition started, we need a new one now
				//TODO
			} else {
				obj.Name = f[1]
				obj.Initialized = true
			}
		case "v":
			if !vertsStarted {
				vertsStarted = true
			}
			//Convert
			vert := *new(GeomVertex)
			for i, b := range f[1:] {
				var tmpv float64
				tmpv, err = strconv.ParseFloat(b, 32)
				if err != nil {
					return
				}
				v := float32(tmpv)
				vert[i] = v
			}
			obj.Verticies = append(obj.Verticies, vert)
		case "vn":
			//Convert
			norm := *new(VertexNormal)
			for i, b := range f[1:] {
				var tmpv float64
				tmpv, err = strconv.ParseFloat(b, 32)
				if err != nil {
					return
				}
				v := float32(tmpv)
				norm[i] = v
			}
			obj.Normals = append(obj.Normals, norm)
		case "f":
			//Split
			if !facesStarted {
				//Determine if quads or tris
				if len(f) == 5 {
					obj.IsQuads = true
				} else {
					obj.IsQuads = false
				}
				//Determine the extent of this object
				slashSplit := strings.Split(f[1], "/")
				switch len(slashSplit) {
				case 0:
					err = errors.New("Error parsing face")
					return
				case 1:
					//Only using verticies, no changes to obj
				case 2:
					//Verticies and uvs, no normals
					obj.HasUVs = true
				case 3:
					//Have to check.
					if slashSplit[1] != "" {
						obj.HasUVs = true
					}
					if slashSplit[2] != "" && slashSplit[2] != "\n" {
						obj.HasNormals = true
					}
				}
				facesStarted = true
			}

			face := *new(FaceTriplet)
			for _, inds := range f[1:] {
				//TODO: MAJOR Better managment for indicies for multiple objects
				slashSplit := strings.Split(inds, "/")
				for i, b := range slashSplit {
					if b == "" || b == "\n" {
						continue //TODO: Error checking for missing values
					}
					conv, err := strconv.ParseUint(b, 10, 0)
					if err != nil {
						err = errors.New("Error converting uint")
					}
					face[i] = uint32(conv - 1)
				}
				obj.FaceIndicies = append(obj.FaceIndicies, face)
			}
		}
	}

	return
}
