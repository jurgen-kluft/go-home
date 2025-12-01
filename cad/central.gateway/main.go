package main

import (
	"github.com/ljanyst/ghostscad/lib/shapes"
	"github.com/ljanyst/ghostscad/sys"

	. "github.com/go-gl/mathgl/mgl64"
	. "github.com/ljanyst/ghostscad/primitive"
)

// TODO
//
// - External Antenna for ESP32
// - Ventilation holes ?
// - Remove USB-C round cutout, not necessary if we use Waveshare board

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
	boxW = PCB_W + 30
	boxL = PCB_L
	boxH = boxBottomH + boxTopH
	boxT = WT
	boxI = 1.25
)

// Bottom of the box
const (
	boxBottomT = boxT
	boxBottomH = 25.0
)

// Top of the box
const (
	boxTopT = boxT
	boxTopH = 15.0
)

// LilyGo ETH PCB board measurements
const (
	LiliGo_PCB_W          = 28.0
	LiliGo_PCB_L          = 59.5
	LiliGo_PCB_MountingHR = 1.5  // Radius of the mounting holes
	LiliGo_PCB_H2HL       = 54.0 // Distance between the mounting holes along the length
	LiliGo_PCB_H2HW       = 23.0 // Distance between the mounting holes along the width
)

// Waveshare ETH PCB board measurements
const (
	WS_PCB_W          = 21.0
	WS_PCB_L          = 72.8
	WS_PCB_MountingHR = 1.6   // Radius of the mounting holes
	WS_PCB_H2HL       = 54.15 // Distance between the mounting holes along the length
	WS_PCB_H2HW       = 18.25 // Distance between the mounting holes along the width
)

// ETH PCB board measurements
const (
	PCB_W          = WS_PCB_W
	PCB_L          = WS_PCB_L
	PCB_MountingHR = WS_PCB_MountingHR
	PCB_H2HL       = WS_PCB_H2HL
	PCB_H2HW       = WS_PCB_H2HW
	PCB_HOFFSET    = 1.5 // Extra height offset of the PCB from the bottom of the box
)

// Inlay magnets for holding the box bottom and top together
// - 5x2 mm Round Magnet, Micro Neodymium Magnet
const (
	magnetDiameter    = 5.0
	magnetR           = magnetDiameter / 2.0
	magnetH           = 2.1
	magnetInlayRadius = magnetR + W1
)

// USB-C hole dimensions
const (
	UsbCHoleDiameter = 11.6 // Radius of the USB-C hole, 1.16 cm
	UsbCHoleRadius   = UsbCHoleDiameter / 2.0
)

const (
	UsbCablePlugWidth  = 11.0
	UsbCablePlugHeight = 5.25
)

// Rotary Encoder cutout dimensions
const (
	RotaryEncoderW = 12.0
	RotaryEncoderH = 12.0
	RotaryEncoderD = 7.0
	RotaryEncoderR = RotaryEncoderD / 2.0
)

// ETHERNET RJ45 hole dimensions
const (
	EthernetHoleW = 16.5 // 14.0
	EthernetHoleH = 13.0 // 11.0
	EthernetHoleO = 5.0  // Height offset from the bottom of the box backside, plus some extra clearance for USB-C plug
)

// Display, 128x128, SH1107 OLED
const (
	sh1107ScreenW       = 37.3 // Actual display width
	sh1107ScreenL       = 34.0 // Actual display length
	sh1107ScreenR       = 0.5  // Actual display corner rounding
	sh1107W             = 47.1 // Overall width including bezel
	sh1107L             = 34.1 // Overall length including bezel
	sh1107MountingW     = 42   // Mounting hole to hole width
	sh1107MountingL     = 29   // Mounting hole to hole length
	sh1107MountingHoleD = 2.0  // Mounting hole diameter
)

