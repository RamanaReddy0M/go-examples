package main

import (
	"fmt"
	"os"
	"strings"
	"tree-problem/tree"
)

func main() {
	cmd := "tree " + strings.Join(os.Args[1:], " ")
	fmt.Println(tree.ListDirAndFiles(tree.ParseCommand(cmd)))
}
