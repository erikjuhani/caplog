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
		file       string
		expectsErr bool
	}{
		{
			file:       "",
			expectsErr: true,
		},
		{
			file:       dir,
			expectsErr: true,
		},
		{
			file:       fmt.Sprintf("%s/%s", dir, "12345_file.log"),
			expectsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			os.Create(tt.file)
			actual := CommitSingleFile(tt.file, "log: entry")
			if (actual != nil) != tt.expectsErr {
				t.Fatalf("expects error %t did not match actual %v", tt.expectsErr, actual)
			}
		})
	}
}
