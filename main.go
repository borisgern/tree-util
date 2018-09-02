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

//"io"
//"path/filepath"
//"strings"

// var dirItems map[string]os.FileInfo
var level int
var levels []levelStatus

func dirTree(out io.Writer, path string, printfiles bool) error {
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
	//fmt.Printf("before for prefix %v\n", prefix)
	for i := level; i > 0; i-- {
		//fmt.Printf("in for i %v, level %v, levels[i].isOpened %v, levels[i].headOnLast %v\n\n", i, level, levels[i].isOpened, levels[i].headOnLast)
		if levels[i-1].isOpened && !levels[i-1].headOnLast {
			prefix = "|\t" + prefix
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
		return fmt.Errorf("No such directory")
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
			newPath := path + "/" + item
			level++
			levels = append(levels, levelStatus{isOpened: true, headOnLast: false})
			dirTreeIn(out, newPath, printfiles)
		}
	}
	levels[level].isOpened = false
	level--
	return nil
}

// var levels map[int][]bool

// var level int

// func dirTree(dir io.Writer, path string, printfiles bool) error {
// 	fmt.Printf("dir %v, path %v, printfiles %v \n", dir, path, printfiles)
// 	dir, err := os.Open(path)
// 	if err != nil {
// 		return fmt.Errorf("No such directory")
// 	}
// 	defer dir.Close()

// 	fileInfos, err := dir.Readdir(-1)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return fmt.Errorf("No items in directory")
// 	}
// 	//fmt.Printf("fileInfos %v\n", fileInfos)
// 	var dirItems []string

// 	for _, fi := range fileInfos {
// 		if !fi.IsDir() && !printfiles {
// 			continue
// 		} else {
// 			dirItems = append(dirItems, fi.Name())
// 		}
// 	}
// 	sort.Strings(dirItems)
// 	levels[level] = []bool{true, true}
// 	if len(dirItems) != 0 {
// 		for i, name := range dirItems {
// 			if i != len(dirItems)-1 {
// 				levels[level][1] = true
// 				if level == 0 {
// 					fmt.Println("├───" + name)
// 				} else if levels[level][0] == levels[level-1][0] {
// 					fmt.Printf("│\t")
// 					fmt.Println("├───" + name)
// 				}
// 			} else {
// 				levels[level][1] = false
// 				if level == 0 {
// 					fmt.Println("└───" + name)
// 				} else if levels[level][0] == levels[level-1][0] {
// 					for i := 0; i < level; i++ {
// 						if levels[i][0] && levels[i][1] {
// 							fmt.Printf("|\t")
// 						} else {
// 							fmt.Printf("\t")
// 						}
// 					}
// 					if i != len(dirItems)-1 {
// 						fmt.Println("├───" + name)
// 					} else {
// 						fmt.Println("└───" + name)
// 					}
// 				}
// 			}
// 			//fmt.Printf("│\t")

// 			newPath := path + "/" + name
// 			level++
// 			err := dirTree(dir, newPath, printfiles)
// 			if err != nil {
// 				return err
// 			}

// 			//}
// 		}

// 		//level--
// 	}
// 	levels[level][0] = false
// 	level--

// 	return nil
// }

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
	//fmt.Printf("out %v %[1]T, path %v, printfiles %v", out, path, printFiles)
}
