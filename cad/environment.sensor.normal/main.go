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

// TODO
//
// What putting the sensor outside of the box?
// We just needs standard cirular holes where sensors can be mounted.
// This way, we can be compatible with many different sensors.
//

// Box dimensions
const (
	boxW = pcbBoardW + 2*20
	boxL = pcbBoardL + 2*10
	boxT = 2.0
	boxI = 1.25
)

// Bottom of the box
const (
	boxBacksideT = boxT
	boxBacksideH = 15.0
)

// Top of the box
const (
	boxFrontsideT = boxT
	boxFrontsideH = 25.0
)

// Inlay magnets for holding the box bottom and top together
// - 5x2 mm Round Magnet, Micro Neodymium Magnet
const (
	magnetDiameter    = 5.0
	magnetR           = magnetDiameter / 2.0
	magnetH           = 2.1
	magnetInlayRadius = magnetR + W1
)

// PCB Board from
const (
	pcbBoardW      = 56.25 + (Rounding / 2)
	pcbBoardL      = 80.25 + (Rounding / 2)
	pcbBoardT      = 1.6                                 // Thickness of the PCB board itself
	pcbBoardMountW = 51                                  // Mounting width, from hole to hole, for the PCB board
	pcbBoardMountL = 75                                  // Mounting length, from hole to hole, for the PCB board
	displayMountHR = 2.0 / 2                             // Mounting hole radius for the PCB board
	pcbBoardBH     = 6.0 - pcbBoardT                     // Height of the tallest component on the bottom of the PCB
	pcbBoardTH     = W1                                  // Thickness of the tallest component on the top of the PCB
	pcbBoardH      = pcbBoardBH + pcbBoardTH + pcbBoardT // Total height of the PCB including components
)

// Display, SH1107 OLED, 128x128
const (
	displayScreenW       = Sh1107ScreenL
	displayScreenL       = Sh1107ScreenW
	displayScreenR       = Sh1107ScreenR
	displayMountingW     = Sh1107MountingW
	displayMountingL     = Sh1107MountingL
	displayMountingHoleR = Sh1107MountingHoleD / 2
)

// Circular hole dimensions for sensor mounting
const (
	sensorMountingHoleD = 10.0 // Diameter of the mounting holes for sensors
	sensorMountingHoleR = sensorMountingHoleD / 2.0
)

// USB-C hole dimensions
const (
	UsbCHoleDiameter     = 11.6 // Radius of the USB-C hole, 1.16 cm
	UsbCHoleRadius       = UsbCHoleDiameter / 2.0
	UsbCHoleRingDiameter = 17.0
)

func newPyramid(w, l, h float64) Primitive {
	baseLT := Vec3{-w / 2, l / 2, -h / 2}
	baseRT := Vec3{w / 2, l / 2, -h / 2}
	baseLB := Vec3{-w / 2, -l / 2, -h / 2}
	baseRB := Vec3{w / 2, -l / 2, -h / 2}
	apex := Vec3{0, 0, h / 2}

	points := []Vec3{
		baseLT,
		baseRT,
		baseRB,
		baseLB,
		apex,
	}

	triangles := []Vec3{
		// Base
		{0, 1, 2},
		{0, 2, 3},
		// Sides
		{0, 1, 4},
		{1, 2, 4},
		{2, 3, 4},
		{3, 0, 4},
	}

	return NewRender(
		10,
		NewPolyhedron(points, triangles),
	)
}

func NewCylinderOnTheXAxis(h, r float64) Primitive {
	return NewRotation(Vec3{0, 90, 0}, NewCylinder(h, r))
}

func NewCylinderOnTheYAxis(h, r float64) Primitive {
	return NewRotation(Vec3{90, 0, 0}, NewCylinder(h, r))
}

