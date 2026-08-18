package agentkit

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev" 

//Define the root command, starting with goclaw 
var rootCmd = &cobra.Command{
	Use: "cloudclaw", 
	Short: "CloudClaw - Agent gateway", 
	Long: "Cloudclaw-  multi agent AI platform with Websocket RPC, tool execution and channel intergration", 
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run gateway function") 
	}, 
}

//Define the version command, starting with version 
var versionCmd = &cobra.Command{
	Use: "version", 
	Short: "Cloudclaw version", 
	Long: "Cloudclaw version information", 
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("")
	}, 
}

// Execute runs the root cobra command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
