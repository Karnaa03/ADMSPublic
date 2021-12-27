package model

type Tree struct {
	Id   int `pg:",pk,notnull,unique"`
	Name string
}
