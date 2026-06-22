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
	pcbBoardW      = 70 + (Rounding / 2)
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
	displayScreenW       = WcsScreenL
	displayScreenL       = WcsScreenW
	displayScreenR       = WcsScreenR
	displayMountingW     = WcsMountingW
	displayMountingL     = WcsMountingL
	displayMountingHoleR = WcsMountingHoleD / 2.0
)

// Circular hole dimensions for sensor mounting
const (
	sensorMountingHoleD = 10.0 // Diameter of the mounting holes for sensors
	sensorMountingHoleR = sensorMountingHoleD / 2.0
)

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
				NewCylinder(h, magnetInlayRadius).SetFn(20),
				NewTranslation(
					Vec3{0, 0, h / 2},
					NewCylinder(magnetH*2, magnetR).SetFn(20),
				),
			),
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
							NewTranslation(
								Vec3{0, 0, boxBacksideH / 2.0},
								newBox(boxW+2*(boxI), boxL+2*(boxI), 2*W1, MainRounding),
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
						NewCylinder(boxBacksideH*2, UsbCHoleRadius).SetFn(20),
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

		// NewTranslation(
		// 	Vec3{0, 0, boxBacksideT},
		// 	newMounting(pcbBoardMountW, pcbBoardMountL, displayMountingHoleR, 4*W1),
		// ),

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
			NewCylinder(lowerHeight, lowerRadius).SetFn(20),
			NewCylinder(lowerHeight+nailHeight, nailRadius).SetFn(20),
		),
		// Cutoff the bottom
		NewTranslation(
			Vec3{0, 0, -(lowerHeight + nailHeight + W1) / 2},
			NewCylinder((lowerHeight+nailHeight+W1), lowerRadius+W2).SetFn(20),
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
						NewTranslation(
							Vec3{0, 0, boxFrontsideH / 2.0},
							newBox(boxW+2*(boxI), boxL+2*(boxI), boxBacksideT*2, MainRounding),
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
						NewCylinder(W2*3, sensorMountingHoleR).SetFn(20),
					),
				),
				NewTranslation(
					Vec3{(boxW / 2), (7 * W2), sensorMountingHoleR / 2.0},
					NewRotation(
						Vec3{0, -90, 0},
						NewCylinder(W2*3, sensorMountingHoleR).SetFn(20),
					),
				),

				// Side
				NewTranslation(
					Vec3{-15, ((boxL / 2) + (boxFrontsideT-boxI)/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR).SetFn(20),
					),
				),
				NewTranslation(
					Vec3{15, ((boxL / 2) + (boxFrontsideT-boxI)/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR).SetFn(20),
					),
				),

				NewTranslation(
					Vec3{-15, -((boxL / 2) + boxFrontsideT/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR).SetFn(20),
					),
				),
				NewTranslation(
					Vec3{15, -((boxL / 2) + boxFrontsideT/2), 0},
					NewRotation(
						Vec3{90, 0, 0},
						NewCylinder(W2*3, sensorMountingHoleR).SetFn(20),
					),
				),

				// Opening on the front for the OLED display
				NewTranslation(
					Vec3{0, 0, -3 * W2},
					newBox(displayScreenW, displayScreenL, 6*W2, displayScreenR),
				),
			),
		),

		// AMOLED display mounting
		NewTranslation(
			Vec3{0, 0, boxFrontsideT},
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
	return NewRender(
		10,
		NewUnion(
			NewTranslation(
				Vec3{boxW/2 + 10, 0, 0},
				newBoxBackside(),
			),
			NewTranslation(
				Vec3{-(boxW/2 + 10), 0, 0},
				newBoxFrontside(),
			),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "environment.sensor.normal", Primitive: newProduct(), Flags: sys.Default},
	})
}
