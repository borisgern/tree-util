package main

import (
	"os"
	"sort"

	"fmt"
)

//"io"
//"path/filepath"
//"strings"

var levels map[int][]bool

var level int

func dirTree(dir *os.File, path string, printfiles bool) error {
	fmt.Printf("dir %v, path %v, printfiles %v \n", dir, path, printfiles)
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("No such directory")
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		//fmt.Println(err)
		return fmt.Errorf("No items in directory")
	}
	//fmt.Printf("fileInfos %v\n", fileInfos)
	var dirItems []string

	for _, fi := range fileInfos {
		if !fi.IsDir() && !printfiles {
			continue
		} else {
			dirItems = append(dirItems, fi.Name())
		}
	}
	sort.Strings(dirItems)
	levels[level] = []bool{true, true}
	if len(dirItems) != 0 {
		for i, name := range dirItems {
			if i != len(dirItems)-1 {
				levels[level][1] = true
				if level == 0 {
					fmt.Println("├───" + name)
				} else if levels[level][0] == levels[level-1][0] {
					fmt.Printf("│\t")
					fmt.Println("├───" + name)
				}
			} else {
				levels[level][1] = false
				if level == 0 {
					fmt.Println("└───" + name)
				} else if levels[level][0] == levels[level-1][0] {
					for i := 0; i < level; i++ {
						if levels[i][0] && levels[i][1] {
							fmt.Printf("|\t")
						} else {
							fmt.Printf("\t")
						}
					}
					if i != len(dirItems)-1 {
						fmt.Println("├───" + name)
					} else {
						fmt.Println("└───" + name)
					}
				}
			}
			//fmt.Printf("│\t")

			newPath := path + "/" + name
			level++
			err := dirTree(dir, newPath, printfiles)
			if err != nil {
				return err
			}

			//}
		}

		//level--
	}
	levels[level][0] = false
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
	levels = make(map[int][]bool, 1)
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	//fmt.Printf("out %v %[1]T, path %v, printfiles %v", out, path, printFiles)
}
