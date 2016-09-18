package main

import (
	"errors"
	"github.com/juu/rekadrage"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	c := rekadrage.Config{Margin: 10, Tolerance: 1000}

	if len(os.Args) != 2 {
		log.Fatal("Expected one argument: fileName")
	}
	fileName := os.Args[1]

	processFile(fileName, c)
}

func processFile(fileName string, c rekadrage.Config) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Printf("Couldn't open file: %v\n", err)
		return
	}
	defer file.Close()

	outPath := "out_" + fileName
	fo, err := os.Create(outPath)
	if err != nil {
		log.Printf("Cannot create output file '", outPath, "':", err)
	}

	img, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Couldn't decode file: %v\n", err)
		return
	}
	log.Printf("File %s opened, type: %s", fileName, format)

	imgOut := rekadrage.Rekadrage(img, c)

	// Finally save the output picture
	switch filepath.Ext(outPath) {
	case ".png":
		err = png.Encode(fo, imgOut)
	case ".jpg":
		err = jpeg.Encode(fo, imgOut, nil)
	default:
		err = errors.New("Unsupported format: " + filepath.Ext(outPath))
	}
	if err != nil {
		log.Println("Error saving output:", err)
		return
	}
	log.Println("Image saved to", outPath)
}
