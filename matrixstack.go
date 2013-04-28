// Matrix stack struct and supporting functions
package go-glutils

import (
	gl "github.com/chsc/gogl/gl33"
)

// MatrixStack - Represents a way to store a sequential series of
// transformations
type MatrixStack struct {
	currMat  *Mat4
	matrices []*Mat4
}

// Creates a default identity matrix as the current matrix
func (ms *MatrixStack) Init() {
	ms.currMat = IdentMat4()
}

// Return pointer to top matrix
func (ms *MatrixStack) Top() *Mat4 {
	return ms.currMat
}

func (ms *MatrixStack) RotateX(deg gl.Float) {
	ms.currMat = ms.currMat.MulM(RotateX(deg))
}

func (ms *MatrixStack) RotateY(deg gl.Float) {
	ms.currMat = ms.currMat.MulM(RotateY(deg))
}

func (ms *MatrixStack) RotateZ(deg gl.Float) {
	ms.currMat = ms.currMat.MulM(RotateZ(deg))
}

func (ms *MatrixStack) Scale(s *Vec4) {
	ms.currMat = ms.currMat.Scale(s)
}

func (ms *MatrixStack) Translate(offset *Vec4) {
	ms.currMat = ms.currMat.Translate(offset)
}

func (ms *MatrixStack) Invert() {
	ms.currMat = ms.currMat.Inverse()
}

func (ms *MatrixStack) MulM(m *Mat4) {
	ms.currMat = ms.currMat.MulM(m)
}

// Create a copy of the current matrix and push
// it onto the stack
func (ms *MatrixStack) Push() {
	copied := ms.currMat.Copy()
	ms.matrices = append(ms.matrices, copied)
}

// Pop the last matrix off the stack and make
// it the current matrix
func (ms *MatrixStack) Pop() {
	if len(ms.matrices) > 0 {
		ms.currMat = ms.matrices[len(ms.matrices)-1]
		ms.matrices = ms.matrices[:len(ms.matrices)-1]
	}
}
