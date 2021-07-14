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
		if strings.HasSuffix(path, ".exe") || !DirExists(path, false) { // if its not a folder or it also doesnt exist
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

func DirExists(path string, log bool) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if log {
			fmt.Printf("%sDirectory \"%s\" does not exist!", colorRed, path)
			fmt.Printf("%s", colorReset)
		}
		return false
	}
	return true
}

func ValidateDir(originalDir []string) {
	if !DirExists(CurDir, true) {
		CurDir = strings.Join(originalDir, "\\")
		Paths = Paths[:len(Paths)-1]
		Paths = append(Paths, strings.TrimSuffix(CurDir, "\n"))
	}
}
