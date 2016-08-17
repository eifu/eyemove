package main

import (
	"./manaco"
	"image"
	"image/png"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	file, err := os.Open("data/test2.jpg")
	defer file.Close()
	if err != nil {
		log.Printf("main open file :%v\n", err)
		os.Exit(1)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		log.Printf("main read file :%v\n", err)
		os.Exit(1)
	}

	nimg := manaco.GaussianFilter(img)

	nimg, _ = manaco.CutoffRGBA(nimg)

	nimg = manaco.Sb(nimg, 1)

	_, w := manaco.Binary(nimg)

	nimg = manaco.Hough(w, img)

	if err := png.Encode(os.Stdout, nimg); err != nil {
		log.Printf("main write file :%v\n", err)
		os.Exit(1)
	}
	log.Printf("Process took %.3fs total\n", time.Since(start).Seconds())

}
