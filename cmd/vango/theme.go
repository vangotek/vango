package vango

import (
	"fmt"
	"os"
	"text/tabwriter"

	"vango/internal/config"
	"vango/internal/theme"

	"github.com/spf13/cobra"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage themes",
	Long:  `Manage themes for your Vango site.`,
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available themes",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load("config.toml")
		themeManager := theme.NewThemeManager(cfg)
		themeManager.LoadThemes()
		themes := themeManager.ListThemes()

		if len(themes) == 0 {
			fmt.Println("No themes found. Create a theme with 'vango theme create <name>'")
			return
		}

		fmt.Println("Available Themes:")
		fmt.Println("")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tAUTHOR\tDESCRIPTION")
		fmt.Fprintln(w, "----\t-------\t------\t-----------")

		for _, theme := range themes {
			active := ""
			if themeManager.GetActiveTheme() != nil && themeManager.GetActiveTheme().Name == theme.Name {
				active = " (active)"
			}

			description := theme.Description
			if len(description) > 50 {
				description = description[:47] + "..."
			}

			fmt.Fprintf(w, "%s%s\t%s\t%s\t%s\n",
				theme.Name, active, theme.Version, theme.Author, description)
		}

		w.Flush()
	},
}

var themeInstallCmd = &cobra.Command{
	Use:   "install [name]",
    Short: "Install a theme from the theme repository",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        cfg, _ := config.Load("config.toml")
        themeManager := theme.NewThemeManager(cfg)
        
        if err := themeManager.InstallTheme(args[0]); err != nil {
            fmt.Fprintf(os.Stderr, "❌ Failed to install theme: %v\n", err)
            os.Exit(1)
        }
        
        fmt.Printf("✅ Theme '%s' installed successfully!\n", args[0])
    },
        
}

var themeUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Set active theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.Load("config.toml")
		themeManager := theme.NewThemeManager(cfg)
		themeManager.LoadThemes()

		if err := themeManager.SetActiveTheme(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully set active theme to: %s\n", args[0])
		fmt.Println("Don't forget to update your config.toml file with:")
		fmt.Printf("theme = \"%s\"\n", args[0])
	},
}

var themeCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new theme",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template, _ := cmd.Flags().GetString("template")
		cfg, _ := config.Load("config.toml")
		themeManager := theme.NewThemeManager(cfg)

		fmt.Printf("Creating theme '%s' with template '%s'\n", args[0], template)

		if err := themeManager.CreateTheme(args[0], template); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Theme '%s' created successfully!\n", args[0])
		fmt.Printf("Theme files are located in: themes/%s/\n", args[0])
		fmt.Println("")
		fmt.Println("Next steps:")
		fmt.Printf("1. Edit themes/%s/theme.json to customize theme metadata\n", args[0])
		fmt.Printf("2. Modify templates in themes/%s/layouts/\n", args[0])
		fmt.Printf("3. Add styles to themes/%s/static/css/style.css\n", args[0])
		fmt.Printf("4. Use the theme with: vango theme use %s\n", args[0])
	},
}



func init() {
	rootCmd.AddCommand(themeCmd)
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeInstallCmd)
	themeCmd.AddCommand(themeUseCmd)
	themeCmd.AddCommand(themeCreateCmd)

	themeCreateCmd.Flags().StringP("template", "t", "basic", "Theme template to use (basic, blog, portfolio, docs)")
}
