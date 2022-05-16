package git

import (
	"fmt"
	"os"
	"testing"
)

func testRepo() (string, func()) {
	dir, _ := os.MkdirTemp("", "caplog")
	runGitCommand("init", "-q", dir)
	return dir, func() { os.RemoveAll(dir) }
}

func TestIsGitRepository(t *testing.T) {
	dir, cleanup := testRepo()
	defer cleanup()

	tests := []struct {
		path     string
		expected bool
	}{
		{
			path:     fmt.Sprintf("%s/not-a-dir", dir),
			expected: false,
		},
		{
			path:     dir,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if tt.expected != isGitRepository(tt.path) {
				t.Fatal("expected did not match output when deciding if directory is a git repository")
			}
		})
	}
}

func TestCommitSingleFile(t *testing.T) {
	dir, _ := os.MkdirTemp("", "caplog")
	defer os.RemoveAll(dir)

	tests := []struct {
		file     string
		expected error
	}{
		{
			file:     fmt.Sprintf("%s/%s", dir, "12345_file.log"),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			os.Create(tt.file)
			if tt.expected != CommitSingleFile(tt.file, "log: entry") {
				t.Fatal("expected did not match the actual error output")
			}
		})
	}
}
