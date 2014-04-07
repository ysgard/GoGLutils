# GoGLutils
---
Simple utilities for use with chsc's excellent  [http://github.com/chsc/gogl](*gogl library*).

*This is still a work in progress!  There are still several bugs to work out.*

## Installation

  go get github.com/Ysgard/goglutils

As of right now, it only works if you have gl33 installed.  I'm working on
letting it use any version of GoGL.

## matrix.go ##

A simple library to provide matrix structs, methods and functions for vector/matrix math.

## matrixstack.go ##

Provides a matrix stack to provide the ability to pop/push 4x4 gl.float matrices.

## shader.go ##

Provides a simple function to load .vert and .frag shaders and return a shader program id.

## mesh.go ***INCOMPLETE*** ##

Provides a simple Mesh struct that keeps track of its vertex arrays and can be loaded via COLLADA (.dae), Object (.obj) and
GLUT mesh (.xml) files. 

## collada.go ## 

Parses a Collada file and provides the raw data in a COLLADA
struct.

## objectloader.go ***INCOMPLETE*** ##

Parses a .obj file and returns the raw data in an OBJECT struct.

## glutmesh.go ***INCOMPLETE*** ##

Reads a GLUT-style XML mesh file and returns a GLUTMESH struct.





PACKAGE DOCUMENTATION ***INCOMPLETE***

Go doc documentation below:



package goglutils
    import "/Users/wyvern/GoCode/src/github.com/Ysgard/goglutils"

    A collection of simple routines and structures to help with matrices as
    well functions for common math operations.

    I tried to keep the structure similar to glm's matrix and vector classes
    Therefore, the matrices are stored in column order:

	v0      v1      v2      v3

    x | { v0x } { v1x } { v2x } { v3x } | y | { v0y } { v1y } { v2y } { v3y
    } | z | { v0z } { v1z } { v2z } { v3z } | w | { v0w } { v1w } { v2w } {
    v3w } |

    So if you were going to scale a matrix, you'd set: mat[0].x = scaleX *
    mat[0].x mat[1].y = scaleY * mat[1].y mat[2].z = scaleZ * mat[2].z

    Must be careful - Can't have vector pointers inside the matrix, as the
    value must be contiguous in memory if we're going to pass them to
    OpenGL.

    This package also contains OpenGL-friendly math functions - these are
    basically Go's math functions wrapped in pointer casts to make them use
    gl.Float (glFloat)

    There is also a simple matrix stack utility, as well as a functions that
    loads and creates shaders from .vert and .grag GLSL files.


CONSTANTS

const Pi = (gl.Float)(math.Pi)


FUNCTIONS

func Clamp(fValue, fMinValue, fMaxValue gl.Float) gl.Float
    Clamp - constrain a value fValue to the range delimited by fMinValue ->
    fMaxValue

func CosGL(Rad gl.Float) gl.Float
    Cosine, in gl.Float

func CreateShader(shaderType gl.Enum, filePath string) gl.Uint
    Create and Compile a shader, and return its object

func CreateShaderProgram(shaderFiles []string) gl.Uint
    CreateShaderProgram - create a shader program and attach the various
    shader objects defined by the files in the slice, then return the
    programID.

func DebugMat(m []gl.Float, s string)
    Pretty-print a []gl.Float slice representing a 16-item transformation
    matrix.

func DegToRad(fAngDeg gl.Float) gl.Float
    Convert degrees to radians

func Ident4() []gl.Float
    Identity matrix, bare

func Invert(m []gl.Float) ([]gl.Float, error)
    Code from the MESA library, adapted for Go

func LerpGL(start, end, ratio gl.Float) gl.Float
    Basic linear interpolation

func ModGL(a, b gl.Float) gl.Float
    Take two gl.Floats and return remainder as a gl.Float

func ReadSourceFile(filename string) (string, error)
    Reads a file and returns its contents as a string.

func SinGL(Rad gl.Float) gl.Float
    Sine, in gl.Float

func TanGL(Rad gl.Float) gl.Float
    Tan, in gl.Float


TYPES

type Mat4 [4]Vec4
    Struct that kinda, sorta represents a glm/glsl 4x4 matrix


