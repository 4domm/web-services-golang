package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

func dirTree(out io.Writer, root string, printFiles bool) error {
	return printDirTree(out, root, printFiles, "")
}

func printDirTree(out io.Writer, root string, printFiles bool, indent string) error {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	var fileInfos []os.FileInfo
	for _, file := range files {
		if file.IsDir() || printFiles {
			fileInfos = append(fileInfos, file)
		}
	}

	for i, fileInfo := range fileInfos {
		var newIndent string
		if i == len(fileInfos)-1 {
			fmt.Fprintf(out, "%s└───%s", indent, fileInfo.Name())
			newIndent = indent + "\t"
		} else {
			fmt.Fprintf(out, "%s├───%s", indent, fileInfo.Name())
			newIndent = indent + "│\t"
		}

		if !fileInfo.IsDir() && printFiles {
			fmt.Fprintf(out, " (%s)", formatSize(fileInfo.Size()))
		}
		fmt.Fprintf(out, "\n")

		if fileInfo.IsDir() {
			if err := printDirTree(out, filepath.Join(root, fileInfo.Name()), printFiles, newIndent); err != nil {
				return err
			}
		}
	}

	return nil
}

func formatSize(size int64) string {
	if size == 0 {
		return "empty"
	}
	return fmt.Sprintf("%db", size)
}
