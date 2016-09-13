package main

import (
	"./manaco"
	"flag"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root, submain)
	if err != nil {
		log.Fatal(err)
	}
}

func submain(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println(err)
		return err
	}
	if info.IsDir() {
		return nil
	}

	log.Println(path + "is loading...")
	infile, err := os.Open(path)
	defer infile.Close()
	if err != nil {
		log.Printf("main open file :%v\n", err)
		return err
	}
	img, _, err := image.Decode(infile)
	if err != nil {
		log.Printf("main read file :%v\n", err)
		return err
	}

	eye_image := manaco.InitEyeImage(&img)

	eye_image.GaussianFilter()

	eye_image.CutoffRGBA()

	eye_image.Sobel(2)

	 _ = eye_image.Binary()

	// nimg = manaco.Hough(w, img)

	rel, err := filepath.Rel("data/images", path)

	outfile, err := os.Create("result/" + "test__" + rel)
	if err != nil {
		return err
	}
	defer outfile.Close()

	if err := png.Encode(outfile, eye_image.MyRGBA); err != nil {
		log.Printf("main write file :%v\n", err)
		return err
	}
	return nil
}
