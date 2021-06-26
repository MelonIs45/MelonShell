package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	_ "strings"
)

var CurDir, _ = os.Getwd()

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorCyan = "\033[36m"
var colorWhite = "\033[37m"

//func init() {
//	if runtime.GOOS == "windows" {
//		colorReset  = ""
//		colorRed    = ""
//		colorGreen  = ""
//		colorYellow = ""
//		colorBlue   = ""
//		colorPurple = ""
//		colorCyan   = ""
//		colorWhite  = ""
//	}
//}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(os.Getenv("Path"))
	var paths []string
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	} else {
		paths = append(paths, path)
	}

	for _, e := range os.Environ() {
		if e == "Path" {
			fmt.Println(e)
		}
	}

	for {
		userName, _ := user.Current()
		hostName, _ := os.Hostname()

		fmt.Print(strings.Split(userName.Username, "\\")[1] + "@" + hostName + ": " + CurDir + "\n> ")
		input, _ := reader.ReadString('\n')

		if input != "" {
			exec := strings.Split(input, " ")

			if strings.HasPrefix(exec[0], "./") {
				ExecutePathProgram(exec[0])
			}
			if strings.HasPrefix(exec[0], "cd") {
				ChangeDirectory(exec[1:])
			}
			if strings.HasPrefix(exec[0], "ls") {
				ListDirectory()
			}
		}

		fmt.Println()
	}
}

func ExecutePathProgram(program string) {
	fmt.Print(strings.Split(program, "./")[1])
}

func ChangeDirectory(path []string) {
	splitDir := strings.Split(CurDir, "\\")

	if strings.HasPrefix(path[0], "../") || strings.HasPrefix(path[0], "..") {
		CurDir = strings.Join(splitDir[:len(splitDir)-1], "\\")
		return
	}

	if strings.HasPrefix(path[0], "~") {
		CurDir, _ = os.Getwd()
		return
	}

	CurDir += "\\" + path[0]
	CurDir = strings.TrimSuffix(CurDir, "\n")
}

func ListDirectory() {
	files, err := ioutil.ReadDir(CurDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), ".") {
			fmt.Print(colorRed, file.Name())
			fmt.Print(colorWhite, "\\\n")
			continue
		}

		if file.IsDir() {
			fmt.Print(colorBlue, file.Name())
			fmt.Print(colorWhite, "\\\n")
			continue
		}

		if !strings.HasPrefix(file.Name(), ".") {
			fmt.Print(colorWhite, "./")
			fmt.Print(colorYellow, file.Name()+"\n")
			continue
		}

		if strings.Contains(file.Name(), ".") {
			fmt.Println(colorGreen, file.Name())
			continue
		}

	}
	fmt.Print(colorReset)
}
