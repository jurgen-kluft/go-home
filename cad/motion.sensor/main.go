package main

import (
	"github.com/ljanyst/ghostscad/lib/shapes"
	"github.com/ljanyst/ghostscad/sys"

	. "github.com/go-gl/mathgl/mgl64"
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

// Global constants
const (
	Unit         = 1.0    // 1 unit = 1 mm
	WT           = 1.6    // Wall Thickness
	W1           = WT     // Thin Wall
	W2           = 2 * WT // Medium Wall
	Rounding     = 0.5    // Small Rounding
	MainRounding = 10     // Main Rounding
)

// Power block dimensions
const (
	powerBlockWidth    = 37.0 + 0.25
	powerBlockLength   = 46.25
	powerBlockHeight   = 24.1 + 0.25
	powerBlockRounding = 2.5
)

// Power box dimensions
const (
	powerBoxWidth  = powerBlockWidth + 2*W1
	powerBoxLength = powerBlockLength + 18.0
	powerBoxHeight = powerBlockHeight * 2
)

// USB-A dimensions
const (
	USBAWidth  = 12.0
	USBALength = 16.0
	USBAHeight = 7.0
)

// Seeed studio XIAO ESP32-C3 dimensions
// USB-C connector
const (
	ESP32C3W = 18.0
	ESP32C3L = 21.0
	ESP32C3H = 5
	ESP32C3T = 1.25 // Thickness of the PCB
)

// The sensor cylindrical connection with the box
const (
	SensorCylindricalR = 8.0 / 2
)

// External antenna dimensions
const (
	AntennaW = 20.0
	AntennaL = 40.0
	AntennaH = 1.5
)

// RD03D mmWave sensor dimensions
const (
	RD03DW  = 15.25
	RD03DL  = 44.5
	RD03DH  = 5.0
	RD03DT  = 3.2 // PCB and things thickness
	RD03DBT = 0.7 // Bottom height (height of components on bottom side
	RD03DTT = 1.0 // Bottom height (height of components on bottom side
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

// newPowerBlock creates the power block primitive that holds the USB power block.
// The power block has cutouts on the top and bottom for heat dissipation
func newPowerBlock() Primitive {
	return NewTranslation(
		Vec3{0, 0, powerBlockHeight/2 + W1},
		NewDifference(
			shapes.NewSmoothedCube(
				Vec3{powerBlockWidth + 2*W1, powerBlockLength + 2*W1, powerBlockHeight + 2*W1},
				powerBlockRounding,
			).Build(),

			// Emptying + Front cutout(for inserting the power block)
			NewTranslation(
				Vec3{0, W1, 0},
				NewBox(powerBlockWidth, powerBlockLength+W1, powerBlockHeight),
			),

			// USB-A cutout
			NewTranslation(
				Vec3{0, -((powerBlockLength / 2) + (USBALength / 2) - W1), 0},
				NewBox(USBAWidth, USBALength, powerBlockHeight-2*W1),
			),
		),
	)
}

// newPowerBox creates the power box primitive that holds the power block, antenna
// and has cutouts for the power block, and sensor.
func newPowerBox() Primitive {
	return NewDifference(
		NewTranslation(
			Vec3{0, 0, powerBoxHeight / 2},
			NewDifference(
				// The box
				NewBox(powerBoxWidth, powerBoxLength, powerBoxHeight),
				// Make the box hollow
				NewBox(powerBoxWidth-2*W1, powerBoxLength-2*W1, powerBoxHeight-2*W1),
			),
		),

		// Power block cutout
		NewTranslation(
			Vec3{0, ((powerBoxLength / 2) - (powerBlockLength / 2)) + W1, powerBlockHeight/2 + W1},
			NewBox(powerBlockWidth, powerBlockLength, powerBlockHeight),
		),

		// Backside cover cutout
		NewTranslation(
			Vec3{0, -((powerBoxLength / 2) - W1), powerBoxHeight / 2},
			//NewBox(powerBoxWidth-2*W2, 2*W1, (powerBoxHeight-2*W2)),
			NewCube(Vec3{powerBoxWidth - 2*W2, 2 * W1, (powerBoxHeight - 2*W2)}),
		),

		// Circular Sensor cutout at the top
		NewTranslation(
			Vec3{0, -(powerBoxLength / 2) + 2*SensorCylindricalR, powerBoxHeight - 4*W1},
			NewCylinder(8*W1, SensorCylindricalR+W1),
		),
	)
}

func newPowerBoxLid() Primitive {
	return NewDifference(
		NewUnion(
			NewTranslation(
				Vec3{0, W1 / 2, 0},
				NewCube(Vec3{powerBoxWidth - W1, W1, powerBoxHeight - W1}),
			),
			NewTranslation(
				Vec3{0, W2 / 2, 0},
				NewCube(Vec3{powerBoxWidth - 2*W2 - 0.2, W2, (powerBoxHeight - 2*W2)}),
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
	cylindricalHeatHole := NewCylinder(powerBlockHeight+2*W2, W2*0.8)
	return NewRender(
		10,
		NewUnion(
			NewDifference(
				NewUnion(
					newPowerBox(),
					NewTranslation(
						Vec3{0, ((powerBoxLength / 2) - (powerBlockLength / 2)) - W1, 0},
						newPowerBlock(),
					),
				),

				// Cutout in the bottom
				NewTranslation(Vec3{0, 0, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{0, 8 * W2, powerBlockHeight / 2}, cylindricalHeatHole),

				NewTranslation(Vec3{(powerBlockWidth/2 - 2*W2), 0, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{-(powerBlockWidth/2 - 2*W2), 0, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{(powerBlockWidth/2 - 3*W2), 2 * W2, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{-(powerBlockWidth/2 - 3*W2), 2 * W2, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{(powerBlockWidth/2 - 2*W2), 4 * W2, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{-(powerBlockWidth/2 - 2*W2), 4 * W2, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{(powerBlockWidth/2 - 3*W2), 6 * W2, powerBlockHeight / 2}, cylindricalHeatHole),
				NewTranslation(Vec3{-(powerBlockWidth/2 - 3*W2), 6 * W2, powerBlockHeight / 2}, cylindricalHeatHole),
			),

			// // DEBUG Antenna shape
			// NewTranslation(
			// 	Vec3{0, SensorCylindricalR, 1.1 * powerBoxHeight},
			// 	NewBox(AntennaW, AntennaL, AntennaH*3),
			// ),
			NewTranslation(
				Vec3{-70, 0, 0},
				NewRotation(
					Vec3{90, 0, 0},
					newPowerBoxLid(),
				),
			),
			// Sensor stick
			NewTranslation(
				Vec3{40, 0, 0},
				newSensorStick(),
			),
			NewTranslation(
				Vec3{-35, 0, 0},
				newSensorStickLid(),
			),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "power_box", Primitive: newProduct(), Flags: sys.Default},
		{Name: "power_box_lid", Primitive: newPowerBoxLid(), Flags: sys.Default},
	})
}
