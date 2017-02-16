package bat

import (
	"image"
	"image/color"
)

func GetImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{253, 253, 253, 0})
	return img
}
