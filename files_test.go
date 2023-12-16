package TBys

import (
	"io/fs"
	"testing"
)

func Test_SearchFile(t *testing.T) {
	sf := TSearchFile{}
	sf.Suffix = ".js"
	sf.Search(`/home/icy/Projects/`, func(path string, info fs.FileInfo) bool {
		t.Log(info)
		return false

	})
}
