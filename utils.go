package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

//func dirContains(s []os.FileInfo, program string) bool {
//	for _, v := range s {
//		if v.Name() == program {
//			return true
//		}
//	}
//	return false
//}

func programInPath(program string) bool {
	for _, path := range Paths {
		if strings.HasSuffix(path, ".exe") || !CheckDir(path, false) { // if its not a folder or it also doesnt exist
			continue
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files { // over each file
			if strings.HasSuffix(file.Name(), ".exe") && file.Name() == program {
				ProgramPath = path
				return true
			}
		}
	}
	return false
}

func CheckDir(path string, log bool) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if log {
			fmt.Println(colorRed, "Directory \""+path+"\" does not exist!")
			fmt.Print(colorReset)
		}
		return false
	}
	return true
}
