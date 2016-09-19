package bat

import (
	"image"
	"image/color"
)

func GetImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(1, 1, color.RGBA{128, 128, 128, 128})
	return img
}
