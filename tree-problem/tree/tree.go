package tree

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Option struct {
	option      string
	value       int
	description string
}

type FileCount struct {
	dirCnt, fileCnt int
}

type TreeConfig struct {
	reqRelPath, reqOnlyDir, reqFilePermsn bool
	level                                 int
	paths                                 []string
}

const (
	//Box Drawing Characters
	BoxVer        = "│"
	BoxHor        = "──"
	BoxVH         = BoxVer + BoxHor
	BoxDowAndRig  = "┌"
	BoxDowAndLef  = "┐"
	BoxUpAndRig   = "└"
	BoxUpAndLef   = "┘"

	OpenBrkt      = "["
	CloseBrkt     = "]"
	Command       = "tree"
	PathSeperator = string(os.PathSeparator)
	NewLine       = "\n"
	Space         = " "
	Spaces3       = "   "
	Spaces4       = "    "
)

func NewTreeConfig() *TreeConfig {
	config := new(TreeConfig)
	return config
}

func ParseCommand(cmd string) TreeConfig {
	regxCmpl := regexp.MustCompile(`\s+`)
	cmd = strings.TrimSpace(regxCmpl.ReplaceAllString(cmd, " "))

	ca := strings.Split(cmd, Space) //args in command

	if strings.TrimSpace(ca[0]) != Command {
		log.Fatalf("command not found: `%v`", ca[0])
	}

	config := NewTreeConfig()

	for i := 1; i < len(ca); i++ {
		v := ca[i]

		op := v //op: option
		switch op {
		case "-d":
			config.reqOnlyDir = true
		case "-f":
			config.reqRelPath = true
		case "-L":
			if len(ca) < i+1 {
				log.Fatal("-L option requires value")
			}
			lVal := parseToInt(ca[i+1])
			if lVal < 1 {
				log.Fatal("-L value greater than 0")
			}
			config.level = lVal
			i++
		case "-p":
			config.reqFilePermsn = true
		default:
			// check for path
			if !strings.HasPrefix(v, "-") {
				config.paths = append(config.paths, v)
				continue
			}
			log.Fatalf("Invalid argument `%v`", op)
		}
	}

	if len(config.paths) < 1 {
		config.paths = []string{"."}
	}
	return *config
}

func ListDirAndFiles(config TreeConfig) string {
	var fc FileCount
	var isNthDirLast []bool
	temp := ""
	
	for _, p := range config.paths {
		fc = FileCount{}
		isNthDirLast = []bool{}
		root := strings.TrimSuffix(p, PathSeperator)
		temp += recListDirAndFiles(root, root+NewLine, 0, &isNthDirLast, &fc, &config)
	}
	return fmt.Sprintf("%v %v %v directory, %v files", temp, NewLine, fc.dirCnt, fc.fileCnt)
}

func ReadDir(root string) []fs.DirEntry {
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return make([]fs.DirEntry, 0)
	}
	return files
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

func recListDirAndFiles(root string, temp string, n int, isNthDirLast *[]bool, fc *FileCount, config *TreeConfig) string {
	files := ReadDir(root)

	if len(files) < 1 || (config.level > 0 && n == config.level) {
		return temp
	}

	lastFile := files[len(files)-1]

	for _, f := range files {
		//ignoring file start with `.`
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		bp := "" // before pipe
		pipe := BoxVH
		ap := Space + f.Name() //after pipe
		isLastFile := lastFile == f

		ap = getAfterPipeVal(root, f, *config, ap)

		if isLastFile {
			pipe = BoxUpAndRig + BoxHor
		}

		bp = getBeforePipeVal(n, *isNthDirLast)

		if !f.IsDir() { // file
			if config.reqOnlyDir {
				continue
			}
			temp += bp + pipe + ap + NewLine
			fc.fileCnt++
			continue
		}
		//directory
		temp += bp + pipe + ap + NewLine
		fc.dirCnt++
		//tracking information(whether directory last or not) from 0 to Nth level directory
		*isNthDirLast = append(*isNthDirLast, isLastFile)
		temp = recListDirAndFiles(root+PathSeperator+f.Name(), temp, n+1, isNthDirLast, fc, config)
	}

	if n > 0 && len(*isNthDirLast) > 0 {
		*isNthDirLast = (*isNthDirLast)[:n-1]
	}
	return temp
}

func getAfterPipeVal(root string, fi fs.DirEntry, config TreeConfig, ap string) string {
	var relPath, fp string // fp: file permission
	if config.reqRelPath {
		relPath = Space + root + PathSeperator + fi.Name()
		ap = relPath
	}

	if config.reqFilePermsn {
		fp = Space + OpenBrkt + getPermsnMode(fi) + CloseBrkt + Space
		ap = fp + fi.Name()
	}

	if config.reqRelPath && config.reqFilePermsn {
		ap = fp + relPath
	}
	return ap
}

func getPermsnMode(f fs.DirEntry) string {
	fi, err := f.Info()
	if err != nil {
		fmt.Println(err)
		return "Unable read mode"
	}
	return fi.Mode().String()
}

func getBeforePipeVal(n int, isNthDirLast []bool) string {
	bp := ""

	if n < 1 || n > len(isNthDirLast) {
		return bp
	}

	for i := 0; i < n; i++ {
		if !isNthDirLast[i] {
			bp += BoxVer + Spaces3
			continue
		}
		bp += Spaces4
	}
	return bp
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

func help() {
	options := make(map[string]Option)
	options["d"] = Option{"-d", 0, "prints only directories"}
	options["f"] = Option{"-f", 0, "prints relative path of each file"}
	options["L"] = Option{"-L", 0, "travese specified nested levels only"}
	options["p"] = Option{"-p", 0, "prints permission along with file name"}
	fmt.Println("help: ", options)
}