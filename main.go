package main

import (
	"github.com/eifu/eyemove/manaco"

	"encoding/json"
	"flag"
	"image"
	"log"
	"os"
	"sync"
)

func main() {

	file, err := os.Open("test1.avi") // For read access.
	if err != nil {
		t.Error(err)
	}

	avi, err := HeadReader(file)

	if err != nil {
		t.Errorf(" %#v\n", err)
	}
	fmt.Printf("%#v \n", avi)
	avi.MOVIReader(40)

	// avi.GetLists[1] is movi_list
	processed := Concurrent(avi.GetLists())

	json_data, _ := json.MarshalIndent(processed, "", "    ")

	f, err := os.Create("test__data.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := f.Write(json_data)
	if err != nil {
		panic(err)
	}
	log.Printf("Wrote %d bytes\n", n)

}

func Concurrent(movi_lists []*avi.ImageChunk) []*manaco.EyeImage {
	wg := new(sync.WaitGroup)
	wg.Add(4)

	final := make([]*manaco.EyeImage, len(movi_lists))

	size := len(movi_lists)
	for i := 0; i < 4; i++ {
		n := movi_lists[i*size/4 : (i+1)*size/4]
		f := final[i*size/4 : (i+1)*size/4]
		go oneQuarter(wg, &n, &f)
	}
	wg.Wait()

	return final
}

func oneQuarter(wg *sync.WaitGroup, names *[]*avi.ImageChunk, result *[]*manaco.EyeImage) {
	wg2 := new(sync.WaitGroup)
	wg2.Add(len(*names))
	var err error

	for i, e := range *names {
		(*result)[i], err = oneFlame(e)
		if err != nil {
			panic(err)
		}
		wg2.Done()
	}
	wg2.Wait()
	wg.Done()
	return
}

func oneFlame(ick *avi.ImageChunk) (*manaco.EyeImage, error) {

	eye_image := manaco.InitEyeImage(ick)

	eye_image.GaussianFilter()

	eye_image.CutoffRGBA()

	eye_image.Sobel(2)

	w := eye_image.Binary()

	eye_image.Hough(w)

	return eye_image, nil
}
