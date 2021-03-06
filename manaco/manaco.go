package manaco

import (
	"image"

	"fmt"
	"github.com/eifu/eyemove/avi"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"

	"sort"
)

const (
	MinEyeR = 10
)

type EyeImage struct {
	MyName          int             `json:"MyName"`
	MyRect          image.Rectangle `json:"-"`
	OriginalImage   *image.RGBA     `json:"-"`
	MyRGBA          *image.RGBA     `json:"-"`
	MyCircle        []Circle        `json:"Mycircle"`
	ValidatedCircle Circle          `json:"ValidtedCircle"`
}

type Circle struct {
	X int
	Y int
	R int
}

func Init(ick *avi.ImageChunk) *EyeImage {

	r := image.Rect(0, 0, 172, 114)
	img := image.NewRGBA(r)
	original := image.NewRGBA(r)

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, color.Gray{uint8(ick.Image[x+y*172])})
			original.Set(x, y, color.Gray{uint8(ick.Image[x+y*172])})
		}
	}

	return &EyeImage{
		MyName:        ick.ImageID,
		MyRect:        r,
		MyRGBA:        img,
		OriginalImage: original,
	}
}

func validateNoise(e0, e1, e2 Circle, e3 []Circle) int {

	var e_array []int = []int{e0.R, e1.R, e2.R}
	sort.Ints(e_array)
	var avg int = e_array[1]

	var diff int = 100
	var val int = 0
	for i, e := range e3 {
		if diff > (avg-e.R)*(avg-e.R) {
			val = i
			diff = (avg - e.R) * (avg - e.R)
		}
	}

	return val
}

func CleanNoise(lei []*EyeImage) {
	// check first 5 frames.
	var rightRindex int

	lei[0].ValidatedCircle = lei[0].MyCircle[0]
	lei[1].ValidatedCircle = lei[1].MyCircle[0]
	lei[2].ValidatedCircle = lei[2].MyCircle[0]

	e0 := lei[0].MyCircle[0]
	e1 := lei[1].MyCircle[0]
	e2 := lei[2].MyCircle[0]

	for lei_i, _ := range lei {

		for y := 0; y < lei[lei_i].MyRect.Max.Y; y++ {
			for x := 0; x < lei[lei_i].MyRect.Max.X; x++ {
				c, _, _, _ := (*lei[lei_i].OriginalImage).At(x, y).RGBA()
				lei[lei_i].MyRGBA.Set(x, y, color.RGBA{uint8(c), uint8(c), uint8(c), 0xFF})
			}
		}

		if len(lei[lei_i].MyCircle) == 0 {
			lei[lei_i].ValidatedCircle = e2

		} else if lei_i < 3 {
			lei[lei_i].ValidatedCircle = lei[lei_i].MyCircle[0]

			lei[lei_i].DrawCircle(0)
		} else {
			rightRindex = validateNoise(e0, e1, e2, lei[lei_i].MyCircle)

			lei[lei_i].ValidatedCircle = lei[lei_i].MyCircle[rightRindex]

			e0 = e1
			e1 = e2
			e2 = lei[lei_i].MyCircle[rightRindex]

			lei[lei_i].DrawCircle(rightRindex)
		}

		fname := fmt.Sprintf("image-id%d.png", lei[lei_i].MyName)
		f, err := os.Create(fname)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		err = png.Encode(f, lei[lei_i].MyRGBA)

		if err != nil {
			panic(err)
		}

	}

}

