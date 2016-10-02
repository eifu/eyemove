package avi

import (
	"bytes"
	"log"
	"testing"
)

func TestNewTestReader(t *testing.T) {
	s := []byte{'\x52', '\x49', '\x46', '\x46'}              // RIFF
	s = append(s, []byte{'\x50', '\x56', '\x62', '\x16'}...) // fileSize
	s = append(s, []byte{'\x41', '\x56', '\x49', '\x20'}...) // AVI
	s = append(s, []byte{'\x4c', '\x49', '\x53', '\x54'}...) // LIST
	s = append(s, []byte{'\x00', '\x01', '\x01', '\x00'}...) // listSize
	s = append(s, []byte{'\x68', '\x64', '\x72', '\x6c'}...) // hdrl
	s = append(s, []byte{'\x61', '\x76', '\x69', '\x68'}...) // avih
	s = append(s, []byte{'\x38', '\x00', '\x00', '\x00'}...) // ckSize
	s = append(s, []byte{'\x35', '\x82', '\x00', '\x00'}...) // dwMicroSecPerFrame
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // dwMaxBytesPerSec
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // dwPaddingGradularity
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // dwFlags
	s = append(s, []byte{'\x10', '\x08', '\x00', '\x00'}...) // dwTotalFrames
	s = append(s, []byte{'\xae', '\x4a', '\x00', '\x00'}...) // dwInitialFrames
	s = append(s, []byte{'\x01', '\x00', '\x00', '\x00'}...) // dwStreams
	s = append(s, []byte{'\x98', '\x4c', '\x00', '\x00'}...) // dwSuggestedBufferSize
	s = append(s, []byte{'\xac', '\x00', '\x00', '\x00'}...) // dwWidth
	s = append(s, []byte{'\x72', '\x00', '\x00', '\x00'}...) // dwHeight
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // dwReserved
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x4c', '\x49', '\x53', '\x54'}...) // LIST
	s = append(s, []byte{'\xa4', '\x04', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x73', '\x74', '\x72', '\x6c'}...) // strl
	s = append(s, []byte{'\x73', '\x74', '\x72', '\x68'}...) // strh
	s = append(s, []byte{'\x38', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x76', '\x69', '\x64', '\x73'}...) // vids
	s = append(s, []byte{'\x44', '\x49', '\x42', '\x20'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x40', '\x42', '\x0f', '\x00'}...)
	s = append(s, []byte{'\x80', '\xc3', '\xc9', '\x01'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\xae', '\x4a', '\x00', '\x00'}...)
	s = append(s, []byte{'\x98', '\x4c', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x73', '\x74', '\x72', '\x66'}...) // strf
	s = append(s, []byte{'\x28', '\x04', '\x00', '\x00'}...)
	s = append(s, []byte{'\x28', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\xac', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x72', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x01', '\x00', '\x08', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x98', '\x4c', '\x00', '\x00'}...)

	avi, err := HeadReader(bytes.NewReader(s))
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}

	//log.Printf("avi  %#v\n", avi)
	_ = avi

	list, err := avi.ListHeadReader()
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}
	//	log.Printf("final: list  %#v  \n error%s\n", list, err)
	_ = list

	avih, err := avi.AVIHeaderReader()
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}
	//	log.Printf("%#v\n", avih)
	_ = avih

	strl, err := avi.ListHeadReader()
	if err != nil {
		t.Errorf("%#v \n", strl)
	}
	log.Printf("%#v\n", strl)

	strh, err := avi.StreamHeaderReader()
	if err != nil {
		t.Errorf("%#v\n", strh)
	}
	log.Printf("%#v\n", strh)

}
