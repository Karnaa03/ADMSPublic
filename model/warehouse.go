package model

type Warehouse struct {
	Row        int   `pg:",unique:wr"`
	Shelf      int   `pg:",unique:wr"`
	ShelfLevel []int `pg:",array,unique:wr"`
}
