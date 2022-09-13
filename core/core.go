package core

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/git"
)

type Meta struct {
	Date time.Time
	Page string
}

func (m Meta) Location() string {
	loc := config.WorkspacePath()

	if len(m.Page) == 0 {
		return loc
	}

	return fmt.Sprintf("%s/%s", loc, m.Page)
}

func (m Meta) String() string {
	// TODO: Refactor
	if len(m.Page) == 0 {
		return fmt.Sprintf(`---
date: %s
---
`, m.Date.Format(metaTimeLayout))
	}

	return fmt.Sprintf(`---
date: %s

page: %s
---
`, m.Date.Format(metaTimeLayout), m.Page)
}

type Log struct {
	Meta
	Data []string
}

const (
	// TODO: Custom time layout
	// WEEKDAY, MONTH DAY, YEAR
	metaTimeLayout = "Monday, January 2, 2006"
	timeFormat     = "15:04"
	timeFileFormat = "02-01-2006"
)

func NewLog(meta Meta, data string, tags []string) Log {
	l := Log{Meta: meta}

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		l.Data = append(l.Data, scanner.Text())
	}

	if len(tags) > 0 {
		l.Data = append(l.Data, fmt.Sprintf("\ntags: %s", strings.Join(tags, ", ")))
	}

	return l
}

func WriteLog(log Log) error {
	if len(log.Data) == 0 {
		return errors.New("no data provided")
	}

	loc := log.Location()

	if _, err := os.Stat(loc); os.IsNotExist(err) {
		if err := os.MkdirAll(loc, os.ModePerm); err != nil {
			return err
		}
	}

	filename := logFilename(log)
	filepath := fmt.Sprintf("%s/%s", loc, filename)

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
	}

	formattedLog := formatLog(log)

	if _, ok := err.(*os.PathError); ok {
		if err := os.WriteFile(filepath, []byte(fmt.Sprintf("%s\n%s", log.Meta.String(), formattedLog)), 0644); err != nil {
			return err
		}
		return git.CommitSingleFile(filepath, formattedLog)
	}

	if _, err := f.WriteString("\n" + formattedLog); err != nil {
		return err
	}

	return git.CommitSingleFile(filepath, formattedLog)
}

func openInEditor(filename string) error {
	executable, err := exec.LookPath(config.Config.Editor)
	if err != nil {
		return err
	}

	command := exec.Command(executable, filename)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func CaptureEditorInput() ([]byte, error) {
	var input []byte

	file, err := os.CreateTemp(os.TempDir(), "caplog")
	if err != nil {
		return input, err
	}

	filename := file.Name()
	defer os.Remove(filename)

	if err := file.Close(); err != nil {
		return input, err
	}

	if err := openInEditor(filename); err != nil {
		return input, err
	}

	return os.ReadFile(filename)
}

// NOTE: Mutates the given slice values with the prefix
func addPrefix(prefix string, s []string) []string {
	for i, v := range s {
		s[i] = prefix + v
	}

	return s
}

func formatLog(log Log) string {
	if len(log.Data) == 0 {
		return ""
	}

	ts := log.Date.Format(timeFormat)
	if len(log.Data) == 1 {
		return fmt.Sprintf("%s\t%s\n", ts, log.Data[0])
	}

	return fmt.Sprintf(
		"%s\t%s\n%s\n",
		ts,
		log.Data[0],
		strings.Join(addPrefix("	", log.Data[1:]), "\n"),
	)
}

func logFilename(log Log) string {
	return fmt.Sprintf("%s.log.md", log.Date.Format(timeFileFormat))
}
