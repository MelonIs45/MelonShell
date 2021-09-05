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

// Current directory, the initial directory that the program starts in
var CurDir, _ = os.Getwd()
var ShellVer = "v0.0.1"

// Path Environment Variable
var Paths = strings.Split(strings.ReplaceAll(os.Getenv("Path"), endLine, ""), ";")
var ProgramPath string

// Color variables
var Red = color.New(color.FgRed)
var Green = color.New(color.FgHiGreen)
var Yellow = color.New(color.FgYellow)
var Cyan = color.New(color.FgCyan)
var Magenta = color.New(color.FgHiMagenta)
var White = color.New(color.FgWhite)
var endLine = "\r\n"

func main() {
	Paths = append(Paths, TrimLineEnd(CurDir)) // Add initial directory to path variable
	reader := bufio.NewReader(os.Stdin)        // Create input stream reader to scan input

	for {
		userName, _ := user.Current()
		hostName, _ := os.Hostname()

		// This chunk is for outputting the line above the input character
		// in the format: user-name@host-name: current-directory
		Green.Printf("%s@%s", strings.Split(userName.Username, "\\")[1], hostName)
		White.Printf(": ")
		Magenta.Printf("%s", CurDir)
		White.Printf("\n> ")

		input, _ := reader.ReadString('\n')

		if input != "" {
			split := strings.Split(input, "\"")
			if len(split) == 1 { // If the length of split is 1, this means that there are
				// no " in the input, so spaces are used to separate the input.
				split = strings.Split(input, " ")
			}

			for i := range split {
				split[i] = TrimLineEnd(split[i])
			}

			//
			// ALL THE FUNCTIONS WITHIN THE PROGRAM
			//

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
	// fileRun is the input after ./
	fileRun := TrimLineEnd(strings.Split(strings.Split(strings.Join(command, " "), "./")[1], " ")[0])
	if !strings.HasSuffix(fileRun, ".exe") {
		fileRun += ".exe"
	}

	if ProgramInPath(TrimLineEnd(fileRun), false) {
		// This chunk initialises the start of a new subprocess
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
	// Directory before anything is changed
	originalDir := strings.Split(CurDir, "\\")
	structureAction := NewStructure(path[0])

	if strings.Contains(path[0], ":") { // Drive directory
		CurDir = strings.Join(path, "\\") // Changes the drive
		ValidateDir(originalDir)
		return
	}

	if strings.Contains(path[0], "\\\\") { // Server directory
		CurDir = strings.Join(path, "\\") // Changes the server
		ValidateDir(originalDir)
		return
	}

	if strings.HasPrefix(path[0], "~") { // Root directory
		CurDir, _ = os.Getwd() // Goes back to the root directory where the program first started
		ValidateDir(originalDir)
		return
	}

	if len(structureAction.dirChanges) == 1 && TrimLineEnd(structureAction.dirChanges[0]) == ".." {
		// This removes the last directory after \ in the current directory
		// acting as a way of moving one directory down.
		CurDir = strings.Join(originalDir[:len(originalDir)-1], "\\")
		Paths = Paths[:len(Paths)-1]
		Paths = append(Paths, TrimLineEnd(CurDir))
		ValidateDir(originalDir)
		return
	}

	// This gets ran when there are more than 1 directory changes in the input
	for _, dir := range structureAction.dirChanges {
		splitDir := strings.Split(CurDir, "\\")
		switch dir {
		case "..":
			// Removes last directory from current directory
			CurDir = strings.Join(splitDir[:len(splitDir)-1], "\\")
			Paths = Paths[:len(Paths)-1]
			Paths = append(Paths, TrimLineEnd(CurDir))
		default:
			if dir != endLine {
				// Adds directory to current directory
				CurDir += "\\" + dir
				Paths = Paths[:len(Paths)-1]
				Paths = append(Paths, TrimLineEnd(CurDir))
			}
		}
	}

	// Adds the newly changed directory to the list of path directories.
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
		if file.IsDir() && strings.Contains(file.Name(), ".") { // If the file is a folder and begins with a '.'
			Red.Printf("%s", file.Name())
			White.Printf("/\n")
			continue
		}

		if file.IsDir() { // If the file is a folder
			Cyan.Printf("%s", file.Name())
			White.Printf("/\n")
			continue
		}

		if strings.HasPrefix(file.Name(), ".") { // If the file begins with a '.'
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
	Red.Printf(name)
	if strings.Contains(name, ".") { // If the name contains a '.', this means its a file
		_, err := os.Create(CurDir + "\\" + name)
		if err != nil {
			return
		}
	} else { // Else its a folder
		err := os.MkdirAll(CurDir+"\\"+name, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func DelItem(pathArr []string) { // Fix
	path := NewStructure(pathArr[0])

	if !DirExists(CurDir+"\\"+strings.Join(path.dirChanges, "\\"), true) {
		return
	}
	err := os.RemoveAll(CurDir + "\\" + strings.Join(path.dirChanges, "\\"))
	if err != nil {
		return
	}
	return
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
	// Sets all the variables to "Unknown"
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
