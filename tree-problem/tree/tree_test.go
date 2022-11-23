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

	//No paths in command test
	cmd = "tree"

	got = ListDirAndFiles(ParseCommand(cmd))
	want = ".\n│── tree.go\n└── tree_test.go\n\n0 directories, 2 files"
	assert.Equal(want, got, "No path in command test")

	//Single file in directory test
	path = "../resources/test-dir/hello/temp"
	cmd = "tree " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir/hello/temp\n└── temp.txt\n\n0 directories, 1 file"
	assert.Equal(want, got, "Single file in directory test")

	//directory with multiple files test
	path = "../resources/test-dir/"
	cmd = "tree " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── empty\n└── hello\n    │── hello.txt\n " +
		"   │── temp\n    │   └── temp.txt\n    └── xelo\n        └── lwlo.rx\n" +
		"\n4 directories, 3 files"
	assert.Equal(want, got, "directory with multiple file test")

	//relative path directories test
	path = "../resources/test-dir/"
	cmd = "tree -f " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── ../resources/test-dir/empty\n└──" +
		" ../resources/test-dir/hello\n    │── ../resources/test-dir/hello/hello.txt\n" +
		"    │── ../resources/test-dir/hello/temp\n    │  " +
		" └── ../resources/test-dir/hello/temp/temp.txt\n " +
		"   └── ../resources/test-dir/hello/xelo\n        " +
		"└── ../resources/test-dir/hello/xelo/lwlo.rx\n\n4 directories, 3 files"
	assert.Equal(want, got, "relative path directories test")

	//Files with permission mode test
	path = "../resources/test-dir/"
	cmd = "tree -p " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── [drwxr-xr-x] empty\n└── [drwxr-xr-x] hello\n" +
		"    │── [-rw-r--r--] hello.txt\n    │── [drwxr-xr-x] temp\n    │   └── [-rw-r--r--]" +
		" temp.txt\n    └── [drwxr-xr-x] xelo\n        └── [-rw-r--r--] lwlo.rx\n\n" +
		"4 directories, 3 files"
	assert.Equal(want, got, "Files with permission mode test")

	//Order files by Modified Time(-t) test
	path = "../resources/test-dir/"
	cmd = "tree -t " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── empty\n└── hello\n    │── hello.txt\n" +
		"    │── temp\n    │   └── temp.txt\n    └── xelo\n        └── lwlo.rx\n\n" +
		"4 directories, 3 files"
	assert.Equal(want, got, "Order files by Modified Time(-t) test")

	//List files with 'No-Indentation' test
	path = "../resources/test-dir/"
	cmd = "tree -i -f " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n ../resources/test-dir/empty\n ../resources/test-dir/hello\n" +
		" ../resources/test-dir/hello/hello.txt\n ../resources/test-dir/hello/temp\n" +
		" ../resources/test-dir/hello/temp/temp.txt\n ../resources/test-dir/hello/xelo\n" +
		" ../resources/test-dir/hello/xelo/lwlo.rx\n\n4 directories, 3 files"
	assert.Equal(want, got, "List files with 'No-Indentation' test")

	//permission mode and relative path test
	path = "../resources/test-dir/"
	cmd = "tree -f -p " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/test-dir\n│── [drwxr-xr-x]  ../resources/test-dir/empty\n└── [drwxr-xr-x]" +
		"  ../resources/test-dir/hello\n    │── [-rw-r--r--]  ../resources/test-dir/hello/hello.txt\n" +
		"    │── [drwxr-xr-x]  ../resources/test-dir/hello/temp\n    │   └── [-rw-r--r--] " +
		" ../resources/test-dir/hello/temp/temp.txt\n    └── [drwxr-xr-x]  " +
		"../resources/test-dir/hello/xelo\n        └── [-rw-r--r--]  " +
		"../resources/test-dir/hello/xelo/lwlo.rx\n\n4 directories, 3 files"
	assert.Equal(want, got, "permission mode and relative path test")

	//Level 5 directories test
	path = "../resources/level-test-dir"
	cmd = "tree -L 5 " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/level-test-dir\n│── META-INF\n│   └── empty\n└── in\n  " +
		"  └── one2n\n        └── tree-prblm\n            └── test-dir\n             " +
		"   │── empty\n                └── hello\n\n8 directories, 0 files"
	assert.Equal(want, got, "Level 5 directories test")

	//only directories upto 7 levels test
	path = "../resources/level-test-dir"
	cmd = "tree -L 7 -d " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/level-test-dir\n│── META-INF\n│   └── empty\n└── in\n" +
		"    └── one2n\n        └── tree-prblm\n            └── test-dir\n " +
		"               │── empty\n                └── hello\n             " +
		"       │── temp\n                    └── xelo\n\n10 directories"
	assert.Equal(want, got, "only directories upto 7 levels test")

	//Parsing command with odd spaces and mutiple args test
	path = "../resources/level-test-dir"
	cmd = "tree  -L            7   -d -t      -p " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/level-test-dir\n│── [drwxr-xr-x] in\n│   └── [drwxr-xr-x] one2n\n│" +
		"       └── [drwxr-xr-x] tree-prblm\n│           └── [drwxr-xr-x] test-dir\n│   " +
		"            │── [drwxr-xr-x] empty\n│               └── [drwxr-xr-x] hello\n│  " +
		"                 │── [drwxr-xr-x] temp\n│                   └── [drwxr-xr-x] xelo\n" +
		"└── [drwxr-xr-x] META-INF\n    └── [drwxr-xr-x] empty\n\n10 directories"
	assert.Equal(want, got, "Parsing command with odd spaces and mutiple args test")

	//List files in multiple paths test
	paths := "../resources/level-test-dir ../resources/test-dir"
	cmd = "tree " + paths

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "../resources/level-test-dir\n│── META-INF\n│   └── empty\n└── in\n    └── one2n\n        " +
		"└── tree-prblm\n            └── test-dir\n                │── empty\n                " +
		"└── hello\n                    │── hello.txt\n                    │── temp\n         " +
		"           │   └── temp.txt\n                    └── xelo\n                        " +
		"└── lwlo.rx\n" +
		"../resources/test-dir\n│── empty\n└── hello\n    │── hello.txt\n    │── temp\n    │   " +
		"└── temp.txt\n    └── xelo\n        └── lwlo.rx\n\n4 directories, 3 files"
	assert.Equal(want, got, "Parsing command with odd spaces and mutiple args test")

	//XML format empty dir test
	path = "../resources/test-dir/empty"
	cmd = "tree -X " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
		"<directory name=\"../resources/test-dir/empty>\n  </directory>\n  <report>\n   " +
		"<directories>0</directories>\n   <files>0</files>\n  </report>\n</tree>"
	assert.Equal(want, got, "XML format empty dir test")

	//XML format single file in directory test
	path = "../resources/test-dir/hello/temp"
	cmd = "tree -X " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
		"<directory name=\"../resources/test-dir/hello/temp>\n    <file name=\"temp.txt\"></file>\n " +
		" </directory>\n  <report>\n   <directories>0</directories>\n   <files>0</files>\n  " +
		"</report>\n</tree>"
	assert.Equal(want, got, "XML format single file in directory test")

	//XML format Level 5 directories test
	path = "../resources/level-test-dir"
	cmd = "tree -L 5 -X " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
		"<directory name=\"../resources/level-test-dir>\n    <directory name=\"META-INF\">\n     " +
		"<directory name=\"empty\">\n     </directory>\n    </directory>\n    <directory name=\"in\">\n" +
		"     <directory name=\"one2n\">\n      <directory name=\"tree-prblm\">\n       " +
		"<directory name=\"test-dir\">\n        <directory name=\"empty\">\n        </directory>\n" +
		"        <directory name=\"hello\">\n        </directory>\n       </directory>\n      " +
		"</directory>\n     </directory>\n    </directory>\n  </directory>\n  <report>\n   " +
		"<directories>8</directories>\n   <files>8</files>\n  </report>\n</tree>"
	assert.Equal(want, got, "XML format Level 5 directories test")

	//Files in XML format with permission mode test
	path = "../resources/test-dir/"
	cmd = "tree -p -X " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
		"<directory name=\"../resources/test-dir>\n    " +
		"<directory name=\"empty\" mode=\"0755\" prot=\"drwxr-xr-x\">\n    </directory>\n    " +
		"<directory name=\"hello\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     " +
		"<file name=\"hello.txt\" mode=\"0644\" prot=\"-rw-r--r--\"></file>\n     " +
		"<directory name=\"temp\" mode=\"0755\" prot=\"drwxr-xr-x\">\n      " +
		"<file name=\"temp.txt\" mode=\"0644\" prot=\"-rw-r--r--\"></file>\n     " +
		"</directory>\n     <directory name=\"xelo\" mode=\"0755\" prot=\"drwxr-xr-x\">\n      " +
		"<file name=\"lwlo.rx\" mode=\"0644\" prot=\"-rw-r--r--\"></file>\n     </directory>\n" +
		"    </directory>\n  </directory>\n  <report>\n   <directories>4</directories>\n   " +
		"<files>4</files>\n  </report>\n</tree>"
	assert.Equal(want, got, "Files in XML format with permission mode test")

	//XML format only directories and permission mode test
	path = "../resources/test-dir/"
	cmd = "tree -p -X -d " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
		"<directory name=\"../resources/test-dir>\n    " +
		"<directory name=\"empty\" mode=\"0755\" prot=\"drwxr-xr-x\">\n    " +
		"</directory>\n    <directory name=\"hello\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     " +
		"<directory name=\"temp\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     </directory>\n     " +
		"<directory name=\"xelo\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     </directory>\n    " +
		"</directory>\n  </directory>\n  <report>\n   <directories>4</directories>\n  </report>\n</tree>"
	assert.Equal(want, got, "XML format only directories and permission mode test")

	//JSON format empty dir test
	path = "../resources/test-dir/empty"
	cmd = "tree -J " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "[\n  {\"type\":\"directory\",\"name\":\"../resources/test-dir/empty\",\"contents\":[\n" +
		"  ]}\n,\n  {\"type\":\"report\",\"directories\":0,\"files\":0}\n]"
	assert.Equal(want, got, "JSON format empty dir test")

	//JSON format single file in directory test
	path = "../resources/test-dir/hello/temp"
	cmd = "tree -J " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "[\n  {\"type\":\"directory\",\"name\":\"../resources/test-dir/hello/temp\",\"contents\":[\n" +
		"    {\"type\":\"file\",\"name\":\"temp.txt\"}\n  ]}\n,\n  " +
		"{\"type\":\"report\",\"directories\":0,\"files\":1}\n]"
	assert.Equal(want, got, "JSON format single file in directory test")

	//Files in JSON format with permission mode test
	path = "../resources/test-dir/"
	cmd = "tree -p -J " + path

	got = ListDirAndFiles(ParseCommand(cmd))
	want = "[\n  {\"type\":\"directory\",\"name\":\"../resources/test-dir\",\"contents\":[\n" +
		"    {\"type\":\"directory\",\"name\":\"empty\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
		"    ]}\n    {\"type\":\"directory\",\"name\":\"hello\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
		"     {\"type\":\"file\",\"name\":\"hello.txt\",\"mode\":\"0644\",\"prot\":\"-rw-r--r--\"},\n" +
		"     {\"type\":\"directory\",\"name\":\"temp\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
		"      {\"type\":\"file\",\"name\":\"temp.txt\",\"mode\":\"0644\",\"prot\":\"-rw-r--r--\"}\n" +
		"     ]}\n     {\"type\":\"directory\",\"name\":\"xelo\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
		"      {\"type\":\"file\",\"name\":\"lwlo.rx\",\"mode\":\"0644\",\"prot\":\"-rw-r--r--\"}\n" +
		"     ]}\n    ]},\n  ]}\n,\n  {\"type\":\"report\",\"directories\":4,\"files\":3}\n]"
	assert.Equal(want, got, "Files in JSON format with permission mode test")

}