func newBox(w, l, h, r float64) Primitive {
	return NewRender(
		10,
		NewDifference(
			shapes.NewSmoothedCube(Vec3{w, l, h + 2*r}, r).Build(),
			// Cutoff the top and bottom to end up with the requested height
			NewTranslation(
				Vec3{0, 0, -((h / 2) + ((h + 2*r) / 2))},
				NewCube(Vec3{w + W2, l + W2, h + 2*r}),
			),
			NewTranslation(
				Vec3{0, 0, ((h / 2) + ((h + 2*r) / 2))},
				NewCube(Vec3{w + W2, l + W2, h + 2*r}),
			),
		),
	)
}

func newWall(wallThickness, outerW, outerL, h, r float64) Primitive {
	return NewRender(
		10,
		NewDifference(
			newBox(outerW, outerL, h, r),
			newBox(outerW-2*wallThickness, outerL-2*wallThickness, 1.2*h, r),
		))
}

func newMagnetInlay(h float64) Primitive {
	return NewRender(
		10,
		NewTranslation(
			Vec3{0, 0, h / 2},
			NewDifference(
				NewCylinder(h, magnetInlayRadius),
				NewTranslation(
					Vec3{0, 0, h / 2},
					NewCylinder(magnetH*2, magnetR),
				),
			),
		),
	)
}

// Insert slide, reusable for different sensors
func newVerticalInsertSlide(w, l, h, t, bh float64) Primitive {
	round := 0.2
	overhang := 1.0
	return NewTranslation(
		Vec3{0, 0, bh / 2},
		NewRender(
			10,
			NewUnion(
				NewTranslation(
					Vec3{0, 0, bh/2 - 2*round + (l+W2)/2},
					NewDifference(
						NewDifference(
							newBox(W1+w+overhang, h, l+W2, round),
							NewTranslation(
								Vec3{0, 0, W1},
								newBox(-overhang+w-overhang, 2*h, l+W1, round),
							),
						),
						NewTranslation(
							Vec3{0, 0, W1},
							newBox(w, t, l+W2, round),
						),
					),
				),
				newBox(W1+w+W1, h, bh, round),
			),
		),
	)
}

func newHorizontalInsertSlide2(w, l, h, t float64) Primitive {
	round := 0.25
	return NewRender(
		10,
		NewDifference(
			NewDifference(
				newBox(W1+w+W1, l+W2, h, round),
				NewTranslation(
					Vec3{0, W2, 0},
					newBox(w-W2, l+W2, h*2, round),
				),
			),
			NewTranslation(
				Vec3{0, W1, 0},
				newBox(w, l+W2, t, round),
			),
		),
	)
}

func newScd41Sensor() Primitive {
	pipeLen := W1 * 3

	// The pipe part, thickness is sensorMountingHoleR and it is hollow for wires to go through
	// Wall thickness is W1
	return NewDifference(
		NewUnion(
			NewTranslation(
				Vec3{0, ((Scd41L + W2) / 2) + W1, 0},
				NewCylinderOnTheYAxis(pipeLen, sensorMountingHoleR),
			),
			// The box part that contains the sensor
			NewTranslation(
				Vec3{0, 0, 0},
				NewDifference(
					NewCube(Vec3{Scd41W + W2, Scd41L + W2, Scd41H + W1}),
					NewTranslation(
						Vec3{0, 0, W1},
						NewCube(Vec3{Scd41W, Scd41L, Scd41H + W1 + W1}),
					),
				),
			),
		),
		NewTranslation(
			Vec3{0, ((Scd41L + W2) / 2) + W1, 0},
			NewCylinderOnTheYAxis(pipeLen*2, sensorMountingHoleR-W1),
		),
	)
}

func newScd41SensorLid() Primitive {
	// The box part that contains the sensor
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewUnion(
				NewCube(Vec3{Scd41W + W2, Scd41L + W2, W1}),
				NewTranslation(
					Vec3{0, 0, W1},
					NewDifference(
						NewCube(Vec3{Scd41W, Scd41L, W1}),
						NewCube(Vec3{Scd41W - W2, Scd41L - W2, W2}),
					),
				),
			),
			// A cutout for the actual CO2 sensor
			NewTranslation(
				Vec3{0, (Scd41L / 2.0) - Scd41BottomToSensorMiddle, 0},
				NewCube(Vec3{Scd41SensorWidth, Scd41SensorLength, W1 + W1}),
			),
		),
	)
}