func (eye *EyeImage) Hough(w []image.Point) {
	rect := eye.MyRect
	var c uint32
	var rad, rf float64
	var x0, x1, x2, x3, tmp, rfsinX, rfcosX int
	var y0, y1, y2, y3 int
	// trigo variable array
	// cosX, sinX
	trigo := make([]float64, 90)
	for i := 0; i < 45; i++ {
		rad = float64(i) * math.Pi / 180.0
		trigo[2*i] = math.Cos(rad)
		trigo[2*i+1] = math.Sin(rad)
	}

	width, height := rect.Max.X, rect.Max.Y
	rmax := height / 2
	acc := make([]int, width*height*(rmax-MinEyeR))

	var p image.Point
	// tranform to 3d space
	for r := 0; r < rmax-MinEyeR; r++ {
		rf = float64(r + MinEyeR)
		for i := 0; i < 45; i++ {
			rfcosX = int(rf * trigo[2*i])
			rfsinX = int(rf * trigo[2*i+1])

			for _, p = range w {

				x0 = p.X + rfcosX
				x1 = p.X + rfsinX
				x2 = p.X - rfsinX
				x3 = p.X - rfcosX

				y0 = p.Y + rfsinX
				y1 = p.Y + rfcosX
				y2 = p.Y - rfcosX
				y3 = p.Y - rfsinX

				// first quadrant 0-45
				if (image.Point{x0, y0}.In(rect)) {
					acc[x0+y0*width+width*height*r] += 1
				}
				// first quadrant 45-90
				if (image.Point{x1, y1}.In(rect)) {
					acc[x1+y1*width+width*height*r] += 1
				}
				// second quadrant 90-135
				if (image.Point{x2, y1}.In(rect)) {
					acc[x2+y1*width+width*height*r] += 1
				}
				// second quadrant 135-180
				if (image.Point{x3, y0}.In(rect)) {
					acc[x3+y0*width+width*height*r] += 1
				}
				// third quadrant 180-225
				if (image.Point{x3, y3}.In(rect)) {
					acc[x3+y3*width+width*height*r] += 1
				}
				// third quadrant 225-270
				if (image.Point{x2, y2}.In(rect)) {
					acc[x2+y2*width+width*height*r] += 1
				}
				// fourth quadrant 270-315
				if (image.Point{x1, y2}.In(rect)) {
					acc[x1+y2*width+width*height*r] += 1
				}
				// fourth quadrant 315-360
				if (image.Point{x0, y3}.In(rect)) {
					acc[x0+y3*width+width*height*r] += 1
				}

			}

		}
	}
	// find maximus value acc in a for each radious
	// maxlist store data max accumulated point for each radious
	maxl := make([]int, rmax-MinEyeR)
	// cntlist store center Point that is maximum acc
	cntl := make([]image.Point, rmax-MinEyeR)

	var cnt image.Point
	for r := 0; r < rmax-MinEyeR; r++ {
		tmp = 0
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if tmp < acc[x+y*width+width*height*r] {
					tmp = acc[x+y*width+width*height*r]
					cnt = image.Point{x, y}
				}
			}
		}
		maxl[r] = tmp
		cntl[r] = cnt
	}

	// second derivative of radious candidates
	// i, i+1, i+2, i+3, i+4
	// \+/  \+/  \-/  \-/
	var cc []int
	for i := 0; i < len(maxl)-4; i++ {
		if maxl[i+1]-maxl[i] < 0 {
			continue
		} else if maxl[i+2]-maxl[i+1] < 0 {
			i += 1
		} else if maxl[i+3]-maxl[i+2] > 0 {
			continue
		} else if maxl[i+4]-maxl[i+3] > 0 {
			i += 2
		} else {
			cc = append(cc, i+2)
			// TODO: i or i+1 is arbitrary
		}
	}

	var dencity float64
	var pixel_in_circle float64
	for _, e := range cc {
		dencity = 0.0
		pixel_in_circle = 0.0
		for y := 0; y < rect.Max.Y; y++ {
			for x := 0; x < rect.Max.X; x++ {
				if (cntl[e].X-x)*(cntl[e].X-x)+(cntl[e].Y-y)*(cntl[e].Y-y) < (e+MinEyeR)*(e+MinEyeR) {
					c, _, _, _ = (*eye.OriginalImage).At(x, y).RGBA()
					dencity += float64(c)
					pixel_in_circle += 1.0
				}
			}
		}
		dencity = dencity / pixel_in_circle
		dencity = dencity / 255.0
		if dencity < 130 {

			eye.MyCircle = append(eye.MyCircle, Circle{cntl[e].X, cntl[e].Y, e + MinEyeR})
		}

	}

	_ = c // this is used for reassign the MyRGBA in the below block

}

func (eye *EyeImage) DrawAllCircle() {

	for i, _ := range eye.MyCircle {
		eye.DrawCircle(i)
	}

}

