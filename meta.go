package main

import "encoding/binary"

const (
	metaDataPageNum = 0
)

type metaPage struct {
	freelistPage pgnum
}

func newEmptyMeta() *metaPage {
	return &metaPage{}
}

//write metaPage data into buffer
func (m *metaPage) serialize(buf []byte) {
	pos := 0
	binary.LittleEndian.PutUint64(buf[pos:], uint64(m.freelistPage))
	pos += pageNumSize
}

//read metaPage data into struct
func (m *metaPage) deserialize(buf []byte) {
	pos := 0
	m.freelistPage = pgnum(binary.LittleEndian.Uint64(buf[pos:]))
	pos += pageNumSize
}
