package main

import (
	"./manaco"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "main open file :%v\n", err)
		os.Exit(1)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "main read file :%v\n", err)
		os.Exit(1)
	}

	nimg := manaco.G_smoothing(img)

	nimg, _ = manaco.CutoffRGBA(nimg)

	nimg = manaco.Sb(nimg, 2)

	_, w := manaco.Binary(nimg)

	nimg = manaco.Hough(w, img)

	if err := png.Encode(os.Stdout, nimg); err != nil {
		fmt.Fprintf(os.Stderr, "main write file :%v\n", err)
		os.Exit(1)
	}
	log.Println("Process took %.2fs total\n", time.Since(start).Seconds())

}
