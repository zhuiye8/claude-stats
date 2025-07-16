package main

import (
	"fmt"
	"os"

	"github.com/claude-stats/claude-stats/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
} 