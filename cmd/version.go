package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	Ver       string
	Commit    string
	BuildDate string
)

func Version() *cobra.Command {
	cmdClient := &cobra.Command{
		Use:   "version",
		Short: "wx-msg-push version",
		Run: func(cmd *cobra.Command, args []string) {
			PrintVersion()
		},
	}
	return cmdClient
}

func PrintVersion() {
	fmt.Println("wx-msg-push", Ver)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Git Commit:", Commit)
	fmt.Println("BuildDate:", BuildDate)
	fmt.Println("Repository: https://github.com/zyh94946/wx-msg-push")
}
