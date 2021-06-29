package main

import (
	"bufio"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
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
var ShellVer = "v0.0.1"
var Paths = strings.Split(strings.ReplaceAll(os.Getenv("Path"), "\n", ""), ";")
var ProgramPath string

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorWhite = "\033[37m"

//func init() {
//	if runtime.GOOS == "windows" {
//		colorReset = ""
//		colorRed = ""
//		colorGreen = ""
//		colorYellow = ""
//		colorBlue = ""
//		colorPurple = ""
//		colorWhite = ""
//	}
//}

func main() {
	reader := bufio.NewReader(os.Stdin)

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
				ExecutePathProgram(split[0], split)
			}
			if strings.HasPrefix(split[0], "cd") {
				ChangeDirectory(split[1:])
			}
			if strings.HasPrefix(split[0], "ls") {
				ListDirectory()
			}
			if strings.HasPrefix(split[0], "db") || strings.HasPrefix(split[0], "debug") {
				GetDebugInfo(split[1:])
			}
			if strings.HasPrefix(split[0], "h") || strings.HasPrefix(split[0], "help") {
				ShowHelp(split[1:])
			}
			if strings.HasPrefix(split[0], "sys") {
				ShowSystemInfo()
			}
		}
		fmt.Println()
	}
}

func ExecutePathProgram(program string, command []string) {
	files, err := ioutil.ReadDir(CurDir)
	if err != nil {
		log.Fatal(err)
	}

	fileRun := strings.TrimSuffix(strings.Split(program, "./")[1], "\n")
	if !strings.HasSuffix(fileRun, ".exe") {
		fileRun += ".exe"
	}

	if dirContains(files, fileRun) {
		cmdInstance := exec.Command(CurDir+"\\"+fileRun, command[1:]...)
		cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
		cmdInstance.Stdout = os.Stdout
		cmdInstance.Stdin = os.Stdin
		cmdInstance.Stderr = os.Stderr
		err = cmdInstance.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cmdInstance.Stdout)
	} else if programInPath(fileRun) {
		cmdInstance := exec.Command(ProgramPath+"\\"+fileRun, command[1:]...)
		cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
		cmdInstance.Stdout = os.Stdout
		cmdInstance.Stdin = os.Stdin
		cmdInstance.Stderr = os.Stderr
		err = cmdInstance.Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(cmdInstance.Stdout)

	} else {
		fmt.Println(colorRed, "File \""+fileRun+"\" is not recognised!")
		fmt.Print(colorReset)
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
	if !CheckDir(CurDir, true) {
		CurDir = strings.Join(originalDir, "\\")
	}
}

func ListDirectory() {
	files, err := ioutil.ReadDir(CurDir)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), ".") {
			fmt.Print(colorRed, file.Name())
			fmt.Print(colorWhite, "/\n")
			continue
		}

		if file.IsDir() {
			fmt.Print(colorBlue, file.Name())
			fmt.Print(colorWhite, "/\n")
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

func GetDebugInfo(prop []string) {
	var curDir, _ = os.Getwd()
	var userName, _ = user.Current()
	var hostName, _ = os.Hostname()

	if len(prop) == 0 {
		fmt.Print(colorGreen, "~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "MelonShell ")
		fmt.Print(colorYellow, ShellVer)
		fmt.Println(colorReset)
		fmt.Print(colorGreen, "~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "User: ")
		fmt.Print(colorYellow, userName.Name)
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Host: ")
		fmt.Print(colorYellow, hostName)
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Shell Location: ")
		fmt.Print(colorYellow, curDir)
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Current Working Directory: ")
		fmt.Print(colorYellow, CurDir)

		return
	}

	switch strings.TrimSuffix(prop[0], "\n") {
	case "cwd":
		fmt.Println(colorYellow, curDir)
	case "ver":
		fmt.Println(colorYellow, ShellVer)
	case "loc":
		fmt.Println(colorYellow, CurDir)
	}
}

func ShowHelp(prop []string) {
	if len(prop) == 0 {
		fmt.Print(colorGreen, "~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "MelonShell Help")
		fmt.Println(colorReset)
		fmt.Print(colorGreen, "~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Change Directory: ")
		fmt.Print(colorYellow, "cd [ .. | ../ | <folder-name> ]")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "List Directory: ")
		fmt.Print(colorYellow, "ls")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Run Program: ")
		fmt.Print(colorYellow, "./[ app-name ]")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Show System Information: ")
		fmt.Print(colorYellow, "sys")
		fmt.Println(colorReset)

		return
	}
}

func ShowSystemInfo() {
	const unknown = "Unknown"
	osName, cpuName, ramAmount, diskAmount := unknown, unknown, unknown, unknown

	hStat, err := host.Info()
	if err == nil {
		osName = hStat.Platform + " " + hStat.PlatformVersion
	}
	// rewrite
	cStats, err := cpu.Info()
	if err == nil && len(cStats) > 0 {
		cpuName = fmt.Sprintf("%s, %d cores", cStats[0].ModelName, cStats[0].Cores)
	}

	mStat, err := mem.VirtualMemory()
	if err == nil {
		ramAmount = humanize.IBytes(mStat.Total)
	}

	dStat, err := disk.Usage("C:\\")
	if err == nil {
		diskAmount = humanize.IBytes(dStat.Total)
	}

	fmt.Print(colorGreen, "~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println(colorReset)
	fmt.Print(colorWhite, "System Information")
	fmt.Println(colorReset)
	fmt.Print(colorGreen, "~~~~~~~~~~~~~~~~~~~~~~")
	fmt.Println(colorReset)
	fmt.Print(colorWhite, "Operating System: ")
	fmt.Print(colorYellow, osName)
	fmt.Println(colorReset)
	fmt.Print(colorWhite, "CPU: ")
	fmt.Print(colorYellow, cpuName)
	fmt.Println(colorReset)
	fmt.Print(colorWhite, "RAM: ")
	fmt.Print(colorYellow, ramAmount)
	fmt.Println(colorReset)
	fmt.Print(colorWhite, "Drive: ")
	fmt.Print(colorYellow, diskAmount)
	fmt.Println(colorReset)
}
