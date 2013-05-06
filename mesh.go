// mesh - a struct representing a mesh
//
//
//



package goglutils

import (
	"encoding/xml"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"os"
)

type Mesh struct {
	name string
	attributes []*Attribute
	indices []*Index
}

type MeshAttribute struct {
	name string
	data []gl.Float
	stride int
}

type MeshIndex struct {
	name string
	data []gl.Uint
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
	return m
}

func AddAttribute