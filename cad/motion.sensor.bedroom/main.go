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
	PowerBoxWidth     = 80.0
	PowerBoxLength    = 120.0
	PowerBoxHeight    = 40.0
	PowerBoxThickness = WT
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

// newSensorBox creates the power box primitive that holds the power block, antenna
// and has cutouts for the power block, and sensor.
func newSensorBox() Primitive {
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
			Vec3{0, -(PowerBoxLength / 2) + 2*SensorCylindricalR, PowerBoxHeight - 4*W1},
			NewCylinder(8*W1, SensorCylindricalR+W1),
		),
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

// newSensorStick creates the cylindrical tube with the rectangular (hollow) part on top where the sensor is located.
func newSensorStick() Primitive {

	cylindricalLength := 10.0
	stickHeight := 2*SensorCylindricalR + W2

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
					NewCylinderOnTheYAxis(2*W2, SensorCylindricalR),
				),
			),

			// Cylindrical part
			NewTranslation(
				Vec3{0, RD03DL/2.0 + cylindricalLength/2.0 + W1, 0},
				NewDifference(
					NewCylinderOnTheYAxis(cylindricalLength+W1, SensorCylindricalR+W1),
					NewCylinderOnTheYAxis(RD03DL, SensorCylindricalR),
				),
			),
		),
	)
}

func newSensorStickLid() Primitive {
	return NewTranslation(
		Vec3{0, 0, W1 / 2.0},
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
	)
}

func newProduct() Primitive {
	cylindricalHeatHole := NewCylinder(PowerBlockHeight+2*W2, W2*0.8)
	return NewRender(
		10,
		NewUnion(
			NewDifference(
				newSensorBox(),

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
				newSensorStick(),
			),
			NewTranslation(
				Vec3{-55, 0, 0},
				newSensorStickLid(),
			),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "sensor_box", Primitive: newProduct(), Flags: sys.Default},
		{Name: "sensor_box_lid", Primitive: newSensorBoxLid(), Flags: sys.Default},
	})
}
