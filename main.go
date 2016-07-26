package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"time"
)

const (
	MinEyeR = 10
)

func hough(img, pimg image.Image) *image.RGBA {
	log.Print("start hough transforming")
	// img is binary image
	rect := img.Bounds()

	var deg, rad float64
	var c uint32
	var x0, y0, x1, y1, tmp int
	width, height := rect.Max.X, rect.Max.Y

	//	rmax := int(math.Min(float64(width), float64(height)))
	rmax := height / 2
	acc := make([]int, width*height*(rmax-MinEyeR))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c, _, _, _ = img.At(int(x), int(y)).RGBA()

			// if pixel is white
			if c&0xFF == 0xFF {
				// tranform to 3d space
				for r := 0; r < rmax-MinEyeR; r++ {
					for deg = 0; deg < 360; deg++ {
						rad = deg * math.Pi / 180.0
						x0 = x + int(float64(r+MinEyeR)*math.Cos(rad))
						y0 = y + int(float64(r+MinEyeR)*math.Sin(rad))
						if 0 <= x0 && x0 < width && 0 <= y0 && y0 < height {
							acc[x0+y0*width+width*height*r] += 1
						}
					}
				}
			}
		}
	}

	// find maximus value acc in a for each radious
	// maxlist store data max accumulated point for each radious
	maxl := make([]int, rmax-MinEyeR)
	// cntlist store center Point that is maximum acc
	cntl := make([]image.Point, rmax-MinEyeR)

	for r := 0; r < rmax-MinEyeR; r++ {
		tmp = 0
		cnt := image.Point{0, 0}
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if tmp < acc[x+y*width+r*width*height] {
					tmp = acc[x+y*width+r*width*height]
					cnt = image.Point{x, y}
				}
			}
		}
		maxl[r] = tmp
		cntl[r] = cnt
	}

	// difference between each value in maxlist
	difm := make([]int, rmax-MinEyeR-1)

	// rectilinear distance between each center pointn in cntlist
	difc := make([]int, rmax-MinEyeR-1)

	for i := 0; i < len(difm)-1; i++ {
		difm[i] = maxl[i] - maxl[i+1]
		difc[i] = cntl[i].X - cntl[i+1].X + cntl[i].Y - cntl[i+1].Y
	}

	// candidates of canditates of radious
	var cc []int
	for i := 1; i < len(difc)-2; i++ {
		if difm[i-1] < 0 && difm[i] < 0 && 0 < difm[i+1] && 0 < difm[i+2] {
			cc = append(cc, i)
			// TODO: i or i+1 is arbitrary
		}
	}

	// accm0, accm1: best 2 accumulation maximums
	// cd0, cd1: best 2 candidates of radious
	var cd0, cd1, accm0, accm1 int

	for _, e := range cc {
		if accm0 < maxl[e] {
			tmp = accm0
			accm0 = maxl[e]
			accm1 = tmp
			tmp = cd0
			cd0 = e
			cd1 = tmp
		} else if accm1 < maxl[e] {
			accm1 = maxl[e]
			cd1 = e
		}
	}

	// determine which one has more black pixels than the other
	var acc0, acc1 float64
	x0, y0, x1, y1 = cntl[cd0].X, cntl[cd0].Y, cntl[cd1].X, cntl[cd1].Y
	cd0 += MinEyeR
	cd1 += MinEyeR
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x-x0)*(x-x0)+(y-y0)*(y-y0) < cd0*cd0 {
				c, _, _, _ = pimg.At(x, y).RGBA()
				acc0 += float64(c & 0xFF)
			}
			if (x-x1)*(x-x1)+(y-y1)*(y-y1) < cd1*cd1 {
				c, _, _, _ = pimg.At(x, y).RGBA()
				acc1 += float64(c & 0xFF)
			}
		}
	}
	dens0, dens1 := acc0/(float64(cd0*cd0)*math.Pi), acc1/(float64(cd1*cd1)*math.Pi)
	if dens0 < dens1 {
		return drawCircle(pimg, cntl[cd0-MinEyeR], cd0)
	}

	return drawCircle(pimg, cntl[cd1-MinEyeR], cd1)
}
func drawCircle(img image.Image, cnt image.Point, r int) *image.RGBA {
	rect := img.Bounds()
	nimg := image.NewRGBA(rect)
	var c uint32
	var deg, rad, x0, y0 float64
	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c, _, _, _ = img.At(x, y).RGBA()
			nimg.Set(x, y, color.RGBA{uint8(c), uint8(c), uint8(c), 0xFF})
		}
	}

	xf, yf, rf := float64(cnt.X), float64(cnt.Y), float64(r)

	for deg = 0; deg < 360; deg++ {
		rad = deg * math.Pi / 180.0
		x0 = xf + rf*math.Cos(rad)
		y0 = yf + rf*math.Sin(rad)
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
	}
	nimg.Set(cnt.X, cnt.Y, color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
	return nimg
}

