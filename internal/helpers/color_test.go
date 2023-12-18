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
		t.Run(tt.in, func(t *testing.T){
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
