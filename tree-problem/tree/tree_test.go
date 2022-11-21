package tree

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestListDirAndFiles(t *testing.T) {
	assert := assert.New(t)

	//empty dir test
	path := "../resources/test-dir/empty"
	cmd := "tree " + path

	got := ListDirAndFiles(ParseCommand(cmd))
	want := "../resources/test-dir/empty\n\n0 directories, 0 files"
	assert.Equal(want, got, "empty directory test")

	//directory with multiple files test
	path = "../resources/test-dir/"
	cmd = "tree " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── empty\n└── hello\n    │── hello.txt\n " + 
	       "   │── temp\n    │   └── temp.txt\n    └── xelo\n        └── lwlo.rx\n"+
		   "\n4 directories, 3 files"
	assert.Equal(want, got, "directory with multiple file test")

	//relative path directories test
	path = "../resources/test-dir/"
	cmd = "tree -f " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── ../resources/test-dir/empty\n└──"+
	       " ../resources/test-dir/hello\n    │── ../resources/test-dir/hello/hello.txt\n"+
		   "    │── ../resources/test-dir/hello/temp\n    │  "+
		   " └── ../resources/test-dir/hello/temp/temp.txt\n "+
		   "   └── ../resources/test-dir/hello/xelo\n        "+
		   "└── ../resources/test-dir/hello/xelo/lwlo.rx\n\n4 directories, 3 files"
	assert.Equal(want, got, "relative path directories test")

	//Level 5 directories test
	path = "../resources/level-test-dir"
	cmd = "tree -L 5 " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/level-test-dir\n│── META-INF\n│   └── empty\n└── in\n  "+
	       "  └── one2n\n        └── tree-prblm\n            └── test-dir\n             "+
	       "   │── empty\n                └── hello\n\n8 directories, 0 files"
	assert.Equal(want, got, "Level 5 directories test")

	//only directories upto 7 levels test
	path = "../resources/level-test-dir"
	cmd = "tree -L 7 -d " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/level-test-dir\n│── META-INF\n│   └── empty\n└── in\n"+
	       "    └── one2n\n        └── tree-prblm\n            └── test-dir\n "+
		   "               │── empty\n                └── hello\n             "+
		   "       │── temp\n                    └── xelo\n\n10 directories"
	assert.Equal(want, got, "only directories upto 7 levels test")
}
