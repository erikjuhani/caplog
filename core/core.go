package core

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/erikjuhani/caplog/config"
	"github.com/erikjuhani/caplog/git"
)

type Log struct {
	Date time.Time
	Data []string
}

const (
	timeFormat     = "15:04"
	timeFileFormat = "2006-01-02T15:04:05"
)

func CreateLog(date time.Time, data string) Log {
	l := Log{Date: date}

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {
		l.Data = append(l.Data, scanner.Text())
	}

	return l
}

func WriteLog(log Log) error {
	if len(log.Data) == 0 {
		return errors.New("no data provided")
	}

	path := config.Get(config.GitLocalRepositoryKey)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}

	filename := generateFilename(log)
	filepath := fmt.Sprintf("%s/%s", path, filename)
	formattedLog := formatLog(log)

	if err := os.WriteFile(filepath, []byte(formattedLog), 0644); err != nil {
		return err
	}

	return git.CommitSingleFile(filepath, formattedLog)
}

func openInEditor(filename string) error {
	executable, err := exec.LookPath(config.Get(config.EditorKey))
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

func formatLog(log Log) string {
	ts := log.Date.Format(timeFormat)
	if len(log.Data) == 1 {
		return fmt.Sprintf("%s\t%s", ts, log.Data[0])
	}
	return fmt.Sprintf("%s\t%s\n%s", ts, log.Data[0], strings.Join(log.Data[1:], "\n"))
}

func generateFilename(log Log) string {
	h := sha1.New()
	h.Write([]byte(log.Date.String()))

	hash := fmt.Sprintf("%x", h.Sum(nil))[0:7]
	date := log.Date.Format(timeFileFormat)

	return fmt.Sprintf("%s_%s.log", date, hash)
}
