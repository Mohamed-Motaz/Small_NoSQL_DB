package main

//Data structure for keeping track of all free
//and empty pages
type freelist struct {
	maxPage       pgnum   //holds the maxiumu page allocated. maxPage * pageSize = fileSize
	releasedPages []pgnum //pages that were allocated but now free
}

const initialPage = 0

func newFreeList() *freelist {
	return &freelist{
		maxPage:       initialPage,
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
