// A detailed description of the format is at
// https://msdn.microsoft.com/en-us/library/ms779636.aspx
package avi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	errMissingPaddingByte     = errors.New("avi: missing padding byte")
	errMissingKeywordHeader   = errors.New("avi: missing keyword")
	errMissingRIFFChunkHeader = errors.New("avi: missing RIFF chunk header")
	errMissingAVIChunkHeader  = errors.New("avi: missing AVI chunk header")
	errMissingLIST            = errors.New("avi: missing LIST keyword")
	errListSubchunkTooLong    = errors.New("avi: list subchunk too long")
	errShortData              = errors.New("avi: short data")
	errStaleReader            = errors.New("avi: stale reader")

	fccRIFF = FOURCC{'R', 'I', 'F', 'F'}       // RIFF is super class of avi file
	fccAVI  = FOURCC{'A', 'V', 'I', ' '}       // AVI is identifier of avi file
	fccLIST = FOURCC{'L', 'I', 'S', 'T'}       // LIST is identifier of LIST type
	fcchdrl = FOURCC{'h', 'd', 'r', 'l'}       // hdrl is header list
	fccavih = FOURCC{'a', 'v', 'i', 'h'}       // avih is AVI header
	fccstrf = FOURCC{'s', 't', 'r', 'f'}       // strf is stream format
	fccstrl = FOURCC{'s', 't', 'r', 'l'}       // strl is stream list
	fccstrh = FOURCC{'s', 't', 'r', 'h'}       // strh is stream header
	fccstrn = FOURCC{'s', 't', 'r', 'n'}       // strn is stream name
	fccvids = FOURCC{'v', 'i', 'd', 's'}       // vids is fccType of stream
	fccmovi = FOURCC{'m', 'o', 'v', 'i'}       // movi
	fccdb   = FOURCC{'\x30', '\x30', 'd', 'b'} // db is uncompressed video frame
	fccrec  = FOURCC{'r', 'e', 'c', ' '}       // rec
	fccindx = FOURCC{'i', 'n', 'd', 'x'}       // indx is optional elememt in List
	fccnnix = FOURCC{'n', 'n', 'i', 'x'}       // nnix is optional element in List
	fccidx1 = FOURCC{'i', 'd', 'x', '1'}       // idx1 is indexer of image files
	fccJUNK = FOURCC{'J', 'U', 'N', 'K'}       // JUNK is data unused.
	fccodml = FOURCC{'o', 'd', 'm', 'l'}       // odml is OpenDML
	fccdmlh = FOURCC{'d', 'm', 'l', 'h'}       // dmlh is OpenDML header
)

// FourCC is a four character code.
type FOURCC [4]byte

// 'RIFF' fileSize 'AVI ' data
// fileSize includes size of 'AVI '(FOURCC), data(io.Reader)
// actual size is fileSize + 8
type AVI struct {
	file  *os.File
	Size  uint32
	lists []*List
	r     io.Reader
}

// 'LIST' listSize listType listData
// listSize includes size of listType(FOURCC), listdata(io.Reader)
// actual size is fileSize + 8
type List struct {
	Size     uint32
	Type     FOURCC
	JunkSize uint32 // JUNK is only in

	lists       []*List
	chunks      []*Chunk
	imagechunks []*ImageChunk
	imageNum    int
}

// ckID ckSize ckData
// ckSize includes size of ckData.
// actual size is ckSize + 8
// The data is always padded to nearest WORD boundary.
type Chunk struct {
	ID   FOURCC
	Size uint32
	Data map[string]uint32
}

type ImageChunk struct {
	ID      FOURCC
	Size    uint32
	Image   []byte
	ImageID int
}

type SuperIndex struct {
	qwOffset   int64
	dwSize     uint32
	dwDuration uint32
}

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
	panic("length must be 4, 2, or 1")
}

func decode(s string) FOURCC {
	return FOURCC{s[0], s[1], s[2], s[3]}
}

func encodeU32(u uint32) *FOURCC {
	return &FOURCC{byte(u >> 0), byte(u >> 8), byte(u >> 16), byte(u >> 24)}
}

