package main

import (
	"github.com/ljanyst/ghostscad/lib/shapes"
	"github.com/ljanyst/ghostscad/sys"

	. "github.com/go-gl/mathgl/mgl64"
	. "github.com/ljanyst/ghostscad/primitive"
)

// History:
// v1.0: Initial version
// v1.1: Adjusted battery holder dimensions to fit better
//       moved battery holder next to sensor box
//       increased sensor box height to fit the ESP8266 board better

// Global constants
const (
	wallThickness = 1.5
	rounding      = 0.4
)

// Cylindrical battery dimensions are 16mm diameter, 34mm length
const (
	batteryWireGap  = 2.0
	batteryDiameter = 16.6 + batteryWireGap
	batteryLength   = 34.0 + 6.0 // 34.0
	batteryRadius   = batteryDiameter / 2.0
	W1              = wallThickness
	W2              = 2 * wallThickness
)

// Switch PCB board + Esp8266
const (
	switchBoardLength = 40.0
	switchBoardWidth  = 21.0
	switchBoardHeight = 15.0 // Including the battery socket and connection
)

// Tunnel between battery holder and sensor box
const (
	tunnelLength = 20.0
	tunnelRadius = 1.2
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
				NewCylinder(W1, batteryRadius+W1),
			),
		),
	)
}

func newBatteryHolder() Primitive {
	h := batteryLength + W2
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewUnion(
				NewCylinder(h, batteryRadius+W1),
				NewTranslation(
					Vec3{-W1 / 2, batteryRadius, 0},
					NewBox(batteryDiameter+W2+W1, switchBoardHeight, h),
				),
			),
			// Make the cylinder hollow, W1 wall thickness
			NewTranslation(Vec3{0, 0, W1},
				NewCylinder(batteryLength+W2, batteryRadius),
			),
		))
}

func newTunnelOnXAxis(l, r float64) Primitive {
	return NewRotation(
		Vec3{0, 90, 0},
		NewCylinder(l, r),
	)
}

func newTunnelOnYAxis(l, r float64) Primitive {
	return NewRotation(
		Vec3{90, 0, 0},
		NewCylinder(l, r),
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
				//NewBox(switchBoardWidth, switchBoardHeight+W2, switchBoardLength-W2),
				NewBox(switchBoardWidth, switchBoardHeight+W2, switchBoardLength),
			),
		),
	)
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

func productSensorBox() Primitive {
	tunnelLength := 4 * W2
	return NewRotation(
		Vec3{0, 0, 180},
		NewDifference(
			NewUnion(
				NewTranslation(
					Vec3{0, -W1, 0},
					newBatteryHolder(),
				),
				NewTranslation(
					//Vec3{0, batteryRadius - 2*W1, -((switchBoardLength+W2)/2 + (batteryLength+W2)/2 + rounding)},
					Vec3{-(switchBoardWidth/2 + batteryRadius + W1 + W1), batteryRadius - 2*W1, 0},
					newSensorBox(),
				),
			),
			// Two wire tunnels for the battery to the ESP8266 board
			NewTranslation(
				Vec3{-(batteryRadius/2 + tunnelLength/2), 0, batteryLength/2 - W1},
				newTunnelOnXAxis(tunnelLength, tunnelRadius),
			),
			NewTranslation(
				Vec3{-(batteryRadius/2 + tunnelLength/2), 0, -(batteryLength/2 - W1)},
				newTunnelOnXAxis(tunnelLength, tunnelRadius),
			),
		),
	)
}

func productBatteryHolderLid() Primitive {
	return NewList(
		NewTranslation(
			Vec3{-30, W1 + W1/2, 0},
			NewRotation(
				Vec3{-90, 0, 0},
				newBatteryHolderLid(),
			),
		),
	)
}

func productSensorBoxLid() Primitive {
	return NewList(
		NewTranslation(
			Vec3{30, W1 + rounding, 0},
			NewRotation(
				Vec3{-90, 0, 0},
				newSensorBoxLid(),
			),
		),
	)
}

func newCombineAll() Primitive {
	// Return a setup of all the parts together for review
	// Move them apart for better viewing
	apart := 20.0
	return NewList(
		NewTranslation(
			Vec3{-apart, 0, 0},
			productBatteryHolderLid(),
		),
		NewTranslation(
			Vec3{apart, 0, 0},
			productSensorBoxLid(),
		),
		NewTranslation(
			Vec3{0, 0, 0},
			productSensorBox(),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "sensorbox", Primitive: productSensorBox(), Flags: sys.Default},
		{Name: "sensorbox lid", Primitive: productSensorBoxLid(), Flags: sys.None},
		{Name: "batteryholder lid", Primitive: productBatteryHolderLid(), Flags: sys.None},
		{Name: "for-review-only", Primitive: newCombineAll(), Flags: sys.None},
	})
}
