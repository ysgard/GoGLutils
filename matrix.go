// A collection of simple routines and structures to help with matrices
// as well functions for common math operations.
//
// I tried to keep the structure similar to glm's matrix and vector classes
// Therefore, the matrices are stored in column order:
//
//       v0      v1      v2      v3
// x | { v0x } { v1x } { v2x } { v3x } |
// y | { v0y } { v1y } { v2y } { v3y } |
// z | { v0z } { v1z } { v2z } { v3z } |
// w | { v0w } { v1w } { v2w } { v3w } |
//
// So if you were going to scale a matrix, you'd set:
// mat[0].X = scaleX * mat[0].X
// mat[1].Y = scaleY * mat[1].Y
// mat[2].Z = scaleZ * mat[2].Z
//
// Must be careful - Can't have vector pointers inside the matrix, as
// the value must be contiguous in memory if we're going to pass
// them to OpenGL.
//
// This package also contains OpenGL-friendly math functions - these
// are basically Go's math functions wrapped in pointer casts to make
// them use gl.Float (glFloat)
//
// There is also a simple matrix stack utility, as well as a functions
// that loads and creates shaders from .vert and .grag GLSL files.
package goglutils

import (
	"errors"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"math"
	"os"
	"strings"
)

// Constants
const degToRad = math.Pi * 2.0 / 360
const Pi = (gl.Float)(math.Pi)

// Change this to change where debug messages get sent
var debugOut = os.Stderr

// ******************************* //
// *     VEC3 - A 3x1 vector     * //
// ******************************* //

// Struct that kinda, sorta represents a vec3 glm/glsl vector
type Vec3 struct {
	X, Y, Z gl.Float
}

func NewVec3(x, y, z gl.Float) *Vec3 {
	return &Vec3{x, y, z}
}

