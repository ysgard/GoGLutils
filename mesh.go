// mesh - a struct representing a mesh
//
//
//

package goglutils

import (
	//"encoding/xml"
	"errors"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"os"
	"strconv"
	"strings"
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

func (mi *MeshIndex) Debug() {
	fmt.Fprintf(os.Stdout, "*** MeshIndex ***\n")
	fmt.Fprintf(os.Stdout, "* name: %s\n", mi.name)
	fmt.Fprintf(os.Stdout, "* data: ")
	for _, val := range mi.data {
		fmt.Fprintf(os.Stdout, " %d ", val)
	}
	fmt.Fprintf(os.Stdout, "\n* primitive: ")
	switch mi.primitive {
	case gl.TRIANGLES:
		fmt.Fprintf(os.Stdout, "%s\n", "triangles")
	case gl.TRIANGLE_FAN:
		fmt.Fprintf(os.Stdout, "%s\n", "triangle_fan")
	case gl.TRIANGLE_STRIP:
		fmt.Fprintf(os.Stdout, "%s\n", "triangle_strip")
	case gl.POINTS:
		fmt.Fprintf(os.Stdout, "%s\n", "points")
	case gl.LINES:
		fmt.Fprintf(os.Stdout, "%s\n", "lines")
	case gl.LINE_LOOP:
		fmt.Fprintf(os.Stdout, "%s\n", "line_loop")
	case gl.LINE_STRIP:
		fmt.Fprintf(os.Stdout, "%s\n", "line_strip")
	default:
		fmt.Fprintf(os.Stdout, "%s\n", "Could not determine primitive!")
	}
}

func (ma *MeshAttribute) Debug() {
	fmt.Fprintf(os.Stdout, "*** MeshAttribute ***\n")
	fmt.Fprintf(os.Stdout, "* name: %s\n", ma.name)
	fmt.Fprintf(os.Stdout, "* data: ")
	for _, val := range ma.data {
		fmt.Fprintf(os.Stdout, " %f ", val)
	}
	fmt.Fprintf(os.Stdout, "\n* stride: %d\n", ma.stride)
}

func NewMeshIndex(name string, data []gl.Uint, primitive gl.Enum) *MeshIndex {
	i := new(MeshIndex)
	i.name = name
	i.data = make([]gl.Uint, len(data))
	copied := copy(i.data, data)
	if copied != len(data) {
		return nil
	}
	i.primitive = primitive
	return i
}

func NewMeshAttribute(name string, data []gl.Float, stride int) *MeshAttribute {
	a := new(MeshAttribute)
	a.name = name
	a.data = make([]gl.Float, len(data))
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
	ma := NewMeshAttribute(name, data, stride)
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
		return nil, errors.New("Mesh:StringToGLUintArray:no data to convert\n")
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
func (m *Mesh) LoadGLUTMesh(file string) error {
	raw, err := LoadGLUTMesh(file)
	if err != nil {
		return err
	}
	for _, attr := range raw.Attribute {
		float_array, err := StringToGLFloatArray(attr.CDATA)
		if err != nil {
			return err
		}
		size, err := strconv.ParseInt(attr.Size, 10, 32)
		if err != nil {
			return err
		}
		err = m.AddMeshAttribute(attr.Index, float_array, int(size))
		if err != nil {
			return err
		}
	}
	for i, indx := range raw.Indices {
		var primitive gl.Enum
		switch indx.Cmd {
		case "triangles":
			primitive = gl.TRIANGLES
		case "tri-fan":
			primitive = gl.TRIANGLE_FAN
		case "tri-strip":
			primitive = gl.TRIANGLE_STRIP
		case "points":
			primitive = gl.POINTS
		case "lines":
			primitive = gl.LINES
		case "line-strip":
			primitive = gl.LINE_STRIP
		case "line-loop":
			primitive = gl.LINE_LOOP
		default:
			continue
		}
		uint_array, err := StringToGLUintArray(indx.CDATA)
		if err != nil {
			return err
		}
		err = m.AddMeshIndex(strconv.FormatInt(int64(i), 10), uint_array, primitive)
		if err != nil {
			return err
		}

	}
	return nil
}

func (m *Mesh) Debug() {
	fmt.Fprintf(os.Stdout, "*** Debug Mesh: %s ***\n", m.name)
	for key, val := range m.attributes {
		fmt.Fprintf(os.Stdout, "* Attribute array: %s\n", key)
		val.Debug()
	}
	for key, val := range m.indices {
		fmt.Fprintf(os.Stdout, "* Index array: %s\n", key)
		val.Debug()
	}
}
