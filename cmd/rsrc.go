/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/akavel/rsrc/rsrc"
	"github.com/spf13/cobra"
)

var fnamein, fnameico, fnameout, arch string

// rsrcCmd represents the rsrc command
var rsrcCmd = &cobra.Command{
	Use:   "rsrc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if fnameout == "" {
			fnameout = "rsrc_windows_" + arch + ".syso"
		}

		err := rsrc.Embed(fnameout, arch, fnamein, fnameico)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(rsrcCmd)

	rsrcCmd.Flags().StringVarP(&fnamein, "manifest", "", "", "path to a Windows manifest file to embed")
	rsrcCmd.Flags().StringVarP(&fnameico, "ico", "", "", "comma-separated list of paths to .ico files to embed")
	rsrcCmd.Flags().StringVarP(&fnameout, "o", "", "", "name of output COFF (.res or .syso) file; if set to empty, will default to 'rsrc_windows_{arch}.syso'")
	rsrcCmd.Flags().StringVarP(&arch, "arch", "", "amd64", "architecture of output file - one of: 386, amd64, [EXPERIMENTAL: arm, arm64]")
}
