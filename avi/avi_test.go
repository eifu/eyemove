package avi

import (
	"fmt"
	"os"
	"testing"
)

func TestRealAVIFile(t *testing.T) {
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
	//avi.AVIPrint()
	fmt.Printf("%#v \n", avi.GetMoviList())
}
