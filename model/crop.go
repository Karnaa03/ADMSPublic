package model

type Crop struct {
	Id   int `pg:",pk,notnull,unique"`
	Name string
}
