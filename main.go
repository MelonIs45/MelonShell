package main

import (
	"bufio"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	_ "strings"
	"syscall"
)

var CurDir, _ = os.Getwd()
var ShellVer = "v0.0.1"
var Paths = strings.Split(strings.ReplaceAll(os.Getenv("Path"), endLine, ""), ";")
var ProgramPath string

var Red = color.New(color.FgRed)
var Green = color.New(color.FgHiGreen)
var Yellow = color.New(color.FgYellow)
var Cyan = color.New(color.FgCyan)
var Magenta = color.New(color.FgHiMagenta)
var White = color.New(color.FgWhite)
var endLine = "\r\n"

func main() {
	Paths = append(Paths, TrimLineEnd(CurDir))
	reader := bufio.NewReader(os.Stdin)

	for {
		userName, _ := user.Current()
		hostName, _ := os.Hostname()

		Green.Printf("%s@%s", strings.Split(userName.Username, "\\")[1], hostName)
		White.Printf(": ")
		Magenta.Printf("%s", CurDir)
		White.Printf("\n> ")

		input, _ := reader.ReadString('\n')

		if input != "" {
			split := strings.Split(input, "\"")
			if len(split) == 1 {
				split = strings.Split(input, " ")
			}

			for i := range split {
				split[i] = TrimLineEnd(split[i])
			}

			if strings.HasPrefix(split[0], "./") {
				ExecutePathProgram(split)
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
			if strings.HasPrefix(split[0], "mk") {
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

func ExecutePathProgram(command []string) {
	fileRun := TrimLineEnd(strings.Split(strings.Split(strings.Join(command, " "), "./")[1], " ")[0])
	if !strings.HasSuffix(fileRun, ".exe") {
		fileRun += ".exe"
	}

	if ProgramInPath(TrimLineEnd(fileRun), false) {
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
		Red.Printf("File \"%s\" is not recognised!", fileRun)
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

	if strings.HasPrefix(path[0], "~") { // root directory
		CurDir, _ = os.Getwd()
		ValidateDir(originalDir)
		return
	}

	if len(structureAction.dirChanges) == 1 && TrimLineEnd(structureAction.dirChanges[0]) == ".." {
		CurDir = strings.Join(originalDir[:len(originalDir)-1], "\\")
		Paths = Paths[:len(Paths)-1]
		Paths = append(Paths, TrimLineEnd(CurDir))
		ValidateDir(originalDir)
		return
	}

	for _, dir := range structureAction.dirChanges {
		splitDir := strings.Split(CurDir, "\\")
		switch dir {
		case "..":
			CurDir = strings.Join(splitDir[:len(splitDir)-1], "\\")
			Paths = Paths[:len(Paths)-1]
			Paths = append(Paths, TrimLineEnd(CurDir))
		default:
			if dir != endLine {
				CurDir += "\\" + dir
				Paths = Paths[:len(Paths)-1]
				Paths = append(Paths, TrimLineEnd(CurDir))
			}
		}
	}

	CurDir = TrimLineEnd(CurDir)
	Paths = Paths[:len(Paths)-1]
	Paths = append(Paths, TrimLineEnd(CurDir))

	ValidateDir(originalDir)
}

func ListDirectory() {
	files, err := ioutil.ReadDir(CurDir + "\\")

	if err != nil {
		Red.Printf("%s", err)
	}

	for _, file := range files {
		if file.IsDir() && strings.Contains(file.Name(), ".") {
			Red.Printf("%s", file.Name())
			White.Printf("/\n")
			continue
		}

		if file.IsDir() {
			Cyan.Printf("%s", file.Name())
			White.Printf("/\n")
			continue
		}

		if strings.HasPrefix(file.Name(), ".") {
			Green.Printf("%s\n", file.Name())
			continue
		} else {
			White.Printf("./")
			Yellow.Printf("%s\n", file.Name())
			continue
		}
	}
}

func MakeItem(name string) {
	if strings.Contains(name, ".") {
		_, err := os.Create(name)
		if err != nil {
			return
		}
	} else {
		err := os.MkdirAll(name, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func DelItem(pathArr []string) {
	path := NewStructure(pathArr[0])

	if len(path.dirChanges) == 1 {
		if !DirExists(CurDir+"\\"+TrimLineEnd(path.dirChanges[0]), true) {
			return
		}
		err := os.RemoveAll(CurDir + "\\" + TrimLineEnd(path.dirChanges[0]))
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
		Green.Printf("~~~~~~~~~~~~~~~~~~~~~~\n")
		White.Printf("MelonShell ")
		Yellow.Printf("%s\n", ShellVer)
		Green.Printf("~~~~~~~~~~~~~~~~~~~~~~\n")
		White.Printf("User: ")
		Yellow.Printf("%s\n", userName.Username)
		White.Printf("Host: ")
		Yellow.Printf("%s\n", hostName)
		White.Printf("Shell Location: ")
		Yellow.Printf("%s\n", curDir)
		White.Printf("Current Directory: ")
		Yellow.Printf("%s\n", CurDir)
		White.Printf("%%PATH%% Environment Variable: ")
		Yellow.Printf("%s\n", Paths)

		return
	}

	switch TrimLineEnd(prop[0]) {
	case "-dir":
		Yellow.Printf("%s\n", CurDir)
	case "-ver":
		Yellow.Printf("%s\n", ShellVer)
	case "-loc":
		Yellow.Printf("%s\n", curDir)
	case "-path":
		Yellow.Printf("%s\n", Paths)
	case "-exes":
		_ = ProgramInPath("", true)
	}
}

func ShowHelp(prop []string) {
	if len(prop) == 0 {
		Green.Printf("~~~~~~~~~~~~~~~~~~~~~~\n")
		White.Printf("MelonShell Help\n")
		White.Printf("[] = mandatory, <> = optional, | = either\n")
		Green.Printf("~~~~~~~~~~~~~~~~~~~~~~\n")
		White.Printf("Change Directory: ")
		Yellow.Printf("cd [ .. | ../ | folder-name ]\n")
		White.Printf("List Directory: ")
		Yellow.Printf("ls\n")
		White.Printf("Run Program: ")
		Yellow.Printf("./[ app-name ]\n")
		White.Printf("Show System Information: ")
		Yellow.Printf("sys\n")
		White.Printf("Show Debug Information:  ")
		Yellow.Printf("db < -dir | -ver | -loc | -path | -exes >\n")
		White.Printf("Make Folder/File: ")
		Yellow.Printf("mk [ name ]\n")
		White.Printf("Delete Folder/File: ")
		Yellow.Printf("rm [ name ]\n")

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

	Green.Printf("~~~~~~~~~~~~~~~~~~~~~~\n")
	White.Printf("System Information\n")
	Green.Printf("~~~~~~~~~~~~~~~~~~~~~~\n")
	White.Printf("Operating System: ")
	Yellow.Printf("%s\n", osName)
	White.Printf("CPU: ")
	Yellow.Printf("%s\n", cpuName)
	White.Printf("RAM Amount: ")
	Yellow.Printf("%s\n", ramAmount)
	White.Printf("Disk Amount: ")
	Yellow.Printf("%s\n", diskAmount)
}

func Melon() {
	fmt.Println("                                                                                \r\n                                                       ((((((///////            \r\n                                                       ((((((///////            \r\n                                                 //////%&&&&&#(((((///////      \r\n                                                 //////%&&&&&#(((((///////      \r\n                                                 //////%&&&&&#(((((///////      \r\n                                           ////////////%&&&&&&&&&&&#(((((///////\r\n                                           ////////////%&&&&&&&&&&&#(((((///////\r\n                                     ////////////.     */////%&&&&&#(((((///////\r\n                                     ////////////.     */////%&&&&&#(((((///////\r\n                               //////////////////,     */////%&&&&&#(((((///////\r\n                               //////////////////////////////%&&&&&#(((((///////\r\n                               //////////////////////////////%&&&&&#(((((///////\r\n                         //////////////////.     *///////////%&&&&&#(((((///////\r\n                         //////////////////.     *///////////%&&&&&#(((((///////\r\n                   //////////////////////////////////////////%&&&&&#(((((///////\r\n                   //////////////////////////////////////////%&&&&&#(((((///////\r\n                   //////////////////////////////////////////%&&&&&#(((((///////\r\n             //////////////////.     ,/////////////////%&&&&&&&&&&&#(((((///////\r\n             //////////////////.     ,/////////////////%&&&&&&&&&&&#(((((///////\r\n       ////////////.     */////////////////////////////%&&&&&#(((((///////      \r\n       ////////////.     */////////////////////////////%&&&&&#(((((///////      \r\n ((((((((((((((((((*.....*/////////////////((((((((((((%%%%%%((((((///////      \r\n ((((((%&&&&&&&&&&&#///////////////////////%&&&&&&&&&&&#(((((/////////////      \r\n ((((((%&&&&&&&&&&&#///////////////////////%&&&&&&&&&&&#(((((/////////////      \r\n ///////(((((%&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&#(((((/////////////            \r\n ///////(((((%&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&#(((((/////////////            \r\n       ///////(((((((((((((((((((((((((((((((((((/////////////                  \r\n       ///////(((((((((((((((((((((((((((((((((((/////////////                  \r\n       ///////(((((((((((((((((((((((((((((((((((/////////////                  \r\n             /////////////////////////////////////                              \r\n             /////////////////////////////////////                              ")
}