// Display, SH1107 OLED, 128x128
const (
	displayScreenW       = sh1107ScreenW
	displayScreenL       = sh1107ScreenL
	displayScreenR       = sh1107ScreenR
	displayMountingW     = sh1107MountingW
	displayMountingL     = sh1107MountingL
	displayMountingHoleR = sh1107MountingHoleD / 2
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

func putOnXAxis(p Primitive) Primitive {
	return NewRotation(Vec3{0, 90, 0}, p)
}

func putOnYAxis(p Primitive) Primitive {
	return NewRotation(Vec3{90, 0, 0}, p)
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

// The backside of the box:
// - PCB inlay
// - Sensor insert-slides
// - USB-C connector cutout
func newBoxBottomPart() Primitive {
	V := Vec3{(boxW / 2), (boxL / 2), 0}
	CurrentLen := V.Len()
	NewLen := CurrentLen - MainRounding + 1

	LTP := V.Mul(NewLen / CurrentLen)
	RTP := Vec3{-LTP.X(), LTP.Y(), 0}
	LBP := Vec3{LTP.X(), -LTP.Y(), 0}
	RBP := Vec3{-LTP.X(), -LTP.Y(), 0}

	return NewUnion(
		NewTranslation(
			Vec3{0, 0, boxBottomH / 2},
			NewRender(
				10,
				NewDifference(
					NewUnion(
						NewDifference(
							newBox(boxW+2*(boxT+boxI), boxL+2*(boxT+boxI), boxBottomH, MainRounding),
							NewTranslation(
								Vec3{0, 0, boxBottomT},
								newBox(boxW, boxL, boxBottomH, MainRounding),
							),
						),
						NewTranslation(
							Vec3{0, 0, boxBottomT},
							newWall(boxI, boxW+2*boxI, boxL+2*boxI, boxBottomH, MainRounding),
						),
					),
					// ETHERNET RJ45 connector cutout
					NewTranslation(
						Vec3{0, boxL / 2, -(boxBottomH / 2) + (EthernetHoleH / 2) + boxBottomT + EthernetHoleO + PCB_HOFFSET},
						putOnYAxis(NewCube(Vec3{EthernetHoleW, EthernetHoleH, 2 * W2})),
					),
					// USB-C cable plug cutout
					NewTranslation(
						Vec3{0, boxL / 2, -(boxBottomH / 2) + (UsbCablePlugHeight / 2) + boxBottomT + PCB_HOFFSET},
						putOnYAxis(NewCube(Vec3{UsbCablePlugWidth, UsbCablePlugHeight, 2 * W2})),
					),
					// TODO Rotary Encoder (circular) cutout
					NewTranslation(
						Vec3{(boxW / 2) - W1 + (3*W2)/2, 0, W1},
						putOnXAxis(NewCylinder(RotaryEncoderH*3, RotaryEncoderR)),
					),
					// USB-C connector cutout
					NewTranslation(
						Vec3{RBP.X(), RBP.Y() + UsbCHoleDiameter + UsbCHoleRadius, 0},
						putOnXAxis(NewCylinder(boxBottomH*2, UsbCHoleRadius)),
					),
				),
				// DEBUG; display the PCB
				//newBox(PCB_W, PCB_L, 2.0, Rounding),
			),
		),

		// Magnet inlays
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				LTP,
				newMagnetInlay(boxBottomH-boxBottomT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				RTP,
				newMagnetInlay(boxBottomH-boxBottomT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				LBP,
				newMagnetInlay(boxBottomH-boxBottomT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxBottomT},
			NewTranslation(
				RBP,
				newMagnetInlay(boxBottomH-boxBottomT),
			),
		),
	)
}

func newMountingNail(lowerRadius, lowerHeight, nailRadius, nailHeight float64) Primitive {
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, lowerHeight / 2},
			NewCylinder(lowerHeight, lowerRadius),
		),
		NewTranslation(
			Vec3{0, 0, lowerHeight + (nailHeight / 2)},
			NewCylinder(lowerHeight+nailHeight, nailRadius),
		),
	)
}

// newMounting creates mounting nails for a PCB board that has mounting holes.
func newMounting(w, l, hr, supportHeight float64) Primitive {
	// Create mounting nails at the 4 corners of the PCB
	// A mounting nail has two parts, the lower part (thicker) and the upper part that has the radius of the mounting hole.
	holeRadius := hr
	mountingWidth := w
	mountingLength := l
	mountingRadius := hr * 1.75
	nailLength := 2 * W1
	return NewUnion(
		NewTranslation(
			Vec3{-mountingWidth / 2, -mountingLength / 2, 0},
			newMountingNail(mountingRadius, supportHeight, holeRadius, nailLength),
		),
		NewTranslation(
			Vec3{mountingWidth / 2, -mountingLength / 2, 0},
			newMountingNail(mountingRadius, supportHeight, holeRadius, nailLength),
		),
		NewTranslation(
			Vec3{-mountingWidth / 2, mountingLength / 2, 0},
			newMountingNail(mountingRadius, supportHeight, holeRadius, nailLength),
		),
		NewTranslation(
			Vec3{mountingWidth / 2, mountingLength / 2, 0},
			newMountingNail(mountingRadius, supportHeight, holeRadius, nailLength),
		),
	)
}

// The frontside of the box, also has 4 sides of a certain height:
// - RD03D sensor insert-slide (upper part)
// - OLED display cutout (lower part)
func newBoxTopPart() Primitive {
	V := Vec3{(boxW / 2), (boxL / 2), 0}
	CurrentLen := V.Len()
	NewLen := CurrentLen - MainRounding + 1

	LTP := V.Mul(NewLen / CurrentLen)
	RTP := Vec3{-LTP.X(), LTP.Y(), LTP.Z()}
	LBP := Vec3{LTP.X(), -LTP.Y(), LTP.Z()}
	RBP := Vec3{-LTP.X(), -LTP.Y(), LTP.Z()}

	return NewUnion(
		NewTranslation(
			Vec3{0, 0, boxTopH / 2},
			NewDifference(
				NewDifference(
					NewDifference(
						newBox(boxW+2*(boxTopT+boxI), boxL+2*(boxTopT+boxI), boxTopH, MainRounding),
						NewTranslation(
							Vec3{0, 0, boxTopT},
							newBox(boxW, boxL, boxTopH, MainRounding),
						),
					),
					NewTranslation(
						Vec3{0, 0, boxTopH / 2},
						newBox(boxW+2*(boxI), boxL+2*(boxI), W2, MainRounding),
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
			Vec3{0, 0, boxTopT},
			newMounting(displayMountingW, displayMountingL, displayMountingHoleR, 0),
		),
		// Magnet inlays
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				LTP,
				newMagnetInlay(boxTopH-boxTopT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				RTP,
				newMagnetInlay(boxTopH-boxTopT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				LBP,
				newMagnetInlay(boxTopH-boxTopT),
			),
		),
		NewTranslation(
			Vec3{0, 0, boxT},
			NewTranslation(
				RBP,
				newMagnetInlay(boxTopH-boxTopT),
			),
		),
	)
}

func newProduct() Primitive {
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, 0},
			newBoxBottomPart(),
		),
		NewTranslation(
			Vec3{0, 0, boxTopH + 50},
			newBoxTopPart(),
		),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "for-review-only", Primitive: newProduct(), Flags: sys.Default},
	})
}
