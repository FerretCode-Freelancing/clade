package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// containerCmd represents the container command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Manage containers over SSH",
	Long:  `Create, list, delete, and set the restart policy for your container.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Please use a subcommand!")
	},
}

func init() {
	rootCmd.AddCommand(containerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// containerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// containerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
