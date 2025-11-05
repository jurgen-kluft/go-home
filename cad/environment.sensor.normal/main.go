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
	WT           = 1.6
	W1           = WT
	W2           = 2 * WT
	Rounding     = 1
	MainRounding = 10
)

// Box dimensions
const (
	boxW = pcbBoardW + 2*20
	boxL = pcbBoardL + 2*15
	boxT = 2.0
	boxI = 1.25
)

// Bottom of the box
const (
	boxBacksideT = boxT
	boxBacksideH = 10.0
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

// PCB Board
const (
	pcbBoardW  = 62 + (Rounding / 2)
	pcbBoardL  = 90 + (Rounding / 2)
	pcbBoardT  = 1.6                                 // Thickness of the PCB board itself
	pcbBoardBH = 6.0 - pcbBoardT                     // Height of the tallest component on the bottom of the PCB
	pcbBoardTH = W1                                  // Thickness of the tallest component on the top of the PCB
	pcbBoardH  = pcbBoardBH + pcbBoardTH + pcbBoardT // Total height of the PCB including components
)

// Display, SH1107 OLED, 128x128
const (
	displayScreenW       = 37.3 // Active display width
	displayScreenL       = 34.0 // Active display length
	displayW             = 47.1 // Overall width including bezel
	displayL             = 34.1 // Overall length including bezel
	displayMountingW     = 42   // Mounting hole to hole width
	displayMountingL     = 29   // Mounting hole to hole length
	displayMountingHoleD = 2.2  // Mounting hole diameter
)

// RD03D
const (
	RD03DW  = 15.25
	RD03DL  = 44.5
	RD03DH  = 5.0
	RD03DT  = 3.2 // PCB and things thickness
	RD03DBT = 0.7 // Bottom height (height of components on bottom side
	RD03DTT = 1.0 // Bottom height (height of components on bottom side
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
	UsbCHoleDiameter = 11.25 // Radius of the USB-C hole, 1 cm
	UsbCHoleRadius   = UsbCHoleDiameter / 2.0
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

func newScd41InsertSlide() Primitive {
	return newVerticalInsertSlide(Scd41W, Scd41L, W1+Scd41T+W1, Scd41T, boxBacksideH-boxBacksideT)
}
func newBh1750InsertSlide() Primitive {
	return newVerticalInsertSlide(Bh1750W, Bh1750L, W1+Bh1750T+W1, Bh1750T, boxBacksideH-boxBacksideT)
}
func newBme280InsertSlide() Primitive {
	return newVerticalInsertSlide(Bme280W, Bme280L, W1+Bme280T+W1, Bme280T, boxBacksideH-boxBacksideT)
}
func newRd03dInsertSlide() Primitive {
	return newHorizontalInsertSlide2(RD03DW, RD03DL, W1+RD03DT+W1, RD03DT)
}

// The backside of the box:
// - PCB inlay
// - Sensor insert-slides
// - USB-C connector cutout
func newBoxBackside() Primitive {
	V := Vec3{(boxW / 2), (boxL / 2), 0}
	CurrentLen := V.Len()
	NewLen := CurrentLen - MainRounding

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
						Vec3{RBP.X(), RBP.Y() + UsbCHoleDiameter + UsbCHoleRadius, 0},
						NewCylinder(boxBacksideH*2, UsbCHoleRadius),
					),
				),
			),
		),
		newPcbInlay(),
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
		NewTranslation(
			Vec3{0, boxL/2 - Scd41H, boxBacksideT},
			newScd41InsertSlide(),
		),
		NewTranslation(
			Vec3{0, -(boxL/2 - W2 - Bh1750H), boxBacksideT},
			newBme280InsertSlide(),
		),
		NewTranslation(
			Vec3{(boxW / 2) - 3*W1, (W2 + RD03DW + W2), boxBacksideT},
			NewRotation(
				Vec3{0, 0, 90},
				newBh1750InsertSlide(),
			),
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
	NewLen := CurrentLen - MainRounding

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
						newBox(boxW+2*(boxI), boxL+2*(boxI), W2, MainRounding),
					),
				),
				// Opening at the top side for inserting the RD03D sensor into the insert-slide
				NewTranslation(
					Vec3{boxW / 2, 0, (-boxFrontsideH / 2) + (RD03DH / 2) + 2*boxFrontsideT},
					NewRotation(
						Vec3{0, 0, 90},
						newBox(RD03DW, RD03DL/2, RD03DH+W2, Rounding),
					),
				),
				// Opening on the front for the OLED display
				NewTranslation(
					Vec3{0, -(boxL / 2) + (displayScreenL / 2) + boxT + W1, -3 * W2},
					newBox(displayScreenW+W2, displayScreenL+W2, 6*W2, Rounding),
				),
			),
		),
		// The insert slide for the RD03D sensor
		NewTranslation(
			Vec3{20, 0, boxFrontsideT + (W1+RD03DT+W1)/2},
			NewRotation(
				Vec3{0, 0, -90},
				newRd03dInsertSlide(),
			),
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
			Vec3{0, 0, 0},
			newBoxBackside(),
		),
		NewTranslation(
			Vec3{0, 0, boxFrontsideH},
			newBoxFrontside(),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "main", Primitive: newProduct(), Flags: sys.Default},
	})
}
