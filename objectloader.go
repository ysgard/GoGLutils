/*
ObjectLoader - load .obj models for use by OpenGL
*/
package main

import (
	"bufio"
	"fmt"
	"github.com/Jragonmiris/mathgl"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	debug = false
)

func debugMsg(msg string) {
	if debug == true {
		fmt.Fprintf(os.Stdout, "%s\n", msg)
	}
}

func debugIndex(d []uint, a []mathgl.Vec3f) {
	if debug == true {
		fmt.Fprintf(os.Stdout, "(((%d elements)))\n", len(d))
		for _, val := range d {
			fmt.Fprintf(os.Stdout, "\t%d\t\t%f %f %f\n", val, a[val][0], a[val][1], a[val][2])
		}
	}
}

// decipherFace - takes a v/uv/n, if a value is missing it 
// returns 0 for that value.
func decipherFace(face string) (uint, uint, uint) {
	var vertexIndice, uvIndice, normalIndice uint64
	var err error
	ar := strings.Split(strings.TrimSpace(face), "/")

	if vertexIndice, err = strconv.ParseUint(ar[0], 0, 32); err != nil {
		vertexIndice = 0
	}
	if uvIndice, err = strconv.ParseUint(ar[1], 0, 32); err != nil {
		uvIndice = 0
	}
	if normalIndice, err = strconv.ParseUint(ar[2], 0, 32); err != nil {
		normalIndice = 0
	}
	return (uint)(vertexIndice), (uint)(uvIndice), (uint)(normalIndice)
}

// Dump the output of LoadOBJ
func dumpOBJ(filePath string) {
	vtx, uv, nrm := loadOBJ(filePath)
	fmt.Printf("*** %d VERTEXES ***\n\n", len(vtx))
	for _, val := range vtx {
		fmt.Printf(" %f\t%f\t%f\n", val[0], val[1], val[2])
	}
	fmt.Printf("*** %d UVS ***\n\n", len(uv))
	for _, val := range uv {
		fmt.Printf(" %f\t%f\n", val[0], val[1])
	}
	fmt.Printf("*** %d NORMALS ***\n\n", len(nrm))
	for _, val := range nrm {
		fmt.Printf(" %f\t%f\t%f\n", val[0], val[1], val[2])
	}
}

func stringsToVec3f(first, second, third string) mathgl.Vec3f {
	var x mathgl.Vec3f
	var a, b, c float64
	var err error = nil
	a, err = strconv.ParseFloat(first, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stringsToVec3f: could not parse %s: %s", first, err)
		a = 0
	}
	b, err = strconv.ParseFloat(second, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stringsToVec3f: could not parse %s: %s", second, err)
		b = 0
	}
	c, err = strconv.ParseFloat(third, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "stringsToVec3f: could not parse %s: %s", third, err)
		c = 0
	}

	x[0] = (float32)(a)
	x[1] = (float32)(b)
	x[2] = (float32)(c)

	return x
}

func stringsToVec2f(first, second string) mathgl.Vec2f {
	var x mathgl.Vec2f
	var a, b float64
	var err error
	if a, err = strconv.ParseFloat(first, 32); err != nil {
		fmt.Fprintf(os.Stderr, "stringsToVec3f: could not parse %s: %s", first, err)
		a = 0
	}
	if b, err = strconv.ParseFloat(second, 32); err != nil {
		fmt.Fprintf(os.Stderr, "stringsToVec3f: could not parse %s: %s", second, err)
		b = 0
	}

	x[0] = (float32)(a)
	x[1] = (float32)(b)

	return x
}

func loadOBJ(filePath string) ([]mathgl.Vec3f, []mathgl.Vec2f, []mathgl.Vec3f) {

	var vertexIndices, uvIndices, normalIndices []uint

	var vertices = make([]mathgl.Vec3f, 1, 50)
	var uvs = make([]mathgl.Vec2f, 1, 50)
	var normals = make([]mathgl.Vec3f, 1, 50)

	// Open the file, and prepare a buffer to read in lines
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not open file %s\n", filePath)
		return nil, nil, nil
	}
	defer fp.Close()
	fileBuf := bufio.NewReader(fp)

	// Loop through the lines, 'till the file is done
	for {
		line, err := fileBuf.ReadString('\n')
		if line == "" && err == io.EOF {
			debugMsg("Reached EOF, no line, quitting")
			break
		} // EOF on line by itself
		debugMsg(fmt.Sprintf("Read %s from file, processing...", line))
		words := strings.Fields(line)
		switch words[0] {
		case "v":
			debugMsg("Found vertex...")
			var vtx = stringsToVec3f(words[1], words[2], words[3])
			vertices = append(vertices, vtx)
		case "vt":
			debugMsg("Found uv coord...")
			var uv = stringsToVec2f(words[1], words[2])
			uvs = append(uvs, uv)
		case "vn":
			debugMsg("Found normal...")
			var normal = stringsToVec3f(words[1], words[2], words[3])
			normals = append(normals, normal)
		case "f":
			debugMsg("Found face...")
			vi1, uvi1, ni1 := decipherFace(words[1])
			vi2, uvi2, ni2 := decipherFace(words[2])
			vi3, uvi3, ni3 := decipherFace(words[3])
			vertexIndices = append(vertexIndices, vi1, vi2, vi3)
			uvIndices = append(uvIndices, uvi1, uvi2, uvi3)
			normalIndices = append(normalIndices, ni1, ni2, ni3)
		default:
			debugMsg("Not v, vt, vn or f, skipping...")
			continue
		}
		if err == io.EOF {
			debugMsg("Found EOF at end of line, quitting.")
			break
		} // EOF at end of line
	}

	// Translate object file vertices into an order more preferable to 
	// OpenGL.  OpenGL expects triplets of vertexes, one after another,
	// defining each triangle, but the OBJ format uses the 'f' (face)
	// line to define a triangle and maps it to a listed vertex.
	realVertices := make([]mathgl.Vec3f, 1, 50)
	realUVs := make([]mathgl.Vec2f, 1, 50)
	realNormals := make([]mathgl.Vec3f, 1, 50)

	debugIndex(vertexIndices, vertices)

	// Now, for each vertex of each triangle
	for i := 0; i < len(vertexIndices); i++ {
		// Get the indices of its attributes
		vertexIndex := vertexIndices[i]
		uvIndex := uvIndices[i]
		normalIndex := normalIndices[i]

		// Get the attributes thanks to the index
		if vertexIndex != 0 {
			realVertices = append(realVertices, vertices[vertexIndex])
		}
		if uvIndex != 0 {
			realUVs = append(realUVs, uvs[uvIndex])
		}
		if normalIndex != 0 {
			realNormals = append(realNormals, normals[normalIndex])
		}
	}
	return realVertices[1:], realUVs[1:], realNormals[1:]

}

//Driver, remove once loader has been tested.
// func main() {
// 	dumpOBJ("art/cylinder.obj")
// }
