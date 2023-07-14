package main

import "errors"

const (
	nodeHeaderSize  = 3
	pageNumSize     = 8 //size of a page number in bytes
	collectionSize  = 16
	magicNumberSize = 4
	counterSize     = 4
)

var writeInsideReadTxErr = errors.New("can't perform a write operation inside a read transaction")
