package main

import (
	"fmt"
	"os"
	"strings"
)

func dirContains(s []os.FileInfo, str string) bool {
	for _, v := range s {
		if v.Name() == str {
			return true
		}
	}
	return false
}

func CheckDir(path string, ogDir []string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(colorRed, "Directory \""+path+"\" does not exist!")
		fmt.Print(colorReset)
		CurDir = strings.Join(ogDir, "\\")
	}
}
