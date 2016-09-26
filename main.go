package main

import (
	"./manaco"
	"flag"
	"sync"
	"image"
	"log"
	"os"
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
	a := Concurrent(names[:12])

	log.Println(a)
}

func Concurrent(names []string) []*manaco.EyeImage {
    wg := new(sync.WaitGroup)
    wg.Add(4)

    name1 := names[:len(names)/4]
    name2 := names[len(names)/4:2*len(names)/4]
    name3 := names[2*len(names)/4:3*len(names)/4]
    name4 := names[3*len(names)/4:]

    final := make([]*manaco.EyeImage, len(names))

    result1 := final[:len(names)]
    result2 := final[len(names)/4:2*len(names)/4]
    result3 := final[2*len(names)/4:3*len(names)/4]
    result4 := final[3*len(names)/4:]

    go onethird(wg, name1, &result1)
    go onethird(wg, name2, &result2)
    go onethird(wg, name3, &result3)
    go onethird(wg, name4, &result4)

    wg.Wait()



    return final
}

func onethird(wg *sync.WaitGroup, names []string, result *[]*manaco.EyeImage) {
	wg2 := new(sync.WaitGroup)
	wg2.Add(len(names))

	for i, e := range names{
		(*result)[i], _ = oneFlame(e)
		wg2.Done()
	}
	wg2.Wait()
	wg.Done()
	return 
}	

func oneFlame(path string) (*manaco.EyeImage, error) {
		log.Println(path + "is loading...")
		path = "data/images/" + path
		infile, err := os.Open(path)
		defer infile.Close()
		if err != nil {
			log.Printf("main open file :%v\n", err)
			return nil, err
		}

		img, _, err := image.Decode(infile)
		if err != nil {
			log.Printf("main read file :%v\n", err)
			return nil, err
		}

		eye_image := manaco.InitEyeImage(&img)

		eye_image.GaussianFilter()

		eye_image.CutoffRGBA()

		eye_image.Sobel(2)

		w := eye_image.Binary()

		eye_image.Hough(w)

		// for i := 0; i < len(eye_image.MyRadius); i++ {
		// 	eye_image.DrawCircle(i)
		// }

		// rel, err := filepath.Rel("data/images", path)

		// outfile, err := os.Create("result/" + "test__" + rel)
		// defer outfile.Close()
		// if err != nil {
		// 	return err
		// }

		// if err := png.Encode(outfile, eye_image.MyRGBA); err != nil {
		// 	log.Printf("main write file :%v\n", err)
		// 	return err
		// }
		return eye_image, nil

}