// Normalize - Vec3 version
func (v *Vec3) Normalize() {
	lenv := (gl.Float)(math.Sqrt((float64)(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	v.X = v.X / lenv
	v.Y = v.Y / lenv
	v.Z = v.Z / lenv
}

// Cross product - Vec3 version, u.Cross(v) = u x v
func (u *Vec3) Cross(v *Vec3) *Vec3 {
	s := Vec3{
		u.Y*v.Z - u.Z*v.Y,
		u.Z*v.X - u.X*v.Z,
		u.X*v.Y - u.Y*v.X,
	}
	return &s
}

// Add together two Vec3's - u.Add(v)
func (u *Vec3) Add(v *Vec3) *Vec3 {
	s := Vec3{
		u.X + v.X,
		u.Y + v.Y,
		u.Z + v.Z,
	}
	return &s
}

// Subtract two Vec3's - u.Sub(v)
func (u *Vec3) Sub(v *Vec3) *Vec3 {
	s := Vec3{
		u.X - v.X,
		u.Y - v.Y,
		u.Z - v.Z,
	}
	return &s
}

// Multiply vector by a scalar
func (u *Vec3) MulS(f gl.Float) *Vec3 {
	s := Vec3{
		u.X * f,
		u.Y * f,
		u.Z * f,
	}
	return &s
}

// ******************************* //
// *     VEC4 - A 4x1 vector     * //
// ******************************* //

// Struct that kinda, sorta represents a glm vec4
type Vec4 struct {
	X, Y, Z, W gl.Float
}

func NewVec4(x, y, z, w gl.Float) *Vec4 {
	return &Vec4{x, y, z, w}
}

// Vec4 from a Vec3
func (v3 *Vec3) To4W(f gl.Float) *Vec4 {
	v4 := Vec4{
		v3.X,
		v3.Y,
		v3.Z,
		f,
	}
	return &v4
}

// Implicit Vec4 from a Vec3, assumes 1.0 for w
func (v3 *Vec3) To4() *Vec4 {
	return v3.To4W(1.0)
}

// Normalize - normalizes a vector, doesn't include w
func (v *Vec4) Normalize() {
	lenv := (gl.Float)(math.Sqrt((float64)(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
	v.X = v.X / lenv
	v.Y = v.Y / lenv
	v.Z = v.Z / lenv
}

// ******************************* //
// *     MAT4 - A 4x4 Matrix     * //
// ******************************* //

// Struct that kinda, sorta represents a glm/glsl 4x4 matrix
type Mat4 [4]Vec4

// Return a Mat4 as a *gl.Float
func (m *Mat4) GetPtr() *gl.Float {
	return &m[0].X
}

// Multiply receiving matrix by given Vec4 and return
// the new Vec4
func (m *Mat4) MulV(v *Vec4) *Vec4 {
	rv := Vec4{0.0, 0.0, 0.0, 0.0}
	rv.X = m[0].X*v.X + m[1].X*v.Y + m[2].X*v.Z + m[3].X*v.W
	rv.Y = m[0].Y*v.X + m[1].Y*v.Y + m[2].Y*v.Z + m[3].Y*v.W
	rv.Z = m[0].Z*v.X + m[1].Z*v.Y + m[2].Z*v.Z + m[3].Z*v.W
	rv.W = m[0].W*v.X + m[1].W*v.Y + m[2].W*v.Z + m[3].W*v.W
	return &rv
}

// Multiply receiving matrix by given Mat4 and return
// the new Mat.
func (m1 *Mat4) MulM(m2 *Mat4) *Mat4 {
	var rm = Mat4{
		{
			m1[0].X*m2[0].X + m1[1].X*m2[0].Y + m1[2].X*m2[0].Z + m1[3].X*m2[0].W,
			m1[0].Y*m2[0].X + m1[1].Y*m2[0].Y + m1[2].Y*m2[0].Z + m1[3].Y*m2[0].W,
			m1[0].Z*m2[0].X + m1[1].Z*m2[0].Y + m1[2].Z*m2[0].Z + m1[3].Z*m2[0].W,
			m1[0].W*m2[0].X + m1[1].W*m2[0].Y + m1[2].W*m2[0].Z + m1[3].W*m2[0].W,
		},
		{
			m1[0].X*m2[1].X + m1[1].X*m2[1].Y + m1[2].X*m2[1].Z + m1[3].X*m2[1].W,
			m1[0].Y*m2[1].X + m1[1].Y*m2[1].Y + m1[2].Y*m2[1].Z + m1[3].Y*m2[1].W,
			m1[0].Z*m2[1].X + m1[1].Z*m2[1].Y + m1[2].Z*m2[1].Z + m1[3].Z*m2[1].W,
			m1[0].W*m2[1].X + m1[1].W*m2[1].Y + m1[2].W*m2[1].Z + m1[3].W*m2[1].W,
		},
		{
			m1[0].X*m2[2].X + m1[1].X*m2[2].Y + m1[2].X*m2[2].Z + m1[3].X*m2[2].W,
			m1[0].Y*m2[2].X + m1[1].Y*m2[2].Y + m1[2].Y*m2[2].Z + m1[3].Y*m2[2].W,
			m1[0].Z*m2[2].X + m1[1].Z*m2[2].Y + m1[2].Z*m2[2].Z + m1[3].Z*m2[2].W,
			m1[0].W*m2[2].X + m1[1].W*m2[2].Y + m1[2].W*m2[2].Z + m1[3].W*m2[2].W,
		},
		{
			m1[0].X*m2[3].X + m1[1].X*m2[3].Y + m1[2].X*m2[3].Z + m1[3].X*m2[3].W,
			m1[0].Y*m2[3].X + m1[1].Y*m2[3].Y + m1[2].Y*m2[3].Z + m1[3].Y*m2[3].W,
			m1[0].Z*m2[3].X + m1[1].Z*m2[3].Y + m1[2].Z*m2[3].Z + m1[3].Z*m2[3].W,
			m1[0].W*m2[3].X + m1[1].W*m2[3].Y + m1[2].W*m2[3].Z + m1[3].W*m2[3].W,
		},
	}
	return &rm
}

// Returns the transpose of a given matrix
func (m *Mat4) Transpose() *Mat4 {
	var rm = Mat4{
		{m[0].X, m[1].X, m[2].X, m[3].X},
		{m[0].Y, m[1].Y, m[2].Y, m[3].Y},
		{m[0].Z, m[1].Z, m[2].Z, m[3].Z},
		{m[0].W, m[1].W, m[2].W, m[3].W},
	}
	return &rm
}

// Scale - Scales a matrix using a passed Vec4, the vec4 should take the
// form { sx, sy, sz, 1.0 }
func (m *Mat4) Scale(s *Vec4) *Mat4 {
	scaleMat := IdentMat4()
	scaleMat[0].X = s.X
	scaleMat[1].Y = s.Y
	scaleMat[2].Z = s.Z
	return m.MulM(scaleMat)
}

// Multiplies a Matrix by a scalar s and returns the new matrix
func (m *Mat4) MulS(s gl.Float) *Mat4 {
	var rm = Mat4{
		{m[0].X * s, m[0].Y * s, m[0].Z * s, m[0].W * s},
		{m[1].X * s, m[1].Y * s, m[1].Z * s, m[1].W * s},
		{m[2].X * s, m[2].Y * s, m[2].Z * s, m[2].W * s},
		{m[3].X * s, m[3].Y * s, m[3].Z * s, m[3].W * s},
	}
	return &rm
}

// Take a translation vector {tx, ty, tz, 1.0} and
// translate the matrix
func (m *Mat4) Translate(offset *Vec4) *Mat4 {
	var tm = IdentMat4()
	tm[3].X = offset.X
	tm[3].Y = offset.Y
	tm[3].Z = offset.Z
	return m.MulM(tm)
}

// ToArray - produce a []gl.Float array from a given struct.
// Perhaps not necessary, doing &Mat4 should be sufficient!
func (m *Mat4) ToArray() []gl.Float {
	arr := make([]gl.Float, 16)
	for i, vec := range m {
		arr[i*4] = vec.X
		arr[i*4+1] = vec.Y
		arr[i*4+2] = vec.Z
		arr[i*4+3] = vec.W
	}
	return arr
}

// FromArray - produce a Mat4 from a []gl.Float.  Basically
// the inverse of ToArray
func FromArray(arr []gl.Float) (*Mat4, error) {
	if len(arr) < 16 {
		return nil, errors.New("Need 16-element float array")
	}
	rm := IdentMat4()
	for i := 0; i < 4; i++ {
		rm[i].X = arr[i*4]
		rm[i].Y = arr[i*4+1]
		rm[i].Z = arr[i*4+2]
		rm[i].W = arr[i*4+3]
	}
	return rm, nil
}

// Return a Mat4 with identity values
func IdentMat4() *Mat4 {
	var m Mat4
	m[0].X = 1.0
	m[1].Y = 1.0
	m[2].Z = 1.0
	m[3].W = 1.0
	return &m
}

// Create a copy of a given mat4
func (m *Mat4) Copy() *Mat4 {
	copy := IdentMat4()
	for i := 0; i < 4; i++ {
		copy[i].X = m[i].X
		copy[i].Y = m[i].Y
		copy[i].Z = m[i].Z
		copy[i].W = m[i].W
	}
	return copy
}

// Returns a Mat4 representing a rotation matrix
// for the angle given in degrees
func RotateX(fAngDeg gl.Float) *Mat4 {
	fAngRad := DegToRad(fAngDeg)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[1].Y = fCos
	theMat[2].Y = -fSin
	theMat[1].Z = fSin
	theMat[2].Z = fCos
	return theMat
}

// Returns a Mat4 representing a rotation matrix
// for the angle given in degree
func RotateY(fAngDeg gl.Float) *Mat4 {
	fAngRad := DegToRad(fAngDeg)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[0].X = fCos
	theMat[2].X = fSin
	theMat[0].Z = -fSin
	theMat[2].Z = fCos
	return theMat
}

// Returns a Mat4 representing a rotation matrix
// for the angle given in degrees
func RotateZ(fAngDeg gl.Float) *Mat4 {
	fAngRad := DegToRad(fAngDeg)
	fCos := CosGL(fAngRad)
	fSin := SinGL(fAngRad)
	theMat := IdentMat4()
	theMat[0].X = fCos
	theMat[1].X = -fSin
	theMat[0].Y = fSin
	theMat[1].Y = fCos
	return theMat
}

// Pretty-prints a Mat4 with an optional header
func (m *Mat4) Print(s string) {
	if s == "" {
		s = "Debugging Matrix"
	}
	slen := len(s) + 2

	var dashes string
	if (58-slen)&1 == 1 {
		// odd-string
		dashes = strings.Repeat("-", (58-slen-1)/2) + " " +
			s + " " + strings.Repeat("-", ((58-slen-1)/2))
	} else {
		// even-string
		dashes = strings.Repeat("-", (58-slen)/2) + " " +
			s + " " + strings.Repeat("-", (58-slen)/2-1)
	}
	fmt.Fprintf(debugOut, "%s\n", dashes)
	fmt.Fprintf(debugOut, "%9.3f       %9.3f       %9.3f       %9.3f\n", m[0].X, m[1].X, m[2].X, m[3].X)
	fmt.Fprintf(debugOut, "%9.3f       %9.3f       %9.3f       %9.3f\n", m[0].Y, m[1].Y, m[2].Y, m[3].Y)
	fmt.Fprintf(debugOut, "%9.3f       %9.3f       %9.3f       %9.3f\n", m[0].Z, m[1].Z, m[2].Z, m[3].Z)
	fmt.Fprintf(debugOut, "%9.3f       %9.3f       %9.3f       %9.3f\n\n", m[0].W, m[1].W, m[2].W, m[3].W)
	//fmt.Fprintf(debugOut, "\t------------------------------------------------------------\n")
}

// Returns a new Mat4 representing the inverse of the Mat4
func (m *Mat4) Inverse() *Mat4 {
	// Convert Mat4 to an array of floats
	inArray := m.ToArray()
	if outArray, err := Invert(inArray); err == nil {
		// Craft Mat4 from the array of floats
		outMat, _ := FromArray(outArray)
		return outMat
	} else {
		return nil
	}
}

// Returns an orthographic projection matrix
func Ortho(left, right, bottom, top, nearVal, farVal gl.Float) *Mat4 {
	m := IdentMat4()
	m[0].X = 2.0 / (right - left)
	m[1].Y = 2.0 / (top - bottom)
	m[2].Z = -2.0 / (farVal - nearVal)
	m[3].X = -(right + left) / (right - left)
	m[3].Y = -(top + bottom) / (top - bottom)
	m[3].Z = -(farVal + nearVal) / (farVal - nearVal)
	return m
}

// Returns a perspective projection matrix
func Perspective(fovy, aspect, zNear, zFar gl.Float) *Mat4 {
	f := 1 / (TanGL(fovy / 2.0))
	m := IdentMat4()
	m[0].X = f / aspect
	m[1].Y = f
	m[2].Z = (zFar + zNear) / (zNear - zFar)
	m[3].W = 0
	m[2].W = -1
	m[3].Z = (2 * zFar * zNear) / (zNear - zFar)
	return m
}

// Returns a frustum
func Frustum(left, right, bottom, top, near, far gl.Float) *Mat4 {

	m := IdentMat4()
	if (right == left) || (top == bottom) || (near == far) || (near < 0.0) || (far < 0.0) {
		fmt.Fprintf(os.Stderr, "Frustum error: Returning identity\n")
		return m
	}

	m[0].X = (2.0 * near) / (right - left)
	m[1].Y = (2.0 * near) / (top - bottom)

	m[2].X = (right + left) / (right - left)
	m[2].Y = (top + bottom) / (top - bottom)
	m[2].Z = -(far + near) / (far - near)
	m[2].W = -1.0

	m[3].Z = -(2.0 * far * near) / (far - near)
	m[2].W = 0.0

	return m
}

// ************************************ //
// *     OpenGL utility functions     * //
// ************************************ //

// Take two gl.Floats and return remainder as a gl.Float
func ModGL(a, b gl.Float) gl.Float {
	return (gl.Float)(math.Mod((float64)(a), (float64)(b)))
}

// Basic linear interpolation
func LerpGL(start, end, ratio gl.Float) gl.Float {
	return start + (end-start)*ratio
}

// Cosine, in gl.Float
func CosGL(Rad gl.Float) gl.Float {
	return (gl.Float)(math.Cos((float64)(Rad)))
}

// Sine, in gl.Float
func SinGL(Rad gl.Float) gl.Float {
	return (gl.Float)(math.Sin((float64)(Rad)))
}

// Tan, in gl.Float
func TanGL(Rad gl.Float) gl.Float {
	return (gl.Float)(math.Tan((float64)(Rad)))
}

// Identity matrix, bare
func Ident4() []gl.Float {
	return []gl.Float{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 1.0,
	}
}

// Convert degrees to radians
func DegToRad(fAngDeg gl.Float) gl.Float {
	return fAngDeg * degToRad
}

// Clamp - constrain a value fValue to the range delimited by
// fMinValue -> fMaxValue
func Clamp(fValue, fMinValue, fMaxValue gl.Float) gl.Float {
	if fValue < fMinValue {
		return fMinValue
	} else if fValue > fMaxValue {
		return fMaxValue
	} else {
		return fValue
	}
}

// Pretty-print a []gl.Float slice representing
// a 16-item transformation matrix.
func DebugMat(m []gl.Float, s string) {
	fmt.Fprintf(debugOut, "\t-----------------------%s-------------------------\n", s)
	for i := 0; i < 4; i++ {
		fmt.Fprintf(debugOut, "\t%f\t%f\t%f\t%f\n", m[i*4], m[i*4+1], m[i*4+2], m[i*4+3])
	}
	fmt.Fprintf(debugOut, "\t--------------------------------------------------------\n")
}

// Code from the MESA library, adapted for Go
func Invert(m []gl.Float) ([]gl.Float, error) {

	//double inv[16], det;
	//int i;
	inv := make([]gl.Float, 16)
	invOut := make([]gl.Float, 16)
	if len(m) != 16 {
		return nil, errors.New("Not a 4x4 matrix, needs 16 elements")
	}

	inv[0] = m[5]*m[10]*m[15] -
		m[5]*m[11]*m[14] -
		m[9]*m[6]*m[15] +
		m[9]*m[7]*m[14] +
		m[13]*m[6]*m[11] -
		m[13]*m[7]*m[10]

	inv[4] = -m[4]*m[10]*m[15] +
		m[4]*m[11]*m[14] +
		m[8]*m[6]*m[15] -
		m[8]*m[7]*m[14] -
		m[12]*m[6]*m[11] +
		m[12]*m[7]*m[10]

	inv[8] = m[4]*m[9]*m[15] -
		m[4]*m[11]*m[13] -
		m[8]*m[5]*m[15] +
		m[8]*m[7]*m[13] +
		m[12]*m[5]*m[11] -
		m[12]*m[7]*m[9]

	inv[12] = -m[4]*m[9]*m[14] +
		m[4]*m[10]*m[13] +
		m[8]*m[5]*m[14] -
		m[8]*m[6]*m[13] -
		m[12]*m[5]*m[10] +
		m[12]*m[6]*m[9]

	inv[1] = -m[1]*m[10]*m[15] +
		m[1]*m[11]*m[14] +
		m[9]*m[2]*m[15] -
		m[9]*m[3]*m[14] -
		m[13]*m[2]*m[11] +
		m[13]*m[3]*m[10]

	inv[5] = m[0]*m[10]*m[15] -
		m[0]*m[11]*m[14] -
		m[8]*m[2]*m[15] +
		m[8]*m[3]*m[14] +
		m[12]*m[2]*m[11] -
		m[12]*m[3]*m[10]

	inv[9] = -m[0]*m[9]*m[15] +
		m[0]*m[11]*m[13] +
		m[8]*m[1]*m[15] -
		m[8]*m[3]*m[13] -
		m[12]*m[1]*m[11] +
		m[12]*m[3]*m[9]

	inv[13] = m[0]*m[9]*m[14] -
		m[0]*m[10]*m[13] -
		m[8]*m[1]*m[14] +
		m[8]*m[2]*m[13] +
		m[12]*m[1]*m[10] -
		m[12]*m[2]*m[9]

	inv[2] = m[1]*m[6]*m[15] -
		m[1]*m[7]*m[14] -
		m[5]*m[2]*m[15] +
		m[5]*m[3]*m[14] +
		m[13]*m[2]*m[7] -
		m[13]*m[3]*m[6]

	inv[6] = -m[0]*m[6]*m[15] +
		m[0]*m[7]*m[14] +
		m[4]*m[2]*m[15] -
		m[4]*m[3]*m[14] -
		m[12]*m[2]*m[7] +
		m[12]*m[3]*m[6]

	inv[10] = m[0]*m[5]*m[15] -
		m[0]*m[7]*m[13] -
		m[4]*m[1]*m[15] +
		m[4]*m[3]*m[13] +
		m[12]*m[1]*m[7] -
		m[12]*m[3]*m[5]

	inv[14] = -m[0]*m[5]*m[14] +
		m[0]*m[6]*m[13] +
		m[4]*m[1]*m[14] -
		m[4]*m[2]*m[13] -
		m[12]*m[1]*m[6] +
		m[12]*m[2]*m[5]

	inv[3] = -m[1]*m[6]*m[11] +
		m[1]*m[7]*m[10] +
		m[5]*m[2]*m[11] -
		m[5]*m[3]*m[10] -
		m[9]*m[2]*m[7] +
		m[9]*m[3]*m[6]

	inv[7] = m[0]*m[6]*m[11] -
		m[0]*m[7]*m[10] -
		m[4]*m[2]*m[11] +
		m[4]*m[3]*m[10] +
		m[8]*m[2]*m[7] -
		m[8]*m[3]*m[6]

	inv[11] = -m[0]*m[5]*m[11] +
		m[0]*m[7]*m[9] +
		m[4]*m[1]*m[11] -
		m[4]*m[3]*m[9] -
		m[8]*m[1]*m[7] +
		m[8]*m[3]*m[5]

	inv[15] = m[0]*m[5]*m[10] -
		m[0]*m[6]*m[9] -
		m[4]*m[1]*m[10] +
		m[4]*m[2]*m[9] +
		m[8]*m[1]*m[6] -
		m[8]*m[2]*m[5]

	det := m[0]*inv[0] + m[1]*inv[4] + m[2]*inv[8] + m[3]*inv[12]

	if det == 0 {
		return nil, errors.New("No inverse for this matrix!")
	}

	det = 1.0 / det

	for i := 0; i < 16; i++ {
		invOut[i] = inv[i] * det
	}

	return invOut, nil
}
