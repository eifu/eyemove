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

	data := make([]byte, 8+8+8+65792+8+8+19608+8+19608+8+19608)

	if _, err = file.Read(data); err != nil {
		t.Error(err)
	}

	avi, err := HeadReader(bytes.NewReader(data))

	if err != nil {
		t.Errorf(" %#v %s", data, err)
	}

	avi.AVIPrint()
}
