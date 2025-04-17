package utils

import (
	"image"
	"fmt"	
	"image/color"
	"strings"
	"image/png"
	"image/jpeg"
	"os"
)

// LoadImage loads an image from the given file path and returns it as an image.Image.
func LoadImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var img image.Image
	if isPNG(filePath) {
		img, err = png.Decode(file)
	} else if isJPEG(filePath) {
		img, err = jpeg.Decode(file)
	} else {
		return nil, fmt.Errorf("unsupported image format")
	}
	return img, err
}

// isPNG checks if the file has a .png extension.
func isPNG(filePath string) bool {
	return strings.HasSuffix(filePath, ".png")
}

// isJPEG checks if the file has a .jpeg or .jpg extension.
func isJPEG(filePath string) bool {
	return strings.HasSuffix(filePath, ".jpeg") || strings.HasSuffix(filePath, ".jpg")
}

// CalculateDifference calculates the difference between two images and returns the number of differing pixels.
func CalculateDifference(img1, img2 image.Image) (int, error) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1 != bounds2 {
		return 0, fmt.Errorf("images must have the same dimensions")
	}

	differenceCount := 0
	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			if !pixelsEqual(img1.At(x, y), img2.At(x, y)) {
				differenceCount++
			}
		}
	}
	return differenceCount, nil
}

// pixelsEqual compares two colors and returns true if they are equal.
func pixelsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}