func (fcc *FOURCC) String() string {
	return string([]byte{fcc[0], fcc[1], fcc[2], fcc[3]})
}
func equal(a, b FOURCC) bool {
	if a[0] != b[0] || a[1] != b[1] || a[2] != b[2] || a[3] != b[3] {
		return false
	}
	return true
}

func (avi *AVI) GetMoviList() []*ImageChunk {
	return avi.lists[1].imagechunks
}

func (avi *AVI) AVIPrint() {
	fmt.Printf("AVI (%d)\n", avi.Size)
	for _, l := range avi.lists {
		l.ListPrint("\t")
	}

}

func (l *List) ListPrint(indent string) {
	fmt.Printf("%sList (%d) %s\n", indent, l.Size, l.Type.String())

	for _, e := range l.chunks {
		e.ChunkPrint(indent + "\t")
	}

	for _, e := range l.lists {
		e.ListPrint(indent + "\t")
	}
	if l.JunkSize != 0 {
		fmt.Printf("\t%sJUNK (%d)\n", indent, l.JunkSize)
	}

	for _, e := range l.imagechunks {
		e.ImageChunkPrint(indent + "\t")
	}
}

func (c *Chunk) ChunkPrint(indent string) {
	fmt.Printf("%s%s(%d)\n", indent, c.ID, c.Size)
	for k, v := range c.Data {
		if k == "fccType" || k == "fccHandler" || k == "dwChunkId" {
			fmt.Printf("%s\t%s: %s\n", indent, k, encodeU32(v))
		} else {
			fmt.Printf("%s\t%s: %d\n", indent, k, v)
		}
	}
}

func (ick *ImageChunk) ImageChunkPrint(indent string) {
	fmt.Printf("%s%s ID: %d\n", indent, ick.ID, ick.ImageID)
}

func (avi *AVI) readData(size uint32) ([]byte, error) {
	data := make([]byte, size)
	if _, err := avi.file.Read(data); err != nil {
		return nil, err
	}
	avi.r = bytes.NewReader(data)

	buf := make([]byte, size)
	if n, err := io.ReadFull(avi.r, buf); err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = errShortData
		}
		fmt.Println(n, " out of  ", size)
		return nil, err
	}
	return buf, nil
}

// NewReader returns the RIFF stream's form type, such as "AVI " or "WAVE", and
// its chunks as a *Reader.
func HeadReader(f *os.File) (*AVI, error) {

	avi := &AVI{file: f}
	buf, err := avi.readData(12)
	if err != nil {
		return nil, err
	}

	// Make sure the first FOURCC lieral is 'RIFF'
	if !equal([4]byte{buf[0], buf[1], buf[2], buf[3]}, fccRIFF) {
		return nil, errMissingRIFFChunkHeader
	}

	// Make sure the 9th to 11th bytes is 'AVI '
	if !equal([4]byte{buf[8], buf[9], buf[10], buf[11]}, fccAVI) {
		return nil, errMissingAVIChunkHeader
	}

	avi.Size = decodeU32(buf[4:8])

	// hdrl
	list, err := avi.ListReader()
	if err != nil {
		return nil, err
	}
	avi.lists = append(avi.lists, list)

	return avi, nil
}

// ListReader returns List type
func (avi *AVI) ListReader() (*List, error) {
	var l List

	buf, err := avi.readData(12)
	if err != nil {
		return nil, err
	}

	// Make sure that first 4th letters are "LIST"
	if !equal(FOURCC{buf[0], buf[1], buf[2], buf[3]}, fccLIST) {
		return nil, errMissingLIST
	}

	l.Size = decodeU32(buf[4:8])
	copy(l.Type[:], buf[8:12])

	switch l.Type {
	case fcchdrl:
		// avih chunk ... 8 + 56 bytes
		if err := avi.ChunkReader(&l); err != nil {
			return nil, err
		}

		// strl List ... 12 + 56
		l2, err := avi.ListReader()
		if err != nil {
			return nil, err
		}
		l.lists = append(l.lists, l2)

		// odml List ... 12 + 40
		l3, err := avi.ListReader()
		if err != nil {
			return nil, err
		}
		l.lists = append(l.lists, l3)

		// JUNK ... 12 + 64496
		if err := avi.JUNKReader(&l); err != nil {
			return nil, err
		}

	case fccstrl:
		// strh 8 + 56
		if err := avi.ChunkReader(&l); err != nil {
			return nil, err
		}

		// strf 8 + 1064
		if err := avi.ChunkReader(&l); err != nil {
			return nil, err
		}

		// indx 8 + 40
		if err := avi.ChunkReader(&l); err != nil {
			return nil, err
		}

	case fccodml:
		// dmlr 8 + 4
		if err := avi.ChunkReader(&l); err != nil {
			return nil, err
		}
	}
	return &l, nil
}

