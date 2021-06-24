package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	_ "strings"
)

var CurDir, _ = os.Getwd()

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

		fmt.Print(userName.Username + "@" + hostName + ": " + CurDir + "\n> ")
		input, _ := reader.ReadString('\n')

		if input != "" {
			exec := strings.Split(input, " ")

			if strings.HasPrefix(exec[0], "./") {
				ExecutePathProgram(exec[0])
			}
			if strings.HasPrefix(exec[0], "cd") {
				ChangeDirectory(exec[1:])
			}
		}
	}
}

func ExecutePathProgram(program string) {
	fmt.Print(strings.Split(program, "./")[1])
}

func ChangeDirectory(path []string) {
	splitDir := strings.Split(CurDir, "\\")
	if strings.HasPrefix(path[0], "../") || strings.HasPrefix(path[0], "..") {
		CurDir = strings.Join(splitDir[:len(splitDir)-1], "\\")
	}

	if strings.HasPrefix(path[0], "~") {
		CurDir, _ = os.Getwd()
	}
}
