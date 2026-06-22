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
// - Xiao Seeed Studio, ESP32-S3
// - 6 axis IMU sensor (e.g. MPU6050)

// Cube dimensions

// Firebeetle DFRobot ESP32-C6 v1.0
const (
	Firebeetle2Length    = 60.0
	Firebeetle2LengthH2H = 56.6
	Firebeetle2Width     = 25.4
	Firebeetle2WidthH2H  = 22.0
)

// Xiao Seeed Studio, ESP32-S3
const (
	XiaoEsp32S3Length       = 23.0
	XiaoEsp32S3Width        = 18.0
	XiaoEsp32S3PCBThickness = 1.2
)

// MicroController dimensions
const (
	MicroControllerLength       = XiaoEsp32S3Length
	MicroControllerWidth        = XiaoEsp32S3Width
	MicroControllerPCBThickness = XiaoEsp32S3PCBThickness
)

const (
	CubeDimension = 3*MicroControllerWidth + 2*W2
	// Pimple Profile Control
	PipRadius  = 3.0 // Width footprint of the pimple
	PipProfile = 1.2
)

// 900mAh battery dimensions
const (
	batteryWireGap  = 2.0
	batteryDiameter = 16.8 + batteryWireGap
	batteryLength   = 32.5 + 6.0 // 34.0
	batteryRadius   = batteryDiameter / 2.0
)

// Tunnel out of the battery holder for wires to exit
const (
	wireTunnelLength = 20.0
	wireTunnelRadius = 1.2
)

// getLowProfileDimples returns a union of cutting tools designed to carve pockets INWARD from Z = 0.
func getLowProfileDimples(number int) Primitive {
	pips := []Primitive{}
	s := CubeDimension * 0.22

	addDimple := func(x, y float64) {
		// A sphere positioned so its crown cuts exactly PipProfile into the wall.
		// It extends slightly into positive space (+1.0) to guarantee a clean breakthrough edge.
		dimpleCutter := NewTranslation(Vec3{0, 0, -PipRadius + PipProfile},
			NewUnion(
				NewSphere(PipRadius).SetFn(20),
				NewCylinder(PipRadius+20.0, 0.5).SetFn(20),
			),
		)
		pips = append(pips, NewTranslation(Vec3{x, y, 0}, dimpleCutter))
	}

	buildDiePattern(number, s, addDimple)
	return NewUnion(pips...)
}

// Helper to reuse the standard die face layout grid across both functions
func buildDiePattern(number int, s float64, addFunc func(x, y float64)) {
	switch number {
	case 1:
		addFunc(0, 0)
	case 2:
		addFunc(-s, -s)
		addFunc(s, s)
	case 3:
		addFunc(-s, -s)
		addFunc(0, 0)
		addFunc(s, s)
	case 4:
		addFunc(-s, -s)
		addFunc(-s, s)
		addFunc(s, -s)
		addFunc(s, s)
	case 5:
		addFunc(-s, -s)
		addFunc(-s, s)
		addFunc(0, 0)
		addFunc(s, -s)
		addFunc(s, s)
	case 6:
		addFunc(-s, -s)
		addFunc(-s, 0)
		addFunc(-s, s)
		addFunc(s, -s)
		addFunc(s, 0)
		addFunc(s, s)
	}
}

