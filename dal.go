package main

import (
	"errors"
	"fmt"
	"os"
)

//Data Access Layer
type dal struct {
	file     *os.File
	pageSize int

	*metaPage
	*freelist
}

type pgnum uint64

type page struct {
	num  pgnum
	data []byte
}

func newDal(path string) (*dal, error) {

	dal := &dal{
		metaPage: newEmptyMeta(),
		pageSize: os.Getpagesize() / 4,
	}

	if _, err := os.Stat(path); err == nil {
		//db file exists
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			dal.close()
			return nil, err
		}

		meta, err := dal.readMetaPage()
		if err != nil {
			dal.close()
			return nil, err
		}
		dal.metaPage = meta

		freelist, err := dal.readFreeList()
		if err != nil {
			dal.close()
			return nil, err
		}
		dal.freelist = freelist

	} else if errors.Is(err, os.ErrNotExist) {
		//db file doesn't exist
		dal.file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			dal.close()
			return nil, err
		}

		dal.freelist = newFreeList()
		dal.freelistPage = dal.getNextPage()
		_, err := dal.writeFreeList()
		if err != nil {
			return nil, err
		}
		dal.writeMetaPage(dal.metaPage)

	} else {
		return nil, err
	}

	return dal, nil
}

func (d *dal) close() error {
	if d.file != nil {
		err := d.file.Close()
		if err != nil {
			return fmt.Errorf("couldn't close the file: %s", err)
		}
		d.file = nil
	}
	return nil
}

func (d *dal) allocateEmptyPage() *page {
	return &page{
		data: make([]byte, d.pageSize),
	}
}

func (d *dal) readPage(pageNum pgnum) (*page, error) {
	p := d.allocateEmptyPage()

	offsetToRead := int64(pageNum) * int64(d.pageSize)

	_, err := d.file.ReadAt(p.data, int64(offsetToRead))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *dal) writePage(p *page) error {
	offsetToWrite := int64(p.num) * int64(d.pageSize)
	_, err := d.file.WriteAt(p.data, int64(offsetToWrite))
	return err
}

func (d *dal) writeMetaPage(meta *metaPage) (*page, error) {
	p := d.allocateEmptyPage()
	p.num = metaDataPageNum
	meta.serialize(p.data)

	err := d.writePage(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (d *dal) readMetaPage() (*metaPage, error) {
	p, err := d.readPage(metaDataPageNum)
	if err != nil {
		return nil, err
	}

	meta := newEmptyMeta()
	meta.deserialize(p.data)
	return meta, nil
}

func (d *dal) writeFreeList() (*page, error) {
	p := d.allocateEmptyPage()
	p.num = d.freelistPage
	d.freelist.serialize(p.data)

	err := d.writePage(p)
	if err != nil {
		return nil, err
	}
	d.freelistPage = p.num //set the freelist page now that p has been written
	return p, nil
}

func (d *dal) readFreeList() (*freelist, error) {
	p, err := d.readPage(d.freelistPage)
	if err != nil {
		return nil, err
	}

	freelist := newFreeList()
	freelist.deserialize(p.data)
	return freelist, nil
}
