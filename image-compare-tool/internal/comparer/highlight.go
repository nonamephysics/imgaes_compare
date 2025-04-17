package comparer

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"errors"
)

var ErrImageSizeMismatch = errors.New("image sizes do not match")

// HighlightDifferences takes two images and a file path to save the highlighted image.
func HighlightDifferences(baseImg, compareImg image.Image, outputPath string) error {
	bounds := baseImg.Bounds()
	if bounds != compareImg.Bounds() {
		return ErrImageSizeMismatch
	}

	highlighted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			basePixel := baseImg.At(x, y)
			comparePixel := compareImg.At(x, y)

			if basePixel != comparePixel {
				// Highlight differences
				highlighted.Set(x, y, color.RGBA{255, 0, 0, 255}) // Red for compare image
			} else {
				highlighted.Set(x, y, basePixel) // Keep base image pixel
			}
		}
	}

	// Save the highlighted image
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, highlighted)
}