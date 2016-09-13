package manaco

import (
	"image"
	"image/draw"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"
	"time"
)

const (
	MinEyeR = 10
)

type eyeImage struct {
	MyRect image.Rectangle
	MyImage *draw.Image
}

func InitEyeImage(img *image.Image) *eyeImage{
	return &eyeImage{
		MyRect: (*img).Bounds(),
		MyImage: img,
	}
}

func process1(out chan<- [3]int, rmax int, trigo []float64) {
	for r := 0; r < rmax-MinEyeR; r++ {
		rf := float64(r + MinEyeR)
		for deg := 0; deg < 45; deg++ {
			rfcosX := int(rf * trigo[2*deg])
			rfsinX := int(rf * trigo[2*deg+1])
			out <- [3]int{rfcosX, rfsinX, r}
		}
	}
	close(out)
}

func process2(out chan<- [][3]int, in <-chan [3]int, w []image.Point, rect image.Rectangle) {
	for v := range in {
		rfcosX := v[0]
		rfsinX := v[1]
		radius := v[2]
		for _, p := range w {

			p_news := make([][3]int, 0, 8)

			if (image.Point{p.X + rfcosX, p.Y + rfsinX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X + rfcosX, p.Y + rfsinX, radius})
			}
			if (image.Point{p.X + rfsinX, p.Y + rfcosX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X + rfsinX, p.Y + rfcosX, radius})
			}
			if (image.Point{p.X - rfsinX, p.Y + rfcosX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X - rfsinX, p.Y + rfcosX, radius})
			}
			if (image.Point{p.X - rfcosX, p.Y + rfsinX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X - rfcosX, p.Y + rfsinX, radius})
			}
			if (image.Point{p.X - rfcosX, p.Y - rfsinX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X - rfcosX, p.Y - rfsinX, radius})
			}
			if (image.Point{p.X - rfsinX, p.Y - rfcosX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X - rfsinX, p.Y - rfcosX, radius})
			}
			if (image.Point{p.X + rfsinX, p.Y - rfcosX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X + rfsinX, p.Y - rfcosX, radius})
			}
			if (image.Point{p.X + rfcosX, p.Y - rfsinX}.In(rect)) {
				p_news = append(p_news, [3]int{p.X + rfcosX, p.Y - rfsinX, radius})
			}
			out <- p_news
		}
	}
	close(out)
}

func process3(in <-chan [][3]int, acc []int, width, height int) {
	for v := range in {
		for _, p_new := range v {
			acc[p_new[0]+p_new[1]*width+width*height*p_new[2]] += 1
		}
	}
}

func Hough2(w []image.Point, pimg image.Image) *image.RGBA {
	rect := pimg.Bounds()

	var rad float64
	var c uint32
	var x0, x1, tmp int
	var y0, y1 int
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
	n := time.Now()

	// tranform to 3d space

	deg_chan := make(chan [3]int)
	p_chan := make(chan [][3]int)
	go process1(deg_chan, rmax, trigo)
	go process2(p_chan, deg_chan, w, rect)
	process3(p_chan, acc, width, height)

	log.Printf("  transform takes %.2f \n", time.Since(n).Seconds())
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

	// TODO: best 2 is arbitrary
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
	var acc0, acc1 uint32
	x0, y0, x1, y1 = cntl[cd0].X, cntl[cd0].Y, cntl[cd1].X, cntl[cd1].Y
	r0, r1 := (cd0+MinEyeR)*(cd0+MinEyeR), (cd1+MinEyeR)*(cd1+MinEyeR)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x-x0)*(x-x0)+(y-y0)*(y-y0) < r0 {
				c, _, _, _ = pimg.At(x, y).RGBA()
				acc0 += c & 0xFF
			}
			if (x-x1)*(x-x1)+(y-y1)*(y-y1) < r1 {
				c, _, _, _ = pimg.At(x, y).RGBA()
				acc1 += c & 0xFF
			}
		}
	}
	dens0, dens1 := float64(acc0)/(float64(r0)*math.Pi), float64(acc1)/(float64(r1)*math.Pi)
	if dens0 < dens1 {
		return DrawCircle(pimg, cntl[cd0], cd0+MinEyeR)
	}
	return DrawCircle(pimg, cntl[cd1], cd1+MinEyeR)
}

func Hough(w []image.Point, pimg image.Image) *image.RGBA {
	rect := pimg.Bounds()

	var rad, rf float64
	var c uint32
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

	// TODO: best 2 is arbitrary
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
	var acc0, acc1 uint32
	x0, y0, x1, y1 = cntl[cd0].X, cntl[cd0].Y, cntl[cd1].X, cntl[cd1].Y
	r0, r1 := (cd0+MinEyeR)*(cd0+MinEyeR), (cd1+MinEyeR)*(cd1+MinEyeR)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x-x0)*(x-x0)+(y-y0)*(y-y0) < r0 {
				c, _, _, _ = pimg.At(x, y).RGBA()
				acc0 += c & 0xFF
			}
			if (x-x1)*(x-x1)+(y-y1)*(y-y1) < r1 {
				c, _, _, _ = pimg.At(x, y).RGBA()
				acc1 += c & 0xFF
			}
		}
	}
	dens0, dens1 := float64(acc0)/(float64(r0)*math.Pi), float64(acc1)/(float64(r1)*math.Pi)
	if dens0 < dens1 {
		return DrawCircle(pimg, cntl[cd0], cd0+MinEyeR)
	}
	return DrawCircle(pimg, cntl[cd1], cd1+MinEyeR)
}

