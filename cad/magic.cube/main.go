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
// Firebeetle DFRobot ESP32-C6 v1.0
const (
	FirebeetleLength    = 60.0
	FirebeetleLengthH2H = 56.6
	FirebeetleWidth     = 25.4
	FirebeetleWidthH2H  = 22.0
)

const (
	CubeDimension = 3*FirebeetleWidth + 2*W2
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
				shapes.NewSmoothedCube(Vec3{CubeDimension, CubeDimension, CubeDimension}, PowerBlockRounding).Build(),

				// Emptying + Top cutout(for inserting the power block)
				NewTranslation(
					Vec3{0, 0, W2},
					NewBox(CubeDimension-2*W2, CubeDimension-2*W2, CubeDimension),
				),

				// 6 sides, drill holes:
				// at side 1 sketching the number 1
				// at side 2 sketching the number 2
				// at side 3 sketching the number 3
				// at side 4 sketching the number 4
				// at side 5 sketching the number 5
				// at side 6 sketching the number 6
			),

			// 4 pins for the Firebeetle board to be attached to the inside of the cube. The holes are 3mm in diameter and 4mm deep, and are located at the corners of a rectangle that is 56.6mm by 22mm (the dimensions of the Firebeetle board). The holes are centered on the Z axis and are located 10mm from the top of the cube.
			NewTranslation(
				Vec3{0, (-CubeDimension / 2) + W1*3, 0},
				NewUnion(
					NewTranslation(Vec3{FirebeetleWidthH2H / 2, 0, FirebeetleLengthH2H / 2},
						NewCylinderOnTheYAxis(W1*3, 1),
					),
					NewTranslation(Vec3{-FirebeetleWidthH2H / 2, 0, FirebeetleLengthH2H / 2},
						NewCylinderOnTheYAxis(W1*3, 1),
					),
					NewTranslation(Vec3{FirebeetleWidthH2H / 2, 0, -FirebeetleLengthH2H / 2},
						NewCylinderOnTheYAxis(W1*3, 1),
					),
					NewTranslation(Vec3{-FirebeetleWidthH2H / 2, 0, -FirebeetleLengthH2H / 2},
						NewCylinderOnTheYAxis(W1*3, 1),
					),
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
				Vec3{0, 0, W2},
				NewCube(Vec3{CubeDimension - 2*W2 - 0.2, (CubeDimension - 2*W2), 2 * W2}),
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
	h := FirebeetleLength
	return NewTranslation(
		Vec3{0, 0, W2 + h/2 - PowerBlockRounding},
		NewDifference(
			NewTranslation(
				Vec3{0, 0, 0},
				NewBox(CubeDimension, FirebeetleWidth, h),
			),

			// Make a hollow cylinder
			NewTranslation(Vec3{0, 0, (FirebeetleLength-batteryLength)/2 + W1},
				NewCylinder(batteryLength+W2, batteryRadius),
			),

			// Make 2 more hollow cylinders, mainly to reduce the amount of support material
			// needed when printing the battery holder.
			NewTranslation(Vec3{-(batteryRadius*2 + W1), 0, W1},
				NewCylinder(batteryLength, batteryRadius),
			),
			NewTranslation(Vec3{batteryRadius*2 + W1, 0, W1},
				NewCylinder(batteryLength, batteryRadius),
			),

			// Wire tunnel top
			NewTranslation(Vec3{0, batteryRadius, h/2 - W2},
				NewCylinderOnTheYAxis(wireTunnelLength, wireTunnelRadius),
			),

			// Wire tunnel bottom
			NewTranslation(Vec3{0, batteryRadius, h/2 - W2 - (batteryLength - 2*W2)},
				NewCylinderOnTheYAxis(wireTunnelLength, wireTunnelRadius),
			),
		),
	)
}

func newCombineAll() Primitive {
	// Return a setup of all the parts together for review
	// Move them apart for better viewing
	apart := 50.0
	return NewList(
		NewTranslation(
			Vec3{-apart, 0, 0},
			newMagicCube(),
		),
		NewTranslation(
			Vec3{apart, 0, 0},
			newMagicCubeLid(),
		),
		NewTranslation(Vec3{-apart, 0, 0}, newBatteryHolder()),
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
