package bot

import (
	"image/color"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	out, err := ParseHexColor("0102AB")
	if err != nil {
		t.Error(err)
	}
	exc := color.RGBA{R: 1, G: 2, B: 171, A: 255}
	if out != exc {
		t.Error("Colors are not equal \n need: ", exc, "\n out: ", out)
	}
}
