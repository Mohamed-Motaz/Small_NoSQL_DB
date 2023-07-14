package main

import "encoding/binary"

const (
	metaPageNum        = 0
	magicNumber uint32 = 0xD00DB00D
)

type meta struct {
	root         pgnum
	freelistPage pgnum
}

func newEmptyMeta() *meta {
	return &meta{}
}

//write meta data into buffer
func (m *meta) serialize(buf []byte) {
	pos := 0

	binary.LittleEndian.PutUint32(buf[pos:], magicNumber)
	pos += magicNumberSize

	binary.LittleEndian.PutUint64(buf[pos:], uint64(m.root))
	pos += pageNumSize

	binary.LittleEndian.PutUint64(buf[pos:], uint64(m.freelistPage))
	pos += pageNumSize
}

//read meta data into struct
func (m *meta) deserialize(buf []byte) {
	pos := 0

	magicNumberRes := binary.LittleEndian.Uint32(buf[pos:])
	pos += magicNumberSize

	if magicNumberRes != magicNumber {
		panic("The file is not a libra db file")
	}

	m.root = pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos += pageNumSize

	m.freelistPage = pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos += pageNumSize
}
