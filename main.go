package main

import (
	"io"
	"os"
	"path/filepath"

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
				source := Source{sourceType: SOURCE_STDIN, input: string(input)}
				RunProgram(source)
				return
			}

			if filename != "" {
				source := Source{sourceType: SOURCE_FILE, fileName: filename}
				RunProgram(source)
			}
		},
	}

	rootCmd.
		Flags().
		StringVarP(&filename, "filename", "f", "", "filename of sql file")

	rootCmd.Execute()
}

var databaseUrl = "postgres://postgres:postgres@localhost:5432/galaxy_dev"

func ExecuteExplain(query string) string {
	pgConn := Connection{
		databaseUrl: databaseUrl,
	}
	pgConn.Connect()
	defer pgConn.Close()
	return pgConn.ExecuteExplain(query)
}

func CreatePgexDir() string {
	workingDir, _ := os.Getwd()
	dirPath := filepath.Join(workingDir, "_pgex")
	os.MkdirAll(dirPath, 0755)
	return dirPath
}
