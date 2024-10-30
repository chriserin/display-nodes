package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var filename string

	var rootCmd = &cobra.Command{
		Use:   "pg_explain",
		Short: "read explain in json format from stdin",
		Long:  `read explain in json format from stdin`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				input, _ := io.ReadAll(os.Stdin)
				explainPlan := Convert(string(input))
				RunProgram(explainPlan)
				return
			}
		},
	}

	rootCmd.
		Flags().
		StringVarP(&filename, "filename", "f", "", "filename of sql file")

	rootCmd.Execute()
}
