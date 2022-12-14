package tree

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type FileCount struct {
	dirCnt, fileCnt int
}

type TreeConfig struct {
	reqRelPath, reqOnlyDir, reqFilePermsn, sortByModTime, noIndent, reqXmlFormat, reqJsonFormat bool
	level                                                                                       int
	paths                                                                                       []string
}

const (
	//Box Drawing Characters
	BoxVer      = "│"
	BoxHor      = "──"
	BoxVH       = BoxVer + BoxHor
	BoxUpAndRig = "└"

	OpenBrkt      = "["
	CloseBrkt     = "]"
	OpenTag       = "<"
	Slash         = "/"
	CloseTag      = ">"
	JSONArrEnd    = "]}"
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
	//remove extra spaces in cmd
	regxCmpl := regexp.MustCompile(`\s+`)
	cmd = strings.TrimSpace(regxCmpl.ReplaceAllString(cmd, Space))
	ca := strings.Split(cmd, Space) //args in command
	if strings.TrimSpace(ca[0]) != Command {
		log.Fatalf("command not found: `%v`", ca[0])
	}

	config := NewTreeConfig()
	for i := 1; i < len(ca); i++ {
		arg := ca[i] //op: option
		switch arg {
		case "-d":
			config.reqOnlyDir = true
		case "-f":
			config.reqRelPath = true
		case "-i":
			config.noIndent = true
		case "-J":
			config.reqJsonFormat = true
			config.reqXmlFormat = false
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
		case "-t":
			config.sortByModTime = true
		case "-X":
			config.reqXmlFormat = true
			config.reqJsonFormat = false
		default:
			// check for path
			if !strings.HasPrefix(arg, "-") {
				config.paths = append(config.paths, arg)
				continue
			}
			log.Fatalf("Invalid argument `%v`", arg)
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
		if config.reqXmlFormat {
			temp += recListDirAndFilesInXML(root, temp, 0, &fc, &config)
			continue
		}

		if config.reqJsonFormat {
			temp += recListDirAndFilesInJSON(root, temp, 0, &fc, &config)
			continue
		}
		temp += recListDirAndFiles(root, root+NewLine, 0, &isNthDirLast, &fc, &config)
	}
	return formatRes(temp, fc, config)
}

func recListDirAndFiles(root string, temp string, n int, isNthDirLast *[]bool, fc *FileCount, config *TreeConfig) string {
	files := GetFiles(root, *config)
	if len(files) < 1 || (config.level > 0 && n == config.level) {
		*isNthDirLast = resizeToNMinus1(n, *isNthDirLast)
		return temp
	}

	lastFile := files[len(files)-1]
	for _, f := range files {
		bp := getBeforePipeVal(n, *isNthDirLast, *config) // before pipe
		isLastFile := lastFile == f
		pipe := getPipeVal(isLastFile, *config) // pipe (│── or └──)
		ap := getAfterPipeVal(root, f, *config) // after pipe
		temp += bp + pipe + ap + NewLine        //line structure in tree

		if !f.IsDir() { // file
			fc.fileCnt++
			continue
		}
		fc.dirCnt++
		//tracking information(whether directory last or not) from 0 to Nth level directory
		*isNthDirLast = append(*isNthDirLast, isLastFile)
		temp = recListDirAndFiles(root+PathSeperator+f.Name(), temp, n+1, isNthDirLast, fc, config)
	}
	*isNthDirLast = resizeToNMinus1(n, *isNthDirLast)
	return temp
}

func recListDirAndFilesInXML(root string, temp string, n int, fc *FileCount, config *TreeConfig) string {
	files := GetFiles(root, *config)

	if n == 0 {
		temp += strings.Repeat(Space, n+2) + OpenTag + "directory name=\"" + root + CloseTag + NewLine
	}

	closeDirTag := OpenTag + Slash + "directory" + CloseTag + NewLine
	if n > 0 && n == config.level {
		return temp + strings.Repeat(Space, n+3) + closeDirTag
	}

	for _, f := range files {
		if !f.IsDir() { // file
			temp += strings.Repeat(Space, n+4) + OpenTag + "file" + getFileAttrsVal(f, *config) + CloseTag +
				OpenTag + Slash + "file" + CloseTag + NewLine
			fc.fileCnt++
			continue
		}
		temp += strings.Repeat(Space, n+4) + OpenTag + "directory" + getFileAttrsVal(f, *config) + CloseTag + NewLine
		fc.dirCnt++
		temp = recListDirAndFilesInXML(root+PathSeperator+f.Name(), temp, n+1, fc, config)
	}

	if n > 0 {
		return temp + strings.Repeat(Space, n+3) + closeDirTag
	}
	return temp + strings.Repeat(Space, n+2) + closeDirTag
}

func recListDirAndFilesInJSON(root string, temp string, n int, fc *FileCount, config *TreeConfig) string {
	files := GetFiles(root, *config)

	if n == 0 {
		temp += strings.Repeat(Space, n+2) + "{\"type\":\"directory\",\"name\":\"" + root + "\",\"contents\":[" + NewLine
	}

	if n > 0 && n == config.level {
		return temp + strings.Repeat(Space, n+3) + JSONArrEnd + NewLine
	}

	var lastFile fs.DirEntry
	if len(files) > 0 {
		lastFile = files[len(files)-1]
	}

	for _, f := range files {
		isLastFile := f == lastFile
		if !f.IsDir() { // file
			temp += strings.Repeat(Space, n+4) + "{\"type\":\"file\"" + getFileAttrsVal(f, *config) + "}"
			if !isLastFile {
				temp += ","
			}
			temp += NewLine
			fc.fileCnt++
			continue
		}
		temp += strings.Repeat(Space, n+4) + "{\"type\":\"directory\"" + getFileAttrsVal(f, *config) + ",\"contents\":[" + NewLine
		fc.dirCnt++
		temp = recListDirAndFilesInJSON(root+PathSeperator+f.Name(), temp, n+1, fc, config)
	}

	if n > 0 {
		temp += strings.Repeat(Space, n+3) + JSONArrEnd
		if len(files) > 1 {
			temp += ","
		}
		return temp + NewLine
	}
	return temp + strings.Repeat(Space, n+2) + JSONArrEnd + NewLine
}

func GetFiles(root string, config TreeConfig) []fs.DirEntry {
	files := IgnoreDotFiles(ReadDir(root))
	if config.reqOnlyDir {
		files = ReadOnlyDir(files)
	}

	if config.sortByModTime {
		SortByModTime(files)
	}
	return files
}

func ReadDir(root string) []fs.DirEntry {
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Println(err)
		return make([]fs.DirEntry, 0)
	}
	return files
}

