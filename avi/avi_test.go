package avi

import (
	"bytes"
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
	s = append(s, []byte{'\x10', '\x08', '\x00', '\x00'}...) // dwFlags

	s = append(s, []byte{'\xae', '\x67', '\x00', '\x00'}...) // dwTotalFrames
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // dwInitialFrames
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
	s = append(s, []byte{'\x38', '\x00', '\x00', '\x00'}...) // size of stream header
	s = append(s, []byte{'\x76', '\x69', '\x64', '\x73'}...) // vids fccType

	s = append(s, []byte{'\x44', '\x49', '\x42', '\x20'}...) // handler
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // wPriority
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // wLanguage
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // dwInitialFrames

	s = append(s, []byte{'\x40', '\x42', '\x0f', '\x00'}...) // dwScale
	s = append(s, []byte{'\x80', '\xc3', '\xc9', '\x01'}...) // dwRate
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\xae', '\x4a', '\x00', '\x00'}...)

	s = append(s, []byte{'\x98', '\x4c', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //

	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) //
	s = append(s, []byte{'\x73', '\x74', '\x72', '\x66'}...) // strf: video stream format
	s = append(s, []byte{'\x28', '\x04', '\x00', '\x00'}...) // size
	s = append(s, []byte{'\x28', '\x00', '\x00', '\x00'}...) // biSize

	s = append(s, []byte{'\xac', '\x00', '\x00', '\x00'}...) // biWidth
	s = append(s, []byte{'\x72', '\x00', '\x00', '\x00'}...) // biHeight
	s = append(s, []byte{'\x01', '\x00', '\x08', '\x00'}...) // biPlanes
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // biBitCount

	s = append(s, []byte{'\x98', '\x4c', '\x00', '\x00'}...) // biCompression
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // biSizeImage
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...) // biXPelsPerMeter
	s = append(s, []byte{'\x00', '\x01', '\x00', '\x00'}...) // biYPelsPerMeter

	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x01', '\x01', '\x01', '\x00'}...)
	s = append(s, []byte{'\x02', '\x02', '\x02', '\x00'}...)

	s = append(s, []byte{'\x03', '\x03', '\x03', '\x00'}...)
	s = append(s, []byte{'\x04', '\x04', '\x04', '\x00'}...)
	s = append(s, []byte{'\x05', '\x05', '\x05', '\x00'}...)
	s = append(s, []byte{'\x06', '\x06', '\x06', '\x00'}...)

	s = append(s, []byte{'\x07', '\x07', '\x07', '\x00'}...)
	s = append(s, []byte{'\x08', '\x08', '\x08', '\x00'}...)
	s = append(s, []byte{'\x09', '\x09', '\x09', '\x00'}...)
	s = append(s, []byte{'\x0a', '\x0a', '\x0a', '\x00'}...)

	s = append(s, []byte{'\x0b', '\x0b', '\x0b', '\x00'}...)
	s = append(s, []byte{'\x0c', '\x0c', '\x0c', '\x00'}...)
	s = append(s, []byte{'\x0d', '\x0d', '\x0d', '\x00'}...)
	s = append(s, []byte{'\x0e', '\x0e', '\x0e', '\x00'}...)

	s = append(s, []byte{'\x0f', '\x0f', '\x0f', '\x00'}...)
	s = append(s, []byte{'\x10', '\x10', '\x10', '\x00'}...)
	s = append(s, []byte{'\x11', '\x11', '\x11', '\x00'}...)
	s = append(s, []byte{'\x12', '\x12', '\x12', '\x00'}...)

	s = append(s, []byte{'\x13', '\x13', '\x13', '\x00'}...)
	s = append(s, []byte{'\x14', '\x14', '\x14', '\x00'}...)
	s = append(s, []byte{'\x15', '\x15', '\x15', '\x00'}...)
	s = append(s, []byte{'\x16', '\x16', '\x16', '\x00'}...)

	s = append(s, []byte{'\x17', '\x17', '\x17', '\x00'}...)
	s = append(s, []byte{'\x18', '\x18', '\x18', '\x00'}...)
	s = append(s, []byte{'\x19', '\x19', '\x19', '\x00'}...)
	s = append(s, []byte{'\x1a', '\x1a', '\x1a', '\x00'}...)

	s = append(s, []byte{'\x1b', '\x1b', '\x1b', '\x00'}...)
	s = append(s, []byte{'\x1c', '\x1c', '\x1c', '\x00'}...)
	s = append(s, []byte{'\x1d', '\x1d', '\x1d', '\x00'}...)
	s = append(s, []byte{'\x1e', '\x1e', '\x1e', '\x00'}...)

	s = append(s, []byte{'\x1f', '\x1f', '\x1f', '\x00'}...)
	s = append(s, []byte{'\x20', '\x20', '\x20', '\x00'}...)
	s = append(s, []byte{'\x21', '\x21', '\x21', '\x00'}...)
	s = append(s, []byte{'\x22', '\x22', '\x22', '\x00'}...)

	s = append(s, []byte{'\x23', '\x23', '\x23', '\x00'}...)
	s = append(s, []byte{'\x24', '\x24', '\x24', '\x00'}...)
	s = append(s, []byte{'\x25', '\x25', '\x25', '\x00'}...)
	s = append(s, []byte{'\x26', '\x26', '\x26', '\x00'}...)

	s = append(s, []byte{'\x27', '\x27', '\x27', '\x00'}...)
	s = append(s, []byte{'\x28', '\x28', '\x28', '\x00'}...)
	s = append(s, []byte{'\x29', '\x29', '\x29', '\x00'}...)
	s = append(s, []byte{'\x2a', '\x2a', '\x2a', '\x00'}...)

	s = append(s, []byte{'\x2b', '\x2b', '\x2b', '\x00'}...)
	s = append(s, []byte{'\x2c', '\x2c', '\x2c', '\x00'}...)
	s = append(s, []byte{'\x2d', '\x2d', '\x2d', '\x00'}...)
	s = append(s, []byte{'\x2e', '\x2e', '\x2e', '\x00'}...)

	s = append(s, []byte{'\x2f', '\x2f', '\x2f', '\x00'}...)
	s = append(s, []byte{'\x30', '\x30', '\x30', '\x00'}...)
	s = append(s, []byte{'\x31', '\x31', '\x31', '\x00'}...)
	s = append(s, []byte{'\x32', '\x32', '\x32', '\x00'}...)

	s = append(s, []byte{'\x33', '\x33', '\x33', '\x00'}...)
	s = append(s, []byte{'\x34', '\x34', '\x34', '\x00'}...)
	s = append(s, []byte{'\x35', '\x35', '\x35', '\x00'}...)
	s = append(s, []byte{'\x36', '\x36', '\x36', '\x00'}...)

	s = append(s, []byte{'\x37', '\x37', '\x37', '\x00'}...)
	s = append(s, []byte{'\x38', '\x38', '\x38', '\x00'}...)
	s = append(s, []byte{'\x39', '\x39', '\x39', '\x00'}...)
	s = append(s, []byte{'\x3a', '\x3a', '\x3a', '\x00'}...)

	s = append(s, []byte{'\x3b', '\x3b', '\x3b', '\x00'}...)
	s = append(s, []byte{'\x3c', '\x3c', '\x3c', '\x00'}...)
	s = append(s, []byte{'\x3d', '\x3d', '\x3d', '\x00'}...)
	s = append(s, []byte{'\x3e', '\x3e', '\x3e', '\x00'}...)

	s = append(s, []byte{'\x3f', '\x3f', '\x3f', '\x00'}...)
	s = append(s, []byte{'\x40', '\x40', '\x40', '\x00'}...)
	s = append(s, []byte{'\x41', '\x41', '\x41', '\x00'}...)
	s = append(s, []byte{'\x42', '\x42', '\x42', '\x00'}...)

	s = append(s, []byte{'\x43', '\x43', '\x43', '\x00'}...)
	s = append(s, []byte{'\x44', '\x44', '\x44', '\x00'}...)
	s = append(s, []byte{'\x45', '\x45', '\x45', '\x00'}...)
	s = append(s, []byte{'\x46', '\x46', '\x46', '\x00'}...)

	s = append(s, []byte{'\x47', '\x47', '\x47', '\x00'}...)
	s = append(s, []byte{'\x48', '\x48', '\x48', '\x00'}...)
	s = append(s, []byte{'\x49', '\x49', '\x49', '\x00'}...)
	s = append(s, []byte{'\x4a', '\x4a', '\x4a', '\x00'}...)

	s = append(s, []byte{'\x4b', '\x4b', '\x4b', '\x00'}...)
	s = append(s, []byte{'\x4c', '\x4c', '\x4c', '\x00'}...)
	s = append(s, []byte{'\x4d', '\x4d', '\x4d', '\x00'}...)
	s = append(s, []byte{'\x4e', '\x4e', '\x4e', '\x00'}...)

	s = append(s, []byte{'\x4f', '\x4f', '\x4f', '\x00'}...)
	s = append(s, []byte{'\x50', '\x50', '\x50', '\x00'}...)
	s = append(s, []byte{'\x51', '\x51', '\x51', '\x00'}...)
	s = append(s, []byte{'\x52', '\x52', '\x52', '\x00'}...)

	s = append(s, []byte{'\x53', '\x53', '\x53', '\x00'}...)
	s = append(s, []byte{'\x54', '\x54', '\x54', '\x00'}...)
	s = append(s, []byte{'\x55', '\x55', '\x55', '\x00'}...)
	s = append(s, []byte{'\x56', '\x56', '\x56', '\x00'}...)

	s = append(s, []byte{'\x57', '\x57', '\x57', '\x00'}...)
	s = append(s, []byte{'\x58', '\x58', '\x58', '\x00'}...)
	s = append(s, []byte{'\x59', '\x59', '\x59', '\x00'}...)
	s = append(s, []byte{'\x5a', '\x5a', '\x5a', '\x00'}...)

	s = append(s, []byte{'\x5b', '\x5b', '\x5b', '\x00'}...)
	s = append(s, []byte{'\x5c', '\x5c', '\x5c', '\x00'}...)
	s = append(s, []byte{'\x5d', '\x5d', '\x5d', '\x00'}...)
	s = append(s, []byte{'\x5e', '\x5e', '\x5e', '\x00'}...)

	s = append(s, []byte{'\x5f', '\x5f', '\x5f', '\x00'}...)
	s = append(s, []byte{'\x60', '\x60', '\x60', '\x00'}...)
	s = append(s, []byte{'\x61', '\x61', '\x61', '\x00'}...)
	s = append(s, []byte{'\x62', '\x62', '\x62', '\x00'}...)

	s = append(s, []byte{'\x63', '\x63', '\x63', '\x00'}...)
	s = append(s, []byte{'\x64', '\x64', '\x64', '\x00'}...)
	s = append(s, []byte{'\x65', '\x65', '\x65', '\x00'}...)
	s = append(s, []byte{'\x66', '\x66', '\x66', '\x00'}...)

	s = append(s, []byte{'\x67', '\x67', '\x67', '\x00'}...)
	s = append(s, []byte{'\x68', '\x68', '\x68', '\x00'}...)
	s = append(s, []byte{'\x69', '\x69', '\x69', '\x00'}...)
	s = append(s, []byte{'\x6a', '\x6a', '\x6a', '\x00'}...)

	s = append(s, []byte{'\x6b', '\x6b', '\x6b', '\x00'}...)
	s = append(s, []byte{'\x6c', '\x6c', '\x6c', '\x00'}...)
	s = append(s, []byte{'\x6d', '\x6d', '\x6d', '\x00'}...)
	s = append(s, []byte{'\x6e', '\x6e', '\x6e', '\x00'}...)

	s = append(s, []byte{'\x6f', '\x6f', '\x6f', '\x00'}...)
	s = append(s, []byte{'\x70', '\x70', '\x70', '\x00'}...)
	s = append(s, []byte{'\x71', '\x71', '\x71', '\x00'}...)
	s = append(s, []byte{'\x72', '\x72', '\x72', '\x00'}...)

	s = append(s, []byte{'\x73', '\x73', '\x73', '\x00'}...)
	s = append(s, []byte{'\x74', '\x74', '\x74', '\x00'}...)
	s = append(s, []byte{'\x75', '\x75', '\x75', '\x00'}...)
	s = append(s, []byte{'\x76', '\x76', '\x76', '\x00'}...)

	s = append(s, []byte{'\x77', '\x77', '\x77', '\x00'}...)
	s = append(s, []byte{'\x78', '\x78', '\x78', '\x00'}...)
	s = append(s, []byte{'\x79', '\x79', '\x79', '\x00'}...)
	s = append(s, []byte{'\x7a', '\x7a', '\x7a', '\x00'}...)

	s = append(s, []byte{'\x7b', '\x7b', '\x7b', '\x00'}...)
	s = append(s, []byte{'\x7c', '\x7c', '\x7c', '\x00'}...)
	s = append(s, []byte{'\x7d', '\x7d', '\x7d', '\x00'}...)
	s = append(s, []byte{'\x7e', '\x7e', '\x7e', '\x00'}...)

	s = append(s, []byte{'\x7f', '\x7f', '\x7f', '\x00'}...)
	s = append(s, []byte{'\x80', '\x80', '\x80', '\x00'}...)
	s = append(s, []byte{'\x81', '\x81', '\x81', '\x00'}...)
	s = append(s, []byte{'\x82', '\x82', '\x82', '\x00'}...)

	s = append(s, []byte{'\x83', '\x83', '\x83', '\x00'}...)
	s = append(s, []byte{'\x84', '\x84', '\x84', '\x00'}...)
	s = append(s, []byte{'\x85', '\x85', '\x85', '\x00'}...)
	s = append(s, []byte{'\x86', '\x86', '\x86', '\x00'}...)

	s = append(s, []byte{'\x87', '\x87', '\x87', '\x00'}...)
	s = append(s, []byte{'\x88', '\x88', '\x88', '\x00'}...)
	s = append(s, []byte{'\x89', '\x89', '\x89', '\x00'}...)
	s = append(s, []byte{'\x8a', '\x8a', '\x8a', '\x00'}...)

	s = append(s, []byte{'\x8b', '\x8b', '\x8b', '\x00'}...)
	s = append(s, []byte{'\x8c', '\x8c', '\x8c', '\x00'}...)
	s = append(s, []byte{'\x8d', '\x8d', '\x8d', '\x00'}...)
	s = append(s, []byte{'\x8e', '\x8e', '\x8e', '\x00'}...)

	s = append(s, []byte{'\x8f', '\x8f', '\x8f', '\x00'}...)
	s = append(s, []byte{'\x90', '\x90', '\x90', '\x00'}...)
	s = append(s, []byte{'\x91', '\x91', '\x91', '\x00'}...)
	s = append(s, []byte{'\x92', '\x92', '\x92', '\x00'}...)

	s = append(s, []byte{'\x93', '\x93', '\x93', '\x00'}...)
	s = append(s, []byte{'\x94', '\x94', '\x94', '\x00'}...)
	s = append(s, []byte{'\x95', '\x95', '\x95', '\x00'}...)
	s = append(s, []byte{'\x96', '\x96', '\x96', '\x00'}...)

	s = append(s, []byte{'\x97', '\x97', '\x97', '\x00'}...)
	s = append(s, []byte{'\x98', '\x98', '\x98', '\x00'}...)
	s = append(s, []byte{'\x99', '\x99', '\x99', '\x00'}...)
	s = append(s, []byte{'\x9a', '\x9a', '\x9a', '\x00'}...)

	s = append(s, []byte{'\x9b', '\x9b', '\x9b', '\x00'}...)
	s = append(s, []byte{'\x9c', '\x9c', '\x9c', '\x00'}...)
	s = append(s, []byte{'\x9d', '\x9d', '\x9d', '\x00'}...)
	s = append(s, []byte{'\x9e', '\x9e', '\x9e', '\x00'}...)

	s = append(s, []byte{'\x9f', '\x9f', '\x9f', '\x00'}...)
	s = append(s, []byte{'\xa0', '\xa0', '\xa0', '\x00'}...)
	s = append(s, []byte{'\xa1', '\xa1', '\xa1', '\x00'}...)
	s = append(s, []byte{'\xa2', '\xa2', '\xa2', '\x00'}...)

	s = append(s, []byte{'\xa3', '\xa3', '\xa3', '\x00'}...)
	s = append(s, []byte{'\xa4', '\xa4', '\xa4', '\x00'}...)
	s = append(s, []byte{'\xa5', '\xa5', '\xa5', '\x00'}...)
	s = append(s, []byte{'\xa6', '\xa6', '\xa6', '\x00'}...)

	s = append(s, []byte{'\xa7', '\xa7', '\xa7', '\x00'}...)
	s = append(s, []byte{'\xa8', '\xa8', '\xa8', '\x00'}...)
	s = append(s, []byte{'\xa9', '\xa9', '\xa9', '\x00'}...)
	s = append(s, []byte{'\xaa', '\xaa', '\xaa', '\x00'}...)

	s = append(s, []byte{'\xab', '\xab', '\xab', '\x00'}...)
	s = append(s, []byte{'\xac', '\xac', '\xac', '\x00'}...)
	s = append(s, []byte{'\xad', '\xad', '\xad', '\x00'}...)
	s = append(s, []byte{'\xae', '\xae', '\xae', '\x00'}...)

	s = append(s, []byte{'\xaf', '\xaf', '\xaf', '\x00'}...)
	s = append(s, []byte{'\xb0', '\xb0', '\xb0', '\x00'}...)
	s = append(s, []byte{'\xb1', '\xb1', '\xb1', '\x00'}...)
	s = append(s, []byte{'\xb2', '\xb2', '\xb2', '\x00'}...)

	s = append(s, []byte{'\xb3', '\xb3', '\xb3', '\x00'}...)
	s = append(s, []byte{'\xb4', '\xb4', '\xb4', '\x00'}...)
	s = append(s, []byte{'\xb5', '\xb5', '\xb5', '\x00'}...)
	s = append(s, []byte{'\xb6', '\xb6', '\xb6', '\x00'}...)

	s = append(s, []byte{'\xb7', '\xb7', '\xb7', '\x00'}...)
	s = append(s, []byte{'\xb8', '\xb8', '\xb8', '\x00'}...)
	s = append(s, []byte{'\xb9', '\xb9', '\xb9', '\x00'}...)
	s = append(s, []byte{'\xba', '\xba', '\xba', '\x00'}...)

	s = append(s, []byte{'\xbb', '\xbb', '\xbb', '\x00'}...)
	s = append(s, []byte{'\xbc', '\xbc', '\xbc', '\x00'}...)
	s = append(s, []byte{'\xbd', '\xbd', '\xbd', '\x00'}...)
	s = append(s, []byte{'\xbe', '\xbe', '\xbe', '\x00'}...)

	s = append(s, []byte{'\xbf', '\xbf', '\xbf', '\x00'}...)
	s = append(s, []byte{'\xc0', '\xc0', '\xc0', '\x00'}...)
	s = append(s, []byte{'\xc1', '\xc1', '\xc1', '\x00'}...)
	s = append(s, []byte{'\xc2', '\xc2', '\xc2', '\x00'}...)

	s = append(s, []byte{'\xc3', '\xc3', '\xc3', '\x00'}...)
	s = append(s, []byte{'\xc4', '\xc4', '\xc4', '\x00'}...)
	s = append(s, []byte{'\xc5', '\xc5', '\xc5', '\x00'}...)
	s = append(s, []byte{'\xc6', '\xc6', '\xc6', '\x00'}...)

	s = append(s, []byte{'\xc7', '\xc7', '\xc7', '\x00'}...)
	s = append(s, []byte{'\xc8', '\xc8', '\xc8', '\x00'}...)
	s = append(s, []byte{'\xc9', '\xc9', '\xc9', '\x00'}...)
	s = append(s, []byte{'\xca', '\xca', '\xca', '\x00'}...)

	s = append(s, []byte{'\xcb', '\xcb', '\xcb', '\x00'}...)
	s = append(s, []byte{'\xcc', '\xcc', '\xcc', '\x00'}...)
	s = append(s, []byte{'\xcd', '\xcd', '\xcd', '\x00'}...)
	s = append(s, []byte{'\xce', '\xce', '\xce', '\x00'}...)

	s = append(s, []byte{'\xcf', '\xcf', '\xcf', '\x00'}...)
	s = append(s, []byte{'\xd0', '\xd0', '\xd0', '\x00'}...)
	s = append(s, []byte{'\xd1', '\xd1', '\xd1', '\x00'}...)
	s = append(s, []byte{'\xd2', '\xd2', '\xd2', '\x00'}...)

	s = append(s, []byte{'\xd3', '\xd3', '\xd3', '\x00'}...)
	s = append(s, []byte{'\xd4', '\xd4', '\xd4', '\x00'}...)
	s = append(s, []byte{'\xd5', '\xd5', '\xd5', '\x00'}...)
	s = append(s, []byte{'\xd6', '\xd6', '\xd6', '\x00'}...)

	s = append(s, []byte{'\xd7', '\xd7', '\xd7', '\x00'}...)
	s = append(s, []byte{'\xd8', '\xd8', '\xd8', '\x00'}...)
	s = append(s, []byte{'\xd9', '\xd9', '\xd9', '\x00'}...)
	s = append(s, []byte{'\xda', '\xda', '\xda', '\x00'}...)

	s = append(s, []byte{'\xdb', '\xdb', '\xdb', '\x00'}...)
	s = append(s, []byte{'\xdc', '\xdc', '\xdc', '\x00'}...)
	s = append(s, []byte{'\xdd', '\xdd', '\xdd', '\x00'}...)
	s = append(s, []byte{'\xde', '\xde', '\xde', '\x00'}...)

	s = append(s, []byte{'\xdf', '\xdf', '\xdf', '\x00'}...)
	s = append(s, []byte{'\xe0', '\xe0', '\xe0', '\x00'}...)
	s = append(s, []byte{'\xe1', '\xe1', '\xe1', '\x00'}...)
	s = append(s, []byte{'\xe2', '\xe2', '\xe2', '\x00'}...)

	s = append(s, []byte{'\xe3', '\xe3', '\xe3', '\x00'}...)
	s = append(s, []byte{'\xe4', '\xe4', '\xe4', '\x00'}...)
	s = append(s, []byte{'\xe5', '\xe5', '\xe5', '\x00'}...)
	s = append(s, []byte{'\xe6', '\xe6', '\xe6', '\x00'}...)

	s = append(s, []byte{'\xe7', '\xe7', '\xe7', '\x00'}...)
	s = append(s, []byte{'\xe8', '\xe8', '\xe8', '\x00'}...)
	s = append(s, []byte{'\xe9', '\xe9', '\xe9', '\x00'}...)
	s = append(s, []byte{'\xea', '\xea', '\xea', '\x00'}...)

	s = append(s, []byte{'\xeb', '\xeb', '\xeb', '\x00'}...)
	s = append(s, []byte{'\xec', '\xec', '\xec', '\x00'}...)
	s = append(s, []byte{'\xed', '\xed', '\xed', '\x00'}...)
	s = append(s, []byte{'\xee', '\xee', '\xee', '\x00'}...)

	s = append(s, []byte{'\xef', '\xef', '\xef', '\x00'}...)
	s = append(s, []byte{'\xf0', '\xf0', '\xf0', '\x00'}...)
	s = append(s, []byte{'\xf1', '\xf1', '\xf1', '\x00'}...)
	s = append(s, []byte{'\xf2', '\xf2', '\xf2', '\x00'}...)

	s = append(s, []byte{'\xf3', '\xf3', '\xf3', '\x00'}...)
	s = append(s, []byte{'\xf4', '\xf4', '\xf4', '\x00'}...)
	s = append(s, []byte{'\xf5', '\xf5', '\xf5', '\x00'}...)
	s = append(s, []byte{'\xf6', '\xf6', '\xf6', '\x00'}...)

	s = append(s, []byte{'\xf7', '\xf7', '\xf7', '\x00'}...)
	s = append(s, []byte{'\xf8', '\xf8', '\xf8', '\x00'}...)
	s = append(s, []byte{'\xf9', '\xf9', '\xf9', '\x00'}...)
	s = append(s, []byte{'\xfa', '\xfa', '\xfa', '\x00'}...)

	s = append(s, []byte{'\xfb', '\xfb', '\xfb', '\x00'}...)
	s = append(s, []byte{'\xfc', '\xfc', '\xfc', '\x00'}...)
	s = append(s, []byte{'\xfd', '\xfd', '\xfd', '\x00'}...)
	s = append(s, []byte{'\xfe', '\xfe', '\xfe', '\x00'}...)

	s = append(s, []byte{'\xff', '\xff', '\xff', '\x00'}...)
	s = append(s, []byte{'\x69', '\x6e', '\x64', '\x78'}...)
	s = append(s, []byte{'\x28', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x04', '\x00', '\x00', '\x00'}...)

	s = append(s, []byte{'\x01', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x30', '\x30', '\x64', '\x62'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)

	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\xe0', '\x55', '\x5b', '\x16'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x90', '\x55', '\x02', '\x00'}...)

	s = append(s, []byte{'\xae', '\x4a', '\x00', '\x00'}...)
	s = append(s, []byte{'\x4c', '\x49', '\x53', '\x54'}...)
	s = append(s, []byte{'\x10', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x6f', '\x64', '\x6d', '\x6c'}...)

	s = append(s, []byte{'\x64', '\x6d', '\x6c', '\x68'}...)
	s = append(s, []byte{'\x04', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\xae', '\x4a', '\x00', '\x00'}...)
	s = append(s, []byte{'\x4a', '\x55', '\x4e', '\x4b'}...)

	s = append(s, []byte{'\xf0', '\xfb', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)
	s = append(s, []byte{'\x00', '\x00', '\x00', '\x00'}...)

	avi, err := HeadReader(bytes.NewReader(s))
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}

	list, err := avi.ListHeadReader()
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}
	_ = list
	avih, err := avi.ChunkReader()
	if err != nil {
		t.Errorf(" %#v %s", s, err)
	}
	avih.ChunkPrint("")

	strl, err := avi.ListHeadReader()
	if err != nil {
		t.Errorf("%#v \n", strl)
	}
	strl.ListPrint("")

	odml, err := avi.ListHeadReader()
	if err != nil {
		t.Error("%#v \n")
	}
	odml.ListPrint("")

}