// newVerticalSliderRails creates left and right mounting blocks
// shifted safely inward toward the center of the box.
func newVerticalSliderRails(width, length float64) Primitive {
	railHeightZ := length // Height of the blocks along Z
	railWidthX := width   // Total thickness of each printed block along X
	railDepthY := width   // How far the blocks stick out from the back wall
	slotDepthX := 2.0     // How deep the PCB slides into the block
	printTolerance := 0.2 // Clearance spacing for a smooth slide fit

	slotThicknessY := MicroControllerPCBThickness + printTolerance

	// 1. Build the Left Rail Block
	// The track slot cuts into the inner right side of this block (facing +X)
	leftBlock := NewBox(railWidthX, railDepthY, railHeightZ)
	leftSlotCutter := NewTranslation(
		Vec3{railWidthX / 2.0, 0, 1.0},
		NewBox(2.0*slotDepthX+1.0, slotThicknessY, railHeightZ+2.0),
	)
	leftRailFinished := NewDifference(leftBlock, leftSlotCutter)

	// 2. Build the Right Rail Block
	// The track slot cuts into the inner left side of this block (facing -X)
	rightBlock := NewBox(railWidthX, railDepthY, railHeightZ)
	rightSlotCutter := NewTranslation(
		Vec3{-railWidthX / 2.0, 0, 1.0},
		NewBox(2.0*slotDepthX+1.0, slotThicknessY, railHeightZ+2.0),
	)
	rightRailFinished := NewDifference(rightBlock, rightSlotCutter)

	// 3. FIX: Position the rails using the internal slot distance rather than the outer span.
	// This pulls the boxes away from the cube's corners and places them directly around the board width.
	railMoveX := (MicroControllerWidth / 2.0) + (railWidthX/2.0 - slotDepthX)
	leftRailPlaced := NewTranslation(
		Vec3{-railMoveX, 0, 0},
		leftRailFinished,
	)
	rightRailPlaced := NewTranslation(
		Vec3{railMoveX, 0, 0},
		rightRailFinished,
	)

	return NewUnion(leftRailPlaced, rightRailPlaced)
}

// newMagicCube creates the power block primitive that holds the USB power block.
// The power block has cutouts on the top and bottom for heat dissipation
func newMagicCube() Primitive {
	offset := (CubeDimension / 2.0)

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

				// Subtracting dimple cutters from the faces
				NewTranslation(Vec3{0, -offset, 0}, NewRotation(Vec3{-90, 0, 0}, getLowProfileDimples(1))),
				NewTranslation(Vec3{0, offset, 0}, NewRotation(Vec3{90, 0, 0}, getLowProfileDimples(2))),
				NewTranslation(Vec3{-offset, 0, 0}, NewRotation(Vec3{0, 90, 0}, getLowProfileDimples(3))),
				NewTranslation(Vec3{offset, 0, 0}, NewRotation(Vec3{0, -90, 0}, getLowProfileDimples(4))),
				NewTranslation(Vec3{0, 0, -offset}, NewRotation(Vec3{0, 0, 0}, getLowProfileDimples(5))),
			),

			// Microcontroller left and right block with slider cutout to insert the microcontroller horizontally
			// from the top of the cube.
			// This is attached to the main cube with a small gap (W2) to allow for some tolerance when sliding the
			// microcontroller in.
			NewTranslation(
				Vec3{0, CubeDimension/2 - W2 - (8.0 / 2), -(((CubeDimension / 2.0) - W2 + (Rounding / 2.0)) - MicroControllerLength/2.0)},
				newVerticalSliderRails(8.0, MicroControllerLength),
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

	// The lid sits flat on the XY plane.
	// The outer top surface of the main lid plate is at Z = W1.

	outerSurfaceHeight := W1

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

		// Translate Side 6 dimple cutters onto the top surface
		// No rotation needed: getLowProfileDimples cuts downward into the surface along the -Z axis natively
		NewTranslation(Vec3{0, 0, outerSurfaceHeight}, getLowProfileDimples(6)),
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
	h := batteryLength
	return NewTranslation(
		Vec3{0, 0, W2 + h/2 - PowerBlockRounding},
		NewDifference(
			NewTranslation(
				Vec3{0, 0, 0},
				NewBox(CubeDimension, batteryDiameter+W2, h),
			),

			// Make a hollow cylinder
			NewTranslation(Vec3{0, 0, W1},
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
			NewTranslation(Vec3{0, -batteryRadius, h/2 - W2},
				NewCylinderOnTheYAxis(wireTunnelLength, wireTunnelRadius),
			),

			// Wire tunnel bottom
			NewTranslation(Vec3{0, batteryRadius, h/2 - W2 - (batteryLength - 2*W2)},
				NewCylinderOnTheYAxis(wireTunnelLength, wireTunnelRadius),
			),
			NewTranslation(Vec3{0, -batteryRadius, h/2 - W2 - (batteryLength - 2*W2)},
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
