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

// The Button has to hold the following inside:
// - Battery
// - ESP8266 microcontroller board with the latch circuit
// - Button is positioned on the top of the case

// Button dimensions
// - Diameter: 50 mm
// - Height: ~20 mm

// battery dimensions, ZCLP800
const (
	batteryWidth    = 30.0
	batteryLength   = 36.6
	batteryHeight   = 8.2
	batteryDiameter = 47.0
	batteryRadius   = batteryDiameter / 2.0
)

const (
	ButtonDiameter    = batteryDiameter
	ButtonRadius      = ButtonDiameter / 2.0
	ButtonHeight      = 15.0
	PressButtonRadius = 15.0
)

// Esp8266 board dimensions
const (
	switchBoardLength = 40.0
	switchBoardWidth  = 21.0
	switchBoardHeight = 6.8 // The battery socket should be externally accessible, since we don't want the extra height
)

// newButton creates the button primitive that holds the ESP8266 board and battery inside
func newButton() Primitive {
	buttonHeight := batteryHeight + W1 + switchBoardHeight + W1
	return NewTranslation(
		Vec3{0, 0, (buttonHeight / 2.0)},
		NewDifference(
			NewCylinder(buttonHeight, ButtonRadius+W2),
			NewTranslation(
				Vec3{0, 0, W1},
				NewCylinder(buttonHeight, ButtonRadius),
			),
		),
	)
}

func newBatterySeparator() Primitive {
	// A separator to keep the battery and the ESP8266 board apart, so they don't short circuit each other
	// The battery is 8.2 mm high, so we can make the separator 10 mm high to give some clearance
	// It is basically a disc (thickness = W1) with some holes for the wires to pass through and some holes
	// and structure for the ESP8266 board to latch onto

	legVector := Vec3{ButtonRadius - (W2 / 2.0) - 0.2, 0, -batteryHeight / 2.0}

	return NewUnion(
		NewTranslation(
			Vec3{0, 0, W1 / 2.0},
			NewDifference(
				NewCylinder(W1, ButtonRadius),
				// Hole for the battery wire to pass through
				NewRotation(
					Vec3{0, 0, 0},
					NewTranslation(
						Vec3{ButtonRadius - (W1 * 6.0), 0, 0},
						NewCylinder(W2, W1*4.0),
					),
				),
			),

			// Structure for the ESP8266 board to latch onto
			NewColor("red",
				NewTranslation(
					Vec3{0, (switchBoardLength / 2.0) - (7.4 / 2.0), W1},
					NewBox(switchBoardWidth, 7.4, W2),
				),
				NewTranslation(
					Vec3{0, -(switchBoardLength / 2.0) + (1.5 / 2.0), W1},
					NewBox(switchBoardWidth, 1.5, W2),
				),
			),

			// 4 Legs (basically a box with height = batteryHeight, width = W1, length = W1) for the separator to stand on
			NewColor("blue",
				NewTranslation(
					legVector,
					NewBox(W2, W2, batteryHeight),
				),
				NewRotation(
					Vec3{0, 0, 90},
					NewTranslation(
						legVector,
						NewBox(W2, W2, batteryHeight),
					),
				),
				NewRotation(
					Vec3{0, 0, 180},
					NewTranslation(
						legVector,
						NewBox(W2, W2, batteryHeight),
					),
				),
				NewRotation(
					Vec3{0, 0, 270},
					NewTranslation(
						legVector,
						NewBox(W2, W2, batteryHeight),
					),
				),
			),
		),
	)
}

func newButtonPress() Primitive {
	ringHeight := ButtonHeight
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, ringHeight},
			NewColor("lightgreen",
				NewIntersection(
					NewDifference(
						NewScale(Vec3{1.0, 1.0, 0.5}, NewSphere(ButtonRadius+W2+W2)),
						NewScale(Vec3{1.0, 1.0, 0.5}, NewSphere(ButtonRadius+W2)),
						NewTranslation(
							Vec3{0, 0, -ButtonRadius / 2.0},
							NewBox(2*ButtonDiameter, 2*ButtonDiameter, ButtonRadius),
						),
					),

					// Circular cutout the top for the actual button
					NewCylinder(ButtonRadius+W2+W2, PressButtonRadius-0.2),
				),
			),
		),
	)
}

func newButtonLid() Primitive {
	//ringHeight := ButtonHeight
	shaftHeight := ((ButtonRadius + W2) * 0.5) - 2.7
	return NewUnion(
		NewTranslation(
			Vec3{0, 0, 0},
			NewColor("yellow",
				NewDifference(
					NewScale(Vec3{1.0, 1.0, 0.5}, NewSphere(ButtonRadius+W2+W2)),
					NewScale(Vec3{1.0, 1.0, 0.5}, NewSphere(ButtonRadius+W2)),
					NewTranslation(
						Vec3{0, 0, -ButtonRadius / 2.0},
						NewBox(2*ButtonDiameter, 2*ButtonDiameter, ButtonRadius),
					),

					// Circular cutout the top for the actual button
					NewCylinder(ButtonRadius+W2+W2, PressButtonRadius),
				),
			),
		),
		NewColor("orange",
			NewTranslation(
				Vec3{0, 0, shaftHeight / 2.0},
				NewRing(PressButtonRadius+W1+W1, PressButtonRadius+W1, shaftHeight),
			),
		),
		// NewTranslation(
		// 	Vec3{0, 0, ringHeight / 2.0},
		// 	NewColor("purple", NewRing(ButtonRadius+W2+W2, ButtonRadius+W2, ringHeight)),
		// ),
	)
}

func newCombineAll() Primitive {
	// Return a setup of all the parts together for review
	// Move them apart for better viewing
	apart := 50.0
	return NewList(
		NewTranslation(
			Vec3{-apart, 0, 0},
			newButton(),
		),
		NewTranslation(
			Vec3{apart, 0, 0},
			newButtonLid(),
		),
		NewTranslation(
			Vec3{0, -apart, apart},
			newButtonPress(),
		),
		NewTranslation(
			Vec3{0, apart, apart},
			newBatterySeparator(),
		),
		// NewTranslation(
		// 	Vec3{0, 0, apart},
		// 	NewColor("blue",
		// 		NewPlate(
		// 			[]Vec3{
		// 				{-switchBoardLength / 2.0, -switchBoardWidth / 2.0, 0},
		// 				{switchBoardLength / 2.0, -switchBoardWidth / 2.0, 0},
		// 				{switchBoardLength / 2.0, switchBoardWidth / 2.0, 0},
		// 				{-switchBoardLength / 2.0, switchBoardWidth / 2.0, 0},
		// 			},
		// 			5.0,
		// 			2.0,
		// 			3.0,
		// 		),
		// 	),
		// ),
		// NewTranslation(
		// 	Vec3{0, 0, -apart},
		// 	NewColor("green", NewBar(batteryLength, batteryWidth, batteryHeight, 5.0)),
		// ),
	)
}

func main() {
	sys.Initialize()
	sys.RenderMultiple([]sys.Shape{
		{Name: "button", Primitive: newButton(), Flags: sys.Default},
		{Name: "button_lid", Primitive: newButtonLid(), Flags: sys.Default},
		{Name: "for-review-only", Primitive: newCombineAll(), Flags: sys.None},
	})
}
