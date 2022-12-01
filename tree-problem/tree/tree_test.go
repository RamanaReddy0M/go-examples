package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	cmd  string
	desc string
	want string
}

func TestListDirAndFiles(t *testing.T) {

	tests := []test{
		{cmd: "tree ../resources/test-dir/empty", desc: "empty dir test",
			want: "../resources/test-dir/empty\n" +
				  "\n0 directories, 0 files"},
		{cmd: "tree ", desc: "No paths in command test",
			want: ".\n"+
			      "│── tree.go\n"+
				  "└── tree_test.go\n\n"+
				  "0 directories, 2 files"},
		{cmd: "tree ../resources/test-dir/hello/temp", desc: "Single file in directory test",
			want: "../resources/test-dir/hello/temp\n"+
			      "└── temp.txt\n\n"+
				  "0 directories, 1 file"},
		{cmd: "tree ../resources/test-dir/", desc: "directory with multiple files test",
			want: "../resources/test-dir\n"+
			      "│── empty\n"+
				  "└── hello\n"+
				  "    │── hello.txt\n" +
				  "    │── temp\n"+
				  "    │   └── temp.txt\n"+
				  "    └── xelo\n"+
				  "        └── lwlo.rx\n" +
			      "\n4 directories, 3 files"},
		{cmd: "tree -f ../resources/test-dir/", desc: "relative path directories test",
			want: "../resources/test-dir\n"+
			      "│── ../resources/test-dir/empty\n"+
				  "└── ../resources/test-dir/hello\n"+
				  "    │── ../resources/test-dir/hello/hello.txt\n" +
				  "    │── ../resources/test-dir/hello/temp\n"+
				  "    │   └── ../resources/test-dir/hello/temp/temp.txt\n"+
				  "    └── ../resources/test-dir/hello/xelo\n"+
				  "        └── ../resources/test-dir/hello/xelo/lwlo.rx\n\n"+
				  "4 directories, 3 files"},
		{cmd: "tree -p ../resources/test-dir/", desc: "Files with permission mode test",
			want: "../resources/test-dir\n"+
			      "│── [drwxr-xr-x] empty\n"+
				  "└── [drwxr-xr-x] hello\n" +
				  "    │── [-rw-r--r--] hello.txt\n"+
				  "    │── [drwxr-xr-x] temp\n"+
				  "    │   └── [-rw-r--r--] temp.txt\n"+
				  "    └── [drwxr-xr-x] xelo\n"+
				  "        └── [-rw-r--r--] lwlo.rx\n\n" +
				 "4 directories, 3 files"},
		{cmd: "tree -t ../resources/test-dir/", desc: "Order files by Modified Time(-t) test",
			want: "../resources/test-dir\n"+
			      "│── empty\n"+
				  "└── hello\n"+
				  "    │── hello.txt\n" +
				  "    │── temp\n"+
				  "    │   └── temp.txt\n"+
				  "    └── xelo\n"+
				  "        └── lwlo.rx\n\n" +
			      "4 directories, 3 files"},
		{cmd: "tree -i -f ../resources/test-dir/", desc: "List files with 'No-Indentation' test",
			want: "../resources/test-dir\n"+
			      " ../resources/test-dir/empty\n"+
				  " ../resources/test-dir/hello\n" +
				  " ../resources/test-dir/hello/hello.txt\n"+
				  " ../resources/test-dir/hello/temp\n" +
				  " ../resources/test-dir/hello/temp/temp.txt\n"+
				  " ../resources/test-dir/hello/xelo\n" +
				  " ../resources/test-dir/hello/xelo/lwlo.rx\n\n"+
				  "4 directories, 3 files"},
		{cmd: "tree -p -f ../resources/test-dir/", desc: "permission mode and relative path test",
			want: "../resources/test-dir\n"+
			      "│── [drwxr-xr-x]  ../resources/test-dir/empty\n"+
				  "└── [drwxr-xr-x]  ../resources/test-dir/hello\n"+
				  "    │── [-rw-r--r--]  ../resources/test-dir/hello/hello.txt\n" +
				  "    │── [drwxr-xr-x]  ../resources/test-dir/hello/temp\n"+
				  "    │   └── [-rw-r--r--]  ../resources/test-dir/hello/temp/temp.txt\n"+
				  "    └── [drwxr-xr-x]  ../resources/test-dir/hello/xelo\n"+
				  "        └── [-rw-r--r--]  ../resources/test-dir/hello/xelo/lwlo.rx\n\n"+
				  "4 directories, 3 files"},
		{cmd: "tree -L 5 ../resources/level-test-dir", desc: "Level 5 directories test",
			want: "../resources/level-test-dir\n"+
			      "│── META-INF\n"+
				  "│   └── empty\n"+
				  "└── in\n"+
				  "    └── one2n\n"+
				  "        └── tree-prblm\n"+
				  "            └── test-dir\n"+
				  "                │── empty\n"+
				  "                └── hello\n\n"+
				  "8 directories, 0 files"},
		{cmd: "tree -L 7 -d ../resources/level-test-dir", desc: "only directories upto 7 levels test",
			want: "../resources/level-test-dir\n"+
			      "│── META-INF\n"+
				  "│   └── empty\n"+
				  "└── in\n" +
				  "    └── one2n\n"+
				  "        └── tree-prblm\n"+
				  "            └── test-dir\n"+
				  "                │── empty\n"+
				  "                └── hello\n"+
				  "                    │── temp\n"+
				  "                    └── xelo\n\n"+
				  "10 directories"},
		{cmd: "tree  -L            7   -d -t      -p ../resources/level-test-dir", desc: "Parsing command with odd spaces and mutiple args test",
			want: "../resources/level-test-dir\n"+
			"│── [drwxr-xr-x] in\n"+
			"│   └── [drwxr-xr-x] one2n\n"+
			"│       └── [drwxr-xr-x] tree-prblm\n"+
			"│           └── [drwxr-xr-x] test-dir\n"+
			"│               │── [drwxr-xr-x] empty\n"+
			"│               └── [drwxr-xr-x] hello\n"+
			"│                   │── [drwxr-xr-x] temp\n"+
			"│                   └── [drwxr-xr-x] xelo\n" +
			"└── [drwxr-xr-x] META-INF\n"+
			"    └── [drwxr-xr-x] empty\n\n"+
			"10 directories"},
		{cmd: "tree ../resources/level-test-dir ../resources/test-dir", desc: "Parsing command with odd spaces and mutiple args test",
			want: "../resources/level-test-dir\n"+
			      "│── META-INF\n"+
				  "│   └── empty\n"+
				  "└── in\n"+
				  "    └── one2n\n"+
				  "        └── tree-prblm\n"+
				  "            └── test-dir\n"+
				  "                │── empty\n"+
				  "                └── hello\n"+
				  "                    │── hello.txt\n"+
				  "                    │── temp\n"+
				  "                    │   └── temp.txt\n"+
				  "                    └── xelo\n"+
				  "                        └── lwlo.rx\n" +
				  "../resources/test-dir\n"+
				  "│── empty\n"+
				  "└── hello\n"+
				  "    │── hello.txt\n"+
				  "    │── temp\n"+
				  "    │   └── temp.txt\n"+
				  "    └── xelo\n"+
				  "        └── lwlo.rx\n\n"+
				  "4 directories, 3 files"},
		{cmd: "tree -X ../resources/test-dir/empty", desc: "XML format empty dir test",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
				"<directory name=\"../resources/test-dir/empty>\n  </directory>\n  <report>\n   " +
				"<directories>0</directories>\n   <files>0</files>\n  </report>\n</tree>"},
		{cmd: "tree -X ../resources/test-dir/hello/temp", desc: "XML format single file in directory test",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
				"<directory name=\"../resources/test-dir/hello/temp>\n    <file name=\"temp.txt\"></file>\n " +
				" </directory>\n  <report>\n   <directories>0</directories>\n   <files>0</files>\n  " +
				"</report>\n</tree>"},
		{cmd: "tree -X -L 5 ../resources/level-test-dir", desc: "XML format Level 5 directories test",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
				"<directory name=\"../resources/level-test-dir>\n    <directory name=\"META-INF\">\n     " +
				"<directory name=\"empty\">\n     </directory>\n    </directory>\n    <directory name=\"in\">\n" +
				"     <directory name=\"one2n\">\n      <directory name=\"tree-prblm\">\n       " +
				"<directory name=\"test-dir\">\n        <directory name=\"empty\">\n        </directory>\n" +
				"        <directory name=\"hello\">\n        </directory>\n       </directory>\n      " +
				"</directory>\n     </directory>\n    </directory>\n  </directory>\n  <report>\n   " +
				"<directories>8</directories>\n   <files>8</files>\n  </report>\n</tree>"},
		{cmd: "tree -p -X ../resources/test-dir/", desc: "Files in XML format with permission mode test",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
				"<directory name=\"../resources/test-dir>\n    " +
				"<directory name=\"empty\" mode=\"0755\" prot=\"drwxr-xr-x\">\n    </directory>\n    " +
				"<directory name=\"hello\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     " +
				"<file name=\"hello.txt\" mode=\"0644\" prot=\"-rw-r--r--\"></file>\n     " +
				"<directory name=\"temp\" mode=\"0755\" prot=\"drwxr-xr-x\">\n      " +
				"<file name=\"temp.txt\" mode=\"0644\" prot=\"-rw-r--r--\"></file>\n     " +
				"</directory>\n     <directory name=\"xelo\" mode=\"0755\" prot=\"drwxr-xr-x\">\n      " +
				"<file name=\"lwlo.rx\" mode=\"0644\" prot=\"-rw-r--r--\"></file>\n     </directory>\n" +
				"    </directory>\n  </directory>\n  <report>\n   <directories>4</directories>\n   " +
				"<files>4</files>\n  </report>\n</tree>"},
		{cmd: "tree -X -p -d ../resources/test-dir/", desc: "XML format only directories and permission mode test",
			want: "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<tree>\n  " +
				"<directory name=\"../resources/test-dir>\n    " +
				"<directory name=\"empty\" mode=\"0755\" prot=\"drwxr-xr-x\">\n    " +
				"</directory>\n    <directory name=\"hello\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     " +
				"<directory name=\"temp\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     </directory>\n     " +
				"<directory name=\"xelo\" mode=\"0755\" prot=\"drwxr-xr-x\">\n     </directory>\n    " +
				"</directory>\n  </directory>\n  <report>\n   <directories>4</directories>\n  </report>\n</tree>"},
		{cmd: "tree -J ../resources/test-dir/empty", desc: "JSON format empty dir test",
			want: "[\n  {\"type\":\"directory\",\"name\":\"../resources/test-dir/empty\",\"contents\":[\n" +
				"  ]}\n,\n  {\"type\":\"report\",\"directories\":0,\"files\":0}\n]"},
		{cmd: "tree -J ../resources/test-dir/hello/temp", desc: "JSON format single file in directory test",
			want: "[\n  {\"type\":\"directory\",\"name\":\"../resources/test-dir/hello/temp\",\"contents\":[\n" +
				"    {\"type\":\"file\",\"name\":\"temp.txt\"}\n  ]}\n,\n  " +
				"{\"type\":\"report\",\"directories\":0,\"files\":1}\n]"},
		{cmd: "tree -J -p ../resources/test-dir/", desc: "Files in JSON format with permission mode test",
			want: "[\n  {\"type\":\"directory\",\"name\":\"../resources/test-dir\",\"contents\":[\n" +
				"    {\"type\":\"directory\",\"name\":\"empty\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
				"    ]}\n    {\"type\":\"directory\",\"name\":\"hello\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
				"     {\"type\":\"file\",\"name\":\"hello.txt\",\"mode\":\"0644\",\"prot\":\"-rw-r--r--\"},\n" +
				"     {\"type\":\"directory\",\"name\":\"temp\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
				"      {\"type\":\"file\",\"name\":\"temp.txt\",\"mode\":\"0644\",\"prot\":\"-rw-r--r--\"}\n" +
				"     ]}\n     {\"type\":\"directory\",\"name\":\"xelo\",\"mode\":\"0755\",\"prot\":\"drwxr-xr-x\",\"contents\":[\n" +
				"      {\"type\":\"file\",\"name\":\"lwlo.rx\",\"mode\":\"0644\",\"prot\":\"-rw-r--r--\"}\n" +
				"     ]}\n    ]},\n  ]}\n,\n  {\"type\":\"report\",\"directories\":4,\"files\":3}\n]"},
	}

	assert := assert.New(t)
	for _, t := range tests {
		got := ListDirAndFiles(ParseCommand(t.cmd))
		assert.Equal(t.want, got, t.desc)
	}
}
