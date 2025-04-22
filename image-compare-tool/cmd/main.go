package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"image-compare-tool/internal/comparer"
)

func main() {
	baseImagePath := flag.String("base", "", "Path to the base image")
	compareImagePath := flag.String("compare", "", "Path to the image to compare")
	tolerance := flag.Float64("tolerance", 0, "Tolerance level in percent")
	clean := flag.Bool("clean", false, "Remove intermediate and output files after comparison")

	flag.Parse()

	if *baseImagePath == "" || *compareImagePath == "" {
		log.Fatal("Both base and compare image paths must be provided.")
	}

	result, err := comparer.CompareImages(*baseImagePath, *compareImagePath, *tolerance)
	if err != nil {
		log.Fatalf("Error comparing images: %v", err)
	}
	fmt.Println("Comparison result:", result)

	// Handle cleanup if the clean flag is set
	if *clean {
		cleanupFiles(*baseImagePath, *compareImagePath, result)
	}
}

func cleanupFiles(baseImagePath, compareImagePath, result string) {
	// Remove converted PNG files if the original files were SVG
	if removeIfConverted(baseImagePath) {
		fmt.Printf("Removed converted PNG for base image: %s\n", baseImagePath)
	}
	if removeIfConverted(compareImagePath) {
		fmt.Printf("Removed converted PNG for compare image: %s\n", compareImagePath)
	}

	// Remove the highlighted differences image if images are different
	if resultContainsHighlight(result) {
		highlightedImagePath := getHighlightedImagePath(baseImagePath)
		if err := os.Remove(highlightedImagePath); err == nil {
			fmt.Printf("Removed highlighted differences image: %s\n", highlightedImagePath)
		} else {
			fmt.Printf("Failed to remove highlighted differences image: %s\n", highlightedImagePath)
		}
	}
}

func removeIfConverted(imagePath string) bool {
	if filepath.Ext(imagePath) == ".svg" {
		convertedPNGPath := imagePath[:len(imagePath)-len(filepath.Ext(imagePath))] + ".png"
		if _, err := os.Stat(convertedPNGPath); err == nil {
			os.Remove(convertedPNGPath)
			return true
		}
	}
	return false
}

func resultContainsHighlight(result string) bool {
	return result != "" && resultContains(result, "highlighted in:")
}

func getHighlightedImagePath(baseImagePath string) string {
	return filepath.Join(filepath.Dir(baseImagePath), "highlighted_differences.png")
}

func resultContains(result, substring string) bool {
	return filepath.Base(result) == substring
}
