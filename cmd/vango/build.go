package vango

import (
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your static site",
	Long:  `Build your static site using the configuration file.

This command processes all markdown files in the content directory,
applies templates, and generates a static website in the public directory.`,
	Example: `  vango build                    # Build with default config
  vango build -c custom.toml      # Build with custom config
  vango build --verbose           # Build with verbose output`,
	Run: func(cmd *cobra.Command, args []string) {
		buildSite(cmd)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

