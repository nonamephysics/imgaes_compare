package comparer

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
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
	log.Printf("Loading image: %s", path)

	// Check if the file is an SVG and convert it to PNG
	if filepath.Ext(path) == ".svg" {
		var err error
		path, err = convertSVGToPNG(path)
		if err != nil {
			return nil, fmt.Errorf("failed to convert SVG to PNG: %w", err)
		}
	}

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

func convertSVGToPNG(svgPath string) (string, error) {
	log.Printf("Converting SVG to PNG using headless browser: %s", svgPath)

	// Define the output PNG path
	pngPath := svgPath[:len(svgPath)-len(filepath.Ext(svgPath))] + ".png"

	// Use chromedp to render the SVG and capture it as a PNG
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for the operation
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("file://"+svgPath),
		chromedp.FullScreenshot(&buf, 100),
	)
	if err != nil {
		return "", fmt.Errorf("failed to render SVG with headless browser: %w", err)
	}

	// Write the PNG data to the output file
	err = ioutil.WriteFile(pngPath, buf, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write PNG file: %w", err)
	}

	log.Printf("SVG successfully converted to PNG: %s", pngPath)
	return pngPath, nil
}
