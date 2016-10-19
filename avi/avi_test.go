package avi

import (
	"bytes"

	"os"
	"testing"
)

func TestRealAVIFile(t *testing.T) {
	file, err := os.Open("test1.avi") // For read access.
	if err != nil {
		t.Error(err)
	}

	data := make([]byte, 12+12+65792+12)

	if _, err = file.Read(data); err != nil {
		t.Error(err)
	}

	avi, err := HeadReader(bytes.NewReader(data))

	if err != nil {
		t.Errorf(" %#v %s", data, err)
	}

	avi.MOVIReader()
	avi.AVIPrint()
}
