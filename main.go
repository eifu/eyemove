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
	Concurrent(names[:120])
}

func mini(wg *sync.WaitGroup, s string){
	oneFlame(s)
	wg.Done()
}

func Concurrent(names []string)  {
    wg := new(sync.WaitGroup)
    wg.Add(len(names))

    name1 := names[:len(names)/3]
    name2 := names[len(names)/3:2*len(names)/3]
    name3 := names[2*len(names)/3:len(names)]

    for i := 0; i < len(names)/3 ; i ++{
        go mini(wg, name1[i])

        go mini(wg, name2[i])

        go mini(wg, name3[i])
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
