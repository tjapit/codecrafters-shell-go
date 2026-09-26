package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type builtin_fn func(cmd string, args []string)
type BuiltinMap map[string]builtin_fn

func main() {
	reader := bufio.NewReader(os.Stdin)
	builtins := BuiltinMap{
		"exit": func(cmd string, args []string) { os.Exit(0) },
		"echo": func(cmd string, args []string) { fmt.Println(strings.Join(args, " ")) },
		"pwd": func(cmd string, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("error executing command: " + cmd)
				return
			}
			fmt.Println(cwd)
		},
		"cd": func(cmd string, args []string) {
			err := os.Chdir(args[0])
			if err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", args[0])
			}
		},
	}
	for {
		fmt.Print("$ ")
		cmd, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		args := strings.Split(strings.TrimSpace(cmd), " ")
		cmd, args = args[0], args[1:]
		if builtin := builtins[cmd]; builtin != nil {
			builtin(cmd, args)
		} else if cmd == "type" {
			lookupCmd := args[0]
			if found := builtins[lookupCmd]; found != nil || lookupCmd == "type" {
				fmt.Printf("%s is a shell builtin\n", lookupCmd)
			} else if fullpath := lookPath(lookupCmd); fullpath != "" {
				fmt.Printf("%s is %s\n", lookupCmd, fullpath)
			} else {
				fmt.Printf("%s: not found\n", lookupCmd)
			}
		} else {
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
