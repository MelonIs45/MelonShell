# A custom console made in go (originally created in 2021).

## To run (tested on go1.23.10):
- use `go run .` in the same directory as the `main.go` file.

This will get you into the main console environment

## List of commands

| Command | Description |
| ------- | ----------- |
| ./ <program> | Runs an executable in the program argument |
| cd <directory> | Changes the console's directory to the one specified, accepts .. |
| ls | Lists all the files in the console's directory |
| db, debug | Lists general shell information |
| h, help | Shows help to do with commands |
| sys | Lists general system information |
| exit | Closes the shell |
| mk [directory] | Makes a folder at the directory argument path |
| rm [directory] | Removes a folder at the directory argument path |