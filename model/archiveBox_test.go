package model

import (
	"testing"
)

func TestArchiveBox_GetFsmGraph(t *testing.T) {
	archiveBox := ArchiveBox{}
	t.Log(archiveBox.GetFsmGraph())
}
