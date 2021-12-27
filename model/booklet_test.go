package model

import (
	"testing"
)

func TestBooklet_GetFsmGraph(t *testing.T) {
	booklet := Booklet{}
	t.Log(booklet.GetFsmGraph())
}