func DrawCircle(img image.Image, cnt image.Point, r int) *image.RGBA {
	rect := img.Bounds()
	nimg := image.NewRGBA(rect)
	var c uint32
	var deg, rad, x0, y0, x1, y1 float64
	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c, _, _, _ = img.At(x, y).RGBA()
			nimg.Set(x, y, color.RGBA{uint8(c), uint8(c), uint8(c), 0xFF})
		}
	}
	xf, yf, rf := float64(cnt.X), float64(cnt.Y), float64(r)
	for deg = 0; deg < 45; deg++ {
		rad = deg * math.Pi / 180.0
		x1 = rf * math.Cos(rad)
		y1 = rf * math.Sin(rad)
		// first quadrant 0-45
		x0 = xf + x1
		y0 = yf + y1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// first quadrant 46-90
		x0 = xf + y1
		y0 = yf + x1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// second quadrant 91-135
		x0 = xf - y1
		y0 = yf + x1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// second quadrant 136-180
		x0 = xf - x1
		y0 = yf + y1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// third quadrant 181-215
		x0 = xf - x1
		y0 = yf - y1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// third quadrant 216-270
		x0 = xf - y1
		y0 = yf - x1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// fourth quadrant 271-315
		x0 = xf + y1
		y0 = yf - x1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
		// fourth quadrant 316-359
		x0 = xf + x1
		y0 = yf - y1
		nimg.Set(int(x0), int(y0), color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
	}
	nimg.Set(cnt.X, cnt.Y, color.RGBA{0x32, 0x7D, 0x7D, 0xFF})
	return nimg
}

func (eye *eyeImage)GaussianFilter() *eyeImage {

	m := image.NewRGBA(eye.MyRect)
	draw.Draw(m, m.Bounds(), *eye.MyImage, image.ZP, draw.Src)

	nimg1 := &eyeImage{
		MyRect:eye.MyRect,
		MyImage: m,
	}
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

	for y := 0; y < eye.MyRect.Max.Y; y++ {

		// store a column of pixel val to int array
		mid = make([]int, eye.MyRect.Max.X)
		for x := 0; x < 3; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
			nimg1.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}
		for x := 3; x < eye.MyRect.Max.X-3; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
		}
		for x := eye.MyRect.Max.X - 3; x < eye.MyRect.Max.X; x++ {
			c0, _, _, _ = eye.MyRGBA.At(x, y).RGBA()
			mid[x] = 4 * int(c0&0xFF)
			nimg1.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
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
			nimg1.MyRGBA.Set(x, y, color.Gray{uint8(c)})
		}
	}

	nimg2 := &eyeImage{
		MyRect: eye.MyRect,
		MyRGBA:		image.NewRGBA(eye.MyRect),
	}
	for x := 0; x < eye.MyRect.Max.X; x++ {

		// store a column of pixel val to int array
		mid = make([]int, eye.MyRect.Max.Y)
		for y := 0; y < 3; y++ {
			c0, _, _, _ = nimg1.MyRGBA.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
			nimg2.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}
		for y := 3; y < eye.MyRect.Max.Y-3; y++ {
			c0, _, _, _ = nimg1.MyRGBA.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
		}
		for y := eye.MyRect.Max.Y - 3; y < eye.MyRect.Max.Y; y++ {
			c0, _, _, _ = nimg1.MyRGBA.At(x, y).RGBA()
			mid[y] = 4 * int(c0&0xFF)
			nimg2.MyRGBA.Set(x, y, color.Gray{uint8(c0)})
		}

		// invoke corresponding floating val to pix array
		for y := 3; y < eye.MyRect.Max.Y-3; y++ {
			c = c_arr[mid[y-3]+3]
			c += c_arr[mid[y-2]+2]
			c += c_arr[mid[y-1]+1]
			c += c_arr[mid[y]]
			c += c_arr[mid[y+1]+1]
			c += c_arr[mid[y+2]+2]
			c += c_arr[mid[y+3]+3]
			nimg2.MyRGBA.Set(x, y, color.Gray{uint8(c)})
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
		log.Fatalf("conv1d2: error invalid length %v\n", a)
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

func Binary(img image.Image) (*image.RGBA, []image.Point) {
	rect := img.Bounds()
	width, height := rect.Max.X, rect.Max.Y
	cl := make([]uint32, width*height)
	nimg := image.NewRGBA(rect)
	var acc, ave, c0 uint32
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c0, _, _, _ = img.At(x, y).RGBA()
			cl[x+y*width] = c0 & 0xFF
			acc = acc + c0&0xFF
		}
	}
	ave = acc / uint32(width*height)

	var w []image.Point

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if cl[x+y*width] > ave {
				w = append(w, image.Point{x, y})
				nimg.Set(x, y, color.Gray{0xFF})
			} else {
				nimg.Set(x, y, color.Gray{0x0})
			}
		}
	}
	return nimg, w
}

func CutoffRGBA(img image.Image) (*image.RGBA, uint32) {
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

func ExpandRGBA(img image.Image) *image.RGBA {
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

func Sobel(img image.Image, w float64) *image.RGBA {
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
	var c uint32
	for j := y - 1; j < y+2; j++ {
		for i := x - 1; i < x+2; i++ {
			if (image.Point{i, j}.In(img.Bounds())) {
				c, _, _, _ = img.At(i, j).RGBA()
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

func Prewitt(img image.Image) image.Image {
	return Sobel(img, 1)
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
