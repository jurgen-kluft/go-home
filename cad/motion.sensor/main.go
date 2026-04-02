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

// We embed a USB power block into the bottom of the box, so we need to
// account for its dimensions when designing the box. We call this the
// power box. The power block needs to be flush so that the box can sit flat
// on the socket.

// On top of the power box we make a circular hole where we can insert a
// cylindrical part connected to a rectangular/ part that holds the sensor.
// This cylindrical+rectangular part is separated from the power box, and
// can be inserted into the power box through a circular slot on the top. The
// wiring from the sensor can go through the cylindrical part into the power box.

// The power box also holds space for the ESP32-C3 and the external antenna.

// Power box dimensions
const (
	PowerBoxWidth  = PowerBlockWidth + 2*W1
	PowerBoxLength = PowerBlockLength + 18.0
	PowerBoxHeight = PowerBlockHeight * 2
)

// newPowerBlock creates the power block primitive that holds the USB power block.
// The power block has cutouts on the top and bottom for heat dissipation
func newPowerBlock() Primitive {
	cylindricalHeatHole := NewCylinder(PowerBlockHeight+2*W2, W2*0.8)

	return NewUnion(
		NewDifference(
			shapes.NewSmoothedCube(
				Vec3{PowerBlockWidth + 2*W1, PowerBlockLength + 2*W1, PowerBlockHeight + 2*W1},
				PowerBlockRounding,
			).Build(),

			// Emptying + Front cutout(for inserting the power block)
			NewTranslation(
				Vec3{0, W1, 0},
				NewBox(PowerBlockWidth, PowerBlockLength+W1, PowerBlockHeight),
			),

			// USB-A cutout
			NewTranslation(
				Vec3{0, -((PowerBlockLength / 2) + (USBALength / 2) - W1), 0},
				NewBox(USBAWidth, USBALength, PowerBlockHeight-2*W1),
			),

			// Heat dissipation holes on the sensor box, especially around the power block area
			NewTranslation(Vec3{0, (PowerBlockLength / 6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{0, 5*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),

			NewTranslation(Vec3{(PowerBlockWidth/2 - 3*W2), 2*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{-(PowerBlockWidth/2 - 3*W2), 2*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{(PowerBlockWidth/2 - 2*W2), 3*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{-(PowerBlockWidth/2 - 2*W2), 3*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{(PowerBlockWidth/2 - 3*W2), 4*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{-(PowerBlockWidth/2 - 3*W2), 4*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{(PowerBlockWidth/2 - 2*W2), 5*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
			NewTranslation(Vec3{-(PowerBlockWidth/2 - 2*W2), 5*(PowerBlockLength/6) - (PowerBlockLength / 2), 0}, cylindricalHeatHole),
		),
	)
}

// newPowerBox creates the power box primitive that holds the power block, antenna
// and has cutouts for the power block, and sensor.
func newPowerBox() Primitive {
	return NewDifference(
		NewTranslation(
			Vec3{0, 0, PowerBoxHeight / 2},
			NewDifference(
				// The box
				NewBox(PowerBoxWidth, PowerBoxLength, PowerBoxHeight),
				// Make the box hollow
				NewBox(PowerBoxWidth-2*W1, PowerBoxLength-2*W1, PowerBoxHeight-2*W1),
			),
		),

		// Power block cutout
		NewTranslation(
			Vec3{0, ((PowerBoxLength / 2) - (PowerBlockLength / 2)) + W1, PowerBlockHeight/2 + W1},
			NewBox(PowerBlockWidth, PowerBlockLength, PowerBlockHeight),
		),

		// Backside cover cutout
		NewTranslation(
			Vec3{0, -((PowerBoxLength / 2) - W1), PowerBoxHeight / 2},
			//NewBox(PowerBoxWidth-2*W2, 2*W1, (PowerBoxHeight-2*W2)),
			NewCube(Vec3{PowerBoxWidth - 2*W2, 2 * W1, (PowerBoxHeight - 2*W2)}),
		),

		// Circular Sensor cutout at the top
		NewTranslation(
			Vec3{0, -(PowerBoxLength / 2) + 3*SensorStickCylindricalRadius, PowerBoxHeight - 4*W1},
			NewCylinder(8*W1, SensorStickCylindricalRadius+W1),
		),
		// Circular Sensor cutout at the top
		NewTranslation(
			Vec3{0, -(PowerBoxLength / 2) + 8*SensorStickCylindricalRadius, PowerBoxHeight - 4*W1},
			NewCylinder(8*W1, SensorStickCylindricalRadius+W1),
		),
	)
}

func newSensorBoxLid() Primitive {
	cylindricalHeatHole := NewRotation(Vec3{90, 0, 0},
		NewTranslation(Vec3{0, 0, 0},
			NewCylinder(3*W2, W2*0.6),
		),
	)
	holeOffsetX1 := PowerBlockWidth/2 - 2*W2
	holeOffsetX2 := PowerBlockWidth/2 - 3*W2
	holeOffsetZ := (PowerBoxHeight / 4)

	return NewRotation(
		Vec3{90, 0, 0},
		NewDifference(
			NewUnion(
				NewTranslation(
					Vec3{0, W1 / 2, 0},
					NewCube(Vec3{PowerBoxWidth - W1, W1, PowerBoxHeight - W1}),
				),
				NewTranslation(
					Vec3{0, W2 / 2, 0},
					NewCube(Vec3{PowerBoxWidth - 2*W2 - 0.2, W2, (PowerBoxHeight - 2*W2)}),
				),
			),

			// Cut some holes for heat dissipation
			NewTranslation(Vec3{0, 0, 0 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{0, 0, 8*W2 - holeOffsetZ}, cylindricalHeatHole),

			NewTranslation(Vec3{holeOffsetX1, 0, 0 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{-holeOffsetX1, 0, 0 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{holeOffsetX2, 0, 2*W2 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{-holeOffsetX2, 0, 2*W2 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{holeOffsetX1, 0, 4*W2 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{-holeOffsetX1, 0, 4*W2 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{holeOffsetX2, 0, 6*W2 - holeOffsetZ}, cylindricalHeatHole),
			NewTranslation(Vec3{-holeOffsetX2, 0, 6*W2 - holeOffsetZ}, cylindricalHeatHole),
		),
	)
}

// newSensorStick creates the cylindrical tube with the rectangular (hollow) part on top where the sensor is located.
func newSensorStick() Primitive {

	cylindricalLength := 10.0
	stickHeight := 2*SensorStickCylindricalRadius + W2

	return NewTranslation(
		Vec3{0, 0, stickHeight / 2.0},
		NewUnion(
			NewDifference(
				NewBox(RD03DW+2*W1, RD03DL+2*W1, stickHeight),
				// Make it hollow
				NewCube(Vec3{RD03DW, RD03DL, stickHeight - 2*W1}),
				// Cutout the top, so that we can insert the sensor easily
				NewTranslation(
					Vec3{0, 0, stickHeight / 2},
					NewCube(Vec3{RD03DW, RD03DL, 2 * W2}),
				),
				// Cutout for the cylindrical part
				NewTranslation(
					Vec3{0, (RD03DL + W1) / 2, 0},
					NewCylinderOnTheYAxis(2*W2, SensorStickCylindricalRadius),
				),
			),

			// Cylindrical part
			NewTranslation(
				Vec3{0, RD03DL/2.0 + cylindricalLength/2.0 + W1, 0},
				NewDifference(
					NewCylinderOnTheYAxis(cylindricalLength+W1, SensorStickCylindricalRadius+W1),
					NewCylinderOnTheYAxis(RD03DL, SensorStickCylindricalRadius),
				),
			),
		),
	)
}

func newSensorStickLid() Primitive {
	cylindricalHeatHole := NewCylinder(10, 2)

	return NewTranslation(
		Vec3{0, 0, W1 / 2.0},
		NewDifference(
			NewUnion(
				NewTranslation(
					Vec3{0, 0, W1 / 2},
					NewBox(RD03DW+2*W1, RD03DL+2*W1, W1),
				),
				NewTranslation(
					Vec3{0, 0, W2/2 + Rounding/2},
					NewDifference(
						NewBox(RD03DW-0.2, RD03DL, W2),
						NewBox(RD03DW-0.2-W2, RD03DL-W2, 2*W2),
					),
				),
			),
			NewTranslation(Vec3{0, 0, 0}, cylindricalHeatHole),
			NewTranslation(Vec3{0, 2 * RD03DL / 6, 0}, cylindricalHeatHole),
			NewTranslation(Vec3{0, -2 * RD03DL / 6, 0}, cylindricalHeatHole),

			NewTranslation(Vec3{W2, RD03DL / 6, 0}, cylindricalHeatHole),
			NewTranslation(Vec3{-W2, RD03DL / 6, 0}, cylindricalHeatHole),
			NewTranslation(Vec3{W2, -RD03DL / 6, 0}, cylindricalHeatHole),
			NewTranslation(Vec3{-W2, -RD03DL / 6, 0}, cylindricalHeatHole),
		),
	)
}

func newSensorBox() Primitive {
	return NewTranslation(
		Vec3{0, 0, PowerBoxLength / 2},
		NewRotation(
			Vec3{-90, 0, 0},
			NewUnion(
				NewDifference(
					newPowerBox(),
					// Cutout the bottom of the box for the powerblock
					NewTranslation(
						Vec3{0, ((PowerBoxLength / 2) - (PowerBlockLength / 2)) - W1, 0},
						shapes.NewSmoothedCube(
							Vec3{PowerBlockWidth, PowerBlockLength, PowerBlockHeight + 2*W1},
							PowerBlockRounding,
						).Build(),
					),
				),
				NewTranslation(
					Vec3{0, ((PowerBoxLength / 2) - (PowerBlockLength / 2)) - W1, (PowerBlockHeight + 2*W1) / 2},
					newPowerBlock(),
				),
			),
		),
	)
}

func newCombineAll() Primitive {
	// Return a setup of all the parts together for review
	// Move them apart for better viewing
	apart := 35.0
	return NewList(
		NewTranslation(
			Vec3{-apart, 0, 0},
			newSensorBox(),
		),
		NewTranslation(
			Vec3{apart, 0, 0},
			newSensorBoxLid(),
		),
		// Sensor stick
		NewTranslation(
			Vec3{-apart * 2, 0, 0},
			newSensorStick(),
		),
		NewTranslation(
			Vec3{apart * 2, 0, 0},
			newSensorStickLid(),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "sensor_box", Primitive: newSensorBox(), Flags: sys.Default},
		{Name: "sensor_box_lid", Primitive: newSensorBoxLid(), Flags: sys.Default},
		{Name: "RD03D_stick", Primitive: newSensorStick(), Flags: sys.Default},
		{Name: "RD03D_stick_lid", Primitive: newSensorStickLid(), Flags: sys.Default},
		{Name: "for-review-only", Primitive: newCombineAll(), Flags: sys.None},
	})
}
