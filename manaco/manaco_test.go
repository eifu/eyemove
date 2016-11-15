package manaco

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestID1(t *testing.T) {

	f, err := os.Open("image-id1.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	imgData, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	r := image.Rect(0, 0, 172, 114)
	img := image.NewRGBA(r)
	original := image.NewRGBA(r)
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, imgData.At(x, y))
			original.Set(x, y, imgData.At(x, y))
		}
	}
	eye_image := &EyeImage{
		MyName:        997,
		MyRect:        r,
		MyRGBA:        img,
		OriginalImage: original,
	}

	eye_image.GaussianFilter()

	eye_image.CutoffRGBA()

	eye_image.Sobel(1)

	w := eye_image.Binary()

	eye_image.Hough(w)

	for y := 0; y < eye_image.MyRect.Max.Y; y++ {
		for x := 0; x < eye_image.MyRect.Max.X; x++ {
			eye_image.MyRGBA.Set(x, y, eye_image.OriginalImage.At(x, y))
		}
	}

	eye_image.DrawCircle(0)
	//	eye_image.DrawCircle(1)

	fname := fmt.Sprintf("test__%d.png", eye_image.MyName)

	out_file, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	defer out_file.Close()
	err = png.Encode(out_file, eye_image.MyRGBA)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", eye_image)
}
