package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var cliOptions struct {
	filename string
	database string
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pg_explain",
		Short: "read explain in json format from stdin",
		Long:  `read explain in json format from stdin`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				input, _ := io.ReadAll(os.Stdin)
				source := Source{sourceType: SOURCE_STDIN, input: string(input)}
				RunProgram(source)
				return
			}

			if cliOptions.filename != "" {
				source := Source{sourceType: SOURCE_FILE, fileName: cliOptions.filename}
				RunProgram(source)
			}
		},
	}

	rootCmd.
		Flags().
		StringVarP(&cliOptions.filename, "filename", "f", "", "filename of sql file")

	rootCmd.
		Flags().
		StringVarP(&cliOptions.database, "database", "d", "", "database")

	rootCmd.Execute()
}

var databaseUrl = "postgres://postgres:postgres@localhost:5432/"

func ExecuteExplain(query string, settings []Setting) string {
	pgConn := Connection{
		databaseUrl: databaseUrl + cliOptions.database,
	}
	pgConn.Connect()
	defer pgConn.Close()
	for _, setting := range settings {
		pgConn.SetSetting(setting)
	}
	return pgConn.ExecuteExplain(query)
}
