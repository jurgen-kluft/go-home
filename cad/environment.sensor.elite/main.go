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
	boxW = pcbBoardW + 2*20
	boxL = pcbBoardL + 2*15
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

// PCB Board
const (
	pcbBoardW      = 90 + (Rounding / 2)
	pcbBoardL      = 90 + (Rounding / 2)
	pcbBoardT      = 1.6                                 // Thickness of the PCB board itself
	pcbBoardMountW = 51.0                                // Mounting width, from hole to hole, for the PCB board
	pcbBoardMountL = 75.0                                // Mounting length, from hole to hole, for the PCB board
	displayMountHR = 2.0 / 2                             // Mounting hole radius for the PCB board
	pcbBoardBH     = 6.0 - pcbBoardT                     // Height of the tallest component on the bottom of the PCB
	pcbBoardTH     = W1                                  // Thickness of the tallest component on the top of the PCB
	pcbBoardH      = pcbBoardBH + pcbBoardTH + pcbBoardT // Total height of the PCB including components
)

// Display
const (
	displayScreenW       = WcsScreenL
	displayScreenL       = WcsScreenW
	displayScreenR       = WcsScreenR
	displayMountingW     = WcsMountingW
	displayMountingL     = WcsMountingL
	displayMountingHoleR = WcsMountingHoleD / 2
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
						Vec3{RBP.X() + 4, RBP.Y() + UsbCHoleDiameter + UsbCHoleRadius, 0},
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

const (
	// --- 1. Dimension Constants (in mm) ---
	pipeLength = 40.0
	pipeOD     = 20.0
	pipeID     = 16.0

	// SEN66 dimensional requirements
	sen66X        = 55.6
	sen66Y        = 26.0
	sen66Z        = 21.7
	wallThickness = 2.0
)

func buildSensorForSEN66() Primitive {

	compODX := sen66X + (wallThickness * 2)
	compODY := sen66Y + (wallThickness * 2)

	// FIX: Make the solid outer block taller by adding an explicit floor buffer
	// into the negative Z quadrant so subtractions cannot delete the floor.
	compODZ := sen66Z + wallThickness

	// --- 2. Solid Outer Compartment Block ---
	cornerRadius := 3.0
	cornerX := (compODX / 2.0) - cornerRadius
	cornerY := (compODY / 2.0) - cornerRadius

	// FIX: Outer hull starts at Z = -2.0 instead of 0 to give a solid structural base plate
	compartmentOuter := NewHull(
		NewTranslation(Vec3{cornerX, cornerY, -2.0}, NewCylinder(compODZ+2.0, cornerRadius)),
		NewTranslation(Vec3{-cornerX, cornerY, -2.0}, NewCylinder(compODZ+2.0, cornerRadius)),
		NewTranslation(Vec3{cornerX, -cornerY, -2.0}, NewCylinder(compODZ+2.0, cornerRadius)),
		NewTranslation(Vec3{-cornerX, -cornerY, -2.0}, NewCylinder(compODZ+2.0, cornerRadius)),
	)

	// --- 3. Inner Pocket via Hull ---
	// Sits at Z = wallThickness (2.0), leaving a clean 4mm solid plastic buffer floor (from Z=-2 to Z=2)
	innerRadius := 1.5
	innerX := (sen66X / 2.0) - innerRadius
	innerY := (sen66Y / 2.0) - innerRadius
	pocketHeight := sen66Z + 10.0 // Extends cleanly out through the ceiling

	sensorPocket := NewHull(
		NewTranslation(Vec3{innerX, innerY, wallThickness}, NewCylinder(pocketHeight, innerRadius)),
		NewTranslation(Vec3{-innerX, innerY, wallThickness}, NewCylinder(pocketHeight, innerRadius)),
		NewTranslation(Vec3{innerX, -innerY, wallThickness}, NewCylinder(pocketHeight, innerRadius)),
		NewTranslation(Vec3{-innerX, -innerY, wallThickness}, NewCylinder(pocketHeight, innerRadius)),
	)

	// --- 4. Wire Feedthrough Hole ---
	// FIX: Controlled diameter hole that cuts through the bottom base plate cleanly
	// without removing the wide shelf that supports the SEN66 housing.
	wireHole := NewTranslation(Vec3{0, 0, -3.0}, NewCylinder(wallThickness+6.0, pipeID/2))

	// --- 5. Side Ventilation Windows ---
	ventWidth := 10.0
	ventLength := compODY + 6.0
	ventHeight := 12.0

	intakeVent := NewTranslation(
		Vec3{-(compODX / 2.0) - 1.0, -ventLength / 2.0, wallThickness},
		NewCube(Vec3{ventWidth, ventLength, ventHeight}),
	)

	exhaustVent := NewTranslation(
		Vec3{(compODX / 2.0) - ventWidth + 1.0, -ventLength / 2.0, wallThickness},
		NewCube(Vec3{ventWidth, ventLength, ventHeight}),
	)

	// --- 6. Assemble the Subtraction Tree ---
	compartment := NewDifference(
		compartmentOuter,
		sensorPocket,
		wireHole,
		intakeVent,
		exhaustVent,
	)

	// --- 7. Hollow Pipe Stalk with Interlocking Flange ---
	// Pipe spans up to Z = 40.0. The inner hollow core terminates cleanly at the flange rim.
	stalk := NewDifference(
		NewCylinder(pipeLength, pipeOD/2),
		NewCylinder(pipeLength+2.0, pipeID/2),
	)

	// --- 8. Structural Assembly Placement ---
	// FIX: Shift the compartment up so its base layer (Z = -2.0) interfaces and buries
	// into the stalk head (Z = 40.0). This overlaps the geometries into a clean, rigid, solid mesh.
	positionedCompartment := NewTranslation(Vec3{0, 0, sen66Z/2.0 + 2.0*wallThickness + pipeLength/2.0}, compartment)

	return NewUnion(stalk, positionedCompartment)
}

func buildCompartmentLidForSEN66() Primitive {
	// Reference dimensions from the compartment layout
	compODX := sen66X + (wallThickness * 2)
	compODY := sen66Y + (wallThickness * 2)

	lidTopThickness := 2.0
	plugDepth := 4.0  // How deep the underside alignment wall goes down into the compartment
	tolerance := 0.25 // The 3D printing safety clearance on all sides

	// Dimensions for the lower inserting plug block
	plugX := sen66X - (tolerance * 2.0)
	plugY := sen66Y - (tolerance * 2.0)

	// --- 1. Top Cover Layer (Matches the main box footprint) ---
	cornerRadius := 3.0
	cornerX := (compODX / 2.0) - cornerRadius
	cornerY := (compODY / 2.0) - cornerRadius

	lidTopPlate := NewHull(
		NewTranslation(Vec3{cornerX, cornerY, 0}, NewCylinder(lidTopThickness, cornerRadius)),
		NewTranslation(Vec3{-cornerX, cornerY, 0}, NewCylinder(lidTopThickness, cornerRadius)),
		NewTranslation(Vec3{cornerX, -cornerY, 0}, NewCylinder(lidTopThickness, cornerRadius)),
		NewTranslation(Vec3{-cornerX, -cornerY, 0}, NewCylinder(lidTopThickness, cornerRadius)),
	)

	// --- 2. Underside Alignment Plug Block ---
	// Centers the cube block underneath the lid plate face

	lidUndersidePlate := NewHull(
		NewTranslation(Vec3{cornerX, cornerY, 0}, NewCylinder(lidTopThickness+1.0, cornerRadius)),
		NewTranslation(Vec3{-cornerX, cornerY, 0}, NewCylinder(lidTopThickness+1.0, cornerRadius)),
		NewTranslation(Vec3{cornerX, -cornerY, 0}, NewCylinder(lidTopThickness+1.0, cornerRadius)),
		NewTranslation(Vec3{-cornerX, -cornerY, 0}, NewCylinder(lidTopThickness+1.0, cornerRadius)),
	)

	// scaling factor to scale it down to the inner dimensions of the compartment, while accounting for the tolerance
	scalar := Vec3{plugX / compODX, plugY / compODY, 1.0}
	scalarInner := Vec3{0.9 * scalar.X(), 0.9 * scalar.Y(), 8.0}
	undersidePlug := NewTranslation(Vec3{0, 0, -plugDepth / 2},
		NewDifference(
			NewScale(scalar, lidUndersidePlate),
			NewScale(scalarInner, lidUndersidePlate),
		),
	)

	// Join them together into a single printable piece
	return NewUnion(
		NewColor("blue", lidTopPlate),
		NewColor("green", undersidePlug),
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
					Vec3{0, (displayScreenL / 2), -3 * W2},
					newBox(displayScreenW, displayScreenL, 6*W2, displayScreenR),
				),
			),
		),

		// NewTranslation(
		// 	Vec3{-4, -((boxW / 2) - (W2 + RD03DW + W2) + (4.2 / 2)) / 2, boxBacksideH},
		// 	NewTranslation(
		// 		LTP,
		// 		newBox(5.0, (boxW/2)-(W2+RD03DW+W2)+(4.2/2), boxBacksideH, 0.1),
		// 	),
		// ),

		// AMOLED display mounting
		NewTranslation(
			Vec3{0, (displayScreenL / 2), boxFrontsideT},
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
		{Name: "SEN66_Sensor", Primitive: NewUnion(
			NewTranslation(
				Vec3{0, 0, pipeLength/2.0 + sen66Z + 2.0 + 2.0*wallThickness},
				buildCompartmentLidForSEN66(),
			),
			buildSensorForSEN66(),
		), Flags: sys.Default},
	})
}
