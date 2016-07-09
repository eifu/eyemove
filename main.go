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
	//	fmt.Println(uint8(r0), uint8(g0), uint8(b0), l0)
	l1 := l2 - l0
	if l1 == 0 {
		return 1.0
	} else {
		return l1
	}
}

func gradient_magnitude(dx, dy float64) float64 {
	//	fmt.Println(dx, dy, int(math.Sqrt(dx*dx+dy*dy)))

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
		fmt.Fprintf(os.Stderr, "main %v\n", err)
		os.Exit(1)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main read file :%v\n", err)
		os.Exit(1)
	}
	img = expandRGBA(img)

	//	png.Encode(os.Stdout, img)

	xmax, ymax := img.Bounds().Max.X, img.Bounds().Max.Y
	//chunk_list := make([]Chunk, (xmax-1)*(ymax-1))
	var chunk_list []Chunk
	var theta, magnitude float64
	for y := 0; y < ymax; y++ {
		for x := 0; x < xmax; x++ {
			dx, dy := dx(img, x, y), dy(img, x, y)
			theta = gradient_theta(dx, dy)
			magnitude = gradient_magnitude(dx, dy)
			chunk_list = append(chunk_list, Chunk{x: x, y: y, theta: theta, magnitude: magnitude})
		}
	}

	rect := img.Bounds()
	nimg := image.NewRGBA(rect)

	for _, chunk := range chunk_list {
		//	fmt.Println(chunk.theta, chunk.magnitude)
		if math.Abs(float64(chunk.theta)) <= 120 && chunk.magnitude >= 5 {
			nimg.Set(chunk.x, chunk.y, color.RGBA{0x0, 0x0, 0x0, 0xFF})
		} else {
			nimg.Set(chunk.x, chunk.y, color.RGBA{0xFF, 0xFF, 0xFF, 0xFF})
		}
	}
	png.Encode(os.Stdout, nimg)

}
