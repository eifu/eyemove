package main

import (
	"./manaco"
	"flag"
	"sync"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)

	f, err := os.Open(root)
	if err != nil {
		log.Fatal(err)
	}

	// names has all the files from the directory
	names, err := f.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	Concurrent(names)
}


func Convolve(names []string)  {
    wg := new(sync.WaitGroup)
    wg.Add(len(names))
    for _, e := range names{
        go func(e string) {
            oneFlame(e)
            wg.Done()
        }(e)
    }
    wg.Wait()
    return
}

func oneFlame(path string) error {
		log.Println(path + "is loading...")
		path = "data/images/" + path
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

		w := eye_image.Binary()

		eye_image.Hough(w)

		for i := 0; i < len(eye_image.MyRadius); i++ {
			eye_image.DrawCircle(i)
		}

		rel, err := filepath.Rel("data/images", path)

		outfile, err := os.Create("result/" + "test__" + rel)
		defer outfile.Close()
		if err != nil {
			return err
		}

		if err := png.Encode(outfile, eye_image.MyRGBA); err != nil {
			log.Printf("main write file :%v\n", err)
			return err
		}
		return nil

}
