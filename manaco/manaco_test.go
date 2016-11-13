package manaco

import (
	"image"
	_ "image/png"
	"os"
	"testing"
)

func main(t *testing.T) {

	f, err := os.Open("id1_withCircle.png")
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
		MyName:        999,
		MyRect:        r,
		MyRGBA:        img,
		OriginalImage: original,
	}

	//	eye_image := manaco.Init(ick)

	eye_image.GaussianFilter()

	eye_image.CutoffRGBA()

	eye_image.Sobel(2)

	w := eye_image.Binary()

	eye_image.Hough(w)

}
