package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func ProgramInPath(program string, logExes bool) bool {
	for _, path := range Paths {
		if strings.HasSuffix(path, ".exe") || !DirExists(path, false) { // if its not a folder or it also doesnt exist
			continue
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files { // over each file
			if strings.HasSuffix(file.Name(), ".exe") {
				if logExes {
					Yellow.Printf("%s\n", file.Name())
				}
				if file.Name() == program {
					ProgramPath = path
					return true
				}
			}
		}
	}
	return false
}

func DirExists(path string, log bool) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if log {
			Red.Printf("Directory \"%s\" does not exist!", path)
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

func TrimLineEnd(str string) string {
	return strings.TrimSuffix(strings.TrimSuffix(str, "\n"), "\r")
}
