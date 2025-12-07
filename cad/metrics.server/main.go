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

// Box dimensions
const (
	boxW = lcdW + 2*20
	boxL = lcdL + 2*15
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

// Waveshare ESP32-S3 LCD board with:
// - TF card slot
// - GPIO header pins

// Display, 320x480 LCD from Waveshare
const (
	lcdScreenW       = 34.79 // Actual display width
	lcdScreenL       = 42.62 // Actual display length
	lcdScreenR       = 9.0   // Actual display corner rounding
	lcdW             = 36.79 // Overall width including bezel
	lcdL             = 52.56 // Overall length including bezel
	lcdMountingW     = 31.73 // Mounting hole to hole width
	lcdMountingL     = 47.56 // Mounting hole to hole length
	lcdMountingHoleD = 1.8   // Mounting hole diameter
)

// Display, SH1107 OLED, 128x128
const (
	displayScreenW       = lcdScreenL
	displayScreenL       = lcdScreenW
	displayScreenR       = lcdScreenR
	displayW             = lcdW
	displayL             = lcdL
	displayMountingW     = lcdMountingW
	displayMountingL     = lcdMountingL
	displayMountingHoleD = lcdMountingHoleD
	displayMountingHoleR = lcdMountingHoleD / 2
)

// USB-C hole dimensions
const (
	UsbCHoleDiameter = 11.25 // Radius of the USB-C hole, 1 cm
	UsbCHoleRadius   = UsbCHoleDiameter / 2.0
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
						Vec3{RBP.X(), RBP.Y() + UsbCHoleDiameter + UsbCHoleRadius, 0},
						NewCylinder(boxBacksideH*2, UsbCHoleRadius),
					),
				),
			),
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
		NewTranslation(
			Vec3{-15, boxL/2 - Scd41H, boxBacksideT},
			newScd41InsertSlide(),
		),
		NewTranslation(
			Vec3{-15, -(boxL/2 - W1 - Bh1750H), boxBacksideT},
			newBme280InsertSlide(),
		),
		NewTranslation(
			Vec3{(boxW / 2) - 3*W1, -(W2 + RD03DW + W2), boxBacksideT},
			NewRotation(
				Vec3{0, 0, 90},
				newBh1750InsertSlide(),
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
						newBox(boxW+2*(boxI), boxL+2*(boxI), W2, MainRounding),
					),
				),

				// Pyramid like opening at the top for BH1750 sensor (incoming light)
				NewTranslation(
					Vec3{(boxW / 2) + 0.25, (W2 + RD03DW + W2), -(boxFrontsideH / 2) + Bh1750BottomToSensorMiddle},
					NewRotation(
						Vec3{0, -90, 0},
						newPyramid(Bh1750SensorHeight*3, Bh1750SensorWidth*3, boxFrontsideT+3*W1),
					),
				),

				// Pyramid like opening at the side for BME280 sensor
				NewTranslation(
					Vec3{-15, ((boxL / 2) + (boxFrontsideT-boxI)/2), (boxFrontsideH / 2) - Bme280BottomToSensorMiddle},
					NewRotation(
						Vec3{90, 0, 0},
						newPyramid(Bme280SensorLength*3, Bme280SensorWidth*3, boxFrontsideT+3*W1),
					),
				),

				// Rectangular opening at the side for Scd41 sensor (CO2)
				NewTranslation(
					Vec3{-15, -((boxL / 2) + boxFrontsideT/2), (boxFrontsideH / 2) - Scd41BottomToSensorMiddle},
					newBox(Scd41SensorWidth, boxFrontsideT+4*W1, Scd41SensorHeight, Rounding),
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

		// The insert slide for the RD03D sensor
		NewTranslation(
			Vec3{28, -30, boxFrontsideT + (W1+RD03DT+W1)/2},
			NewRotation(
				Vec3{0, 0, 90},
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
			Vec3{0, 0, boxFrontsideH + 50},
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