func newScd41SensorReview() Primitive {
	return NewUnion(
		newScd41Sensor(),
		NewTranslation(
			Vec3{0, 0, Scd41BottomToSensorMiddle - (Scd41SensorHeight / 2)},
			newScd41SensorLid(),
		),
	)
}

// BH1750 sensor (light)
func newBh1750Sensor() Primitive {
	pipeLen := W1 * 3
	sW := Bh1750W + Bh1750Spacing
	sL := Bh1750L + Bh1750Spacing
	sH := W2 + W1

	// The pipe part, thickness is sensorMountingHoleR and it is hollow for wires to go through
	// Wall thickness is W1
	return NewDifference(
		NewUnion(
			NewTranslation(
				//Vec3{0, ((sL + W2) / 2) + W1, 0},
				Vec3{0, sL/2 + W1 - sensorMountingHoleR, -(pipeLen / 2) - sH/2},
				NewCylinder(pipeLen, sensorMountingHoleR),
			),
			// The box part that contains the sensor
			NewDifference(
				NewCube(Vec3{sW + W2, sL + W2, sH}),
				NewTranslation(
					Vec3{0, 0, (sH+W1)/2.0 - sH/2.0 + W1},
					NewCube(Vec3{sW, sL, sH + W1}),
				),
			),
		),
		NewTranslation(
			Vec3{0, sL/2 + W1 - sensorMountingHoleR, -(pipeLen / 2) - sH/2},
			NewCylinder(pipeLen*4, sensorMountingHoleR-W1),
		),
	)
}

func newBh1750SensorLid() Primitive {
	// The box part that contains the sensor
	sW := Bh1750W + Bh1750Spacing
	sL := Bh1750L + Bh1750Spacing
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewUnion(
				NewCube(Vec3{sW + W2, sL + W2, W1}),
				NewTranslation(
					Vec3{0, 0, W1},
					NewDifference(
						NewCube(Vec3{sW - 0.2, sL, W1}),
						NewCube(Vec3{sW - W2, sL - W2, W2}),
					),
				),
			),
			// A cutout for light to reach the BH1750 sensor
			NewTranslation(
				Vec3{0, (sL / 2.0) - Bh1750BottomToSensorMiddle, 0},
				NewCube(Vec3{Bh1750SensorWidth * 1.5, Bh1750SensorHeight * 1.5, W1 + W1}),
			),
		),
	)
}

func newBh1750SensorReview() Primitive {
	sW := Bh1750W + 4*W1
	sH := W2 + W1
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, sH / 2},
			NewRotation(
				Vec3{180, 0, 0},
				newBh1750Sensor(),
			),
		),
		NewTranslation(
			Vec3{sW, 0, W1 / 2.0},
			newBh1750SensorLid(),
		),
	)
}

func newBme280Sensor() Primitive {
	pipeLen := W1 * 3
	pipeTranslation := Vec3{0, -(Bme280L/2.0 + W1 - sensorMountingHoleD/2.0), -pipeLen/2 - Bme280H/2}

	// The pipe part, thickness is sensorMountingHoleR and it is hollow for wires to go through
	// Wall thickness is W1
	return NewDifference(
		NewUnion(
			NewTranslation(
				pipeTranslation,
				NewCylinder(pipeLen, sensorMountingHoleR),
			),
			// The box part that contains the sensor
			NewTranslation(
				Vec3{0, 0, 0},
				NewDifference(
					NewCube(Vec3{Bme280W + W2, Bme280L + W2, Bme280H + W1}),
					NewTranslation(
						Vec3{0, 0, W1},
						NewCube(Vec3{Bme280W, Bme280L, Bme280H + W1 + W1}),
					),
				),
			),
		),
		NewTranslation(
			pipeTranslation,
			NewCylinder(pipeLen*2, sensorMountingHoleR-W1),
		),
	)
}

