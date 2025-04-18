package main

import (
	"flag"
	"fmt"
	"log"

	"image-compare-tool/internal/comparer"
)

func main() {
	baseImagePath := flag.String("base", "", "Path to the base image")
	compareImagePath := flag.String("compare", "", "Path to the image to compare")
	tolerance := flag.Float64("tolerance", 0, "Tolerance level in percent")

	flag.Parse()

	if *baseImagePath == "" || *compareImagePath == "" {
		log.Fatal("Both base and compare image paths must be provided.")
	}

	result, err := comparer.CompareImages(*baseImagePath, *compareImagePath, *tolerance)
	if err != nil {
		log.Fatalf("Error comparing images: %v", err)
	}
	fmt.Println("Comparison result:", result)
}
