package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	var size int
	var dest string
	flag.IntVar(&size, "size", 64, "width and height of the PNG in pixels")
	flag.StringVar(&dest, "dest", "transparent.png", "destination file")
	flag.Parse()

	if size <= 0 {
		fmt.Fprintln(os.Stderr, "--size must be greater than zero")
		os.Exit(2)
	}

	file, err := os.Create(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create transparent.png: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if err := png.Encode(file, image.NewNRGBA(image.Rect(0, 0, size, size))); err != nil {
		fmt.Fprintf(os.Stderr, "encode transparent.png: %v\n", err)
		os.Exit(1)
	}
}
