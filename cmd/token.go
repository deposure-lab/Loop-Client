package cmd

import (
	"fmt"

	"github.com/agglabs/loop-client/config"
	"github.com/spf13/cobra"
)

var addTokenCmd = &cobra.Command{
	Use:   "add-token [token]",
	Short: "Save authentication token",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfiguration()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		cfg.Token = args[0]
		if err := config.SaveConfiguration(cfg); err != nil {
			fmt.Println("Failed to save token:", err)
			return
		}

		fmt.Println("Token updated successfully in", config.GetConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(addTokenCmd)
}
