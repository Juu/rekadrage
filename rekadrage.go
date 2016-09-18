package rekadrage

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
)

type Config struct {
	Margin    int // Margin to be added in the output picture
	Tolerance int // Tolerance admitted for the difference with control pixel. Higher tolerance will result in a narrower cropping.
}

func Rekadrage(img image.Image, c Config) image.Image {

	frame := img.Bounds()
	log.Println("Image size:", frame)

	// The pixel at 0-0 is used as control value to check color
	control := img.At(0, 0)

	// Optimization attempt to avoid parsing the whole pixels. May be useful with a big picture.
	frameHeightFromTop(img, &frame, control, c.Tolerance)
	frameHeightFromBottom(img, &frame, control, c.Tolerance)
	frameWidthFromLeft(img, &frame, control, c.Tolerance)
	frameWidthFromRight(img, &frame, control, c.Tolerance)

	log.Println("Found picture at positions:", frame)
	log.Println("Adding margin:", c.Margin)
	frame.Min.X -= c.Margin
	frame.Min.Y -= c.Margin
	frame.Max.X += c.Margin
	frame.Max.Y += c.Margin
	frame = frame.Intersect(img.Bounds())
	log.Println("Frame with margins:", frame)
	imgOut := image.NewRGBA(frame)
	draw.Draw(imgOut, frame, img, image.Point{frame.Min.X, frame.Min.Y}, 0)

	return imgOut
}

func frameHeightFromTop(img image.Image, f *image.Rectangle, c color.Color, t int) bool {
	for y := f.Min.Y; y < f.Max.Y; y++ {
		for x := f.Min.X; x < f.Max.X; x++ {
			if !checkColorMatch(img.At(x, y), c, t) {
				f.Min.Y = y
				return true
			}
		}
	}
	return false
}

func frameHeightFromBottom(img image.Image, f *image.Rectangle, c color.Color, t int) bool {
	for y := f.Max.Y - 1; y >= f.Min.Y; y-- {
		for x := f.Min.X; x < f.Max.X; x++ {
			if !checkColorMatch(img.At(x, y), c, t) {
				f.Max.Y = y + 1
				return true
			}
		}
	}
	return false
}

func frameWidthFromLeft(img image.Image, f *image.Rectangle, c color.Color, t int) bool {
	for x := f.Min.X; x < f.Max.X; x++ {
		for y := f.Min.Y; y < f.Max.Y; y++ {
			if !checkColorMatch(img.At(x, y), c, t) {
				f.Min.X = x
				return true
			}
		}
	}
	return false
}

func frameWidthFromRight(img image.Image, f *image.Rectangle, c color.Color, t int) bool {
	for x := f.Max.X - 1; x >= f.Min.X; x-- {
		for y := f.Min.Y; y < f.Max.Y; y++ {
			if !checkColorMatch(img.At(x, y), c, t) {
				f.Max.X = x + 1
				return true
			}
		}
	}
	return false
}

func checkColorMatch(c1, c2 color.Color, t int) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	diff := math.Abs(float64(r1)-float64(r2)) + math.Abs(float64(g1)-float64(g2)) + math.Abs(float64(b1)-float64(b2))
	return int(diff) <= t
}
