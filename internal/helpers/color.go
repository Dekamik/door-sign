package helpers

import (
	"fmt"
	"image/color"
)

func ParseHexToColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)

	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits
		c.R *= 17
		c.G *= 17
		c.B *= 17

	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}

	return
}

func HexString(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func Lerp(min color.RGBA, max color.RGBA, value float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(min.R)*(1-value) + float64(max.R)*value),
		G: uint8(float64(min.G)*(1-value) + float64(max.G)*value),
		B: uint8(float64(min.B)*(1-value) + float64(max.B)*value),
		A: 0xff,
	}
}

func LerpHexString(min string, max string, value float64) (string, error) {
	minC, err := ParseHexToColor(min)
	if err != nil {
		return "", err
	}

	maxC, err := ParseHexToColor(max)
	if err != nil {
		return "", err
	}

	c := Lerp(minC, maxC, value)

	return HexString(c), nil
}
