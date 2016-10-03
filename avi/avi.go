// A detailed description of the format is at
// https://msdn.microsoft.com/en-us/library/ms779636.aspx
package avi

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math"
)

var (
	errMissingPaddingByte     = errors.New("avi: missing padding byte")
	errMissingKeywordHeader   = errors.New("avi: missing keyword")
	errMissingRIFFChunkHeader = errors.New("avi: missing RIFF chunk header")
	errMissingAVIChunkHeader  = errors.New("avi: missing AVI chunk header")
	errListSubchunkTooLong    = errors.New("avi: list subchunk too long")
	errShortChunkData         = errors.New("avi: short chunk data")
	errShortChunkHeader       = errors.New("avi: short chunk header")
	errShortListData          = errors.New("avi: short list data")
	errShortListHeader        = errors.New("avi: short list header")
	errStaleReader            = errors.New("avi: stale reader")
)

// u32 decodes the first four bytes of b as a little-endian integer.
func decodeU32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func encodeU32(u uint32) string {
	return string([]byte{
		byte(u >> 0),
		byte(u >> 8),
		byte(u >> 16),
		byte(u >> 2)})
}

func encode(chunkID, contents string) string {
	n := len(contents)
	if n&1 == 1 {
		contents += "\x00"
	}
	return chunkID + encodeU32(uint32(n)) + contents
}

// FourCC is a four character code.
type FOURCC [4]byte

// 'RIFF' fileSize fileType (data)
// fileSize includes size of fileType, data but not include size of fileSize, 'RIFF'
type AVI struct {
	fileSize [4]byte
	data     io.Reader
}

type AVIHeader struct {
	fcc                   FOURCC
	cb                    uint32
	dwMicroSecPerFrame    uint32
	dwMaxBytesPerSec      uint32
	dwPaddingGranularity  uint32
	dwFlags               uint32
	dwTotalFrames         uint32
	dwInitialFrames       uint32
	dwStreams             uint32
	dwSuggestedBufferSize uint32
	dwWidth               uint32
	dwHeight              uint32
	dwReserved            uint32
}

type StrHeader struct {
	fcc                   FOURCC
	cb                    uint32
	fccType               FOURCC
	fccHandler            FOURCC
	dwFlags               uint32
	wPriority             uint32
	wLanguage             uint32
	dwInitialFrames       uint32
	dwScale               uint32
	dwRate                uint32
	dwStart               uint32
	dwLength              uint32
	dwSuggestedBufferSize uint32
	dwQuality             uint32
	dwSampleSize          uint32
	rcFrame               [4]uint32
}

// StrFormat is the same as BITFORMAT
type StrFormat struct {
}

// ckID ckSize ckData
//
// ckID FOUCC, ckSize uint32, ckData io.Reader
// ckSize includes size of ckData, but not include size of padding, ckID, ckSize
// The data is always padded to nearest WORD boundary.
type Chunk struct {
	ckID   FOURCC
	ckSize uint32
	ckData map[string]uint32
}

// 'LIST' listSize listType listData
// listSize includes size of listType, listdata, but not include 'LIST', listSize
type List struct {
	listSize [4]byte
	listType FOURCC
	listData io.Reader
}

func (l List) PrettyPrint() (string, error) {
	for {

	}
}

var (
	fccRIFF                 = FOURCC{'R', 'I', 'F', 'F'} // RIFF is super class of avi file
	fccAVI                  = FOURCC{'A', 'V', 'I', ' '} // AVI is identifier of avi file
	fccLIST                 = FOURCC{'L', 'I', 'S', 'T'} // LIST is identifier of LIST type
	fcchdrl                 = FOURCC{'h', 'd', 'r', 'l'} // hdrl is header list
	fccavih                 = FOURCC{'a', 'v', 'i', 'h'} // avih is AVI header
	fccstrf                 = FOURCC{'s', 't', 'r', 'f'} // strf is stream format
	fccstrl                 = FOURCC{'s', 't', 'r', 'l'} // strl is stream list
	fccstrh                 = FOURCC{'s', 't', 'r', 'h'} // strh is stream header
	fccstrn                 = FOURCC{'s', 't', 'r', 'n'} // strn is stream name
	fccvids                 = FOURCC{'v', 'i', 'd', 's'} // vids is fccType of stream
	fccmovi                 = FOURCC{'m', 'o', 'v', 'i'} //
	fccrec                  = FOURCC{'r', 'e', 'c', ' '} //
	fccidx1                 = FOURCC{'i', 'd', 'x', '1'} // idx1 is indexer of image files
	fccJUNK                 = FOURCC{'J', 'U', 'N', 'K'} // JUNK is data unused.
	fccMap  map[FOURCC]bool = map[FOURCC]bool{fccRIFF: true, fccAVI: true, fccLIST: true, fcchdrl: true, fccavih: true, fccstrl: true, fccstrh: true, fccstrn: true, fccvids: true, fccmovi: true, fccrec: true, fccidx1: true}
)

