package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/agglabs/loop-client/config"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [app]",
	Short: "Stop a background application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetName := args[0]
		pidFile := config.GetPidFilePath(targetName)

		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Printf("Application '%s' does not seem to be running (PID file not found).\n", targetName)
			return
		}

		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			fmt.Println("Invalid PID in file:", err)
			return
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("Failed to find process %d: %v\n", pid, err)
			return
		}

		err = process.Kill()
		if err != nil {
			fmt.Printf("Warning: Failed to kill process %d (it might have already exited).\n", pid)
		} else {
			fmt.Printf("Successfully stopped application '%s' (PID: %d).\n", targetName, pid)
		}

		os.Remove(pidFile)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
