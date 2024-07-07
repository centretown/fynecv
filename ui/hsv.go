package ui

import (
	"fmt"
	"image/color"
	"math"

	"gopkg.in/yaml.v3"
)

const (
	HueRed     float32 = 0.0
	HueYellow  float32 = 60.0
	HueGreen   float32 = 120.0
	HueCyan    float32 = 180.0
	HueBlue    float32 = 240.0
	HueMagenta float32 = 300.0
	HueMax     float32 = 360.0
)

type HSV struct {
	Hue        float32 `yaml:"hue" json:"hue"`
	Saturation float32 `yaml:"saturation" json:"saturation"`
	Value      float32 `yaml:"value" json:"value"`
}

func (hsv *HSV) FromColor(c color.Color) {
	r, g, b, a := c.RGBA()
	rgba := color.NRGBA{
		uint8(r),
		uint8(g),
		uint8(b),
		uint8(a),
	}
	hsv.FromRGB(rgba)
}

func (hsv *HSV) FromRGB(color color.NRGBA) {
	var red float32 = float32(color.R) / 255.0
	var green float32 = float32(color.G) / 255.0
	var blue float32 = float32(color.B) / 255.0

	major := max(red, green, blue)
	delta := major - min(red, green, blue)

	var hue float64
	if delta != 0.0 {
		switch major {
		case red:
			hue = 60.0 * float64((green-blue)/delta)
		case green:
			hue = 60.0 * (float64((blue-red)/delta) + 2.0)
		case blue:
			hue = 60.0 * (float64((red-green)/delta) + 4.0)
		}
	}

	if hue < 0 {
		hue += 360.0
	}

	hsv.Hue = float32(hue)

	if major == 0.0 {
		hsv.Saturation = 0.0
	} else {
		hsv.Saturation = delta / major
	}

	hsv.Value = major
}

func (hsv *HSV) ToRGB() color.NRGBA {
	var major float64 = float64(hsv.Value * hsv.Saturation)
	var minor float64 = major *
		(1.0 -
			math.Abs(
				math.Mod(float64(hsv.Hue)/60.0, 2.0)-1.0))
	var red, green, blue float64
	switch {
	case hsv.Hue < 60.0:
		red = major
		green = minor
	case hsv.Hue < 120.0:
		red = minor
		green = major
	case hsv.Hue < 180.0:
		green = major
		blue = minor
	case hsv.Hue < 240.0:
		green = minor
		blue = major
	case hsv.Hue < 300.0:
		red = minor
		blue = major
	default:
		red = major
		blue = minor
	}

	intensity := float64(hsv.Value) - major
	var color color.NRGBA
	color.R = uint8(math.Round((red + intensity) * 255.0))
	color.G = uint8(math.Round((green + intensity) * 255.0))
	color.B = uint8(math.Round((blue + intensity) * 255.0))
	color.A = 255
	return color
}

func (hsv *HSV) ToGradient(target HSV, index uint16, length uint16) HSV {
	if hsv.Hue > target.Hue {
		target.Hue += 360.0
	}
	ratio := float32(index) / float32(length)
	hue := hsv.Hue + (target.Hue-hsv.Hue)*ratio

	hue = float32(math.Mod(float64(hue), float64(HueMax)))

	saturation := hsv.Saturation + (target.Saturation-hsv.Saturation)*ratio
	value := hsv.Value + (target.Value-hsv.Value)*ratio
	return HSV{hue, saturation, value}
}

func (hsv HSV) MarshalYAML() (interface{}, error) {
	type alias HSV
	var mod alias = alias(hsv)
	mod.Saturation *= 100
	mod.Value *= 100

	node := yaml.Node{}
	err := node.Encode(mod)

	if err != nil {
		return nil, err
	}
	return node, err
}

func (hsv *HSV) UnmarshalYAML(node *yaml.Node) (err error) {
	type alias HSV
	var mod *alias = (*alias)(hsv)
	err = node.Decode(mod)
	if err != nil {
		return err
	}
	mod.Saturation /= 100
	mod.Value /= 100
	return nil
}

func (hsv *HSV) MakeCode() string {
	s := fmt.Sprintf("{%d,%d,%d}",
		int(hsv.Hue*1530/360), int(hsv.Saturation*255), int(hsv.Value*255))
	return s
}
