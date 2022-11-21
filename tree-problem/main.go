package main

import (
	"fmt"
	"tree-problem/tree"
)

func main() {
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker"
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker-examples/micronaut-postgres/build/resources"
	//path := "/Users/ramana/java-workspace"

	path := "./resources/level-test-dir/"
	cmd := "tree -L 7 -d " + path
	fmt.Println(tree.ListDirAndFiles(tree.ParseCommand(cmd)))
	want := "../resources/level-test-dir\n│── META-INF\n│   └── empty\n"+
	       "└── in\n    └── one2n\n        └── tree-prblm\n            └── test-dir\n    "+
	       "            │── empty\n                └── hello\n\n8 directories, 0 files"
	fmt.Printf("%v", want)
}
