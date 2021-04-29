package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zyh94946/wx-msg-push/cmd"
)

var (
	v string
	c string
	d string
)

func main() {
	cmd.Ver = v
	cmd.Commit = c
	cmd.BuildDate = d

	rootCmd := &cobra.Command{Use: "wx-msg-push"}
	rootCmd.AddCommand(cmd.Server())
	rootCmd.AddCommand(cmd.Version())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("rootCmd.Execute failed", err.Error())
	}
}
