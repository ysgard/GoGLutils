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
	"unsafe"
)

type Mesh struct {
	name string
	//attributes []*MeshAttribute
	attributes []*MeshAttribute
	//indices    []*MeshIndex
	indices  []*MeshIndex
	glvao    gl.Uint
	glbuffer gl.Uint
}

type MeshAttribute struct {
	desc   string
	data   []gl.Float
	stride int
}

type MeshIndex struct {
	desc      string
	data      []gl.Uint
	primitive gl.Enum
	ref       *MeshAttribute
}

func (mi *MeshIndex) Debug() {
	fmt.Fprintf(os.Stdout, "*** MeshIndex ***\n")
	fmt.Fprintf(os.Stdout, "* desc: %s\n", mi.desc)
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
		fmt.Fprintf(os.Stdout, "%s\n", "Could not determine primitive!\n")
	}
	fmt.Fprintf(os.Stdout, "* ref: %p\n", mi.ref)
}

func (ma *MeshAttribute) Debug() {
	fmt.Fprintf(os.Stdout, "*** MeshAttribute ***\n")
	fmt.Fprintf(os.Stdout, "* desc: %s\n", ma.desc)
	fmt.Fprintf(os.Stdout, "* data: ")
	for _, val := range ma.data {
		fmt.Fprintf(os.Stdout, " %f ", val)
	}
	fmt.Fprintf(os.Stdout, "\n* stride: %d\n", ma.stride)
}

func NewMeshIndex(desc string, data []gl.Uint, primitive gl.Enum, ref *MeshAttribute) *MeshIndex {
	i := new(MeshIndex)
	i.desc = desc
	i.data = make([]gl.Uint, len(data))
	copied := copy(i.data, data)
	if copied != len(data) {
		return nil
	}
	i.primitive = primitive
	i.ref = ref
	return i
}

func NewMeshAttribute(desc string, data []gl.Float, stride int) *MeshAttribute {
	a := new(MeshAttribute)
	a.desc = desc
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
	m.attributes = []*MeshAttribute{}
	m.indices = []*MeshIndex{}
	m.glvao = 0
	m.glbuffer = 0
	return m
}

// Add an attribute array to a mesh
func (m *Mesh) AddMeshAttribute(desc string, data []gl.Float, stride int) error {
	ma := NewMeshAttribute(desc, data, stride)
	if ma == nil {
		return errors.New("Mesh:AddMeshAttribute: Could not allocate new attribute array\n")
	}
	m.attributes = append(m.attributes, ma)
	return nil
}

// Add an index to the mesh
func (m *Mesh) AddMeshIndex(desc string, data []gl.Uint, primitive gl.Enum, ref *MeshAttribute) error {
	mi := NewMeshIndex(desc, data, primitive, ref)
	if mi == nil {
		return errors.New("Mesh:AddMeshIndex: Could not allocate new index array\n")
	}
	m.indices = append(m.indices, mi)
	return nil
}

// Set the attribute array this index points to
func (mi *MeshIndex) SetAttribute(ref *MeshAttribute) {
	mi.ref = ref
}

// Set the attribute array a particular index points to
func (m *Mesh) SetAttributeReference(index int, ref *MeshAttribute) {
	m.indices[index].SetAttribute(ref)
}

// Set the primitive type this index draws
func (mi *MeshIndex) SetPrimitive(primitive gl.Enum) {
	mi.primitive = primitive
}

// Set the primitive type for a specific index
func (m *Mesh) SetPrimitive(index int, primitive gl.Enum) {
	m.indices[index].SetPrimitive(primitive)
}

// Get a specific Attribute
func (m *Mesh) Attribute(index int) *MeshAttribute {
	return m.attributes[index]
}

// Helper function, splits a string of floats - like
// "23.3 0.0 2323.0" to a []gl.Float
func StringToGLFloatArray(data string) ([]gl.Float, error) {
	fields := strings.Fields(data)
	if len(fields) == 0 {
		return nil, errors.New("Mesh:StringToGLFloatArray: No data to convert\n")
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
		return nil, errors.New("Mesh:StringToGLUintArray: No data to convert\n")
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
	if len(m.attributes) == 0 {
		return errors.New("Mesh:LoadGLUTMesh: Could not load GLUT mesh, no attribute arrays")
	}

	// GLUT meshes don't single out a particular attribute array for their indexes (annoyingly),
	// so we assume the first one.  Can always reset an index reference later if needed.
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
		err = m.AddMeshIndex(strconv.FormatInt(int64(i), 10), uint_array, primitive, m.attributes[0])
		if err != nil {
			return err
		}

	}
	return nil
}

func (m *Mesh) Debug() {
	fmt.Fprintf(os.Stdout, "*** Debug Mesh: %s ***\n", m.name)
	for _, val := range m.attributes {
		val.Debug()
	}
	for _, val := range m.indices {
		val.Debug()
	}
}

// Set the Mesh's render context
func (m *Mesh) SetRenderContext(vao, buf gl.Uint) {
	m.glvao = vao
	m.glbuffer = buf
}

// Render the mesh using the provided context.
func (m *Mesh) Render() error {

	// Bind the VAO
	gl.BindVertexArray(m.glvao)
	// Bind the buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, m.glbuffer)

	// If we don't have a valid context, we return an error.
	if m.glvao == 0 || m.glbuffer == 0 {
		return errors.New(fmt.Sprintf("Mesh:Render: Mesh context invalid: vao=%d, buffer=%d", m.glvao, m.glbuffer))
	}
	if gl.IsBuffer(m.glbuffer) == gl.FALSE {
		return errors.New("Mesh:Render: Invalid OpenGL buffer!")
	}
	if gl.IsVertexArray(m.glvao) == gl.FALSE {
		return errors.New("Mesh:Render: Invalid OpenGL VAO!")
	}
	// If we don't have any indices, we don't have anything to do, return
	if len(m.indices) == 0 {
		return nil
	}

	// Buffer the vertex positions
	bufferLen := unsafe.Sizeof(gl.Float(0)) * (uintptr)(len(m.attributes[0].data))
	gl.BufferData(gl.ARRAY_BUFFER,
		gl.Sizeiptr(bufferLen),
		gl.Pointer(&m.attributes[0].data[0]),
		gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// For each index, draw that part of the mesh
	for _, indx := range m.indices {
		gl.DrawElements(indx.primitive,
			gl.Sizei(len(indx.data)),
			gl.UNSIGNED_INT, nil)
	}
	// Unbind the VAO
	gl.BindVertexArray(0)
	return nil
}

// Helper function, specify a VAO/Buffer when rendering
func (m *Mesh) RenderVB(v, b gl.Uint) error {
	// Save the current vao, buffer
	tmpvao := m.glvao
	tmpbuf := m.glbuffer
	m.SetRenderContext(v, b)
	err := m.Render()
	m.SetRenderContext(tmpvao, tmpbuf)
	return err
}
