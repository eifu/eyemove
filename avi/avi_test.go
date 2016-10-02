package avi

import (
	"bytes"
	"log"
	"testing"
)

func TestNewTestReader(t *testing.T) {
	s := []byte("\x52\x49\x46\x46")  // RIFF
	s = append(s, []byte{'\x50', '\x56', '\x62', '\x16'}...)  // fileSize
	s = append(s, []byte{'\x41', '\x56', '\x49', '\x20'}...)  // AVI 
	s = append(s, []byte{'\x4c', '\x49', '\x53', '\x54'}...)  // listSize
	s = append(s, []byte{'\x00', '\x01', '\x01', '\x00'}...)  // listdata
	s = append(s, []byte{'\x68', '\x64', '\x72', '\x6c'}...)
	s = append(s, []byte{'\x61', '\x76', '\x69', '\x68'}...)
	s = append(s, []byte{'\x38', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x35', '\x82', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x10', '\x08', '\x00', '\x00'}...)
	avi, r, err := HeadReader(bytes.NewReader(s))
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}

	log.Printf("avi  %#v  ioReader %#v\n", avi, r)

	list, r, err := avi.ListHeadReader(r)

	log.Printf("final: list  %#v  \nioReader %#v\n error%s\n", list, r, err)
}
