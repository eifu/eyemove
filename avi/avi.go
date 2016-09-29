// A detailed description of the format is at
// https://msdn.microsoft.com/en-us/library/ms779636.aspx
package avi

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
)

var (
	errMissingPaddingByte     = errors.New("riff: missing padding byte")
	errMissingRIFFChunkHeader = errors.New("riff: missing RIFF chunk header")
	errListSubchunkTooLong    = errors.New("riff: list subchunk too long")
	errShortChunkData         = errors.New("riff: short chunk data")
	errShortChunkHeader       = errors.New("riff: short chunk header")
	errStaleReader            = errors.New("riff: stale reader")
)

// u32 decodes the first four bytes of b as a little-endian integer.
func decodeU32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

const chunkHeaderSize = 8

// FourCC is a four character code.
type FOURCC [4]byte

type avi struct {
	fileSize ckSize
	fileType FOURCC
	data     []byte
}

type ckID FOURCC    // ckID is a FOURCC that identifies the data contained in the chunk.
type ckSize [4]byte // ckSize is a 4-byte value giving the size of the data in ckData.
type ckData []byte  // ckData is zero or more bytes of data.

type listSize [4]byte
type listType FOURCC
type listData []byte

var (
	AVI  = FOURCC{'A', 'V', 'I', ' '}
	LIST = FOURCC{'L', 'I', 'S', 'T'}
	hdrl = FOURCC{'h', 'd', 'r', 'l'}
	avih = FOURCC{'a', 'v', 'i', 'h'}
	strl = FOURCC{'s', 't', 'r', 'l'}
	strh = FOURCC{'s', 't', 'r', 'h'}
	strn = FOURCC{'s', 't', 'r', 'n'}
	movi = FOURCC{'m', 'o', 'v', 'i'}
	rec  = FOURCC{'r', 'e', 'c', ' '}
	idx1 = FOURCC{'i', 'd', 'x', '1'}
)

func equal(a, b FOURCC) bool {
	if a[0] != b[0] || a[1] != b[1] || a[2] != b[2] || a[3] != b[3]{
		return false
	}
	return true
}

// NewReader returns the RIFF stream's form type, such as "AVI " or "WAVE", and
// its chunks as a *Reader.
func HeadReader(r io.Reader) (FOURCC, *Reader, error) {
	buf := make(FOURCC)

	// Make sure that io.Reader has enough stuff to read.
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errMissingRIFFChunkHeader
		}
		return FOURCC{}, nil, err
	}
	// Make sure the first FOURCC lieral is 'RIFF'
	if buf[0] != 'R' || buf[1] != 'I' || buf[2] != 'F' || buf[3] != 'F' {
		return FOURCC{}, nil, errMissingRIFFChunkHeader
	}

	return ListReader(decodeU32(buf[4:]), r) // return the size of data and the rest of data.
}

// NewListReader returns a LIST chunk's list type, such as "movi" or "wavl",
// and its chunks as a *Reader.
func ListReader(chunkLen uint32, chunkData io.Reader) (FOURCC, *Reader, error) {
	if chunkLen < 4 {
		return FOURCC{}, nil, errShortChunkData
	}
	data := &Reader{r: chunkData}

	if _, err := io.ReadFull(chunkData, data.buf[:4]); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortChunkData
		}
		return FOURCC{}, nil, err
	}
	listType := FOURCC{data.buf[0], data.buf[1], data.buf[2], data.buf[3]}

	data.totalLen = chunkLen - 4 // totalLen is the size of data after listType('avi '')
	return listType, data, nil
}

// Reader reads chunks from an underlying io.Reader.
type Reader struct {
	r   io.Reader
	err error

	totalLen uint32
	chunkLen uint32

	chunkReader *chunkReader
	buf         [chunkHeaderSize]byte
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
	if z.padded { // what is padded??
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
	if z.totalLen < chunkHeaderSize {
		return FOURCC{}, 0, nil, errShortChunkHeader
	}
	z.totalLen -= chunkHeaderSize

	if _, z.err = io.ReadFull(z.r, z.buf[:chunkHeaderSize]); z.err != nil {
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
