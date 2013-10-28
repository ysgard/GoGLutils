package goglutils

import (
	"testing"
	"math"
)


// Test the Vec3 functions

func TestVec3Normalize(t *testing.T) {
	veclen := math.Sqrt(4.5 * 4.5 + 5.5 * 5.5 + 3.4 * 3.4)
	expected := 4.5 / veclen + 5.5 / veclen + 3.4 / veclen
	vi := NewVec3(4.5, 5.5, 3.4)
	vi.Normalize()
	out := float64(vi.X + vi.Y + vi.Z)
	if (out != expected) {
		t.Errorf("Normalize yields %v, want %v", out, expected)
	}
}

func TestVec3Cross(t *testing.T) {

}

func TestVec3Add(t *testing.T) {

}

func TestVec3Sub(t *testing.T) {

}

func TestVec3MulS(t *testing.T) {

}