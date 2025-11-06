package imgutil

import (
	"image"
	"image/color"
)

func isWhite(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r > 150 || g > 150 || b > 150
}

func GetTwoColumns(img image.Image, column int) uint16 {
	var columns uint16

	for i := 0; i < 7; i++ {
		if isWhite(img.At(column, i)) {
			columns = columns | 1<<i
		}

		if isWhite(img.At(column+1, i)) {
			columns = columns | 1<<(i+7)
		}
	}

	return columns
}
