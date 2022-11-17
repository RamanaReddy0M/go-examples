package main

import (
	"fmt"
	"tree-problem/tree"
)

func main() {
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker"
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker-examples/micronaut-postgres/build/resources"
	//path := "/Users/ramana/java-workspace"
	path := "/Users/ramana/go-workspace/go-examples/"
	cmd := "tree -d " + path
	//tree.ParseCommand(cmd)
	//app: 18715 directory, 71516 files
	//tree: 1864 directories, 7416 files
	// f, _ := os.ReadDir(path)
	// fmt.Println(f[len(f)-1])

	fmt.Println(tree.ListDirAndFiles(cmd))
	// tree.FileWalk(path)
}
