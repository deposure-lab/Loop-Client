package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "aggloop",
	Short: "AGG Loop - Edge Infrastructure Tunneling Client",
	Long:  `Kernel-level ZTNA engineered in Go. Maximum throughput, zero enterprise bloat.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func ensureSudo() {
	if runtime.GOOS != "windows" && os.Geteuid() != 0 {
		fmt.Println("[ ERROR ] This command must be run as sudo/root.")
		os.Exit(1)
	}
}
