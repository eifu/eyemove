package avi

import (
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

	avi.MOVIReader()
	avi.AVIPrint()
}
