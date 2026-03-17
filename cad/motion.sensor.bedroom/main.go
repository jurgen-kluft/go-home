package main

import (
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

// Components:
// - USB-C socket
// - Main ESP8266
//   - APDS9960 sensor (color, gesture, proximity) (sensor stick)
// - Relay board
//   - Secondary ESP8266
//   - RD03D Sensor (sensor stick)
//   - OLED display

// Size of the sensor box
const (
	PowerBoxWidth     = 2.0 + ESP8266W + 5.0 + RelayBoardW + 5.0 + ESP8266W + 2.0
	PowerBoxLength    = 5.0 + max(ESP8266L, RelayBoardL, ESP8266L) + 5.0
	PowerBoxHeight    = 40.0
	PowerBoxThickness = WT
)

// newSensorBox creates the power box primitive that holds the power block, antenna
// and has cutouts for the power block, and sensor.
func newSensorBox() Primitive {
	cylindricalHeatHole := NewCylinder(PowerBlockHeight+2*W2, W2*0.8)

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
			Vec3{0, -(PowerBoxLength / 2) + 2*SensorStickCylindricalRadius, PowerBoxHeight - 4*W1},
			NewCylinder(8*W1, SensorStickCylindricalRadius+W1),
		),

		// Cutout in the bottom
		NewTranslation(Vec3{0, 0, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{0, 8 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),

		NewTranslation(Vec3{(PowerBlockWidth/2 - 2*W2), 0, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{-(PowerBlockWidth/2 - 2*W2), 0, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{(PowerBlockWidth/2 - 3*W2), 2 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{-(PowerBlockWidth/2 - 3*W2), 2 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{(PowerBlockWidth/2 - 2*W2), 4 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{-(PowerBlockWidth/2 - 2*W2), 4 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{(PowerBlockWidth/2 - 3*W2), 6 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),
		NewTranslation(Vec3{-(PowerBlockWidth/2 - 3*W2), 6 * W2, PowerBlockHeight / 2}, cylindricalHeatHole),
	)
}

func newSensorBoxLid() Primitive {
	return NewDifference(
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
		NewTranslation(Vec3{10, 0, 0}, NewRotation(Vec3{90, 0, 0}, NewCylinder(8*W2, W2*0.6))),
		NewTranslation(Vec3{-10, 0, 0}, NewRotation(Vec3{90, 0, 0}, NewCylinder(8*W2, W2*0.6))),

		NewTranslation(Vec3{0, 0, 4 * W2}, NewRotation(Vec3{90, 0, 0}, NewCylinder(8*W2, W2*0.6))),
		NewTranslation(Vec3{0, 0, -4 * W2}, NewRotation(Vec3{90, 0, 0}, NewCylinder(8*W2, W2*0.6))),
	)
}

func newProduct() Primitive {
	return NewRender(
		10,
		NewUnion(
			NewTranslation(
				Vec3{0, 0, 0},
				newSensorBox(),
			),

			// // DEBUG Antenna shape
			// NewTranslation(
			// 	Vec3{0, SensorCylindricalR, 1.1 * PowerBoxHeight},
			// 	NewBox(AntennaW, AntennaL, AntennaH*3),
			// ),
			NewTranslation(
				Vec3{-120, 0, 0},
				NewRotation(
					Vec3{90, 0, 0},
					newSensorBoxLid(),
				),
			),
			// Sensor stick
			NewTranslation(
				Vec3{60, 0, 0},
				NewSensorStick(RD03DW, RD03DL, RD03DH, 20.0),
			),
			NewTranslation(
				Vec3{-55, 0, 0},
				NewSensorStickLid(RD03DW, RD03DL, RD03DH),
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
			Vec3{0, 0, 0},
			newSensorBox(),
		),
		NewTranslation(
			Vec3{0, apart, PowerBoxHeight / 2},
			newSensorBoxLid(),
		),
		// Sensor stick
		NewTranslation(
			Vec3{0, apart * 2, 0},
			NewSensorStick(RD03DW, RD03DL, RD03DH, 20.0),
		),
		NewTranslation(
			Vec3{0, -apart * 2, 0},
			NewSensorStickLid(RD03DW, RD03DL, RD03DH),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "sensor_box", Primitive: newProduct(), Flags: sys.Default},
		{Name: "sensor_box_lid", Primitive: newSensorBoxLid(), Flags: sys.Default},
		{Name: "for-review-only", Primitive: newCombineAll(), Flags: sys.None},
	})
}
