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

	// decode image data
	imgData, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	// load image data from decoded data
	r := image.Rect(0, 0, 172, 114)
	img := image.NewRGBA(r)
	original := image.NewRGBA(r)
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, imgData.At(x, y))
			original.Set(x, y, imgData.At(x, y))
		}
	}

	// create an instance of EyeImage struct
	eye_image := &EyeImage{
		MyName:        997,
		MyRect:        r,
		MyRGBA:        img,
		OriginalImage: original,
	}

	// process gaussian filter to make it smooth
	eye_image.GaussianFilter()

	// cut off half of whiter pixel to avarage color
	eye_image.CutoffRGBA()

	// Sobel algorithm to find edge
	eye_image.Sobel(1)

	// make binary image based on color.
	// create an array of image.Point of 'white' pixel.
	// from binary image.
	w := eye_image.Binary()

	// process Hough transform algrithm
	// eyeImage will get myCircle attribute
	// if it finds circle in the image.
	eye_image.Hough(w)

	// re-load original image to create circles-embedded image
	for y := 0; y < eye_image.MyRect.Max.Y; y++ {
		for x := 0; x < eye_image.MyRect.Max.X; x++ {
			eye_image.MyRGBA.Set(x, y, eye_image.OriginalImage.At(x, y))
		}
	}

	eye_image.DrawCircle(0)
	//	eye_image.DrawCircle(1)

	// make file name based on ID number of eyeImage
	fname := fmt.Sprintf("test__%d.png", eye_image.MyName)

	// make png file of circles-embedded image
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
