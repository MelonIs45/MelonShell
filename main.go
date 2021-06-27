package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	_ "strings"
	"syscall"
)

var CurDir, _ = os.Getwd()

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorWhite = "\033[37m"

//func init() {
//	if runtime.GOOS == "windows" {
//		colorReset  = ""
//		colorRed    = ""
//		colorGreen  = ""
//		colorYellow = ""
//		colorBlue   = ""
//		colorPurple = ""
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

		fmt.Print(colorGreen, strings.Split(userName.Username, "\\")[1]+"@"+hostName)
		fmt.Print(colorWhite, ": ")
		fmt.Print(colorPurple, CurDir)
		fmt.Print(colorWhite, "\n> ")
		fmt.Print(colorReset)

		input, _ := reader.ReadString('\n')

		if input != "" {
			split := strings.Split(input, " ")

			if strings.HasPrefix(split[0], "./") {
				ExecutePathProgram(split[0])
			}
			if strings.HasPrefix(split[0], "cd") {
				ChangeDirectory(split[1:])
			}
			if strings.HasPrefix(split[0], "ls") {
				ListDirectory()
			}
		}

		fmt.Println()
	}
}

func ExecutePathProgram(program string) {
	files, err := ioutil.ReadDir(CurDir)
	if err != nil {
		log.Fatal(err)
	}

	fileRun := strings.TrimSuffix(strings.Split(program, "./")[1], "\n")
	if !strings.HasSuffix(fileRun, ".exe") {
		fileRun += ".exe"
	}

	fmt.Println(CurDir + "\\" + fileRun)

	if dirContains(files, fileRun) {
		cmdInstance := exec.Command(CurDir+"\\"+fileRun, "")
		cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
		err = cmdInstance.Start()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ChangeDirectory(path []string) {
	originalDir := strings.Split(CurDir, "\\")
	structureAction := NewStructure(path[0])

	if strings.HasPrefix(path[0], "~") {
		CurDir, _ = os.Getwd()
		return
	}

	if len(structureAction.dirChanges) == 1 && strings.TrimSuffix(structureAction.dirChanges[0], "\n") == ".." {
		CurDir = strings.Join(originalDir[:len(originalDir)-1], "\\")
		return
	}

	for _, dir := range structureAction.dirChanges {
		splitDir := strings.Split(CurDir, "\\")
		switch dir {
		case "..":
			CurDir = strings.Join(splitDir[:len(splitDir)-1], "\\")
		default:
			if !(dir == "\n") {
				CurDir += "\\" + dir
			}
		}
	}

	CurDir = strings.TrimSuffix(CurDir, "\n")
	CheckDir(CurDir, originalDir)
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

		if strings.HasPrefix(file.Name(), ".") {
			fmt.Println(colorGreen, file.Name())
			continue
		} else {
			fmt.Print(colorWhite, "./")
			fmt.Print(colorYellow, file.Name()+"\n")
			continue
		}

	}
	fmt.Print(colorReset)
}
