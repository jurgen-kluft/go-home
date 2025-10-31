package main

import (
	"github.com/ljanyst/ghostscad/lib/shapes"
	"github.com/ljanyst/ghostscad/sys"

	. "github.com/go-gl/mathgl/mgl64"
	. "github.com/ljanyst/ghostscad/primitive"
)

// Conventions (the X-Z plane is the ground plane):
// - X axis is width
// - Y axis is height
// - Z axis is length/depth

// Global constants
const (
	WT       = 1.5
	W1       = WT
	W2       = 2 * WT
	Rounding = 0.4
)

// Box dimensions
const (
	boxW = 60.0
	boxL = 80.0
	boxH = 25.0
)

// Box bottom
const (
	boxBottomH = 10.0
)

// Inlay magnets for holding the box bottom and top together
const (
	magnetR = 5.0
	magnetH = 2.0
)

// PCB Board
const (
	pcbBoardW = 16.0
	pcbBoardL = 16.0
	pcbBoardH = 15.0
)

// Scd41
const (
	Scd41W = 27.0
	Scd41L = 27.0
	Scd41H = 7.8
)

// Bh1750
const (
	Bh1750W = 20.0
	Bh1750L = 20.0
	Bh1750H = 7.0
)

// Bme280
const (
	Bme280W = 16.0
	Bme280L = 22.0
	Bme280H = 3.0
)

// USB-C hole dimensions
const (
	UsbCHoleR = 8
)

// The X-Z plane is the ground planes
func newBox(w, l, h float64) Primitive {
	return shapes.NewSmoothedCube(Vec3{w, l, h}, Rounding).Build()
}

// The bottom part of the box:
// - PCB inlay
// -
func newBoxBottom() Primitive {
	return newBox(boxW, boxBottomH, boxL)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "main", Primitive: newBoxBottom(), Flags: sys.Default},
	})
}
