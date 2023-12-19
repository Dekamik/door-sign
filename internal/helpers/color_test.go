package helpers

import (
	"errors"
	"image/color"
	"testing"
)

var hexToColorTests = []struct {
	in  string
	out color.RGBA
	err error
}{
	{"#112233", color.RGBA{R: 17, G: 34, B: 51, A: 255}, nil},
	{"#123", color.RGBA{R: 17, G: 34, B: 51, A: 255}, nil},
	{"#000233", color.RGBA{R: 0, G: 2, B: 51, A: 255}, nil},
	{"#023", color.RGBA{R: 0, G: 34, B: 51, A: 255}, nil},
	{"invalid", color.RGBA{R: 0, G: 0, B: 0, A: 255}, errors.New("input does not match format")},
	{"#abcd", color.RGBA{R: 0, G: 0, B: 0, A: 255}, errors.New("invalid length, must be 7 or 4")},
	{"#-12", color.RGBA{R: 0, G: 0, B: 0, A: 255}, errors.New("expected integer")},
}

func TestParseHexToColor(t *testing.T) {
	for _, tt := range hexToColorTests {
		t.Run(tt.in, func(t *testing.T) {
			c, err := ParseHexToColor(tt.in)
			if err == nil && c != tt.out {
				t.Errorf("got %q want %q", c, tt.out)
			}
			if err != nil && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("got %q want %q", err, tt.err)
			}
		})
	}
}

var hexStringTests = []struct {
	in  color.RGBA
	out string
}{
	{color.RGBA{R: 17, G: 34, B: 51, A: 255}, "#112233"},
	{color.RGBA{R: 0, G: 2, B: 51, A: 255}, "#000233"},
	{color.RGBA{R: 0, G: 34, B: 51, A: 255}, "#002233"},
}

func TestHexString(t *testing.T) {
	for _, tt := range hexStringTests {
		t.Run(tt.out, func(t *testing.T) {
			s := HexString(tt.in)
			if s != tt.out {
				t.Errorf("got %q want %q", s, tt.out)
			}
		})
	}
}

var lerpTests = []struct {
	min      color.RGBA
	max      color.RGBA
	value    float64
	expected color.RGBA
}{
	{
		color.RGBA{R: 0, G: 0, B: 0, A: 0},
		color.RGBA{R: 255, G: 255, B: 255, A: 255},
		0.5,
		color.RGBA{R: 127, G: 127, B: 127, A: 255},
	},
	{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		0.5,
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
	},
	{
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
		0.5,
		color.RGBA{R: 0, G: 0, B: 0, A: 255},
	},
}

func TestLerp(t *testing.T) {
	for _, tt := range lerpTests {
		t.Run("Lerp", func(t *testing.T) {
			c := Lerp(tt.min, tt.max, tt.value)
			if c != tt.expected {
				t.Errorf("got %q want %q", c, tt.expected)
			}
		})
	}
}

var lerpHexStringTests = []struct {
	min      string
	max      string
	value    float64
	expected string
}{
	{"#000000", "#ffffff", 0.5, "#7f7f7f"},
	{"#000", "#fff", 0.5, "#7f7f7f"},
	{"#000000", "#000000", 0.5, "#000000"},
	{"#fff", "#000", 0.5, "#7f7f7f"},
}

func TestLerpHexString(t *testing.T) {
	for _, tt := range lerpHexStringTests {
		t.Run(tt.expected, func(t *testing.T) {
			s, _ := LerpHexString(tt.min, tt.max, tt.value)
			if s != tt.expected {
				t.Errorf("got %q want %q", s, tt.expected)
			}
		})
	}
}
