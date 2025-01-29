package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "testdir")
	assert.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	assert.True(t, IsDir(tempDir), "Expected true for existing directory")

	nonExistentPath := filepath.Join(tempDir, "does_not_exist")
	assert.False(t, IsDir(nonExistentPath), "Expected false for non-existent path")

	tempFile := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(tempFile, []byte("test content"), 0644)
	assert.NoError(t, err, "Failed to create test file")
	assert.False(t, IsDir(tempFile), "Expected false for a file path")

	specialDir := filepath.Join(tempDir, "dir_with_#@!")
	err = os.Mkdir(specialDir, 0755)
	assert.NoError(t, err, "Failed to create special character directory")
	assert.True(t, IsDir(specialDir), "Expected true for special character directory")

	symlinkDir := filepath.Join(tempDir, "symlink_to_dir")
	err = os.Symlink(specialDir, symlinkDir)
	assert.NoError(t, err, "Failed to create symlink to directory")
	assert.True(t, IsDir(symlinkDir), "Expected true for symlink to a directory")

	symlinkFile := filepath.Join(tempDir, "symlink_to_file")
	err = os.Symlink(tempFile, symlinkFile)
	assert.NoError(t, err, "Failed to create symlink to file")
	assert.False(t, IsDir(symlinkFile), "Expected false for symlink to a file")
}