func (avi *AVI) MOVIReader(num int) {
	var movi_list List

	buf, err := avi.readData(12)
	if err != nil {
		return
	}

	// Make sure that first 4 letters are "LIST"
	if !equal(FOURCC{buf[0], buf[1], buf[2], buf[3]}, fccLIST) {
		return
	}

	// Make sure that third 4 letters are "movi"
	if !equal(FOURCC{buf[8], buf[9], buf[10], buf[11]}, fccmovi) {
		return
	}

	movi_list.Type = fccmovi
	movi_list.Size = decodeU32(buf[4:8])

	for i := 0; i < num; i++ {
		avi.ImageChunkReader(&movi_list)
	}

	avi.lists = append(avi.lists, &movi_list)
}

func (avi *AVI) ImageChunkReader(l *List) error {

	ick := ImageChunk{}

	l.imageNum += 1
	ick.ImageID = l.imageNum

	buf, err := avi.readData(8)
	if err != nil {
		return err
	}

	ick.ID = FOURCC{buf[0], buf[1], buf[2], buf[3]}

	buf, err = avi.readData(decodeU32(buf[4:8]))
	if err != nil {
		return err
	}

	ick.Image = buf

	l.imagechunks = append(l.imagechunks, &ick)

	return nil
}

func (avi *AVI) ChunkReader(l *List) error {

	buf, err := avi.readData(8)
	if err != nil {
		return err
	}

	ck := Chunk{}

	copy(ck.ID[:], buf[:4])
	ck.Size = decodeU32(buf[4:])
	switch ck.ID {
	case fccavih:
		ck.Data, err = avi.AVIHeaderReader(ck.Size)
	case fccstrh:
		ck.Data, err = avi.StreamHeaderReader(ck.Size)
	case fccstrf:
		ck.Data, err = avi.StreamFormatReader(ck.Size)
	case fccindx:
		ck.Data, err = avi.MetaIndexReader(ck.Size)
	case fccdmlh:
		ck.Data, err = avi.ExtendedAVIHeaderReader(ck.Size)
	}
	if err != nil {
		return err
	}

	l.chunks = append(l.chunks, &ck) // add chunk object ck to l.chunks
	return nil
}

func (avi *AVI) JUNKReader(l *List) error {

	buf, err := avi.readData(8)
	if err != nil {
		return err
	}

	if !equal(FOURCC{buf[0], buf[1], buf[2], buf[3]}, fccJUNK) {
		return errMissingKeywordHeader
	}
	l.JunkSize = decodeU32(buf[4:8])

	buf, err = avi.readData(l.JunkSize)

	return nil
}

func (avi *AVI) AVIHeaderReader(size uint32) (map[string]uint32, error) {

	buf, err := avi.readData(size)
	if err != nil {
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

	buf, err := avi.readData(size)
	if err != nil {
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

	buf, err := avi.readData(size)
	if err != nil {
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

	buf, err := avi.readData(size)
	if err != nil {
		return nil, err
	}

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
		m["qwOffset1"] = decodeU32(buf[24:28])
		m["qwOffset2"] = decodeU32(buf[28:32])
		m["dwSize"] = decodeU32(buf[32:36])
		m["dwDuration"] = decodeU32(buf[36:40])
	}

	// TODO: aIndex[] might store multiple items.

	return m, nil
}

func (avi *AVI) ExtendedAVIHeaderReader(size uint32) (map[string]uint32, error) {
	buf, err := avi.readData(size)
	if err != nil {
		return nil, err
	}

	m := make(map[string]uint32)
	m["dwTotalFrames"] = decodeU32(buf[:4])
	return m, nil
}
