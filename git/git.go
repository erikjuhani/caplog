package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func isGitRepository(path string) bool {
	// Check if git directory exists
	if _, err := os.Stat(fmt.Sprintf("%s/.git", path)); os.IsNotExist(err) {
		return false
	}

	// TODO: think of a better command to validate git repository
	if err := runGitCommand("-C", path, "rev-parse"); err != nil {
		return false
	}

	return true
}

func CommitSingleFile(path string, msg string) error {
	dirpath := filepath.Dir(path)
	if !isGitRepository(dirpath) {
		if err := runGitCommand("init", "-q", "-b", "trunk", dirpath); err != nil {
			return err
		}
	}

	if err := os.Chdir(dirpath); err != nil {
		return err
	}

	if err := runGitCommand("add", path); err != nil {
		return err
	}

	return runGitCommand("commit", "-m", msg, path)
}

func commandExists(command string) bool {
	if _, err := exec.LookPath(command); err == nil {
		return true
	}

	return false
}

func runGitCommand(args ...string) error {
	if !commandExists("git") {
		return errors.New("git executable not found in path")
	}

	git := exec.Command("git", args...)
	git.Stdin = os.Stdin
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr

	return git.Run()
}
