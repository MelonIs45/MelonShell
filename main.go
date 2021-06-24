package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	_ "strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

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
		fmt.Print("$> ")
		input, _ := reader.ReadString('\n')

		if input != "" {
			fmt.Println(0)

			exec := strings.Split(input, " ")

			if strings.HasPrefix(exec[0], "#/") {
				ExecutePathProgram(exec[0])
			}

		}
	}
}

func ExecutePathProgram(program string) {
	fmt.Print(program)
}
