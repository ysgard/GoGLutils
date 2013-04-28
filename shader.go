/* Loads fragment and vertex shader code from the supplied files. */

package go-glutils

import (
	"bufio"
	"bytes"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"io"
	"os"
	"path/filepath"
)

// Reads a file and returns its contents as a string.
func ReadSourceFile(filename string) (string, error) {

	fp, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ReadSourceFile: Could not open %s!\n", filename)
		fmt.Fprintf(os.Stderr, "os.Open: %e\n", err)
		return "", err
	}
	defer fp.Close()

	r := bufio.NewReaderSize(fp, 4*1024)
	var buffer bytes.Buffer
	for {
		line, err := r.ReadString('\n')
		buffer.WriteString(line)
		if err == io.EOF {
			// We've read the last string. Make sure there's a null byte.
			buffer.WriteByte('\000')
			break
		}
	}
	return buffer.String(), nil

}

// Create and Compile a shader, and return its object
func CreateShader(shaderType gl.Enum, filePath string) gl.Uint {

	// Start by creating the shader object
	if (shaderType != gl.VERTEX_SHADER) && (shaderType != gl.FRAGMENT_SHADER) {
		fmt.Fprintf(os.Stderr, "User error - not a supported shader type passed to CreateShader\n")
		return 0
	}
	shaderId := gl.CreateShader(shaderType)

	// Load the GLSL source code from the shader file
	shaderCode, err := ReadSourceFile(filePath)
	if err != nil {
		return 0
	}

	// Compile the shader
	var result gl.Int = gl.TRUE
	var infoLogLength gl.Int
	fmt.Fprintf(os.Stdout, "Compiling shader: %s\n", filePath)
	glslCode := gl.GLStringArray(shaderCode)
	defer gl.GLStringArrayFree(glslCode)
	gl.ShaderSource(shaderId, gl.Sizei(len(glslCode)), &glslCode[0], nil)
	gl.CompileShader(shaderId)

	// Check the status of the compile - did it work?
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &result)
	gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		errorMsg := gl.GLStringAlloc(gl.Sizei(infoLogLength))
		defer gl.GLStringFree(errorMsg)
		gl.GetShaderInfoLog(shaderId, gl.Sizei(infoLogLength), nil, errorMsg)
		fmt.Fprintf(os.Stdout, "Shader info for %s: %s", filePath, gl.GoString(errorMsg))
	}
	if result == gl.FALSE {
		fmt.Fprintf(os.Stderr, "Vertex shader compile for %s failed!\n", filePath)
		return 0
	}

	return shaderId
}

// CreateShaderProgram - create a shader program and attach the various shader objects
// defined by the files in the slice, then return the programID.
func CreateShaderProgram(shaderFiles []string) gl.Uint {

	// Create the Program object
	var ProgramID gl.Uint = gl.CreateProgram()

	// For each attached shader, figure out its extension, and load a shader of 
	// that type.
	var sid gl.Uint = 0
	for _, shader := range shaderFiles {
		sid = 0
		switch extension := filepath.Ext(shader); extension {
		case ".vertexshader", ".vert":
			sid = CreateShader(gl.VERTEX_SHADER, shader)

		case ".fragmentshader", ".frag":
			sid = CreateShader(gl.FRAGMENT_SHADER, shader)

		default:
			fmt.Fprintf(os.Stderr, "Don't understand extension %s\n", extension)
			fmt.Fprintf(os.Stderr, "Accepted extensions: .fragmentshader, .vertexshader")
		}
		if sid != 0 {
			gl.AttachShader(ProgramID, sid)
			defer gl.DeleteShader(sid)
		}
	}

	// Link the program
	gl.LinkProgram(ProgramID)

	// Check the program
	var result gl.Int = gl.TRUE
	var infoLogLength gl.Int
	gl.GetProgramiv(ProgramID, gl.LINK_STATUS, &result)
	gl.GetProgramiv(ProgramID, gl.INFO_LOG_LENGTH, &infoLogLength)
	if infoLogLength > 0 {
		programErrorMsg := gl.GLStringAlloc(gl.Sizei(infoLogLength))
		gl.GetProgramInfoLog(ProgramID, gl.Sizei(infoLogLength), nil, programErrorMsg)
		fmt.Fprintf(os.Stdout, "Program Info: %s\n", gl.GoString(programErrorMsg))
	}

	fmt.Fprintf(os.Stdout, "\nLoadShader completed, ProgramID: %d\n", ProgramID)
	return ProgramID
}
