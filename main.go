package main

import (
	"./manaco"
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			return submain(rel)
		})
	if err != nil {
		log.Fatal(err)
	}
}

func submain(filepath string) error {
	log.Println(filepath + "is loading...")
	infile, err := os.Open("data/images/" + filepath)
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

	nimg := manaco.GaussianFilter(img)

	nimg, _ = manaco.CutoffRGBA(nimg)

	nimg = manaco.Sobel(nimg, 1)

	_, w := manaco.Binary(nimg)

	nimg = manaco.Hough(w, img)

	outfile, err := os.Create("result/" + "test__" + filepath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	if err := jpeg.Encode(outfile, nimg, nil); err != nil {
		log.Printf("main write file :%v\n", err)
		return err
	}
	return nil
}
