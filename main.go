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
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	_ "strings"
	"syscall"
)

var CurDir, _ = os.Getwd()
var ShellVer = "v0.0.1"
var Paths = strings.Split(strings.ReplaceAll(os.Getenv("Path"), endLine, ""), ";")
var ProgramPath string

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"
var colorBlue = "\033[34m"
var colorPurple = "\033[35m"
var colorWhite = "\033[37m"
var endLine = "\n"
var env = "jet"

func init() {
	if runtime.GOOS == "windows" {
		colorReset = ""
		colorRed = ""
		colorGreen = ""
		colorYellow = ""
		colorBlue = ""
		colorPurple = ""
		colorWhite = ""
		endLine = "\r\n"
		env = "dos"
	}
}

func main() {
	Paths = append(Paths, strings.TrimSuffix(CurDir, endLine))
	reader := bufio.NewReader(os.Stdin)

	for {
		userName, _ := user.Current()
		hostName, _ := os.Hostname()

		fmt.Printf("%s%s@%s", colorGreen, strings.Split(userName.Username, "\\")[1], hostName)
		fmt.Printf("%s: ", colorWhite)
		fmt.Printf("%s%s", colorPurple, CurDir)
		fmt.Printf("%s\r\n> ", colorWhite)
		fmt.Printf("%s", colorReset)

		input, _ := reader.ReadString('\n')

		if input != "" {
			split := strings.Split(input, " ")

			for i := range split {
				split[i] = strings.TrimSuffix(split[i], endLine)
			}

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
				ShowDebugInfo(split[1:])
			}
			if strings.HasPrefix(split[0], "h") || strings.HasPrefix(split[0], "help") {
				ShowHelp(split[1:])
			}
			if strings.HasPrefix(split[0], "sys") {
				ShowSystemInfo()
			}
			if strings.HasPrefix(split[0], "exit") {
				os.Exit(1)
			}
			if strings.HasPrefix(split[0], "mkdir") {
				MakeDir(split[1])
			}
			if strings.HasPrefix(split[0], "make") {
				MakeItem(split[1])
			}
			if strings.HasPrefix(split[0], "rm") {
				DelItem(split[1:])
			}
			if strings.HasPrefix(split[0], "melon") {
				Melon()
			}
		}
		fmt.Println()
	}
}

func ExecutePathProgram(program string, command []string) {
	fileRun := strings.TrimSuffix(strings.Split(program, "./")[1], endLine)
	if !strings.HasSuffix(fileRun, ".exe") {
		fileRun += ".exe"
	}

	if programInPath(fileRun) {
		cmdInstance := exec.Command(ProgramPath+"\\"+fileRun, command[1:]...)
		cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
		cmdInstance.Stdout = os.Stdout
		cmdInstance.Stdin = os.Stdin
		cmdInstance.Stderr = os.Stderr
		cmdInstance.Dir = CurDir
		err := cmdInstance.Run()
		if err != nil {
			fmt.Println(err)
		}

	} else {
		fmt.Printf("%s, File\"%s\" is not recognised!", colorRed, fileRun)
		fmt.Printf("%s", colorReset)
	}
}

func ChangeDirectory(path []string) {
	originalDir := strings.Split(CurDir, "\\")
	structureAction := NewStructure(path[0])

	if strings.Contains(path[0], ":") { // drive directory
		CurDir = strings.Join(path, "\\")
		ValidateDir(originalDir)
		return
	}

	if strings.Contains(path[0], "\\\\") { // server directory
		CurDir = strings.Join(path, "\\")
		ValidateDir(originalDir)
		return
	}

	if strings.HasPrefix(path[0], "~") {
		CurDir, _ = os.Getwd()
		ValidateDir(originalDir)
		return
	}

	if len(structureAction.dirChanges) == 1 && strings.TrimSuffix(structureAction.dirChanges[0], endLine) == ".." {
		CurDir = strings.Join(originalDir[:len(originalDir)-1], "\\")
		Paths = Paths[:len(Paths)-1]
		Paths = append(Paths, strings.TrimSuffix(CurDir, endLine))
		ValidateDir(originalDir)
		return
	}

	for _, dir := range structureAction.dirChanges {
		splitDir := strings.Split(CurDir, "\\")
		switch dir {
		case "..":
			CurDir = strings.Join(splitDir[:len(splitDir)-1], "\\")
			Paths = Paths[:len(Paths)-1]
			Paths = append(Paths, strings.TrimSuffix(CurDir, endLine))
		default:
			if dir != endLine {
				CurDir += "\\" + dir
				Paths = Paths[:len(Paths)-1]
				Paths = append(Paths, strings.TrimSuffix(CurDir, endLine))
			}
		}
	}

	CurDir = strings.TrimSuffix(CurDir, endLine)
	Paths = Paths[:len(Paths)-1]
	Paths = append(Paths, strings.TrimSuffix(CurDir, endLine))

	ValidateDir(originalDir)
}

func ListDirectory() {
	files, err := ioutil.ReadDir(CurDir + "\\")

	if err != nil {
		fmt.Printf("%s%s", colorRed, err)
		fmt.Printf("%s", colorReset)
	}

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), ".") {
			fmt.Printf("%s%s", colorRed, file.Name())
			fmt.Printf("%s/\r\n", colorWhite)
			continue
		}

		if file.IsDir() {
			fmt.Printf("%s%s", colorBlue, file.Name())
			fmt.Printf("%s/\r\n", colorWhite)
			continue
		}

		if strings.HasPrefix(file.Name(), ".") {
			fmt.Printf("%s%s\r\n", colorGreen, file.Name())
			continue
		} else {

			fmt.Printf("%s./", colorWhite)
			fmt.Printf("%s%s\r\n", colorYellow, file.Name())
			continue
		}
	}
	fmt.Printf("%s", colorReset)
}

