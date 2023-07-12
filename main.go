package main

func main() {
	//initialize the db
	dal, _ := newDal("db.db")

	p := dal.allocateEmptyPage()
	p.num = dal.getNextPage()
	copy(p.data, "data")

	dal.writePage(p)
	dal.writeFreeList()

	dal.close()

	dal, _ = newDal("db.db")
	p = dal.allocateEmptyPage()
	p.num = dal.getNextPage()
	copy(p.data, "data2")

	dal.writePage(p)

	pageNum := dal.getNextPage()
	dal.releasePage(pageNum)

	dal.writeFreeList()
}
