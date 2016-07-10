package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
)

func expandRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()
	xmax, ymax := rect.Max.X, rect.Max.Y
	var min, max float64 = 255.0, 0.0
	var minr, ming, minb uint8
	var r, g, b, a uint8
	var r32, g32, b32, a32 uint32
	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			r32, g32, b32, _ = img.At(x, y).RGBA()
			r, g, b = uint8(r32), uint8(g32), uint8(b32)
			if min > luminosity(r, g, b) {
				min = luminosity(r, g, b)
				minr, ming, minb = r, g, b
			}
			if max < luminosity(r, g, b) {
				max = luminosity(r, g, b)
			}
		}
	}

	var ratio float64 = 255 / (max - min)
	nrgba := image.NewRGBA(rect)

	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			r32, g32, b32, a32 = img.At(x, y).RGBA()
			r, g, b, a = uint8(r32), uint8(g32), uint8(b32), uint8(a32)
			nrgba.Set(x, y, color.RGBA{uint8(float64(r-minr) * ratio), uint8(float64(g-ming) * ratio), uint8(float64(b-minb) * ratio), uint8(a)})

			//			fmt.Println(uint8(r), uint8(g), uint8(b), uint8(a), " ~~>", minr, ming, minb)
		}
	}

	return nrgba
}

func luminosity(r, g, b uint8) float64 {
	return float64(r)*0.2126 + float64(g)*0.7152 + float64(b)*0.0722
}

func dy(img image.Image, x, y int) float64 {
	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y
	if x == 0 || x == xmax-1 || y == 0 || y == ymax-1 {
		return 1.0
	}

	r0, g0, b0, _ := img.At(x, y-1).RGBA() // returns uint32
	r2, g2, b2, _ := img.At(x, y+1).RGBA() // returns uint32

	l0 := luminosity(uint8(r0), uint8(g0), uint8(b0))
	l2 := luminosity(uint8(r2), uint8(g2), uint8(b2))

	l1 := l2 - l0
	if l1 == 0 {
		return 1.0
	} else {
		return l1
	}
}

func dx(img image.Image, x, y int) float64 {
	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y
	if x == 0 || x == xmax-1 || y == 0 || y == ymax-1 {
		return 1.0
	}
	r0, g0, b0, _ := img.At(x-1, y).RGBA()
	r2, g2, b2, _ := img.At(x+1, y).RGBA()

	l0 := luminosity(uint8(r0), uint8(g0), uint8(b0))
	l2 := luminosity(uint8(r2), uint8(g2), uint8(b2))

	l1 := l2 - l0
	if l1 == 0 {
		return 1.0
	} else {
		return l1
	}
}

func gradient_magnitude(dx, dy float64) float64 {
	return math.Sqrt(dx*dx + dy*dy)
}

func gradient_theta(dx, dy float64) float64 {
	if dx == 1.0 && dy == 1.0 {
		return -200.0
	}

	th := math.Atan2(dy, dx) * (180 / math.Pi)
	return th
}

type Chunk struct {
	x, y      int
	theta     float64
	magnitude float64
}

func main() {
	file, err := os.Open("data/test1.jpg")
	defer file.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "main open file :%v\n", err)
		os.Exit(1)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main read file :%v\n", err)
		os.Exit(1)
	}

	img = expandRGBA(img) // settle it black(0x00) and white(0xFF)

	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y

	var row_array = make([]int, xmax-1)

	rect := img.Bounds()
	nimg := image.NewRGBA(rect)

	var max, maxx, maxy int
	var min, minx, miny int

	// preprocess
	for y := 0; y < ymax; y++ {

		for x := 1; x < xmax; x++ { // x: 1 ~ xmax

			r0, g0, b0, _ := img.At(x-1, y).RGBA() // returns uint32
			r1, g1, b1, _ := img.At(x, y).RGBA()

			row_array[x-1] = int(luminosity(uint8(r0), uint8(g0), uint8(b0)) - luminosity(uint8(r1), uint8(g1), uint8(b1)))

			if max < int(row_array[x-1]) {
				max = int(row_array[x-1])
				maxx, maxy = x, y
			}
			if min > int(row_array[x-1]) {
				min = int(row_array[x-1])
				minx, miny = x, y
			}

			r1, g1, b1, a1 := img.At(x, y).RGBA()
			nimg.Set(x, y, color.RGBA{uint8(r1), uint8(g1), uint8(b1), uint8(a1)})
		}
		nimg.Set(maxx, maxy, color.RGBA{0xFF, 0x0, 0xFF, 0xFF})
		nimg.Set(minx, miny, color.RGBA{0xFF, 0x0, 0x0, 0xFF})
		max = 0
		min = 255
	}

	png.Encode(os.Stdout, nimg)

}
