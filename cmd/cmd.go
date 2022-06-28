package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	port string

	rootCmd = &cobra.Command{
		Use:   "go-chat",
		Short: "Terminal based chat made with go",
		Long:  "Terminal based chat made with go",
	}

	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start server",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {

	rootCmd.PersistentFlags().StringVar(&port, "port", "", "set port to communicate with each other")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", "8000")
	rootCmd.AddCommand(serverCmd)
}