func equal(a, b FOURCC) bool {
	if a[0] != b[0] || a[1] != b[1] || a[2] != b[2] || a[3] != b[3] {
		return false
	}
	return true
}

// NewReader returns the RIFF stream's form type, such as "AVI " or "WAVE", and
// its chunks as a *Reader.
func HeadReader(r io.Reader) (*AVI, error) {
	buf := make([]byte, 12)

	// Make sure that io.Reader has enough stuff to read.
	if _, err := io.ReadFull(r, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errMissingKeywordHeader
		}
		return nil, err
	}
	// Make sure the first FOURCC lieral is 'RIFF'
	if !equal([4]byte{buf[0], buf[1], buf[2], buf[3]}, fccRIFF) {
		return nil, errMissingRIFFChunkHeader
	}

	var fileSize [4]byte = [4]byte{buf[4], buf[5], buf[6], buf[7]}

	// Make sure the 9th to 11th bytes is 'AVI '
	if !equal([4]byte{buf[8], buf[9], buf[10], buf[11]}, fccAVI) {
		return nil, errMissingAVIChunkHeader
	}

	return &AVI{fileSize: fileSize, data: r}, nil
}

// ListReader returns List type
func (avi *AVI) ListHeadReader() (*List, error) {
	var l List
	var buf = make([]byte, 12)

	r := avi.data

	if _, err := io.ReadFull(r, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortListHeader
		}
		return nil, err
	}

	// Make sure that first 4th letters are "LIST"
	if !equal(FOURCC{buf[0], buf[1], buf[2], buf[3]}, fccLIST) {
		return nil, errShortListHeader
	}

	copy(l.listSize[:], buf[4:8])

	copy(l.listType[:], buf[8:])

	l.listData = r

	return &l, nil
}

func (avi *AVI) ChunkReader() (*Chunk, error) {
	buf := make([]byte, 8)
	ck := Chunk{}
	var err error
	if _, err = io.ReadFull(avi.data, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortListHeader
		}
		return nil, err
	}

	copy(ck.ckID[:], buf[:4])
	ck.ckSize = decodeU32(buf[4:])
	switch ck.ckID {
	case fccavih:
		ck.ckData, err = avi.AVIHeaderReader(ck.ckSize)
	case fccstrh:
		ck.ckData, err = avi.StreamHeaderReader(ck.ckSize)
	case fccstrf:
		ck.ckData, err = avi.StreamFormatReader(ck.ckSize)
	}

	return &ck, nil
}

func (avi *AVI) AVIHeaderReader(size uint32) (map[string]uint32, error) {
	log.Printf("avih size: %d\n", size)
	buf := make([]byte, size)
	if _, err := io.ReadFull(avi.data, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortListHeader
		}
		return nil, err
	}

	m := make(map[string]uint32)

	m["dwMicroSecPerFrame"] = decodeU32(buf[:4])
	m["dwMaxBytesPerSec"] = decodeU32(buf[4:8])
	m["dwPaddingGranularity"] = decodeU32(buf[8:12])
	m["dwFlags"] = decodeU32(buf[12:16])
	m["dwTotalFrames"] = decodeU32(buf[16:20])
	m["dwInitialFrames"] = decodeU32(buf[20:24])
	m["dwStreams"] = decodeU32(buf[24:28])
	m["dwSuggestedBufferSize"] = decodeU32(buf[28:32])
	m["dwWidth"] = decodeU32(buf[32:36])
	m["dwHeight"] = decodeU32(buf[36:40])
	m["dwReserved"] = decodeU32(buf[40:44])
	return m, nil
}

func (avi *AVI) StreamHeaderReader(size uint32) (map[string]uint32, error) {
	buf := make([]byte, size)
	log.Printf("strh size: %d\n", size)
	if _, err := io.ReadFull(avi.data, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortListHeader
		}
		return nil, err
	}
	m := make(map[string]uint32)
	m["fccType"] = decodeU32(buf[:4])
	m["handler"] = decodeU32(buf[4:8])
	m["dwFlags"] = decodeU32(buf[8:12])
	m["wPriority"] = decodeU32(buf[12:16])
	m["wLanguage"] = decodeU32(buf[16:20])
	m["dwInitialFrames"] = decodeU32(buf[20:24])
	m["dwScale"] = decodeU32(buf[24:28])
	m["dwRate"] = decodeU32(buf[28:32])
	m["dwStart"] = decodeU32(buf[32:36])
	m["dwLength"] = decodeU32(buf[36:40])
	m["dwSuggestedBufferSize"] = decodeU32(buf[40:44])
	m["dwQuality"] = decodeU32(buf[44:48])
	m["dwSampleSize"] = decodeU32(buf[48:52])
	m["rcFrame1"] = uint32(buf[48])
	m["rcFrame2"] = uint32(buf[49])
	m["rcFrame3"] = uint32(buf[50])
	m["rcFrame4"] = uint32(buf[51])

	return m, nil
}