func (eye *EyeImage) DrawCircle(i int) {
	rect := eye.MyRect
	temp := image.NewRGBA(rect)
	var red, g, b uint32
	var deg, rad, x0, y0, x1, y1 float64

	x := eye.MyCircle[i].X
	y := eye.MyCircle[i].Y
	r := eye.MyCircle[i].R
	fmt.Println("myname: ", eye.MyName, " x ", x, " y ", y, " r ", r)
	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			red, g, b, _ = eye.MyRGBA.At(x, y).RGBA()
			temp.Set(x, y, color.RGBA{uint8(red), uint8(g), uint8(b), 0xFF})
		}
	}
	xf, yf, rf := float64(x), float64(y), float64(r)
	for deg = 0; deg < 45; deg++ {
		rad = deg * math.Pi / 180.0
		x1 = rf * math.Cos(rad)
		y1 = rf * math.Sin(rad)
		// first quadrant 0-45
		x0 = xf + x1
		y0 = yf + y1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// first quadrant 46-90
		x0 = xf + y1
		y0 = yf + x1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// second quadrant 91-135
		x0 = xf - y1
		y0 = yf + x1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// second quadrant 136-180
		x0 = xf - x1
		y0 = yf + y1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// third quadrant 181-215
		x0 = xf - x1
		y0 = yf - y1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// third quadrant 216-270
		x0 = xf - y1
		y0 = yf - x1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// fourth quadrant 271-315
		x0 = xf + y1
		y0 = yf - x1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// fourth quadrant 316-359
		x0 = xf + x1
		y0 = yf - y1
		temp.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
	}
	temp.Set(x, y, color.RGBA{0x32, 0x7D, 0x7D, 0xFF})

	eye.MyRGBA = temp

}

// GaussianFilter filters an image data in eye and make it smooth.
// Current implementation is optimized by memoization
// This Gaussian filter uses covonlutional computation
func (eye *EyeImage) GaussianFilter() {

	temp := &EyeImage{
		MyRect:        eye.MyRect,
		OriginalImage: eye.OriginalImage,
		MyRGBA:        image.NewRGBA(eye.MyRect),
	}

	var c0 uint32
	var c float64
	var mid []int

	// store floating val to int array from 0 to 255 for memoization
	c_arr := make([]float64, 256*4)
	for i := 0; i < 256; i++ {
		c_arr[4*i] = float64(i) * 0.383
		c_arr[4*i+1] = float64(i) * 0.242
		c_arr[4*i+2] = float64(i) * 0.061
		c_arr[4*i+3] = float64(i) * 0.006
	}

	// vertical filter
	for y := 0; y < eye.MyRect.Max.Y; y++ {

		// store a column of pixel val to int array
		mid = make([]int, eye.MyRect.Max.X)
		for x := 0; x < 3; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
			temp.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}
		for x := 3; x < eye.MyRect.Max.X-3; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
		}
		for x := eye.MyRect.Max.X - 3; x < eye.MyRect.Max.X; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
			temp.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}

		// invoke corresponding floating val to pix array
		for x := 3; x < eye.MyRect.Max.X-3; x++ {
			c = c_arr[mid[x-3]+3]
			c += c_arr[mid[x-2]+2]
			c += c_arr[mid[x-1]+1]
			c += c_arr[mid[x]]
			c += c_arr[mid[x+1]+1]
			c += c_arr[mid[x+2]+2]
			c += c_arr[mid[x+3]+3]
			temp.MyRGBA.Set(x, y, color.Gray{uint8(c)})
		}
	}

	// horizontal filter
	for x := 0; x < temp.MyRect.Max.X; x++ {

		// store a column of pixel val to int array
		mid = make([]int, temp.MyRect.Max.Y)
		for y := 0; y < 3; y++ {
			c0, _, _, _ = temp.MyRGBA.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
			eye.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}
		for y := 3; y < temp.MyRect.Max.Y-3; y++ {
			c0, _, _, _ = temp.MyRGBA.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
		}
		for y := temp.MyRect.Max.Y - 3; y < temp.MyRect.Max.Y; y++ {
			c0, _, _, _ = temp.MyRGBA.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
			eye.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}

		// invoke corresponding floating val to pix array
		for y := 3; y < temp.MyRect.Max.Y-3; y++ {
			c = c_arr[mid[y-3]+3]
			c += c_arr[mid[y-2]+2]
			c += c_arr[mid[y-1]+1]
			c += c_arr[mid[y]]
			c += c_arr[mid[y+1]+1]
			c += c_arr[mid[y+2]+2]
			c += c_arr[mid[y+3]+3]
			eye.MyRGBA.Set(x, y, color.Gray{uint8(c)})
		}
	}
}

