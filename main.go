package main

import (
	"os"
	"sort"

	"fmt"
)

//"io"
//"path/filepath"
//"strings"

var level int

func dirTree(dir *os.File, path string, printfiles bool) error {
	//fmt.Printf("dir %v, path %v, printfiles %v \n", dir, path, printfiles)
	dir, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("No such directory")
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
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

	if len(dirItems) != 0 {
		for i, name := range dirItems {
			//if level == 0 {
			if i != len(dirItems)-1 {
				fmt.Println("├───" + name)
			} else {
				fmt.Println("└───" + name)
			}
			newPath := path + "/" + name
			level++
			err := dirTree(dir, newPath, printfiles)
			if err != nil {
				return err
			}
			//}
		}
		level--
	} else {
		level--
	}
	// if len(dirItems) != 0 {
	// 	for i, name := range dirItems {
	// 		//fmt.Printf("dir is %v\n", dirItems)
	// 		newPath := path + "/" + name
	// 		if level == 0 {
	// 			// fmt.Printf("├───" + name + "\n")
	// 			fmt.Printf("├───%v\n", name)
	// 			level++
	// 		} else {
	// 			indent := "│\t"
	// 			// for i := 0; i < level; i++ {
	// 			// 	indent += "\t"
	// 			// }
	// 			fmt.Printf(indent + "└───" + name + "\n")
	// 			fmt.Println(name)
	// 			level++
	// 		}
	// 		err := dirTree(dir, newPath, printfiles)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	level--
	// } else {
	// 	level--
	// 	return nil
	// }
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
	//fmt.Printf("out %v %[1]T, path %v, printfiles %v", out, path, printFiles)
}
