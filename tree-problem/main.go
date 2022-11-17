package main

import (
	"fmt"
	"tree-problem/tree"
)

func main() {
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker"
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker-examples/micronaut-postgres/build/resources"
	//path := "/Users/ramana/java-workspace"


	path := ".."
	cmd := "tree -d -f -p -L 1 " + path
	fmt.Println(tree.ListDirAndFiles(cmd))
}
