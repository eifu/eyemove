package main

import (
	"github.com/eifu/eyemove/avi"

	"fmt"
	"os"
	"testing"
)

func TestRealAVIFile(t *testing.T) {
	file, err := os.Open("avi/test1.avi") // For read access.
	if err != nil {
		t.Error(err)
	}

	avi_, err := avi.HeadReader(file)

	if err != nil {
		t.Errorf(" %#v\n", err)
	}
	fmt.Printf("%#v \n", avi_)
	avi_.MOVIReader(40)

	aaa := avi_.GetLists()[1]

	fmt.Printf("%#v \n", aaa)
}
