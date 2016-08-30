package main

import (
	"./manaco"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	for i := 1; i < 10; i++ {
		submain(i)
	}
}

func submain(filename int) {
	start := time.Now()
	infile, err := os.Open("data/images/" + "AVI__Eye_378_1_4002900439_024.avi__0000" + strconv.Itoa(filename) + ".jpg")
	defer infile.Close()
	if err != nil {
		log.Printf("main open file :%v\n", err)
		os.Exit(1)
	}
	img, _, err := image.Decode(infile)
	if err != nil {
		log.Printf("main read file :%v\n", err)
		os.Exit(1)
	}

	nimg := manaco.GaussianFilter(img)

	nimg, _ = manaco.CutoffRGBA(nimg)

	nimg = manaco.Sb(nimg, 1)

	_, w := manaco.Binary(nimg)

	nimg = manaco.Hough(w, img)

	outfile, err := os.Create("result/" + strconv.Itoa(filename) + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	if err := png.Encode(outfile, nimg); err != nil {
		log.Printf("main write file :%v\n", err)
		os.Exit(1)
	}
	log.Printf("Process took %.3fs total\n", time.Since(start).Seconds())
}
