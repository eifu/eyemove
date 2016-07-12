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

func cutoffRGBA(img image.Image) (*image.RGBA, uint32) {
	rect := img.Bounds()
	nimg := image.NewRGBA(rect)
	var acc, ave, c0 uint32

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			acc = acc + c0&0xFF
		}
	}
	ave = acc / uint32(rect.Max.X*rect.Max.Y)

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			if c0&0xFF > ave {
				nimg.Set(x, y, color.Gray{uint8(ave)})
			} else {
				nimg.Set(x, y, color.Gray{uint8(c0)})
			}
		}
	}
	return nimg, ave
}

func expandRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()

	var min, max uint8 = 0xFF, 0
	var c0 uint32
	var c1 uint8

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			if min > uint8(c0) {
				min = uint8(c0)
			}
			if max < uint8(c0) {
				max = uint8(c0)
			}
		}
	}
	var ratio float64 = 0xFF / float64(max-min)
	nrgba := image.NewRGBA(rect)

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			c1 = uint8(float64(uint8(c0)-min) * ratio)
			nrgba.Set(x, y, color.Gray{c1})
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

func rowIterate(img image.Image, ave uint32) ([]int, []int) {
	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y

	var maxlist = make([]int, ymax)
	var minlist = make([]int, ymax)
	var row_array = make([]int, xmax-1)

	var max, maxx, min, minx int
	var c0, c1 uint32
	for y := 0; y < ymax; y++ {
		max, min = 0, 0xFF
		for x := 1; x < xmax; x++ {
			c0, _, _, _ = img.At(x-1, y).RGBA()
			c1, _, _, _ = img.At(x, y).RGBA()
			row_array[x-1] = int(c1) - int(c0)
			if c0&0xFF != ave && c1&0xFF != ave {
				if max < row_array[x-1] {
					max, maxx = row_array[x-1], x
				}
				if min > row_array[x-1] {
					min, minx = row_array[x-1], x
				}
			}

		}
		maxlist[y], minlist[y] = maxx, minx
	}
	return maxlist, minlist
}

func colIterate(img image.Image, ave uint32) ([]int, []int) {
	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y

	var maxlist = make([]int, xmax)
	var minlist = make([]int, xmax)
	var col_array = make([]int, ymax-1)

	var max, maxy, min, miny int
	var c0, c1 uint32
	for x := 0; x < xmax; x++ {
		max, min = 0, 0xFF
		for y := 1; y < ymax; y++ {
			c0, _, _, _ = img.At(x, y-1).RGBA()
			c1, _, _, _ = img.At(x, y).RGBA()
			col_array[y-1] = int(c1) - int(c0)
			if c0&0xFF != ave && c1&0xFF != ave {
				if max < col_array[y-1] {
					max, maxy = col_array[y-1], y
				}
				if min > col_array[y-1] {
					min, miny = col_array[y-1], y
				}
			}
		}
		maxlist[x], minlist[x] = maxy, miny
	}
	return maxlist, minlist
}

type Chunk struct {
	x, y      int
	theta     float64
	magnitude float64
}

type Point struct {
	x, y int
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
	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y

	// settle it black(0x00) and white(0xFF)
	nimg := expandRGBA(img)

	// cut off pixels below the average color
	nimg, ave := cutoffRGBA(nimg)

	// iterate lines
	rowmaxlist, rowminlist := rowIterate(nimg, ave)
	colmaxlist, colminlist := colIterate(nimg, ave)

	// put colors to image every column
	for x := 0; x < xmax; x++ {
		colmaxy, colminy := colmaxlist[x], colminlist[x]
		nimg.Set(x, colmaxy, color.RGBA{0x3F, 0xBF, 0x3F, 0xFF})
		nimg.Set(x, colminy, color.RGBA{0x19, 0x4C, 0x19, 0xFF})
	}

	// put colors to image every row
	for y := 0; y < ymax; y++ {
		rowmaxx, rowminx := rowmaxlist[y], rowminlist[y]
		nimg.Set(rowmaxx, y, color.RGBA{0x4B, 0x4B, 0xD1, 0xFF})
		nimg.Set(rowminx, y, color.RGBA{0x1A, 0x1A, 0x65, 0xFF})
	}

	png.Encode(os.Stdout, nimg)

}
