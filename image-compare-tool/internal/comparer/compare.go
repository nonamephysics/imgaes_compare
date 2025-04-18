package comparer

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func CompareImages(basePath string, comparePath string, tolerance float64) (string, error) {
	baseImg, err := loadImage(basePath)
	if err != nil {
		return "", err
	}

	compareImg, err := loadImage(comparePath)
	if err != nil {
		return "", err
	}

	if baseImg.Bounds() != compareImg.Bounds() {
		return "Images have different dimensions", nil
	}

	differenceCount := 0
	totalPixels := baseImg.Bounds().Dx() * baseImg.Bounds().Dy()

	for y := 0; y < baseImg.Bounds().Dy(); y++ {
		for x := 0; x < baseImg.Bounds().Dx(); x++ {
			if !pixelsEqual(baseImg.At(x, y), compareImg.At(x, y)) {
				differenceCount++
			}
		}
	}

	differencePercentage := (float64(differenceCount) / float64(totalPixels)) * 100

	if differencePercentage == 0 {
		return "Images are equal", nil
	} else if differencePercentage <= tolerance {
		return "Images are equal with tolerance: " + formatFloat(differencePercentage), nil
	} else {
		outputPath := filepath.Join(filepath.Dir(basePath), "highlighted_differences.png")
		err := highlightDifferences(baseImg, compareImg, outputPath)
		if err != nil {
			return "", err
		}
		return "Images are not equal. Differences highlighted in: " + outputPath, nil
	}
}

func pixelsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func loadImage(path string) (image.Image, error) {
	log.Printf("Attempting to load image: %s", path)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var img image.Image
	if filepath.Ext(path) == ".png" {
		img, err = png.Decode(file)
	} else if filepath.Ext(path) == ".jpeg" || filepath.Ext(path) == ".jpg" {
		img, err = jpeg.Decode(file)
	} else {
		return nil, errors.New("unsupported image format")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}
	return img, nil
}

func highlightDifferences(baseImg, compareImg image.Image, outputPath string) error {
	bounds := baseImg.Bounds()
	diffImg := image.NewRGBA(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			basePixel := baseImg.At(x, y)
			comparePixel := compareImg.At(x, y)

			if !pixelsEqual(basePixel, comparePixel) {
				// Highlight differences in red
				diffImg.Set(x, y, color.RGBA{255, 0, 0, 255})
			} else {
				// Copy the original pixel
				diffImg.Set(x, y, basePixel)
			}
		}
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	return png.Encode(outputFile, diffImg)
}
