// mesh - a struct representing a mesh

package goglutils

import (
	"encoding/xml"
	"fmt"
	//gl "github.com/chsc/gogl/gl33"
	"os"
)

type RawMesh struct {
	XMLName   xml.Name       `xml:"mesh"`
	Namespace string         `xml:"xmlns,attr"`
	Attribute []RawAttribute `xml:"attribute"`
	Indices   []RawIndices   `xml:"indices"`
}

func (m *RawMesh) Debug() {
	fmt.Fprintf(os.Stdout, "Mesh/Namespace -- %s\n", m.Namespace)
	for _, val := range m.Attribute {
		val.Debug()
	}
	for _, val := range m.Indices {
		val.Debug()
	}
}

type RawAttribute struct {
	XMLName xml.Name `xml:"attribute"`
	Index   string   `xml:"index,attr"`
	Type    string   `xml:"type,attr"`
	Size    string   `xml:"size,attr"`
	CDATA   string   `xml:",chardata"`
}

func (a *RawAttribute) Debug() {
	fmt.Fprintf(os.Stdout, "Attribute/Index -- %s\n", a.Index)
	fmt.Fprintf(os.Stdout, "Attribute/Type -- %s\n", a.Type)
	fmt.Fprintf(os.Stdout, "Attribute/Size -- %s\n", a.Size)
	fmt.Fprintf(os.Stdout, "Attribute/CDATA -- %s\n", a.CDATA)
}

type RawIndices struct {
	XMLName xml.Name `xml:"indices"`
	Cmd     string   `xml:"cmd,attr"`
	Type    string   `xml:"type,attr"`
	CDATA   string   `xml:",chardata"`
}

func (i *RawIndices) Debug() {
	fmt.Fprintf(os.Stdout, "Indices/Cmd -- %s\n", i.Cmd)
	fmt.Fprintf(os.Stdout, "Indices/Type -- %s\n", i.Type)
	fmt.Fprintf(os.Stdout, "Indices/CDATA -- %s\n", i.CDATA)
}

func LoadMesh(filename string) *RawMesh {
	file, err := os.Open(filename)
	fi, err := file.Stat()
	filelen := fi.Size()
	buf := make([]byte, filelen)
	read, err := file.Read(buf)
	if read != int(filelen) {
		fmt.Fprintf(os.Stderr, "Could not read complete contents of file: %d read vs %d size", read, filelen)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}

	m := new(RawMesh)
	err = xml.Unmarshal(buf, m)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	return m
}

// func main() {
// 	if len(os.Args) < 2 {
// 		fmt.Fprintf(os.Stderr, "syntax: Mesh <meshfile>\n")
// 		return
// 	}
// 	filename := os.Args[1]
// 	m := LoadMesh(filename)
// 	m.Debug()
// }
