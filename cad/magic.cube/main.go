package main

import (
	"github.com/ljanyst/ghostscad/lib/shapes"
	"github.com/ljanyst/ghostscad/sys"

	. "github.com/go-gl/mathgl/mgl64"
	. "github.com/jurgen-kluft/go-home/cad/lib"
	. "github.com/ljanyst/ghostscad/primitive"
)

// Conventions:
// For the final STL output, the X-Z plane is the ground plane:
//     - X axis is width
//     - Y axis is height
//     - Z axis is length/depth

// However, when designing the parts, we will use OpenSCAD conventions
// where the X-Y plane is the ground plane:
//     - X axis is width
//     - Y axis is length/depth
//     - Z axis is height

// The Magic Cube has to hold the following inside:
// - Battery
// - Firebeetle ESP32
// - 6 axis IMU sensor (e.g. MPU6050)

// Cube dimensions
const (
	CubeDimension = 75.0
)

// Firebeetle DFRobot ESP32-C6 v1.0
const (
	FirebeetleLength    = 60
	FirebeetleLengthH2H = 56.6
	FirebeetleWidth     = 25.4
	FirebeetleWidthH2H  = 22.0
)

// 900mAh battery dimensions
const (
	batteryWireGap  = 2.0
	batteryDiameter = 16.6 + batteryWireGap
	batteryLength   = 34.0 + 6.0 // 34.0
	batteryRadius   = batteryDiameter / 2.0
)

// Tunnel out of the battery holder for wires to exit
const (
	wireTunnelLength = 20.0
	wireTunnelRadius = 1.2
)

// newMagicCube creates the power block primitive that holds the USB power block.
// The power block has cutouts on the top and bottom for heat dissipation
func newMagicCube() Primitive {
	return NewTranslation(
		Vec3{0, 0, CubeDimension / 2},
		NewUnion(
			NewDifference(
				shapes.NewSmoothedCube(
					Vec3{CubeDimension, CubeDimension, CubeDimension},
					PowerBlockRounding,
				).Build(),

				// Emptying + Top cutout(for inserting the power block)
				NewTranslation(
					Vec3{0, 0, W1},
					NewBox(CubeDimension-2*W2, CubeDimension-2*W2, CubeDimension+W1),
				),
			),
		),
	)
}

func newMagicCubeLid() Primitive {
	// cylindricalHeatHole := NewTranslation(Vec3{0, 0, 0},
	// 	NewCylinder(3*W2, W2*0.6),
	// )

	// holeOffsetX1 := CubeDimension/2 - 2*W2
	// holeOffsetX2 := CubeDimension/2 - 3*W2
	// holeOffsetZ := (CubeDimension / 4)

	return NewDifference(
		NewUnion(
			NewTranslation(
				Vec3{0, 0, W1 / 2},
				NewCube(Vec3{CubeDimension, CubeDimension, W1}),
			),
			NewTranslation(
				Vec3{0, 0, W2 / 2},
				NewCube(Vec3{CubeDimension - 2*W2 - 0.2, (CubeDimension - 2*W2), W2}),
			),
		),

		// Cut some holes for heat dissipation
		// NewTranslation(Vec3{0, 0 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{0, 8*W2 - holeOffsetZ, 0}, cylindricalHeatHole),

		// NewTranslation(Vec3{holeOffsetX1, 0 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{-holeOffsetX1, 0 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{holeOffsetX2, 2*W2 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{-holeOffsetX2, 2*W2 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{holeOffsetX1, 4*W2 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{-holeOffsetX1, 4*W2 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{holeOffsetX2, 6*W2 - holeOffsetZ, 0}, cylindricalHeatHole),
		// NewTranslation(Vec3{-holeOffsetX2, 6*W2 - holeOffsetZ, 0}, cylindricalHeatHole),
	)
}

// The lid has to be wider than the holder to fit over it so that we can
// close it properly.
func newBatteryHolderLid() Primitive {
	return NewTranslation(
		Vec3{0, 0, 1.5 * W1},
		NewUnion(
			NewDifference(
				NewCylinder(3*W1, batteryRadius),
				NewCylinder(4*W1, batteryRadius-W1),
			),
			// Another disc attached
			NewTranslation(Vec3{0, 0, -W1},
				NewCylinder(W1, batteryRadius+W1),
			),
		),
	)
}

func newBatteryHolder() Primitive {
	h := batteryLength + W2
	return NewTranslation(
		Vec3{0, 0, h / 2},
		NewDifference(
			NewUnion(
				NewCylinder(h, batteryRadius+W1),
				NewTranslation(
					Vec3{-W1 / 2, batteryRadius, 0},
					NewBox(batteryDiameter+W2+W1, batteryDiameter, h),
				),
			),
			// Make the cylinder hollow, W1 wall thickness
			NewTranslation(Vec3{0, 0, W1},
				NewCylinder(batteryLength+W2, batteryRadius),
			),

			// Wire tunnel top
			NewTranslation(Vec3{0, batteryRadius, h/2 - wireTunnelRadius - W1},
				NewCylinderOnTheYAxis(wireTunnelLength, wireTunnelRadius),
			),

			// Wire tunnel bottom
			NewTranslation(Vec3{0, batteryRadius, -h/2 + wireTunnelRadius + W1},
				NewCylinderOnTheYAxis(wireTunnelLength, wireTunnelRadius),
			),
		),
	)
}

func newCombineAll() Primitive {
	// Return a setup of all the parts together for review
	// Move them apart for better viewing
	apart := 40.0
	return NewList(
		NewTranslation(
			Vec3{-apart, 0, 0},
			newMagicCube(),
		),
		NewTranslation(
			Vec3{apart, 0, 0},
			newMagicCubeLid(),
		),
		NewTranslation(Vec3{-3 * apart, 0, 0}, newBatteryHolder()),
		NewTranslation(Vec3{3 * apart, 0, 0}, newBatteryHolderLid()),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "magic_cube", Primitive: newMagicCube(), Flags: sys.Default},
		{Name: "magic_cube_lid", Primitive: newMagicCubeLid(), Flags: sys.Default},
		{Name: "for-review-only", Primitive: newCombineAll(), Flags: sys.None},
	})
}
