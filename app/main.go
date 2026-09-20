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
			lookup := args[0]
			if found := slices.Index(builtin, lookup); found != -1 {
				fmt.Printf("%s is a shell builtin\n", lookup)
			} else {
				path, err := exec.LookPath(lookup)
				if err != nil {
					fmt.Printf("%s: not found\n", lookup)
					continue
				}
				fmt.Printf("%s is %s\n", lookup, path)
			}
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
