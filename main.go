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

func hough(img image.Image) (image.Point, int) {
	// img is binary image
	rect := img.Bounds()

	var rx, ry, r0, r, x0, x, y0, y, deg, rad float64
	var c uint32

	width, height, area := float64(rect.Max.X), float64(rect.Max.Y), rect.Max.X*rect.Max.Y

	rmax := int(math.Min(float64(width), float64(height)))
	acc := make([]int, area*(rmax-MinEyeR))

	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			c, _, _, _ = img.At(int(x), int(y)).RGBA()

			// r0 is min distance from horizontal or vertical edges
			rx = math.Min(x, width-x)
			ry = math.Min(y, height-y)
			r0 = math.Min(rx, ry)

			// if min distance is larger than MinEyeR
			// and if pixel is white
			if r0 > MinEyeR && c&0xFF == 0xFF {
				// tranform to 3d space
				for r = MinEyeR; r < r0; r++ {
					for deg = 0; deg < 360; deg++ {
						rad = deg * math.Pi / 180.0
						x0 = x + r*math.Cos(rad)
						y0 = y + r*math.Sin(rad)
						acc[int(x0+y0*width)+(int(r)-MinEyeR)*area] += 1
					}
				}
			}
		}
	}

	// find maximus value acc in a for each radious
	// store data to two arrays
	maxlist := make([]int, rmax-MinEyeR)
	cntlist := make([]image.Point, rmax-MinEyeR)
	for r := 0; r < rmax-MinEyeR; r++ {
		max := 0
		cntP := image.Point{0, 0}
		for y = 0; y < height; y++ {
			for x = 0; x < width; x++ {
				if max < acc[int(x+y*width)+int(r)*area] {
					max = acc[int(x+y*width)+int(r)*area]
					cntP = image.Point{int(x), int(y)}
				}
			}
		}
		maxlist[r] = max
		cntlist[r] = cntP
	}

	// find a local maximam of maxlist
	diff := make([]int, len(maxlist)-1)
	for i := 0; i < len(diff)-1; i++ {
		diff[i] = maxlist[i] - maxlist[i+1]
	}

	// the local maximam with good hat-shape in
	// [-1, 0, 1, 2] scope
	closeto0, diffi := 100, 0
	for i := 1; i < len(diff)-2; i++ {
		tmp := diff[i-1] + diff[i] + diff[i+1] + diff[i+2]
		if closeto0 > tmp*tmp {
			closeto0 = tmp * tmp
			diffi = i + 1
			// TODO: +1 is arbitrary need to check i or i+1 somehow
		}
	}
	return cntlist[diffi], diffi + MinEyeR
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
	var c0 uint32
	var c float64
	var mid []int

	// store floating val to int array from 0 to 255
	c_arr := make([]float64, 256*4)
	for i := 0; i < 256; i++ {
		c_arr[4*i] = float64(i) * 0.383
		c_arr[4*i+1] = float64(i) * 0.242
		c_arr[4*i+2] = float64(i) * 0.061
		c_arr[4*i+3] = float64(i) * 0.006
	}

	for y := 0; y < rect.Max.Y; y++ {

		// store a column of pixel val to int array
		mid = make([]int, rect.Max.X)
		for x := 0; x < 3; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
			nimg1.Set(x, y, color.Gray{uint8(c0)})
		}
		for x := 3; x < rect.Max.X-3; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
		}
		for x := rect.Max.X - 3; x < rect.Max.X; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
			nimg1.Set(x, y, color.Gray{uint8(c0)})
		}

		// invoke corresponding floating val to pix array
		for x := 3; x < rect.Max.X-3; x++ {
			c = c_arr[mid[x-3]+3]
			c += c_arr[mid[x-2]+2]
			c += c_arr[mid[x-1]+1]
			c += c_arr[mid[x]]
			c += c_arr[mid[x+1]+1]
			c += c_arr[mid[x+2]+2]
			c += c_arr[mid[x+3]+3]
			nimg1.Set(x, y, color.Gray{uint8(c)})
		}
	}
	_ = img

	nimg2 := image.NewRGBA(rect)
	for x := 0; x < rect.Max.X; x++ {

		// store a column of pixel val to int array
		mid = make([]int, rect.Max.Y)
		for y := 0; y < 3; y++ {
			c0, _, _, _ = nimg1.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
			nimg2.Set(x, y, color.Gray{uint8(c0)})
		}
		for y := 3; y < rect.Max.Y-3; y++ {
			c0, _, _, _ = nimg1.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
		}
		for y := rect.Max.Y - 3; y < rect.Max.Y; y++ {
			c0, _, _, _ = nimg1.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
			nimg2.Set(x, y, color.Gray{uint8(c0)})
		}

		// invoke corresponding floating val to pix array
		for y := 3; y < rect.Max.Y-3; y++ {
			c = c_arr[mid[y-3]+3]
			c += c_arr[mid[y-2]+2]
			c += c_arr[mid[y-1]+1]
			c += c_arr[mid[y]]
			c += c_arr[mid[y+1]+1]
			c += c_arr[mid[y+2]+2]
			c += c_arr[mid[y+3]+3]
			nimg2.Set(x, y, color.Gray{uint8(c)})
		}
	}
	return nimg2
}

func conv1d(c0, c1, c2, c3, c4, c5, c6 uint32) uint8 {
	f0 := float64(c0&0xFF) * 0.006
	f0 += float64(c1&0xFF) * 0.061
	f0 += float64(c2&0xFF) * 0.242
	f0 += float64(c3&0xFF) * 0.383
	f0 += float64(c4&0xFF) * 0.242
	f0 += float64(c5&0xFF) * 0.061
	f0 += float64(c6&0xFF) * 0.006
	return uint8(f0)
}

func conv1d2(a []uint32) uint8 {
	if len(a) != 7 {
		fmt.Fprintf(os.Stderr, "conv1d2: error invalid length %v\n", a)
		os.Exit(1)
	}
	f0 := float64(a[0]&0xFF) * 0.006
	f0 += float64(a[1]&0xFF) * 0.061
	f0 += float64(a[2]&0xFF) * 0.242
	f0 += float64(a[3]&0xFF) * 0.383
	f0 += float64(a[4]&0xFF) * 0.242
	f0 += float64(a[5]&0xFF) * 0.061
	f0 += float64(a[6]&0xFF) * 0.006
	return uint8(f0)
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
	//	nimg, _ = cutoffRGBA(nimg)

	// settle it black(0x00) and white(0xFF)
	//	img = expandRGBA(img)

	// sobel algorithm for edging
	//	nimg = sb(nimg, 2)

	// prewitt algorithm
	//	img = pw(img)

	// binary conversion
	//	nimg = binary(nimg)

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

	//cnt, r := hough(nimg)

	//nimg = drawCircle(img, cnt, r)

	err = png.Encode(os.Stdout, nimg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main read file :%v\n", err)
		os.Exit(1)
	}
	log.Printf("Process took %s", time.Since(start))
}
