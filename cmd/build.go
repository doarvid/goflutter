/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/doarvid/goflutter/internal/flutter"
	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		proj_path, err := getProjectPath(args)
		if err != nil {
			return
		}
		proj, err := flutter.NewFlutterProject(proj_path)
		if err != nil {
			return
		}
		proj.BuildGoApp(false)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
