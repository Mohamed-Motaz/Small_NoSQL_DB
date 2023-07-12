package main

import (
	"fmt"
	"os"
)

//Data Access Layer
type dal struct {
	file     *os.File
	pageSize uint64
}

type pgnum uint64

type page struct {
	num  pgnum
	data []byte
}

func newDal(path string, pageSize uint64) (*dal, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	dal := &dal{
		file,
		pageSize,
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

	offsetToRead := uint64(pageNum) * d.pageSize

	_, err := d.file.ReadAt(p.data, int64(offsetToRead))
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (d *dal) writePage(p *page) error {
	offsetToWrite := uint64(p.num) * d.pageSize
	_, err := d.file.WriteAt(p.data, int64(offsetToWrite))
	return err
}