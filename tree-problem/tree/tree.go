package tree

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
)

type Option struct {
	option      string
	value       int
	description string
}

const (
	//Box Drawing Characters
	BoxVer       = "│"
	BoxHor       = "──"
	BoxVH        = BoxVer + BoxHor
	BoxDowAndRig = "┌"
	BoxDowAndLef = "┐"
	BoxUpAndRig  = "└"
	BoxUpAndLef  = "┘"
)

type FileCount struct {
	dirCnt, fileCnt int
}

func ParseCommand(cmd string) map[string]Option {
	options := make(map[string]Option)
	ca := strings.Split(cmd, " ")
	if len(ca) < 2 {
		return options
	}
	for i, op := range ca[1:] {
		switch op {
		case "-d":
			options["d"] = Option{"-d", 0, "prints only directories"}
		case "-f":
			options["f"] = Option{"-f", 0, "prints relative path of each file"}
		case "-L":
			if len(ca) < i+1 {
				log.Fatal("-L option requires value")
			}
			lVal := parseToInt(ca[i+1])
			if lVal < 1 {
				log.Fatal("-L value greater than 0")
			}
			options["L"] = Option{"-L", lVal, "travese specified nested levels only"}
		case "-p":
			options["p"] = Option{"-p", 0, "prints permission along with file name"}
		default:
			//log.Fatalf("Invalid argument `%v`", op)
		}
	}
	return options
}

func parseToInt(input string) int {
	num, err := strconv.ParseInt(input, 10, 32)
	//check for error
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	return int(num)
}

func ListDirAndFiles(cmd string) string {
	fc := FileCount{dirCnt: 0, fileCnt: 0}
	args := strings.Split(cmd, " ")
	root := args[len(args)-1]

	ops := ParseCommand(cmd)

	empOpt := Option{}
	relPath := ops["f"] != empOpt
	permsn := ops["p"] != empOpt
	reqOnlyDir := ops["d"] != empOpt
	level := 0
	if ops["L"] != empOpt {
		level = ops["L"].value
		fmt.Println("level: ", level)
	}

	temp := recListDirAndFiles(root, root+"\n", 0, &fc, false, relPath, permsn, reqOnlyDir, level)

	return fmt.Sprintf("%v \n %v directory, %v files", temp, fc.dirCnt, fc.fileCnt)
}

func recListDirAndFiles(root string, temp string, n int, fc *FileCount, isLastDir bool, reqRelPath bool, permsn bool, reqOnlyDir bool, level int) string {
	files := ReadDir(root)

	if len(files) < 1 || (level > 0 && n == level) {
		return temp
	}

	lastFile := files[len(files)-1]

	for _, val := range files {
		if val.Name() == ".git" {
			continue
		}

		bp := "" // before pipe
		pipe := BoxVH
		ap := " " + val.Name() //after pipe
		isLastFile := lastFile == val
		var relPath, fp string // fp: file permission

		if reqRelPath {
			relPath = " " + root + "/" + val.Name()
			ap = relPath
		}
		if permsn {
			fp = " [" + getPermsnMode(val) + "] "
			ap = fp + val.Name()
		}
		if reqRelPath && permsn {
			ap = fp + relPath
		}

		if isLastFile {
			pipe = BoxUpAndRig + BoxHor
		}

		if n > 0 {
			bp = BoxVer + "   " + strings.Repeat("    ", n-1)
		}

		if n == 1 && isLastDir {
			bp = strings.Repeat("    ", n)
		}

		if n > 1 && !isLastDir {
			bp = BoxVer + "   " + strings.Repeat("    ", n-2) + BoxVer + "   "
		}

		if n > 1 && isLastDir {
			bp = BoxVer + "   " + strings.Repeat("    ", n-1)
		}

		if !val.IsDir() { // file
			if reqOnlyDir {
				continue
			}
			temp += bp + pipe + ap + "\n"
			fc.fileCnt++
			continue
		}

		isLastDir = isLastFile && lastFile.IsDir()
		fc.dirCnt++
		temp += bp + pipe + ap + "\n"
		temp = recListDirAndFiles(root+"/"+val.Name(), temp, n+1, fc, isLastDir, reqRelPath, permsn, reqOnlyDir, level)
	}

	return temp
}

func ReadDir(root string) []fs.DirEntry {
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return make([]fs.DirEntry, 0)
	}
	return files
}

func getPermsnMode(f fs.DirEntry) string {
	fi, err := f.Info()
	if err != nil {
		fmt.Println(err)
		return "Unable read mode"
	}
	return fi.Mode().String()
}

func OrderDirAndFiles(files []fs.DirEntry) []fs.DirEntry {
	if len(files) < 1 {
		return []fs.DirEntry{}
	}

	df := make(map[string][]fs.DirEntry)
	df["directories"] = []fs.DirEntry{}
	df["files"] = []fs.DirEntry{}

	for _, f := range files {
		if f.IsDir() {
			df["directories"] = append(df["directories"], f)
			continue
		}
		df["files"] = append(df["files"], f)
	}
	return append(df["directories"], df["files"]...)
}
