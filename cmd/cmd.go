package cmd

import (
	"fmt"

	"github.com/cothi/chat-go/server"
	"github.com/cothi/chat-go/ui"
	"github.com/cothi/chat-go/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	port       string
	serverPort string

	rootCmd = &cobra.Command{
		Use:   "go-chat",
		Short: "Terminal based chat made with go",
		Long:  "Terminal based chat made with go",
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of chat-go",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("chat-go version v0.1 -- HEAD")
		},
	}

	serverCmd = &cobra.Command{
		Use:   "server [OPTIONS]",
		Short: "Start chat server",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start(utils.PortSet(port))
		},
	}

	clientCmd = &cobra.Command{
		Use:   "client [OPTIONS]",
		Short: "Start client chat",
		Run: func(cmd *cobra.Command, args []string) {
			ui.StartClient(utils.PortSet(serverPort))
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&port, "port", "8000", "set port to communicate with each other")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", "8000")

	clientCmd.Flags().StringVarP(&serverPort, "serverPort", "p", "", "server port for communication")
	clientCmd.MarkFlagRequired("serverPort")

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(clientCmd)
}