func FromArray(arr []gl.Float) (*Mat4, error)
    FromArray - produce a Mat4 from a []gl.Float. Basically the inverse of
    ToArray


func IdentMat4() *Mat4
    Return a Mat4 with identity values


func RotateX(fAngDeg gl.Float) *Mat4
    Returns a Mat4 representing a rotation matrix for the angle given in
    degrees


func RotateY(fAngDeg gl.Float) *Mat4
    Returns a Mat4 representing a rotation matrix for the angle given in
    degree


func RotateZ(fAngDeg gl.Float) *Mat4
    Returns a Mat4 representing a rotation matrix for the angle given in
    degrees


func (m *Mat4) Copy() *Mat4
    Create a copy of a given mat4

func (m *Mat4) Inverse() *Mat4
    Returns a new Mat4 representing the inverse of the Mat4

func (m1 *Mat4) MulM(m2 *Mat4) *Mat4
    Multiply receiving matrix by given Mat4 and return the new Mat.

func (m *Mat4) MulS(s gl.Float) *Mat4
    Multiplies a Matrix by a scalar s and returns the new matrix

func (m *Mat4) MulV(v *Vec4) *Vec4
    Multiply receiving matrix by given Vec4 and return the new Vec4

func (m *Mat4) Print(s string)
    Pretty-prints a Mat4 with an optional header

func (m *Mat4) Scale(s *Vec4) *Mat4
    Scale - Scales a matrix using a passed Vec4, the vec4 should take the
    form { sx, sy, sz, 1.0 }

func (m *Mat4) ToArray() []gl.Float
    ToArray - produce a []gl.Float array from a given struct. Perhaps not
    necessary, doing &Mat4 should be sufficient!

func (m *Mat4) Translate(offset *Vec4) *Mat4
    Take a translation vector {tx, ty, tz, 1.0} and translate the matrix

func (m *Mat4) Transpose() *Mat4
    Returns the transpose of a given matrix


type MatrixStack struct {
    // contains filtered or unexported fields
}
    MatrixStack - Represents a way to store a sequential series of
    transformations


func (ms *MatrixStack) Init()
    Creates a default identity matrix as the current matrix

func (ms *MatrixStack) Invert()
    Inverts the topmost matrix on the stack

func (ms *MatrixStack) MulM(m *Mat4)
    Multiplies by another mat4 the topmost matrix on the stack

func (ms *MatrixStack) Pop()
    Pop the last matrix off the stack and make it the current matrix

func (ms *MatrixStack) Push()
    Create a copy of the current matrix and push it onto the stack

func (ms *MatrixStack) RotateX(deg gl.Float)
    X-rotates the topmost matrix on the stack

func (ms *MatrixStack) RotateY(deg gl.Float)
    Y-rotates the topmost matrix on the stack

func (ms *MatrixStack) RotateZ(deg gl.Float)
    Z-rotates the topmost matrix on the stack

func (ms *MatrixStack) Scale(s *Vec4)
    Scales the topmost matrix on the stack

func (ms *MatrixStack) Top() *Mat4
    Return pointer to top matrix

func (ms *MatrixStack) Translate(offset *Vec4)
    Translates the topmost matrix on the stack


type Vec3 struct {
    // contains filtered or unexported fields
}
    Struct that kinda, sorta represents a vec3 glm/glsl vector


func (u *Vec3) Add(v *Vec3) *Vec3
    Add together two Vec3's - u.Add(v)

func (u *Vec3) Cross(v *Vec3) *Vec3
    Cross product - Vec3 version, u.Cross(v) = u x v

func (u *Vec3) MulS(f gl.Float) *Vec3
    Multiply vector by a scalar

func (v *Vec3) Normalize()
    Normalize - Vec3 version

func (u *Vec3) Sub(v *Vec3) *Vec3
    Subtract two Vec3's - u.Sub(v)

func (v3 *Vec3) V3to4(f gl.Float) *Vec4
    Vec4 from a Vec3


type Vec4 struct {
    // contains filtered or unexported fields
}
    Struct that kinda, sorta represents a glm vec4


func (v *Vec4) Normalize()
    Normalize - normalizes a vector, doesn't include w



