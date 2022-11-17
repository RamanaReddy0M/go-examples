package main

import (
	"fmt"
	"tree-problem/tree"
)

func main() {
	//path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker"
	path := "/Users/ramana/java-workspace/fast-dev/hot-reload-inside-docker-examples/micronaut-postgres/build/resources"
	//path := "/Users/ramana/java-workspace"

	//path := ""
	cmd := "tree " + path
	fmt.Println(tree.ListDirAndFiles(tree.ParseCommand(cmd)))
}
