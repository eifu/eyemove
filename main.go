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

	xmax, ymax := rect.Max.X, rect.Max.Y
	var acc, ave uint32
	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			acc = acc + r&0xFF
		}
	}
	ave = acc / uint32(xmax*ymax)

	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r&0xFF > ave {
				nimg.Set(x, y, color.RGBA{uint8(ave), uint8(ave), uint8(ave), 0xFF})
			} else {
				nimg.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), 0xFF})
			}
		}
	}
	return nimg, ave
}

func expandRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()

	var min, max float64 = 255.0, 0.0
	var l float64
	var minr, ming, minb uint32
	var r, g, b, a uint32
	var r1, g1, b1 uint8

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			r, g, b, _ = img.At(x, y).RGBA()
			l = luminosity(uint8(r), uint8(g), uint8(b))

			if min > l {
				min = l
				minr, ming, minb = r, g, b
			}
			if max < l {
				max = l
			}
		}
	}

	var ratio float64 = 0xFF / (max - min)

	nrgba := image.NewRGBA(rect)

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			r, g, b, a = img.At(x, y).RGBA()
			r1 = uint8(float64(r&0xFF-minr&0xFF) * ratio)
			g1 = uint8(float64(g&0xFF-ming&0xFF) * ratio)
			b1 = uint8(float64(b&0xFF-minb&0xFF) * ratio)

			nrgba.Set(x, y, color.RGBA{r1, g1, b1, uint8(a)})
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

	var rowmax, rowmaxx, rowmin, rowminx int
	var rowr0, rowg0, rowb0 uint32
	for y := 0; y < ymax; y++ {
		for x := 1; x < xmax; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			l := luminosity(uint8(r), uint8(g), uint8(b))

			rowr0, rowg0, rowb0, _ = img.At(x-1, y).RGBA()
			rowl0 := luminosity(uint8(rowr0), uint8(rowg0), uint8(rowb0))

			row_array[x-1] = int(l - rowl0)

			if rowr0&0xFF != ave && r&0xFF != ave {
				if rowmax < row_array[x-1] {
					rowmax, rowmaxx = row_array[x-1], x
				}
				if rowmin > row_array[x-1] {
					rowmin, rowminx = row_array[x-1], x
				}
			}

		}
		maxlist[y], minlist[y] = rowmaxx, rowminx

		rowmax = 0x0
		rowmin = 0xFF
	}

	return maxlist, minlist
}

func colIterate(img image.Image, ave uint32) ([]int, []int) {
	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y

	var maxlist = make([]int, xmax)
	var minlist = make([]int, xmax)
	var col_array = make([]int, ymax-1)

	var colmax, colmaxy, colmin, colminy int
	var colr0, colg0, colb0 uint32
	for x := 0; x < xmax; x++ {
		for y := 1; y < ymax; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			l := luminosity(uint8(r), uint8(g), uint8(b))

			colr0, colg0, colb0, _ = img.At(x, y-1).RGBA()
			coll0 := luminosity(uint8(colr0), uint8(colg0), uint8(colb0))
			col_array[y-1] = int(l - coll0)

			if colr0&0xFF != ave && r&0xFF != ave {
				if colmax < col_array[y-1] {
					colmax, colmaxy = col_array[y-1], y
				}
				if colmin > col_array[y-1] {
					colmin, colminy = col_array[y-1], y
				}
			}
		}
		maxlist[x], minlist[x] = colmaxy, colminy

		colmax = 0x0
		colmin = 0xFF
	}
	return maxlist, minlist
}

func eight(img image.Image) *image.RGBA {

	bounds := img.Bounds()
	//	xmax, ymax := bounds.Max.X/8, bounds.Max.Y/8

	nimg := image.NewRGBA(bounds)

	return nimg
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

	nimg := expandRGBA(img) // settle it black(0x00) and white(0xFF)

	nimg, _ = cutoffRGBA(nimg)

	//	nimg = expandRGBA(nimg)
	nimg, ave := cutoffRGBA(nimg)

	// iterate lines
	rowmaxlist, rowminlist := rowIterate(nimg, ave)
	colmaxlist, colminlist := colIterate(nimg, ave)

	for x := 0; x < xmax; x++ {
		colmaxy, colminy := colmaxlist[x], colminlist[x]

		nimg.Set(x, colmaxy, color.RGBA{0x3F, 0xBF, 0x3F, 0xFF})
		nimg.Set(x, colminy, color.RGBA{0x19, 0x4C, 0x19, 0xFF})

	}

	for y := 0; y < ymax; y++ {
		rowmaxx, rowminx := rowmaxlist[y], rowminlist[y]

		nimg.Set(rowmaxx, y, color.RGBA{0x4B, 0x4B, 0xD1, 0xFF})
		nimg.Set(rowminx, y, color.RGBA{0x1A, 0x1A, 0x65, 0xFF})
	}

	png.Encode(os.Stdout, nimg)

}
