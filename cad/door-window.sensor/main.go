package main

import (
	"github.com/ljanyst/ghostscad/lib/shapes"
	"github.com/ljanyst/ghostscad/sys"

	. "github.com/go-gl/mathgl/mgl64"
	. "github.com/ljanyst/ghostscad/primitive"
)

// Global constants
const (
	wallThickness = 1.5
	rounding      = 0.4
)

// Cylindrical battery dimensions are 16mm diameter, 34mm length
const (
	batteryDiameter = 16.0
	batteryLength   = 34.0
	batteryRadius   = batteryDiameter / 2.0
	W1              = wallThickness
	W2              = 2 * wallThickness
)

// Switch PCB board + Esp8266
const (
	switchBoardLength = 40.0
	switchBoardWidth  = 21.0
	switchBoardHeight = 8.0
)

// Tunnel between battery holder and sensor box
const (
	tunnelLength = 20.0
	tunnelRadius = 1.6
)

// The X-Z plane is the ground planes

func NewBox(w, l, h float64) Primitive {
	return shapes.NewSmoothedCube(Vec3{w, l, h}, rounding).Build()
}

// The lid has to be wider than the holder to fit over it so that we can
// close it properly.
func newBatteryHolderLid() Primitive {
	return NewRotation(
		Vec3{0, 180, 0},
		NewUnion(
			NewDifference(
				NewCylinder(3*W1, batteryRadius),
				NewCylinder(4*W1, batteryRadius-W1),
			),
			// Another disc attached
			NewTranslation(Vec3{0, 0, W1},
				NewCylinder(W1, batteryRadius+W2),
			),
		),
	)
}

func newBatteryHolder() Primitive {
	h := batteryLength + W2 + rounding*2
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewUnion(
				NewCylinder(h, batteryRadius+W1),
				NewTranslation(
					Vec3{0, batteryRadius, 0},
					NewBox(batteryDiameter, switchBoardHeight, h),
				),
			),
			// Make the cylinder hollow, 1mm wall thickness
			NewTranslation(Vec3{0, 0, W1},
				NewCylinder(batteryLength+W2, batteryRadius),
			),
		))
}

func newWireTunnel() Primitive {
	return NewTranslation(
		Vec3{0, batteryRadius / 2, -batteryLength / 2},
		NewCylinder(tunnelLength, tunnelRadius),
	)

}

func newSensorBox() Primitive {
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewDifference(
				NewBox(switchBoardWidth+W2, switchBoardHeight+W2, switchBoardLength+W2),
				NewBox(switchBoardWidth, switchBoardHeight, switchBoardLength),
			),
			NewTranslation(
				Vec3{0, -W1, 0},
				NewBox(switchBoardWidth, switchBoardHeight+W2, switchBoardLength-W2),
			),
		))
}

func newSensorBoxLid() Primitive {
	return NewUnion(
		NewDifference(
			NewBox(switchBoardWidth, switchBoardLength, W2),
			NewBox(switchBoardWidth-W2, switchBoardLength-W2, 2*W2),
		),
		NewTranslation(
			Vec3{0, 0, -W1 + rounding},
			NewBox(switchBoardWidth+W2, switchBoardLength+W2, W1),
		))
}

func newCombineAll() Primitive {
	return NewRotation(
		Vec3{0, 0, 180},
		NewDifference(
			NewUnion(
				NewTranslation(
					Vec3{0, -W1, 0},
					newBatteryHolder(),
				),
				NewTranslation(
					Vec3{0, batteryRadius - 2*W1, -((switchBoardLength+W2)/2 + batteryLength/2 + (W1 - rounding))},
					newSensorBox(),
				),
			),
			newWireTunnel(),
		))
}

func newProduct() Primitive {
	return NewList(
		NewTranslation(
			Vec3{0, batteryRadius + W1, batteryLength / 2},
			newCombineAll(),
		),
		NewTranslation(
			Vec3{-30, W1 + W1/2, 0},
			NewRotation(
				Vec3{-90, 0, 0},
				newBatteryHolderLid(),
			),
		),
		NewTranslation(
			Vec3{30, W1 + rounding, 0},
			NewRotation(
				Vec3{-90, 0, 0},
				newSensorBoxLid(),
			),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "main", Primitive: newProduct(), Flags: sys.Default},
	})
}
