package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtin = []string{"echo", "exit", "type"}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		cmd = strings.TrimSpace(cmd)
		args := strings.Split(cmd, " ")
		cmd, args = args[0], args[1:]
		switch cmd {
		case "exit":
			return
		case "echo":
			fmt.Println(strings.Join(args, " "))
		case "type":
			lookupCmd := args[0]
			if found := slices.Index(builtin, lookupCmd); found != -1 {
				fmt.Printf("%s is a shell builtin\n", lookupCmd)
			} else if fullpath := lookPath(lookupCmd); fullpath != "" {
				fmt.Printf("%s is %s\n", lookupCmd, fullpath)
			} else {
				fmt.Printf("%s: not found\n", lookupCmd)
			}
		default:
			if fullpath := lookPath(cmd); fullpath != "" {
				execCmd := exec.Command(cmd, args...)
				out, err := execCmd.Output()
				if err != nil {
					fmt.Println("error executing command: " + cmd)
					return
				}
				for _, v := range out {
					fmt.Printf("%v", string(v))
				}
			} else {
				fmt.Printf("%s: command not found\n", cmd)
			}
		}
	}
}

// lookPath returns the fullpath of the executable if found in the environmment
// PATH variable, empty string otherwise
func lookPath(cmd string) string {
	if envPaths, ok := os.LookupEnv("PATH"); ok {
		fullpath := ""
		for _, path := range strings.Split(envPaths, string(os.PathListSeparator)) {
			fullpath = findCmd(path, cmd)
			if fullpath != "" {
				return fullpath
			}
		}
	}
	return ""
}

// findCmd returns the fullpath of the executable in the given path if found,
// empty string otherwise
func findCmd(path, cmd string) string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return ""
		}
		perm := info.Mode().Perm()
		if entry.Name() == cmd && perm&0111 != 0 {
			return path + "/" + cmd
		}
	}
	return ""
}
