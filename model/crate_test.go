package model

import (
	"testing"
)

func TestCrate_GetFsmGraph(t *testing.T) {
	crate := Crate{}
	t.Log(crate.getFsmGraph())
}
