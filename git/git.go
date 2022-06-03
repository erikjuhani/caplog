package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	ErrNoPathProvided     = errors.New("no path provided")
	ErrExecNotFoundInPath = errors.New("executable not found in path")
)

func hasGitRemote() bool {
	if err := runGitCommand("ls-remote"); err != nil {
		return false
	}

	return true
}

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
	if len(path) == 0 {
		return ErrNoPathProvided
	}

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

	if err := runGitCommand("commit", "-m", msg, path); err != nil {
		return err
	}

	if hasGitRemote() {
		// TODO: adjust with flag or configuration
		// TODO: think about detached process ordering and composition
		if err := runGitCommand("pull", "--rebase=merges"); err != nil {
			return err
		}

		if err := runDetachedGitCommand("push", "--force-with-lease"); err != nil {
			return err
		}
	}

	return nil
}

func commandExists(command string) bool {
	if _, err := exec.LookPath(command); err == nil {
		return true
	}

	return false
}

func execCommand(cmd string, args ...string) (*exec.Cmd, error) {
	if !commandExists(cmd) {
		return nil, fmt.Errorf("%s %w", cmd, ErrExecNotFoundInPath)
	}

	return exec.Command(cmd, args...), nil
}

func runDetachedGitCommand(args ...string) error {
	git, err := execCommand("git", args...)
	if err != nil {
		return err
	}

	if err := git.Start(); err != nil {
		return err
	}

	go func() {
		git.Process.Release()
	}()

	return nil
}

func runGitCommand(args ...string) error {
	git, err := execCommand("git", args...)
	if err != nil {
		return err
	}

	return git.Run()
}
