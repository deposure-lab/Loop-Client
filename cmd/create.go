package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/agglabs/loop-client/config"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new application interactively",
	Run: func(cmd *cobra.Command, args []string) {
		ensureSudo()

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Application name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Print("Application ID (UUID): ")
		appId, _ := reader.ReadString('\n')
		appId = strings.TrimSpace(appId)

		fmt.Print("Port: ")
		port, _ := reader.ReadString('\n')
		port = strings.TrimSpace(port)

		fmt.Print("Protocol (http/tcp/udp) [http]: ")
		scheme, _ := reader.ReadString('\n')
		scheme = strings.TrimSpace(scheme)
		if scheme == "" {
			scheme = "http"
		}

		cfg, err := config.LoadConfiguration()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		cfg.Applications[name] = config.AppConfig{
			AppID:   appId,
			Addr:    port,
			Scheme:  scheme,
			Inspect: true,
		}

		if err := config.SaveConfiguration(cfg); err != nil {
			fmt.Println("Failed to save configuration:", err)
			return
		}

		fmt.Printf("Application \"%s\" added successfully to %s\n", name, config.GetConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