func newBme280SensorLid() Primitive {
	// The box part that contains the sensor
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewUnion(
				NewCube(Vec3{Bme280W + W2, Bme280L + W2, W1}),
				NewTranslation(
					Vec3{0, 0, W1},
					NewDifference(
						NewCube(Vec3{Bme280W, Bme280L, W1}),
						NewCube(Vec3{Bme280W - W2, Bme280L - W2, W2}),
					),
				),
			),
			// A cutout for the actual BME280 sensor
			NewTranslation(
				Vec3{0, (Bme280L / 2.0) - Bme280BottomToSensorMiddle, 0},
				NewCube(Vec3{Bme280SensorWidth, Bme280SensorLength, W1 + W1 + Bme280H}),
			),
		),
	)
}

func newBme280SensorReview() Primitive {
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, Bme280H/2 + W1/2},
			NewRotation(
				Vec3{180, 0, 0},
				newBme280Sensor(),
			),
		),
		NewTranslation(
			Vec3{20, 0, W1 / 2},
			newBme280SensorLid(),
		),
	)
}

func newRD03DSensor() Primitive {

	length := RD03DL
	height := sensorMountingHoleD

	pipeLen := W1 * 3
	pipeTranslation := Vec3{0, -(length/2.0 + W1 + pipeLen/2), 0}

	// The pipe part, thickness is sensorMountingHoleR and it is hollow for wires to go through
	// Wall thickness is W1
	return NewDifference(
		NewUnion(
			NewTranslation(
				pipeTranslation,
				NewCylinderOnTheYAxis(pipeLen, sensorMountingHoleR),
			),
			// The box part that contains the sensor
			NewTranslation(
				Vec3{0, 0, 0},
				NewDifference(
					NewCube(Vec3{RD03DW + W2, length + W2, height}),
					NewTranslation(
						Vec3{0, 0, W1},
						NewCube(Vec3{RD03DW, length, height}),
					),
				),
			),
		),
		NewTranslation(
			pipeTranslation,
			NewCylinderOnTheYAxis(pipeLen*2, sensorMountingHoleR-W1),
		),
	)
}

func newRD03DSensorLid() Primitive {
	length := RD03DL

	// The box part that contains the sensor
	return NewTranslation(
		Vec3{0, 0, 0},
		NewDifference(
			NewUnion(
				NewCube(Vec3{RD03DW + W2, length + W2, W1}),
				NewTranslation(
					Vec3{0, 0, W1},
					NewDifference(
						NewCube(Vec3{RD03DW, length, W1}),
						NewCube(Vec3{RD03DW - W2, length - W2, W2}),
					),
				),
			),
		),
	)
}

func newRD03DSensorReview() Primitive {
	height := sensorMountingHoleD
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, height / 2},
			newRD03DSensor(),
		),
		NewTranslation(
			Vec3{25, 0, W1 / 2},
			newRD03DSensorLid(),
		),
	)
}