func MakeDir(folder string) {
	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return
	}
}

func MakeItem(name string) {
	_, err := os.Create(name)
	if err != nil {
		return
	}
}

func DelItem(pathArr []string) {
	path := NewStructure(pathArr[0])

	if len(path.dirChanges) == 1 {
		if !DirExists(CurDir+"\\"+strings.TrimSuffix(path.dirChanges[0], endLine), true) {
			return
		}
		err := os.RemoveAll(CurDir + "\\" + strings.TrimSuffix(path.dirChanges[0], endLine))
		if err != nil {
			return
		}
		return
	} else {
		if !DirExists(CurDir+"\\"+strings.Join(path.dirChanges, "\\"), true) {
			return
		}
		err := os.RemoveAll(CurDir + "\\" + strings.Join(path.dirChanges, "\\"))
		if err != nil {
			return
		}
		return
	}
}

func ShowDebugInfo(prop []string) {
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
		fmt.Print(colorWhite, "Current Directory: ")
		fmt.Print(colorYellow, CurDir)
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "%PATH% Environment Variable: ")
		fmt.Print(colorYellow, Paths)
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Terminal Behaviour: ")
		fmt.Print(colorYellow, env)

		return
	}

	switch strings.TrimSuffix(prop[0], "\n") {
	case "dir":
		fmt.Println(colorYellow, CurDir)
	case "ver":
		fmt.Println(colorYellow, ShellVer)
	case "loc":
		fmt.Println(colorYellow, curDir)
	case "path":
		fmt.Println(colorYellow, Paths)
	case "switch":
		if env == "dos" {
			colorReset = "\033[0m"
			colorRed = "\033[31m"
			colorGreen = "\033[32m"
			colorYellow = "\033[33m"
			colorBlue = "\033[34m"
			colorPurple = "\033[35m"
			colorWhite = "\033[37m"
			endLine = "\n"
			env = "jet"
		} else if env == "dos" {
			colorReset = ""
			colorRed = ""
			colorGreen = ""
			colorYellow = ""
			colorBlue = ""
			colorPurple = ""
			colorWhite = ""
			endLine = "\r\n"
			env = "dos"
		}

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
		fmt.Print(colorYellow, "cd [ .. | ../ | < folder-name > ]")
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
		fmt.Print(colorWhite, "Show Debug Information: ")
		fmt.Print(colorYellow, "db < dir | ver | loc | path >")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Make New Directory: ")
		fmt.Print(colorYellow, "mkdir [ name ]")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Delete Folder/File: ")
		fmt.Print(colorYellow, "rm [ name ]")
		fmt.Println(colorReset)
		fmt.Print(colorWhite, "Delete Folder/File: ")
		fmt.Print(colorYellow, "rm [ name ]")
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

func Melon() {
	fmt.Println("                                                                                \r\n                                                       ((((((///////            \r\n                                                       ((((((///////            \r\n                                                 //////%&&&&&#(((((///////      \r\n                                                 //////%&&&&&#(((((///////      \r\n                                                 //////%&&&&&#(((((///////      \r\n                                           ////////////%&&&&&&&&&&&#(((((///////\r\n                                           ////////////%&&&&&&&&&&&#(((((///////\r\n                                     ////////////.     */////%&&&&&#(((((///////\r\n                                     ////////////.     */////%&&&&&#(((((///////\r\n                               //////////////////,     */////%&&&&&#(((((///////\r\n                               //////////////////////////////%&&&&&#(((((///////\r\n                               //////////////////////////////%&&&&&#(((((///////\r\n                         //////////////////.     *///////////%&&&&&#(((((///////\r\n                         //////////////////.     *///////////%&&&&&#(((((///////\r\n                   //////////////////////////////////////////%&&&&&#(((((///////\r\n                   //////////////////////////////////////////%&&&&&#(((((///////\r\n                   //////////////////////////////////////////%&&&&&#(((((///////\r\n             //////////////////.     ,/////////////////%&&&&&&&&&&&#(((((///////\r\n             //////////////////.     ,/////////////////%&&&&&&&&&&&#(((((///////\r\n       ////////////.     */////////////////////////////%&&&&&#(((((///////      \r\n       ////////////.     */////////////////////////////%&&&&&#(((((///////      \r\n ((((((((((((((((((*.....*/////////////////((((((((((((%%%%%%((((((///////      \r\n ((((((%&&&&&&&&&&&#///////////////////////%&&&&&&&&&&&#(((((/////////////      \r\n ((((((%&&&&&&&&&&&#///////////////////////%&&&&&&&&&&&#(((((/////////////      \r\n ///////(((((%&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&#(((((/////////////            \r\n ///////(((((%&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&#(((((/////////////            \r\n       ///////(((((((((((((((((((((((((((((((((((/////////////                  \r\n       ///////(((((((((((((((((((((((((((((((((((/////////////                  \r\n       ///////(((((((((((((((((((((((((((((((((((/////////////                  \r\n             /////////////////////////////////////                              \r\n             /////////////////////////////////////                              ")
}
