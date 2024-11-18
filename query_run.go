package main

import (
	"example/display-nodes/sqlsplit"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/x/ansi"
)

var extension string = ".pgex"

type QueryRun struct {
	query            string
	result           string
	originalFilename string
	pgexPointer      string
	settings         []Setting
}

var defaultPgexDir = "_pgex"

func CreatePgexDir() string {
	workingDir, _ := os.Getwd()
	dirPath := filepath.Join(workingDir, defaultPgexDir)
	os.MkdirAll(dirPath, 0755)
	return dirPath
}

func (q QueryRun) previousQueryRun() QueryRun {
	pgexFiles := getQueryRunEntries()
	var currentIndex int
	for i, pgexFile := range pgexFiles {
		if strings.Contains(pgexFile, q.pgexPointer) {
			currentIndex = i
		}
	}

	if currentIndex-1 >= 0 {
		return loadQueryRun(pgexFiles[currentIndex-1])
	} else {
		return q
	}
}

func (q QueryRun) nextQueryRun() QueryRun {
	pgexFiles := getQueryRunEntries()
	var currentIndex int
	for i, pgexFile := range pgexFiles {
		if strings.Contains(pgexFile, q.pgexPointer) {
			currentIndex = i
		}
	}

	if currentIndex+1 < len(pgexFiles) {
		return loadQueryRun(pgexFiles[currentIndex+1])
	} else {
		return q
	}
}

func loadQueryRun(pgexFile string) QueryRun {
	body, err := os.ReadFile(pgexFile)
	if err != nil {
		log.Fatal(err)
	}

	contents := string(body)
	if !strings.Contains(contents, sqlDivider) {
		panic("wrong pgex format no settings above divider")
	}

	if !strings.Contains(contents, explainDivider) {
		panic("wrong pgex format no sql above divider")
	}

	settingsAbove := strings.Split(contents, sqlDivider)

	settingsContent := settingsAbove[0]

	settingsStrings := strings.Split(settingsContent, "\n")

	settings := make([]Setting, 0, len(settingsStrings))
	for _, settingStr := range settingsStrings {
		if ansi.StringWidth(strings.Trim(settingStr, " ")) > 0 {
			settings = append(settings, SettingUnmarshal(settingStr))
		}
	}

	sqlAbove := strings.Split(settingsAbove[1], explainDivider)

	sql := sqlAbove[0]
	plan := sqlAbove[1]

	_, file := path.Split(pgexFile)

	return QueryRun{query: sql, result: plan, pgexPointer: file, settings: settings}
}

func getQueryRunEntries() []string {
	dirEntries, err := os.ReadDir("_pgex/")
	if err != nil {
		log.Fatal(err)
	}

	pgexFiles := make([]string, 0, len(dirEntries))
	wd, _ := os.Getwd()
	for _, d := range dirEntries {
		pgexFile := regexp.MustCompile(`[0-9]{14}_.*\.pgex`)
		if pgexFile.Match([]byte(d.Name())) {
			result := filepath.Join(wd, "_pgex/", d.Name())
			pgexFiles = append(pgexFiles, result)
		}
	}

	return pgexFiles
}

func NewQueryRun(filename string) QueryRun {
	if _, err := os.Stat(filename); err == nil {
		body, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}

		sqls := sqlsplit.Split(string(body))

		if len(sqls) != 1 {
			log.Fatal("too many sql statements in provided file")
		} else {
			return QueryRun{
				query:            sqls[0],
				originalFilename: filename,
			}
		}
	} else {
		log.Fatal(err)
	}
	return QueryRun{}
}

func (q *QueryRun) SetResult(result string) {
	q.result = result
}

func (q *QueryRun) WritePgexFile(pgexDir string) {
	fileName := q.pgexFilename()
	fullFilePath := filepath.Join(pgexDir, fileName)
	contentBytes := []byte(q.pgexFileContent())

	os.WriteFile(fullFilePath, contentBytes, 0666)
	q.pgexPointer = fileName
}

func (q QueryRun) DisplayName() string {
	_, file := path.Split(q.originalFilename)
	return file
}

func (q QueryRun) pgexFilename() string {
	user, _ := user.Current()
	filePath := strings.Replace(q.originalFilename, "~", user.HomeDir, 1)

	_, file := path.Split(filePath)
	name := strings.Split(file, ".")[0]

	formattedNow := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s_%s%s", formattedNow, name, extension)
}

var explainDivider = "---------------- SQL ABOVE / EXPLAIN JSON BELOW ----------------"
var sqlDivider = "---------------- SETTINGS ABOVE / SQL BELOW ----------------"

func (q QueryRun) pgexFileContent() string {
	var buf strings.Builder
	for _, setting := range q.settings {
		buf.WriteString(setting.Marshal())
		buf.WriteString("\n")
	}
	buf.WriteString("\n\n")
	buf.WriteString(sqlDivider)
	buf.WriteString("\n\n")
	buf.WriteString(q.query)
	buf.WriteString("\n\n")
	buf.WriteString(explainDivider)
	buf.WriteString("\n\n")
	buf.WriteString(q.result)
	buf.WriteString("\n\n")
	return buf.String()
}

func (q QueryRun) WithExplain() string {
	explainSegment := `explain (
		format json
	) `

	return explainSegment + q.query
}

func (q QueryRun) WithExplainAnalyze() string {
	explainSegment := `explain (
		settings,
		format json,
		buffers,
		analyze
	) `

	return explainSegment + q.query
}
