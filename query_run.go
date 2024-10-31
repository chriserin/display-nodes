package main

import (
	"example/display-nodes/sqlsplit"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var extension string = ".pgex"

type QueryRun struct {
	query            string
	result           string
	originalFilename string
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

func (q QueryRun) WritePgexFile(pgexDir string) {
	fileName := q.pgexFilename()
	fullFilePath := filepath.Join(pgexDir, fileName)
	contentBytes := []byte(q.pgexFileContent())

	os.WriteFile(fullFilePath, contentBytes, 0666)
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
	return fmt.Sprintf("%s_%s%s", name, formattedNow, extension)
}

func (q QueryRun) pgexFileContent() string {
	content := q.query + "\n\n--------------\n\n" + q.result + "\n"
	return content
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
