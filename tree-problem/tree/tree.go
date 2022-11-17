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

func NewTreeConfig() *TreeConfig {
	config := new(TreeConfig)
	return config
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

	Command = "tree"
	Spaces3 = "	  "
	Spaces4 = Spaces3 + " "
)

func ParseCommand(cmd string) TreeConfig {
	regxCmpl := regexp.MustCompile(`\s+`)
	cmd = strings.TrimSpace(regxCmpl.ReplaceAllString(cmd, " "))

	ca := strings.Split(cmd, " ") //args in command

	if strings.TrimSpace(ca[0]) != Command {
		log.Fatalf("command not found: `%v`", ca[0])
	}

	config := NewTreeConfig()

	for i := 1; i < len(ca); i++ {
		v := ca[i]

		// check for path
		if !strings.HasPrefix(v, "-") && !regexp.MustCompile(`\d`).MatchString(v) {
			config.paths = append(config.paths, v)
			continue
		}

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
			log.Fatalf("Invalid argument `%v`", op)
		}
	}

	if len(config.paths) < 1 {
		config.paths = []string{"."}
	}
	return *config
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

func ListDirAndFiles(config TreeConfig) string {
	var fc FileCount
	temp := ""
	for _, p := range config.paths {
		fc = FileCount{}
		root := strings.TrimSuffix(p, "/")
		temp += recListDirAndFiles(root, root+"\n", 0, false, &fc, &config)
	}

	return fmt.Sprintf("%v \n %v directory, %v files", temp, fc.dirCnt, fc.fileCnt)
}

func recListDirAndFiles(root string, temp string, n int, isLastDir bool, fc *FileCount, config *TreeConfig) string {
	files := ReadDir(root)

	if len(files) < 1 || (config.level > 0 && n == config.level) {
		return temp
	}

	lastFile := files[len(files)-1]

	for _, val := range files {
		//ignoring file start with `.`
		if strings.HasPrefix(val.Name(), ".") {
			continue
		}

		bp := "" // before pipe
		pipe := BoxVH
		ap := " " + val.Name() //after pipe
		isLastFile := lastFile == val
		var relPath, fp string // fp: file permission

		if config.reqRelPath {
			relPath = " " + root + "/" + val.Name()
			ap = relPath
		}
		if config.reqFilePermsn {
			fp = " [" + getPermsnMode(val) + "] "
			ap = fp + val.Name()
		}
		if config.reqRelPath && config.reqFilePermsn {
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
			if config.reqOnlyDir {
				continue
			}
			temp += bp + pipe + ap + "\n"
			fc.fileCnt++
			continue
		}

		isLastDir = isLastFile && lastFile.IsDir()
		fc.dirCnt++
		temp += bp + pipe + ap + "\n"
		temp = recListDirAndFiles(root+"/"+val.Name(), temp, n+1, isLastDir, fc, config)
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

func help() {
	options := make(map[string]Option)
	options["d"] = Option{"-d", 0, "prints only directories"}
	options["f"] = Option{"-f", 0, "prints relative path of each file"}
	options["L"] = Option{"-L", 0, "travese specified nested levels only"}
	options["p"] = Option{"-p", 0, "prints permission along with file name"}
	fmt.Println("help: ", options)
}
