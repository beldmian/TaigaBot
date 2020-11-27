package main

import (
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	_, err = fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	return
}

func AddLabel(img *image.RGBA, x, y int, label string, face font.Face) {
	col := color.White
	length := font.MeasureString(face, label)

	point := fixed.Point26_6{X: fixed.Int26_6(x * 64)-length/2, Y: fixed.Int26_6((y+8) * 64)}


	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  point,
	}
	d.DrawString(label)
}
