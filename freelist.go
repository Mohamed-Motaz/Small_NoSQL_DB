package main

import "encoding/binary"

//Data structure for keeping track of all free
//and empty pages
type freelist struct {
	maxPage       pgnum   //holds the maximum page allocated. maxPage * pageSize = fileSize
	releasedPages []pgnum //pages that were allocated but now free
}

const metaPageNum = 0

func newFreeList() *freelist {
	return &freelist{
		maxPage:       metaPageNum,
		releasedPages: []pgnum{},
	}
}

func (fr *freelist) getNextPage() pgnum {
	if len(fr.releasedPages) != 0 {
		//get the last page
		pgID := fr.releasedPages[len(fr.releasedPages)-1]
		//remove it from the released pages
		fr.releasedPages = fr.releasedPages[:len(fr.releasedPages)-1]
		return pgID
	}
	fr.maxPage++
	return fr.maxPage
}

func (fr *freelist) releasePage(page pgnum) {
	fr.releasedPages = append(fr.releasedPages, page)
}

func (fr *freelist) serialize(buf []byte) {
	pos := 0

	//maxPage num
	binary.LittleEndian.PutUint16(buf[pos:], uint16(fr.maxPage))
	pos += 2

	//releasedPages num
	binary.LittleEndian.PutUint16(buf[pos:], uint16(len(fr.releasedPages)))
	pos += 2

	for _, pageNum := range fr.releasedPages {
		binary.LittleEndian.PutUint64(buf[pos:], uint64(pageNum))
		pos += pageNumSize
	}
}

func (fr *freelist) deserialize(buf []byte) {
	pos := 0

	//maxPage num
	fr.maxPage = pgnum(binary.LittleEndian.Uint16(buf[pos:]))
	pos += 2

	//releasePages num
	realeasedPagesNum := int(binary.LittleEndian.Uint16(buf[pos:]))
	pos += 2

	//todo, I believe I should first empty the fr.releasedPages array befr
	fr.releasedPages = make([]pgnum, 0)
	for i := 0; i < realeasedPagesNum; i++ {
		releasedPageNum := binary.LittleEndian.Uint64(buf[pos:])
		fr.releasedPages = append(fr.releasedPages, pgnum(releasedPageNum))
		pos += pageNumSize
	}
}