func IgnoreDotFiles(files []fs.DirEntry) []fs.DirEntry {
	fs := make([]fs.DirEntry, 0)
	for _, f := range files {
		if !strings.HasPrefix(f.Name(), ".") {
			fs = append(fs, f)
		}
	}
	return fs
}

func ReadOnlyDir(files []fs.DirEntry) []fs.DirEntry {
	dirs := make([]fs.DirEntry, 0)
	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, f)
		}
	}
	return dirs
}

func SortByModTime(files []fs.DirEntry) {
	sort.Slice(files, func(i, j int) bool {
		return getFileInfo(files[i]).ModTime().Unix() < getFileInfo(files[j]).ModTime().Unix()
	})
}

func resizeToNMinus1(n int, sl []bool) []bool {
	if n > 0 && len(sl) > 0 {
		sl = sl[:n-1]
	}
	return sl
}

func getFileInfo(fs fs.DirEntry) fs.FileInfo {
	fi, err := fs.Info()
	if err != nil {
		fmt.Println(OpenBrkt + err.Error() + CloseBrkt)
		return fi
	}
	return fi
}

func getBeforePipeVal(n int, isNthDirLast []bool, config TreeConfig) string {
	bp := ""
	if config.noIndent {
		return bp
	}

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

func getPipeVal(isLastFile bool, config TreeConfig) string {
	pipe := BoxVH //
	if isLastFile {
		pipe = BoxUpAndRig + BoxHor // └──
	}
	if config.noIndent {
		pipe = ""
	}
	return pipe
}

func getAfterPipeVal(root string, fi fs.DirEntry, config TreeConfig) string {
	ap := Space + fi.Name() //after pipe
	var relPath, fp string  // fp: file permission

	if config.reqRelPath {
		relPath = Space + root + PathSeperator + fi.Name()
		ap = relPath
	}

	if config.reqFilePermsn {
		fp = Space + OpenBrkt + getPermsnMode(fi, false) + CloseBrkt + Space
		ap = fp + fi.Name()
	}

	if config.reqRelPath && config.reqFilePermsn {
		ap = fp + relPath
	}
	return ap
}

func getPermsnMode(f fs.DirEntry, inOctal bool) string {
	fi, err := f.Info()
	if err != nil {
		fmt.Println(err)
		return "Unable read mode"
	}
	perm := fi.Mode()
	res := perm.String()
	if inOctal {
		res = fmt.Sprintf("%#o", perm.Perm())
	}
	return res
}

func getFileAttrsVal(file fs.DirEntry, config TreeConfig) string {
	attrs := ""
	if config.reqXmlFormat {
		attrs = " name=\"" + file.Name() + "\""
		if config.reqFilePermsn {
			attrs += " mode=\"" + getPermsnMode(file, true) + "\""
			attrs += " prot=\"" + getPermsnMode(file, false) + "\""
		}
	}
	if config.reqJsonFormat {
		attrs = ",\"name\":\"" + file.Name() + "\""
		if config.reqFilePermsn {
			attrs += ",\"mode\":\"" + getPermsnMode(file, true) + "\""
			attrs += ",\"prot\":\"" + getPermsnMode(file, false) + "\""
		}
	}
	return attrs
}

func formatRes(temp string, fc FileCount, config TreeConfig) string {
	op := ""
	dirStr := ""
	fileStr := ""
	if !config.reqXmlFormat && !config.reqJsonFormat {
		dirStr = fmt.Sprintf("%v directories", fc.dirCnt)
		if fc.dirCnt == 1 {
			dirStr = fmt.Sprintf("%v directory", fc.dirCnt)
		}

		fileStr = fmt.Sprintf("%v files", fc.fileCnt)
		if fc.fileCnt == 1 {
			fileStr = fmt.Sprintf("%v file", fc.fileCnt)
		}
		op = fmt.Sprintf("%v%v%v", temp, NewLine, dirStr)
		if !config.reqOnlyDir {
			op = fmt.Sprintf("%v, %v", op, fileStr)
		}
	}

	if config.reqXmlFormat {
		header := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
		dirStr = fmt.Sprintf("<directories>%v</directories>", fc.dirCnt)
		fileStr = fmt.Sprintf("<files>%v</files>", fc.dirCnt)
		report := fmt.Sprintf("  <report>\n   %v\n   %v\n  </report>", dirStr, fileStr)
		if config.reqOnlyDir {
			report = fmt.Sprintf("  <report>\n   %v\n  </report>", dirStr)
		}
		op = fmt.Sprintf("%v<%v>\n%v%v\n</%v>", header, Command, temp, report, Command)
	}

	if config.reqJsonFormat {
		op = fmt.Sprintf("[\n%v,\n  {\"type\":\"report\",\"directories\":%v", temp, fc.dirCnt)
		if !config.reqOnlyDir {
			op = fmt.Sprintf("%v,\"files\":%v}\n]", op, fc.fileCnt)
		}
	}

	return op
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