func (avi *AVI) StreamFormatReader(size uint32) (map[string]uint32, error) {
	log.Printf("strf size: %d\n", size)
	buf := make([]byte, size)
	if _, err := io.ReadFull(avi.data, buf); err != nil {
		return nil, err
	}
	m := make(map[string]uint32)
	m["biSize"] = decodeU32(buf[:4])
	m["biWidth"] = decodeU32(buf[4:8])
	m["biHeight"] = decodeU32(buf[8:12])
	m["biPlanes"] = decodeU32(buf[12:16])
	m["biBitCount"] = decodeU32(buf[16:20])
	m["biCompression"] = decodeU32(buf[20:24])
	m["biSizeImage"] = decodeU32(buf[24:28])
	m["biXPelsPerMeter"] = decodeU32(buf[28:32])
	m["biYPelsPerMeter"] = decodeU32(buf[32:36])
	m["biClrUsed"] = decodeU32(buf[36:40])
	m["biClrImportant"] = decodeU32(buf[40:44])

	return m, nil
}

// Reader reads chunks from an underlying io.Reader.
type Reader struct {
	r   io.Reader
	err error

	totalLen uint32
	chunkLen uint32

	chunkReader *chunkReader
	buf         [8]byte
	padded      bool
}

// Next returns the next chunk's ID, length and data. It returns io.EOF if there
// are no more chunks. The io.Reader returned becomes stale after the next Next
// call, and should no longer be used.
//
// It is valid to call Next even if all of the previous chunk's data has not
// been read.
func (z *Reader) Next() (FOURCC, uint32, io.Reader, error) {
	if z.err != nil {
		return FOURCC{}, 0, nil, z.err
	}

	// Drain the rest of the previous chunk.
	if z.chunkLen != 0 {
		want := z.chunkLen
		var got int64
		got, z.err = io.Copy(ioutil.Discard, z.chunkReader)
		if z.err != nil {
			return FOURCC{}, 0, nil, z.err
		}
		if uint32(got) != want {
			z.err = errShortChunkData
		}
	}
	z.chunkReader = nil
	if z.padded {
		if z.totalLen == 0 {
			z.err = errListSubchunkTooLong
			return FOURCC{}, 0, nil, z.err
		}
		z.totalLen--

		if _, z.err = io.ReadFull(z.r, z.buf[:1]); z.err != nil {
			if z.err == io.EOF { // are there any case that z.err == io.ErrUnexpectedEOF??
				z.err = errMissingPaddingByte // what is padding byte??
			}
			return FOURCC{}, 0, nil, z.err
		}
	}

	// We are done if we have no more data.
	if z.totalLen == 0 {
		return FOURCC{}, 0, nil, io.EOF
	}

	// Read the next chunk header.
	if z.totalLen < 8 {
		return FOURCC{}, 0, nil, errShortChunkHeader
	}
	z.totalLen -= 8

	if _, z.err = io.ReadFull(z.r, z.buf[:8]); z.err != nil {
		if z.err == io.EOF || z.err == io.ErrUnexpectedEOF {
			z.err = errShortChunkHeader
		}
		return FOURCC{}, 0, nil, z.err
	}
	chunkID := FOURCC{z.buf[0], z.buf[1], z.buf[2], z.buf[3]}
	z.chunkLen = decodeU32(z.buf[4:])
	if z.chunkLen > z.totalLen {
		return FOURCC{}, 0, nil, errListSubchunkTooLong
	}
	z.padded = z.chunkLen&1 == 1
	z.chunkReader = &chunkReader{z}
	return chunkID, z.chunkLen, z.chunkReader, nil
}

type chunkReader struct {
	z *Reader
}

func (c *chunkReader) Read(p []byte) (int, error) {
	if c != c.z.chunkReader {
		return 0, errStaleReader
	}
	z := c.z
	if z.err != nil {
		if z.err == io.EOF {
			return 0, errStaleReader
		}
		return 0, z.err
	}

	n := int(z.chunkLen)
	if n == 0 {
		return 0, io.EOF
	}
	if n < 0 {
		// Converting uint32 to int overflowed.
		n = math.MaxInt32
	}
	if n > len(p) {
		n = len(p)
	}
	n, err := z.r.Read(p[:n])
	z.totalLen -= uint32(n)
	z.chunkLen -= uint32(n)
	if err != io.EOF {
		z.err = err
	}
	return n, err
}
