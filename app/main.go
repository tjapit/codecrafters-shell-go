package main

import (
	"bufio"
	"fmt"
	"os"
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
			if found := slices.Index(builtin, args[0]); found != -1 {
				fmt.Printf("%s is a shell builtin\n", args[0])
			} else {
				fmt.Printf("%s: not found\n", args[0])
			}
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
