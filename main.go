package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type levelStatus struct {
	isOpened   bool
	headOnLast bool
}

var level int
var levels []levelStatus

func dirTree(out io.Writer, path string, printfiles bool) error {
	level = 0
	levels = make([]levelStatus, 1)
	err := dirTreeIn(&out, path, printfiles)
	if err != nil {
		return err
	}
	return nil
}

func writeStroke(out *io.Writer, prefix, str string, index int, dirItems map[string]os.FileInfo) error {
	var size string
	if !dirItems[str].IsDir() {
		if dirItems[str].Size() == 0 {
			size = "empty"
		} else {
			size = strconv.Itoa(int(dirItems[str].Size())) + "b"
		}
		size = fmt.Sprintf(" (%s)", size)
	}
	for i := level; i > 0; i-- {
		if levels[i-1].isOpened && !levels[i-1].headOnLast {
			prefix = "│\t" + prefix
		} else {
			prefix = "\t" + prefix
		}
	}
	_, err := io.WriteString(*out, prefix+str+size+"\n")
	if err != nil {
		return err
	}
	return nil
}

func sortKeys(m map[string]os.FileInfo) (s []string) {
	for k := range m {
		s = append(s, k)
	}
	sort.Strings(s)
	return
}

func dirTreeIn(out *io.Writer, path string, printfiles bool) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	dirItems := make(map[string]os.FileInfo)
	for _, fi := range fileInfos {
		if !fi.IsDir() && !printfiles {
			continue
		} else {
			dirItems[fi.Name()] = fi
		}
	}
	if len(dirItems) != 0 {
		sortedKeys := sortKeys(dirItems)
		var prefix string
		levels[level].isOpened = true
		for i, item := range sortedKeys {

			if i < len(sortedKeys)-1 {
				prefix = "├───"
			} else {
				levels[level].headOnLast = true
				prefix = "└───"
			}
			writeStroke(out, prefix, item, i, dirItems)
			if dirItems[item].IsDir() {
				newPath := path + "/" + item
				level++
				levels = append(levels, levelStatus{isOpened: true, headOnLast: false})
				dirTreeIn(out, newPath, printfiles)
			}
		}
	}
	levels[level].isOpened = false
	levels[level].headOnLast = false
	level--
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
