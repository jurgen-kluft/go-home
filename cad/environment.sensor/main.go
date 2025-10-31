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

// Global constants
const (
	Unit         = 1.0 // 1 unit = 1 mm
	WT           = 1.5
	W1           = WT
	W2           = 2 * WT
	Rounding     = 1
	MainRounding = 10
)

// Box dimensions
const (
	boxW = pcbBoardW + 2*10
	boxL = pcbBoardL + 2*12
	boxH = 25.0
	boxT = 2.0
	boxI = 1.25
)

// Box bottom
const (
	boxBottomT = W1
	boxBottomH = 10.0
)

// Inlay magnets for holding the box bottom and top together
// - 5x2 mm Round Magnet, Micro Neodymium Magnet
const (
	magnetDiameter    = 5.0
	magnetR           = magnetDiameter / 2.0
	magnetH           = 2.0
	magnetInlayHeight = 10.0
	magnetInlayRadius = magnetR + W1
)

// PCB Board
const (
	pcbBoardW  = 62 + (Rounding / 2)
	pcbBoardL  = 90 + (Rounding / 2)
	pcbBoardT  = 1.6                                 // Thickness of the PCB board itself
	pcbBoardBH = 6.0 - pcbBoardT                     // Height of the tallest component on the bottom of the PCB
	pcbBoardTH = W1                                  // Thickness of the tallest component on the top of the PCB
	pcbBoardH  = pcbBoardBH + pcbBoardTH + pcbBoardT // Total height of the PCB including components
)

// SH1107 OLED Display, 128x128
const (
	OledW = 14.0
	OledL = 14.0
	OledH = 2.5
)

// RD03D
const (
	RD03DW  = 15.25
	RD03DL  = 44.5
	RD03DH  = 5.0
	RD03DT  = 1.5 // PCB thickness
	RD03DBT = 0.7 // Bottom height (height of components on bottom side
)

// Scd41 (CO2, temp, humidity)
const (
	Scd41SW = 8    // Width of the sensor body
	Scd41W  = 13.2 // Width of the PCB
	Scd41L  = 22.0 // Length of the PCB
	Scd41H  = 8
	Scd41T  = 1.6
)

// Bh1750 (light)
const (
	Bh1750W = 14.2
	Bh1750L = 18.8
	Bh1750H = 3.0
	Bh1750T = 1.6
)

// Bme280 (temperature, humidity, pressure)
const (
	Bme280W = 10.5
	Bme280L = 13.2
	Bme280H = 3.5
	Bme280T = 1.6
)

// USB-C hole dimensions
const (
	UsbCHoleR = 11.25 // Radius of the USB-C hole, 1 cm
)

func newBox(w, l, h, r float64) Primitive {
	return NewRender(
		20,
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
		20,
		NewDifference(
			newBox(outerW, outerL, h, r),
			newBox(outerW-2*wallThickness, outerL-2*wallThickness, 1.2*h, r),
		))
}

func newMagnetInlay() Primitive {
	return NewRender(
		20,
		NewTranslation(
			Vec3{0, 0, magnetInlayHeight / 2},
			NewDifference(
				NewCylinder(magnetInlayHeight, magnetInlayRadius),
				NewTranslation(
					Vec3{0, 0, magnetInlayHeight/2 - magnetH/2},
					NewCylinder(magnetH*2, magnetR),
				),
			),
		),
	)
}

// Insert slide, reusable for different sensors
func newInsertSlide(w, l, h, t, bh float64) Primitive {
	round := 0.25
	return NewTranslation(
		Vec3{0, 0, bh / 2},
		NewRender(
			200,
			NewUnion(
				NewTranslation(
					Vec3{0, 0, bh/2 - 2*round + (l+W2)/2},
					NewDifference(
						NewDifference(
							newBox(w+W2+W2, h, l+W2, round),
							NewTranslation(
								Vec3{0, 0, W1},
								newBox(w, 2*h, l+W1, round),
							),
						),
						NewTranslation(
							Vec3{0, 0, W1},
							newBox(w+W2, t, l+W2, round),
						),
					),
				),
				newBox(w+W2+W2, h, bh, round),
			),
		),
	)
}

func newScd41InsertSlide() Primitive {
	return newInsertSlide(Scd41W, Scd41L, 2*W2, Scd41T, boxBottomH-boxBottomT)
}
func newBh1750InsertSlide() Primitive {
	return newInsertSlide(Bh1750W, Bh1750L, 2*W2, Bh1750T, boxBottomH-boxBottomT)
}
func newBme280InsertSlide() Primitive {
	return newInsertSlide(Bme280W, Bme280L, 2*W2, Bme280T, boxBottomH-boxBottomT)
}
func newRd03dInsertSlide() Primitive {
	return newInsertSlide(RD03DW, RD03DL, 2*W2, RD03DT+RD03DBT, boxBottomH-boxBottomT)
}

// The frontside of the box, also has 4 sides of a certain height:
// - RD03D sensor insert-slide (upper part)
// - OLED display cutout (lower part)
func newBoxFrontside() Primitive {
	return newBox(boxW, boxL, boxH, MainRounding)
}

// The backside of the box:
// - PCB inlay
// - Sensor insert-slides
// - USB-C connector cutout
func newBoxBackside() Primitive {
	V := Vec3{-(boxW / 2), (boxL / 2), 0}
	CurrentLen := V.Len()
	NewLen := CurrentLen - MainRounding

	LTP := V.Mul(NewLen / CurrentLen)
	RTP := Vec3{-LTP.X(), LTP.Y(), LTP.Z()}
	LBP := Vec3{LTP.X(), -LTP.Y(), LTP.Z()}
	RBP := Vec3{-LTP.X(), -LTP.Y(), LTP.Z()}

	return NewUnion(
		NewTranslation(
			Vec3{0, 0, boxBottomH / 2},
			NewRender(
				20,
				NewDifference(
					newBox(boxW+2*(boxT+boxI), boxL+2*(boxT+boxI), boxBottomH, MainRounding),
					NewTranslation(
						Vec3{0, 0, boxT},
						newBox(boxW, boxL, boxBottomH, MainRounding),
					),
				),
			),
			NewTranslation(
				Vec3{0, 0, W1},
				newWall(boxI, boxW+2*boxI, boxL+2*boxI, boxBottomH, MainRounding),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				LTP,
				newMagnetInlay(),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				RTP,
				newMagnetInlay(),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				LBP,
				newMagnetInlay(),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				RBP,
				newMagnetInlay(),
			),
		),
		NewTranslation(
			Vec3{0, boxL/2 - Scd41H, boxBottomT},
			newScd41InsertSlide(),
		),
		NewTranslation(
			Vec3{0, -(boxL/2 - W2 - Bh1750H), boxBottomT},
			newBh1750InsertSlide(),
		),
		NewTranslation(
			Vec3{(-boxW / 2) + 3*W1, 0, boxBottomT},
			NewRotation(
				Vec3{0, 0, 90},
				newBme280InsertSlide(),
			),
		),
	)
}

func newPcbInlay() Primitive {
	return NewTranslation(
		Vec3{0, 0, boxT},
		NewRender(
			20,
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

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "main", Primitive: NewList(newBoxBackside(), newPcbInlay()), Flags: sys.Default},
	})
}