// The backside of the box:
// - PCB inlay
// - Sensor insert-slides
// - USB-C connector cutout
func newBoxBackside() Primitive {
	V := Vec3{(boxW / 2), (boxL / 2), 0}
	CurrentLen := V.Len()
	NewLen := CurrentLen - MainRounding + 1

	LTP := V.Mul(NewLen / CurrentLen)
	RTP := Vec3{-LTP.X(), LTP.Y(), 0}
	LBP := Vec3{LTP.X(), -LTP.Y(), 0}
	RBP := Vec3{-LTP.X(), -LTP.Y(), 0}

	return NewUnion(
		NewTranslation(
			Vec3{0, 0, boxBacksideH / 2},
			NewRender(
				10,
				NewDifference(
					NewUnion(
						NewDifference(
							newBox(boxW+2*(boxT+boxI), boxL+2*(boxT+boxI), boxBacksideH, MainRounding),
							NewTranslation(
								Vec3{0, 0, boxBacksideT},
								newBox(boxW, boxL, boxBacksideH, MainRounding),
							),
						),
						NewTranslation(
							Vec3{0, 0, boxBacksideT},
							newWall(boxI, boxW+2*boxI, boxL+2*boxI, boxBacksideH, MainRounding),
						),
					),
					// USB-C connector cutout
					NewTranslation(
						Vec3{-(boxW / 2) + 2 + UsbCHoleRingDiameter/2, -(boxL / 2.0) + UsbCHoleRingDiameter, 0},
						NewCylinder(boxBacksideH*2, UsbCHoleRadius),
					),
				),
			),
		),

		// // Distance test for BH1750 sensor
		// NewTranslation(
		// 	Vec3{-4, ((boxW / 2) - (W2 + RD03DW + W2) + 4.2/2) / 2, boxBacksideH},
		// 	NewTranslation(
		// 		LBP,
		// 		newBox(5.0, (boxW/2)-(W2+RD03DW+W2)+4.2/2, boxBacksideH, 0.1),
		// 	),
		// ),

		NewTranslation(
			Vec3{0, 0, boxBacksideT},
			newMounting(pcbBoardMountW, pcbBoardMountL, displayMountingHoleR, 4*W1),
		),
		NewTranslation(
			Vec3{0, 0, boxBacksideT},
			NewTranslation(
				LTP,
				newMagnetInlay(boxBacksideH-boxBacksideT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBacksideT},
			NewTranslation(
				RTP,
				newMagnetInlay(boxBacksideH-boxBacksideT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBacksideT},
			NewTranslation(
				LBP,
				newMagnetInlay(boxBacksideH-boxBacksideT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBacksideT},
			NewTranslation(
				RBP,
				newMagnetInlay(boxBacksideH-boxBacksideT),
			),
		),
	)
}

func newMountingNail(lowerRadius, lowerHeight, nailRadius, nailHeight float64) Primitive {
	return NewDifference(
		NewUnion(
			NewCylinder(lowerHeight, lowerRadius),
			NewCylinder(lowerHeight+nailHeight, nailRadius),
		),
		// Cutoff the bottom
		NewTranslation(
			Vec3{0, 0, -(lowerHeight + nailHeight + W1) / 2},
			NewCylinder((lowerHeight+nailHeight+W1), lowerRadius+W2),
		),
	)
}

// newMounting creates mounting nails for a PCB board that has mounting holes.
func newMounting(w, l, hr, supportHeight float64) Primitive {
	// Create mounting nails at the 4 corners of the PCB
	// A mounting nail has two parts, the lower part (thicker) and the upper part that has the radius of the mounting hole.
	mountingWidth := w
	mountingLength := l
	holeRadius := hr
	nailLength := 3 * W1
	return NewUnion(
		NewTranslation(
			Vec3{-mountingWidth / 2, -mountingLength / 2, 0},
			newMountingNail(holeRadius+W1, supportHeight, holeRadius, nailLength),
		),
		NewTranslation(
			Vec3{mountingWidth / 2, -mountingLength / 2, 0},
			newMountingNail(holeRadius+W1, supportHeight, holeRadius, nailLength),
		),
		NewTranslation(
			Vec3{-mountingWidth / 2, mountingLength / 2, 0},
			newMountingNail(holeRadius+W1, supportHeight, holeRadius, nailLength),
		),
		NewTranslation(
			Vec3{mountingWidth / 2, mountingLength / 2, 0},
			newMountingNail(holeRadius+W1, supportHeight, holeRadius, nailLength),
		),
	)
}

func newPcbInlay() Primitive {
	return NewTranslation(
		Vec3{0, 0, boxT},
		NewRender(
			10,
			NewUnion(
				NewTranslation(
					Vec3{0, 0, pcbBoardH / 2},
					NewDifference(
						newBox(pcbBoardW+W2, pcbBoardL+W2, pcbBoardH, Rounding),
						newBox(pcbBoardW, pcbBoardL, pcbBoardH+2*W2, Rounding),
					),
				),
				NewTranslation(
					Vec3{0, 0, pcbBoardBH + W1/2},
					newWall(W2, pcbBoardW, pcbBoardL, W1, Rounding),
				),
			),
		),
	)
}

// The frontside of the box, also has 4 sides of a certain height:
// - RD03D sensor insert-slide (upper part)
// - OLED display cutout (lower part)
func newBoxFrontside() Primitive {
	V := Vec3{(boxW / 2), (boxL / 2), 0}
	CurrentLen := V.Len()
	NewLen := CurrentLen - MainRounding + 1

	LTP := V.Mul(NewLen / CurrentLen)
	RTP := Vec3{-LTP.X(), LTP.Y(), LTP.Z()}
	LBP := Vec3{LTP.X(), -LTP.Y(), LTP.Z()}
	RBP := Vec3{-LTP.X(), -LTP.Y(), LTP.Z()}

	return NewUnion(
		NewTranslation(
			Vec3{0, 0, boxFrontsideH / 2},
			NewDifference(
				NewDifference(
					NewDifference(
						newBox(boxW+2*(boxFrontsideT+boxI), boxL+2*(boxFrontsideT+boxI), boxFrontsideH, MainRounding),
						NewTranslation(
							Vec3{0, 0, boxFrontsideT},
							newBox(boxW, boxL, boxFrontsideH, MainRounding),
						),
					),
					NewTranslation(
						Vec3{0, 0, boxFrontsideH / 2},
					),
				),

				// Circular openings for sensors

				// TOP
				NewTranslation(
					Vec3{(boxW / 2), -(7 * W2), sensorMountingHoleR / 2.0},
					NewRotation(
						Vec3{0, -90, 0},
						NewCylinder(W2*3, sensorMountingHoleR),
					),
				),
				NewTranslation(
					Vec3{(boxW / 2), (7 * W2), sensorMountingHoleR / 2.0},
					NewRotation(
						Vec3{0, -90, 0},
						NewCylinder(W2*3, sensorMountingHoleR),
					),
				),

				// Side
				NewTranslation(
					Vec3{-15, ((boxL / 2) + (boxFrontsideT-boxI)/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR),
					),
				),
				NewTranslation(
					Vec3{15, ((boxL / 2) + (boxFrontsideT-boxI)/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR),
					),
				),

				NewTranslation(
					Vec3{-15, -((boxL / 2) + boxFrontsideT/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR),
					),
				),
				NewTranslation(
					Vec3{15, -((boxL / 2) + boxFrontsideT/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR),
					),
				),

				// Opening on the front for the OLED display
				NewTranslation(
					Vec3{15, (displayScreenL / 2), -3 * W2},
					newBox(displayScreenW, displayScreenL, 6*W2, displayScreenR),
				),
			),
		),

		// AMOLED display mounting
		NewTranslation(
			Vec3{15, (displayScreenL / 2), boxFrontsideT},
			newMounting(displayMountingL, displayMountingW, displayMountingHoleR, 1),
		),

		// Magnet inlays
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				LTP,
				newMagnetInlay(boxFrontsideH-boxFrontsideT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				RTP,
				newMagnetInlay(boxFrontsideH-boxFrontsideT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				LBP,
				newMagnetInlay(boxFrontsideH-boxFrontsideT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				RBP,
				newMagnetInlay(boxFrontsideH-boxFrontsideT),
			),
		),
	)
}

func newProduct() Primitive {
	return NewUnion(
		NewTranslation(
			Vec3{boxW/2 + 10, 0, 0},
			newBoxBackside(),
		),
		NewTranslation(
			Vec3{-(boxW/2 + 10), 0, 0},
			newBoxFrontside(),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "esp32_enclosure", Primitive: newProduct(), Flags: sys.Default},
		{Name: "scd41_sensor", Primitive: newScd41SensorReview(), Flags: sys.Default},
		{Name: "bh1750_sensor", Primitive: newBh1750SensorReview(), Flags: sys.Default},
		{Name: "bme280_sensor", Primitive: newBme280SensorReview(), Flags: sys.Default},
		{Name: "rd03d_sensor", Primitive: newRD03DSensorReview(), Flags: sys.Default},
	})
}
