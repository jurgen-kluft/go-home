package cadlib

import (
	. "github.com/go-gl/mathgl/mgl64"
	. "github.com/ljanyst/ghostscad/primitive"
)

// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------
// Sensor Stick

// The sensor cylindrical connection with the box
const (
	SensorStickCylindricalRadius = 8.0 / 2
)

// NewSensorStick creates the cylindrical tube with the rectangular (hollow) part on top where the sensor is located.
func NewSensorStick(sensorW, sensorL, sensorH, cylinderL float64) Primitive {

	// if the sensor height is less than the cylinder diameter (including the walls), then we will
	// mount the cylindrical part at the bottom of the bottom of the sensor box and not on the end.

	return NewTranslation(
		Vec3{0, 0, sensorH / 2.0},
		NewUnion(
			NewDifference(
				NewBox(sensorW+2*W1, sensorL+2*W1, sensorH),
				// Make it hollow
				NewCube(Vec3{sensorW, sensorL, sensorH - 2*W1}),
				// Cutout the top, so that we can insert the sensor easily
				NewTranslation(
					Vec3{0, 0, sensorH / 2},
					NewCube(Vec3{sensorW, sensorL, 2 * W2}),
				),
				// Cutout for the cylindrical part
				NewTranslation(
					Vec3{0, (sensorL + W1) / 2, 0},
					NewCylinderOnTheYAxis(2*W2, SensorStickCylindricalRadius),
				),
			),

			// Cylindrical part
			NewTranslation(
				Vec3{0, sensorL/2.0 + cylinderL/2.0 + W1, 0},
				NewDifference(
					NewCylinderOnTheYAxis(cylinderL+W1, SensorStickCylindricalRadius+W1),
					NewCylinderOnTheYAxis(sensorL, SensorStickCylindricalRadius),
				),
			),
		),
	)
}

func NewSensorStickLid(sensorW, sensorL, sensorH float64) Primitive {
	return NewTranslation(
		Vec3{0, 0, W1 / 2.0},
		NewUnion(
			NewTranslation(
				Vec3{0, 0, W1 / 2},
				NewBox(sensorW+2*W1, sensorL+2*W1, W1),
			),
			NewTranslation(
				Vec3{0, 0, W2/2 + Rounding/2},
				NewDifference(
					NewBox(sensorW-0.2, sensorL, W2),
					NewBox(sensorW-0.2-W2, sensorL-W2, 2*W2),
				),
			),
		),
	)
}
