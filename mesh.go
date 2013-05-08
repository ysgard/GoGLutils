// mesh - a struct representing a mesh
//
//
//

package goglutils

import (
	//"encoding/xml"
	"fmt"
	"errors"
	gl "github.com/chsc/gogl/gl33"
	"os"
	"strings"
	"strconv"
)

type Mesh struct {
	name string
	//attributes []*MeshAttribute
	attributes map[string]*MeshAttribute
	//indices    []*MeshIndex
	indices map[string]*MeshIndex
}

type MeshAttribute struct {
	name   string
	data   []gl.Float
	stride int
}

type MeshIndex struct {
	name      string
	data      []gl.Uint
	primitive gl.Enum
}

func NewMeshIndex(name string, data []gl.Uint, primitive gl.Enum) *MeshIndex {
	i := new(MeshIndex)
	copied := copy(i.data, data)
	if copied != len(data) {
		return nil
	}
	i.primitive = primitive
	return i
}

func NewMeshAttribute(data []gl.Float, stride int) *MeshAttribute {
	a := new(MeshAttribute)
	copied := copy(a.data, data)
	if copied != len(data) {
		return nil
	}
	a.stride = stride
	return a
}

func NewMesh(name string) *Mesh {
	m := new(Mesh)
	m.name = name
	m.attributes = map[string]*MeshAttribute{}
	m.indices = map[string]*MeshIndex{}
	return m
}

// Add an attribute array to a mesh
func (m *Mesh) AddMeshAttribute(name string, data []gl.Float, stride int) error {
	ma := NewMeshAttribute(data, stride)
	if ma == nil {
		return errors.New("Mesh:AddMeshAttribute:Could not allocate new attribute array")
	}
	m.attributes[name] = ma
	return nil
}

// Add an index to the mesh
func (m *Mesh) AddMeshIndex(name string, data []gl.Uint, primitive gl.Enum) error {
	mi := NewMeshIndex(name, data, primitive)
	if mi == nil {
		return errors.New("Mesh:AddMeshIndex:Could not allocate new index array")
	}
	m.indices[name] = mi
	return nil
}

// Helper function, splits a string of floats - like 
// "23.3 0.0 2323.0" to a []gl.Float
func StringToGLFloatArray(data string) ([]gl.Float, error) {
	fields := strings.Fields(data)
	if len(fields) == 0 {
		return nil, errors.New("No data to convert")
	}
	returnArray := make([]gl.Float, len(fields))
	for i, val := range fields {
		floatVal, err := strconv.ParseFloat(val, 32)
		if err == nil {
			returnArray[i] = gl.Float(floatVal)
		}
	} 
	return returnArray, nil
}

// Helper function - takes a string of unsigned ints, like
// "23 2 1 3 53 3" and converts them into a []gl.Uint
func StringToGLUintArray(data string) ([]gl.Uint, error) {
	fields := strings.Fields(data)
	if len(fields) == 0 {
		return nil, errors.New("Mesh:StringToGLUintArray:no data to convert")
	}
	returnArray := make([]gl.Uint, len(fields))
	for i, val := range fields {
		intVal, err := strconv.ParseUint(val, 10, 32)
		if err == nil {
			returnArray[i] = gl.Uint(intVal)
		} else {
			returnArray[i] = 0
		}
	}
	return returnArray, nil
}



// Load a mesh from a GLUT mesh file (.xml file)
func (m *Mesh) LoadGLUTMesh(file string) {
	raw, err := LoadGLUTMesh(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load GLUT mesh: %s\n", err)
	} else {
		for i, attr := range raw.Attribute {
			m.AddMeshAttribute(attr.Index, attr.CDATA, attr.)
		}
	}

}
