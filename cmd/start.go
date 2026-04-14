package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/agglabs/loop-client/config"
	"github.com/agglabs/loop-client/daemon"
	"github.com/spf13/cobra"
)

var (
	runInBackground bool
	autoRetry       bool
)

var startCmd = &cobra.Command{
	Use:   "start [app]",
	Short: "Start application by name",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfiguration()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		if cfg.Token == "" {
			fmt.Println("[ ERROR ] Authentication token is missing. Add it using: aggloop add-token <token>")
			return
		}

		var targetName string
		if len(args) > 0 {
			targetName = args[0]
		} else {
			for name := range cfg.Applications {
				targetName = name
				break
			}
		}

		if targetName == "" {
			fmt.Println("No applications configured to start.")
			return
		}

		app, exists := cfg.Applications[targetName]
		if !exists {
			fmt.Printf("Application '%s' not found in config.\n", targetName)
			return
		}

		if runInBackground {
			if runtime.GOOS == "windows" {
				fmt.Println("[ ERROR ] Background mode (-b) is not supported on Windows. Please run in foreground.")
				return
			}

			ensureSudo()

			fmt.Printf("Starting tunnel for %s in the background...\n", targetName)

			exe, err := os.Executable()
			if err != nil {
				fmt.Println("Failed to get executable path:", err)
				return
			}

			cmdArgs := []string{"start", targetName, "--auto-retry=" + strconv.FormatBool(autoRetry)}
			bgCmd := exec.Command(exe, cmdArgs...)

			err = bgCmd.Start()
			if err != nil {
				fmt.Println("Failed to start background process:", err)
				return
			}

			pidFile := config.GetPidFilePath(targetName)
			err = os.WriteFile(pidFile, []byte(strconv.Itoa(bgCmd.Process.Pid)), 0644)
			if err != nil {
				fmt.Println("Warning: Failed to write PID file:", err)
			}

			fmt.Printf("Process started successfully. PID: %d\n", bgCmd.Process.Pid)
			fmt.Printf("Use 'aggloop stop %s' to terminate it.\n", targetName)
			os.Exit(0)
		}

		fmt.Printf("Starting tunnel for %s...\n", targetName)

		err = daemon.RunBackgroundDaemon(cfg.Token, app.AppID, app.Scheme, app.Addr, app.Inspect, autoRetry)
		if err != nil {
			fmt.Println("Failed to start connection:", err)
		}
	},
}

func init() {
	startCmd.Flags().BoolVarP(&runInBackground, "background", "b", false, "Run in background as a daemon (Requires Root, Linux/macOS only)")
	startCmd.Flags().BoolVarP(&autoRetry, "auto-retry", "r", true, "Auto-retry connection on failure")

	rootCmd.AddCommand(startCmd)
}