// Binary changes myRGB in eye to Binary image.
// and returns an array of 'white' pixel in the binary image
// currently the way it takes binary image is
// if a pixel is blighter than average color, turn it to white
// if a pixel is darker than average color, turn it to black
func (eye *EyeImage) Binary() []image.Point {
	rect := eye.MyRect
	width, height := rect.Max.X, rect.Max.Y
	clr := make([]uint32, width*height)
	temp := image.NewRGBA(rect)
	var acc, ave, c0 uint32
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			clr[x+y*width] = c0 & 0xFF
			acc += c0 & 0xFF
		}
	}
	ave = acc / uint32(width*height)

	var w []image.Point

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if clr[x+y*width] > ave {
				w = append(w, image.Point{x, y})
				temp.Set(x, y, color.Gray{0xFF})
			} else {
				temp.Set(x, y, color.Gray{0x0})
			}
		}
	}
	eye.MyRGBA = temp
	return w
}

// CutoffRGBA changes myRGBA and turn 'blighter-than-average' pixel
// to average color.
// The purpose of this function is to remove the flash light
// in the image.
func (eye *EyeImage) CutoffRGBA() {
	rect := eye.MyRect
	temp := image.NewRGBA(rect)

	var acc, ave, c0 uint32

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			acc = acc + c0&0xFF
		}
	}
	ave = acc / uint32(rect.Max.X*rect.Max.Y)

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			if c0&0xFF > ave {
				temp.Set(x, y, color.Gray{uint8(ave)})
			} else {
				temp.Set(x, y, color.Gray{uint8(c0)})
			}
		}
	}
	eye.MyRGBA = temp
}

// This is normalization of pixel value.
func NormalizeRGBA(img image.Image) *image.RGBA {
	rect := img.Bounds()

	var min, max uint8 = 0xFF, 0
	var c0 uint32
	var c1 uint8

	for y := 1; y < rect.Max.Y-1; y++ {
		for x := 1; x < rect.Max.X-1; x++ {
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

// Sobel finds the edge in an image.
// Sobel Algorithm: Given a image, you extract 3 by 3 pixcels matrix.
//
// a b c
// d e f
// g h i
//
// After the Sobel algorithm, the value of e must be as following.
//
// e = sqrt{ (a + w*b + c - g - w*h - i)^2 + (a + w*d + g - c - w*f - i)^2 }
//
// where w is a positive integer argument
func (eye *EyeImage) Sobel(w float64) {
	rect := eye.MyRect
	temp := image.NewRGBA(rect)
	var sum, gx, gy float64
	for j := 0; j < rect.Max.Y; j++ {
		for i := 0; i < rect.Max.X; i++ {
			gy, gx = eye.sb_helper(i, j, w)
			sum = math.Sqrt(gx*gx + gy*gy)
			if sum > 255 {
				sum = 255
			}
			temp.Set(i, j, color.Gray{uint8(sum)})
		}
	}
	eye.MyRGBA = temp
}

func (eye *EyeImage) sb_helper(x, y int, w float64) (float64, float64) {
	rect := eye.MyRect
	var accY, accX float64
	var c uint32
	for j := y - 1; j < y+2; j++ {
		for i := x - 1; i < x+2; i++ {
			if (image.Point{i, j}.In(rect)) {
				c, _, _, _ = eye.MyRGBA.At(i, j).RGBA()
				switch {
				case i == x-1 && j != y:
					accY -= float64(c & 0xFF)
				case i == x-1 && j == y:
					accY -= w * float64(c&0xFF)
				case i == x+1 && j != y:
					accY += float64(c & 0xFF)
				case i == x+1 && j == y:
					accY += w * float64(c&0xFF)
				}
				switch {
				case i != x && j == y-1:
					accX -= float64(c & 0xFF)
				case i == x && j == y-1:
					accX -= w * float64(c&0xFF)
				case i != x && j == y+1:
					accX += float64(c & 0xFF)
				case i == x && j == y+1:
					accX += w * float64(c&0xFF)
				}
			}
		}
	}
	return accY, accX
}

func (eye *EyeImage) Prewitt() {
	eye.Sobel(1)
}
