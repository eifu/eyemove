// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package avi

import (
	"bytes"
	"testing"
	"log"
)

func encodeU32(u uint32) []byte {
	return []byte{
		byte(u >> 0),
		byte(u >> 8),
		byte(u >> 16),
		byte(u >> 24),
	}
}

func TestNewTestReader(t *testing.T){
	s := []byte("\x52\x49\x46\x46")
	s = append(s, []byte{50,56,62,16}...)
	s = append(s, []byte{41,56,49,20}...)

	fileType, r, err := HeadReader(bytes.NewReader(s))
	if err != nil{
		t.Errorf(" %#v %s",s,err)
	}

	log.Printf("filetype  %s   reader %#v\n",fileType, r)

}