// A detailed description of the format is at
// https://msdn.microsoft.com/en-us/library/ms779636.aspx
package avi

import (
	"errors"
	"fmt"
	"io"
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
	switch len(b) {
	case 4:
		return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	case 2:
		return uint32(b[0]) | uint32(b[1])<<8
	case 1:
		return uint32(b[0])
	}
	return 0
}

func decode(s string) FOURCC {
	return FOURCC{s[0], s[1], s[2], s[3]}

}

func encodeU32(u uint32) string {
	return string([]byte{
		byte(u >> 0),
		byte(u >> 8),
		byte(u >> 16),
		byte(u >> 24)})
}

// FourCC is a four character code.
type FOURCC [4]byte

// 'RIFF' fileSize fileType data
// fileSize includes size of fileType, data but not include size of fileSize, 'RIFF'
type AVI struct {
	fileSize [4]byte
	lists    []List
	opts     []Opt
	r        io.Reader
}

// 'LIST' listSize listType listData
// listSize includes size of listType, listdata, but not include 'LIST', listSize
type List struct {
	listSize uint32
	listType FOURCC
	lists    []*List
	chunks   []*Chunk
}

// ckID ckSize ckData
// ckSize includes size of ckData, but not include size of padding, ckID, ckSize
// The data is always padded to nearest WORD boundary.
type Chunk struct {
	ckID   FOURCC
	ckSize uint32
	ckData map[string]uint32
}

type SuperIndex struct {
	qwOffset   int64
	dwSize     uint32
	dwDuration uint32
}

// Optional elements are placed in brackets: [ optional element ]
type Opt struct {
	elem uint32
}

var (
	fccRIFF = FOURCC{'R', 'I', 'F', 'F'} // RIFF is super class of avi file
	fccAVI  = FOURCC{'A', 'V', 'I', ' '} // AVI is identifier of avi file
	fccLIST = FOURCC{'L', 'I', 'S', 'T'} // LIST is identifier of LIST type
	fcchdrl = FOURCC{'h', 'd', 'r', 'l'} // hdrl is header list
	fccavih = FOURCC{'a', 'v', 'i', 'h'} // avih is AVI header
	fccstrf = FOURCC{'s', 't', 'r', 'f'} // strf is stream format
	fccstrl = FOURCC{'s', 't', 'r', 'l'} // strl is stream list
	fccstrh = FOURCC{'s', 't', 'r', 'h'} // strh is stream header
	fccstrn = FOURCC{'s', 't', 'r', 'n'} // strn is stream name
	fccvids = FOURCC{'v', 'i', 'd', 's'} // vids is fccType of stream
	fccmovi = FOURCC{'m', 'o', 'v', 'i'} //
	fccrec  = FOURCC{'r', 'e', 'c', ' '} //
	fccindx = FOURCC{'i', 'n', 'd', 'x'} // indx is optional elememt in List
	fccnnix = FOURCC{'n', 'n', 'i', 'x'} // nnix is optional element in List
	fccidx1 = FOURCC{'i', 'd', 'x', '1'} // idx1 is indexer of image files
	fccJUNK = FOURCC{'J', 'U', 'N', 'K'} // JUNK is data unused.
	fccodml = FOURCC{'o', 'd', 'm', 'l'} // odml is OpenDML
)

func (fcc *FOURCC) String() string {
	return string([]byte{fcc[0], fcc[1], fcc[2], fcc[3]})
}
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

	return &AVI{fileSize: fileSize, r: r}, nil
}

// ListReader returns List type
func (avi *AVI) ListHeadReader() (*List, error) {
	var l List
	var buf = make([]byte, 12)

	if _, err := io.ReadFull(avi.r, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortListHeader
		}
		return nil, err
	}

	// Make sure that first 4th letters are "LIST"
	if !equal(FOURCC{buf[0], buf[1], buf[2], buf[3]}, fccLIST) {
		return nil, errShortListHeader
	}

	l.listSize = decodeU32(buf[4:8])
	copy(l.listType[:], buf[8:])

	switch l.listType {
	case fccstrl:
		c, err := avi.ChunkReader()
		if err != nil {
			return nil, err
		}
		l.chunks = append(l.chunks, c)
		c, err = avi.ChunkReader()
		if err != nil {
			return nil, err
		}
		l.chunks = append(l.chunks, c)

		c, err = avi.ChunkReader()
		if err != nil {
			return nil, err
		}
		l.chunks = append(l.chunks, c)
	}

	return &l, nil
}

func (l *List) ListPrint(indent string) {
	fmt.Printf("%sList (%d) %s\n", indent, l.listSize, l.listType.String())
	for _, e := range l.lists {
		e.ListPrint(indent + "\t")
	}
	for _, e := range l.chunks {
		e.ChunkPrint(indent + "\t")
	}
}

func (avi *AVI) ChunkReader() (*Chunk, error) {
	buf := make([]byte, 8)
	ck := Chunk{}
	var err error
	if _, err = io.ReadFull(avi.r, buf); err != nil {
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
	case fccindx:
		ck.ckData, err = avi.MetaIndexReader(ck.ckSize)
	}

	return &ck, nil
}

func (c *Chunk) ChunkPrint(indent string) {
	fmt.Printf("%sckID: %s\n", indent, c.ckID.String())
	fmt.Printf("%sckSize: %d\n", indent, c.ckSize)
	for k, v := range c.ckData {
		if k == "fccType" || k == "fccHandler" || k == "dwChunkId" {
			fmt.Printf("%s\t%s: %s\n", indent, k, encodeU32(v))
		} else {
			fmt.Printf("%s\t%s: %d\n", indent, k, v)
		}
	}
}

func (avi *AVI) AVIHeaderReader(size uint32) (map[string]uint32, error) {
	buf := make([]byte, size)

	if _, err := io.ReadFull(avi.r, buf); err != nil {
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
	if _, err := io.ReadFull(avi.r, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortListHeader
		}
		return nil, err
	}
	m := make(map[string]uint32)
	m["fccType"] = decodeU32(buf[:4])
	m["fccHandler"] = decodeU32(buf[4:8])
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
	buf := make([]byte, size)
	if _, err := io.ReadFull(avi.r, buf); err != nil {
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

func (avi *AVI) MetaIndexReader(size uint32) (map[string]uint32, error) {
	buf := make([]byte, size)
	if _, err := io.ReadFull(avi.r, buf); err != nil {
		return nil, err
	}
	fmt.Printf("buf size: %d,  content:%#v", len(buf), buf)
	m := make(map[string]uint32)
	m["wLongsPerEntry"] = decodeU32(buf[:2])
	m["bIndexSubType"] = decodeU32(buf[2:3])
	m["bIndexType"] = decodeU32(buf[3:4])
	m["nEntriesInUse"] = decodeU32(buf[4:8])
	m["dwChunkId"] = decodeU32(buf[8:12])
	m["dwReserved1"] = decodeU32(buf[12:16])
	m["dwReserved2"] = decodeU32(buf[16:20])
	m["dwReserved3"] = decodeU32(buf[20:24])

	// aIndex[] part
	switch m["bIndexType"] {
	case 0x0:
		m["qwOffset"] = decodeU32(buf[24:32])
		m["dwSize"] = decodeU32(buf[32:36])
		m["dwDuration"] = decodeU32(buf[36:40])
	}

	// TODO: aIndex[] might store multiple items.

	return m, nil
}