func g_smoothing(img image.Image) *image.RGBA {
	log.Print("start gaussian smoothing")
	rect := img.Bounds()
	nimg1 := image.NewRGBA(rect)
	// convolution algorithm
	var c0, c1, c2, c3, c4, c5, c6 uint32

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < 3; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			nimg1.Set(x, y, color.Gray{uint8(c0)})
		}

		for x := 3; x < rect.Max.X-3; x++ {
			c0, _, _, _ = img.At(x-3, y).RGBA()
			c1, _, _, _ = img.At(x-2, y).RGBA()
			c2, _, _, _ = img.At(x-1, y).RGBA()
			c3, _, _, _ = img.At(x, y).RGBA()
			c4, _, _, _ = img.At(x+1, y).RGBA()
			c5, _, _, _ = img.At(x+2, y).RGBA()
			c6, _, _, _ = img.At(x+3, y).RGBA()
			c := conv1d(c0, c1, c2, c3, c4, c5, c6)
			nimg1.Set(x, y, color.Gray{c})
		}

		for x := rect.Max.X - 3; x < rect.Max.X; x++ {
			c0, _, _, _ := img.At(x, y).RGBA()
			nimg1.Set(x, y, color.Gray{uint8(c0)})
		}
	}
	_ = img

	nimg2 := image.NewRGBA(rect)
	for x := 0; x < rect.Max.X; x++ {
		for y := 0; y < 3; y++ {
			c0, _, _, _ := nimg1.At(x, y).RGBA()
			nimg2.Set(x, y, color.Gray{uint8(c0)})
		}
		for y := 3; y < rect.Max.Y-3; y++ {
			c0, _, _, _ = nimg1.At(x, y-3).RGBA()
			c1, _, _, _ = nimg1.At(x, y-2).RGBA()
			c2, _, _, _ = nimg1.At(x, y-1).RGBA()
			c3, _, _, _ = nimg1.At(x, y).RGBA()
			c4, _, _, _ = nimg1.At(x, y+1).RGBA()
			c5, _, _, _ = nimg1.At(x, y+2).RGBA()
			c6, _, _, _ = nimg1.At(x, y+3).RGBA()
			c := conv1d(c0, c1, c2, c3, c4, c5, c6)
			nimg2.Set(x, y, color.Gray{c})
		}
		for y := rect.Max.Y - 3; y < rect.Max.Y; y++ {
			c0, _, _, _ := nimg1.At(x, y).RGBA()
			nimg2.Set(x, y, color.Gray{uint8(c0)})
		}
	}
	return nimg2
}

func conv1d(c0, c1, c2, c3, c4, c5, c6 uint32) uint8 {
	f0 := float64(c0&0xFF) * 0.006
	f1 := float64(c1&0xFF) * 0.061
	f2 := float64(c2&0xFF) * 0.242
	f3 := float64(c3&0xFF) * 0.383
	f4 := float64(c4&0xFF) * 0.242
	f5 := float64(c5&0xFF) * 0.061
	f6 := float64(c6&0xFF) * 0.006
	return uint8(f0 + f1 + f2 + f3 + f4 + f5 + f6)
}

func binary(img image.Image) *image.RGBA {
	log.Print("start settling black or white")
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
				nimg.Set(x, y, color.Gray{255})
			} else {
				nimg.Set(x, y, color.Gray{0})
			}
		}
	}
	return nimg
}

func cutoffRGBA(img image.Image) (*image.RGBA, uint32) {
	log.Print("start cut off below average")
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
	log.Print("start expanding")
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

func sb(img image.Image, w float64) *image.RGBA {
	log.Print("start sobel algorithm...")
	rect := img.Bounds()
	nimg := image.NewRGBA(rect)
	var sum, gx, gy float64
	for j := 0; j < rect.Max.Y; j++ {
		for i := 0; i < rect.Max.X; i++ {
			gy, gx = sb_helper(img, i, j, w)
			sum = math.Sqrt(gx*gx + gy*gy)
			if sum > 255 {
				sum = 255
			}
			nimg.Set(i, j, color.Gray{uint8(sum)})
		}
	}
	return nimg
}

func sb_helper(img image.Image, x, y int, w float64) (float64, float64) {
	var accY, accX float64
	for j := y - 1; j < y+2; j++ {
		for i := x - 1; i < x+2; i++ {
			if (image.Point{i, j}.In(img.Bounds())) {
				c, _, _, _ := img.At(i, j).RGBA()
				switch {
				case i == x-1 && j != y:
					accY = accY - float64(c&0xFF)
				case i == x-1 && j == y:
					accY = accY - w*float64(c&0xFF)
				case i == x+1 && j != y:
					accY = accY + float64(c&0xFF)
				case i == x+1 && j == y:
					accY = accY + w*float64(c&0xFF)
				}
			}
			if (image.Point{i, j}.In(img.Bounds())) {
				c, _, _, _ := img.At(i, j).RGBA()
				switch {
				case i != x && j == y-1:
					accX = accX - float64(c&0xFF)
				case i == x && j == y-1:
					accX = accX - w*float64(c&0xFF)
				case i != x && j == y+1:
					accX = accX + float64(c&0xFF)
				case i == x && j == y+1:
					accX = accX + w*float64(c&0xFF)
				}
			}

		}
	}
	return accY, accX
}

func pw(img image.Image) image.Image {
	return sb(img, 1)
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

func row_iterate(img image.Image, ave uint32) ([]int, []int) {
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

func col_iterate(img image.Image, ave uint32) ([]int, []int) {
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

func main() {
	start := time.Now()
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

	// gaussian smoothing function
	nimg := g_smoothing(img)

	// cut off pixels below the average color
	nimg, _ = cutoffRGBA(nimg)

	// settle it black(0x00) and white(0xFF)
	//	img = expandRGBA(img)

	// sobel algorithm for edging
	nimg = sb(nimg, 5)

	// prewitt algorithm
	//	img = pw(img)

	// binary conversion
	nimg = binary(nimg)

	/*
		// iterate lines
		rowmaxlist, rowminlist := row_iterate(nimg, ave)
		colmaxlist, colminlist := col_iterate(nimg, ave)

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
	*/

	nimg = hough(nimg, img)

	err = png.Encode(os.Stdout, nimg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main read file :%v\n", err)
		os.Exit(1)
	}
	log.Printf("Process took %s", time.Since(start))

}
