package goglutils

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

// X-rotates the topmost matrix on the stack
func (ms *MatrixStack) RotateX(deg gl.Float) {
	ms.currMat = ms.currMat.MulM(RotateX(deg))
}

// Y-rotates the topmost matrix on the stack
func (ms *MatrixStack) RotateY(deg gl.Float) {
	ms.currMat = ms.currMat.MulM(RotateY(deg))
}

// Z-rotates the topmost matrix on the stack
func (ms *MatrixStack) RotateZ(deg gl.Float) {
	ms.currMat = ms.currMat.MulM(RotateZ(deg))
}

// Scales the topmost matrix on the stack
func (ms *MatrixStack) Scale(s *Vec4) {
	ms.currMat = ms.currMat.Scale(s)
}

// Translates the topmost matrix on the stack
func (ms *MatrixStack) Translate(offset *Vec4) {
	ms.currMat = ms.currMat.Translate(offset)
}

// Inverts the topmost matrix on the stack
func (ms *MatrixStack) Invert() {
	ms.currMat = ms.currMat.Inverse()
}

// Multiplies by another mat4 the topmost matrix on the stack
func (ms *MatrixStack) MulM(m *Mat4) {
	ms.currMat = ms.currMat.MulM(m)
}

func (ms *MatrixStack) Ortho(left, right, bottom, top, nearVal, farVal gl.Float) {
	ms.currMat = ms.currMat.MulM(Ortho(left, right, bottom, top, nearVal, farVal))
}

func (ms *MatrixStack) Perspective(fov, aspect, zNear, zFar gl.Float) {
	ms.currMat = ms.currMat.MulM(Perspective(fov, aspect, zNear, zFar))
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
