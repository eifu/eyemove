package main

import (
	"./manaco"
	"flag"
	"image"
	"log"
	"os"
	"encoding/json"
	"sync"
)

func main() {
	flag.Parse()
	root := flag.Arg(0)

	dir_data, err := os.Open(root)
	if err != nil {
		log.Fatal(err)
	}

	// names has all the files from the directory
	names, err := dir_data.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	a := Concurrent(names[:13])

	log.Println(a)

	json_data, _ := json.Marshal(a)

	log.Println(string(json_data))

	f, err := os.Create("hi.json")
	if err != nil{
		panic(err)
	}
	defer f.Close()

	n, err := f.Write(json_data)
	if err != nil{
		panic(err)
	}
	log.Printf("Wrote %d bytes\n", n)


}

func Concurrent(names []string) []*manaco.EyeImage {
	wg := new(sync.WaitGroup)
	wg.Add(4)

	final := make([]*manaco.EyeImage, len(names))
	
	size := len(names)
	for i := 0; i < 4; i++ {
		n := names[i*size/4:(i+1)*size/4]
		f := final[i*size/4:(i+1)*size/4]
		go oneQuarter(wg, &n, &f)
	}

	wg.Wait()

	return final
}

func oneQuarter(wg *sync.WaitGroup, names *[]string, result *[]*manaco.EyeImage) {
	wg2 := new(sync.WaitGroup)
	wg2.Add(len(*names))

	for i, e := range *names {
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
