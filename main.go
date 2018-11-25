package main

import (
	"github.com/bino7/gmcc/data"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "gmcc"}
	rootCmd.AddCommand(data.CMDS...)
	rootCmd.Execute()
}
