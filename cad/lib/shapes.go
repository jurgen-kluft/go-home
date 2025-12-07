package cadlib

import (
	. "github.com/go-gl/mathgl/mgl64"
	"github.com/ljanyst/ghostscad/lib/shapes"
	. "github.com/ljanyst/ghostscad/primitive"
)

func NewBox(w, l, h float64) Primitive {
	return shapes.NewSmoothedCube(Vec3{w, l, h}, Rounding).Build()
}

func NewCylinderOnTheXAxis(h, r float64) Primitive {
	return NewRotation(Vec3{0, 90, 0}, NewCylinder(h, r))
}

func NewCylinderOnTheYAxis(h, r float64) Primitive {
	return NewRotation(Vec3{90, 0, 0}, NewCylinder(h, r))
}

func NewPyramid(w, l, h float64) Primitive {
	baseLT := Vec3{-w / 2, l / 2, -h / 2}
	baseRT := Vec3{w / 2, l / 2, -h / 2}
	baseLB := Vec3{-w / 2, -l / 2, -h / 2}
	baseRB := Vec3{w / 2, -l / 2, -h / 2}
	apex := Vec3{0, 0, h / 2}

	points := []Vec3{
		baseLT,
		baseRT,
		baseRB,
		baseLB,
		apex,
	}

	triangles := []Vec3{
		// Base
		{0, 1, 2},
		{0, 2, 3},
		// Sides
		{0, 1, 4},
		{1, 2, 4},
		{2, 3, 4},
		{3, 0, 4},
	}

	return NewRender(
		10,
		NewPolyhedron(points, triangles),
	)
}
