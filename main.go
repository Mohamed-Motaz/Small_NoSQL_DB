package main

import (
	"os"
)

func main() {
	//initialize the db
	dal, _ := newDal("db.db", uint64(os.Getpagesize()))

	p := dal.allocateEmptyPage()
	p.num = dal.getNextPage()
	copy(p.data, "data")
	dal.writePage(p)
}
