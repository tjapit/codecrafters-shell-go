package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
		if cmd == "exit" {
			return
		}
		fmt.Printf("%s: command not found\n", cmd)
	}
}
