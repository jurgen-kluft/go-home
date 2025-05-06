package yeel

import (
	"fmt"
	"strings"
	"time"
)

type Keyframe interface {
	Expression() (string, error)
}

type Animation []Keyframe

func (a Animation) Expression() (string, error) {
	var exprstrs []string

	for _, v := range a {
		exp, err := v.Expression()
		if err != nil {
			return "", err
		}
		exprstrs = append(exprstrs, exp)
	}
	expstr := strings.Join(exprstrs, ",")
	return expstr, nil
}

// RawKeyframe .
type RawKeyframe struct {
	RawExpression string
}

func (r RawKeyframe) Expression() (string, error) {

	return r.RawExpression, nil
}

type RGBKeyframe struct {
	Duration   time.Duration
	Brightness int
	RGB        RGB
}

func (r RGBKeyframe) Expression() (string, error) {
	if err := validateDuration("duration", r.Duration); err != nil {
		return "", err
	}
	b := Brightness(r.Brightness)
	if err := b.Validate(); err != nil {
		return "", err
	}

	res := fmt.Sprintf("%d,1,%d,%d", durationMillis(r.Duration), r.RGB.Int(), b.Int())
	return res, nil
}

type WhiteKeyframe struct {
	Duration   time.Duration
	Brightness int
}

func (r WhiteKeyframe) Expression() (string, error) {
	if err := validateDuration("duration", r.Duration); err != nil {
		return "", err
	}
	b := Brightness(r.Brightness)
	if err := b.Validate(); err != nil {
		return "", err
	}
	rgb := NewRGBNorm(1, 1, 1)
	res := fmt.Sprintf("%d,1,%d,%d", durationMillis(r.Duration), rgb.Int(), b.Int())
	return res, nil
}

type ColorTempratureKeyframe struct {
	Duration        time.Duration
	Brightness      int
	ColorTemprature int
}

func (r ColorTempratureKeyframe) Expression() (string, error) {
	if err := validateDuration("duration", r.Duration); err != nil {
		return "", err
	}
	b := Brightness(r.Brightness)
	if err := b.Validate(); err != nil {
		return "", err
	}
	ct := ColorTemprature(r.ColorTemprature)
	if err := ct.Validate(); err != nil {
		return "", err
	}

	res := fmt.Sprintf("%d,2,%d,%d", durationMillis(r.Duration), r.ColorTemprature, r.Brightness)
	return res, nil
}

type SleepKeyframe struct {
	Duration time.Duration
}

func (r SleepKeyframe) Expression() (string, error) {
	if err := validateDuration("duration", r.Duration); err != nil {
		return "", err
	}

	res := fmt.Sprintf("%d,7,1,1", durationMillis(r.Duration))
	return res, nil
}
