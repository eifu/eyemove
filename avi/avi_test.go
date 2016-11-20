package avi

import (
	"fmt"
	"os"
	"testing"
)

func TestRealAVIFile(t *testing.T) {
	// Read avi file
	file, err := os.Open("test1.avi")
	if err != nil {
		t.Error(err)
	}

	// Construct avi struct, and read avi header.
	avi, err := HeadReader(file)
	if err != nil {
		t.Errorf(" %#v\n", err)
	}

	// Print out what it has in the avi struct.
	fmt.Printf("%#v \n", avi)

	// Read MOVI list in the avi file.
	// the argument is the num of frames this function reads
	avi.MOVIReader(40)

	// this prints out the entire structure of avi struct
	avi.AVIPrint()
}
