// mesh - a struct representing a mesh
//
//
//

package goglutils

import (
	//"encoding/xml"
	//"fmt"
	"errors"
	gl "github.com/chsc/gogl/gl33"
	//"os"
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
