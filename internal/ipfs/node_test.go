package ipfs

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testRepoPath = "./test_repo"
	testBindPort = 4001
)

func TestCreateRepo(t *testing.T) {
	// Clean up any previous test repo
	os.RemoveAll(testRepoPath)

	// Test creating a new repo
	err := CreateRepo(testRepoPath, testBindPort, []string{})
	require.NoError(t, err)

	// Check if the repo directory was created
	_, err = os.Stat(testRepoPath)
	require.NoError(t, err)

	// Test creating a repo with invalid path
	err = CreateRepo("", testBindPort, []string{})
	require.Error(t, err)

	// Test creating a repo with invalid port
	err = CreateRepo(testRepoPath, 0, []string{})
	require.Error(t, err)

	// Clean up
	os.RemoveAll(testRepoPath)
}

func TestCreateNode(t *testing.T) {
	// Clean up any previous test repo
	os.RemoveAll(testRepoPath)

	// Create a new repo
	err := CreateRepo(testRepoPath, testBindPort, []string{})
	require.NoError(t, err)

	// Create a new node
	ctx := context.Background()
	node, api, err := CreateNode(ctx, testRepoPath)
	require.NoError(t, err)
	require.NotNil(t, node)
	require.NotNil(t, api)

	// Test creating a node with an invalid repo path
	_, _, err = CreateNode(ctx, "./invalid_path")
	require.Error(t, err)

	// Clean up
	node.Close()
	os.RemoveAll(testRepoPath)
}

func TestGarbageCollect(t *testing.T) {
	// Clean up any previous test repo
	os.RemoveAll(testRepoPath)

	// Create a new repo
	err := CreateRepo(testRepoPath, testBindPort, []string{})
	require.NoError(t, err)

	// Create a new node
	ctx := context.Background()
	node, _, err := CreateNode(ctx, testRepoPath)
	require.NoError(t, err)

	// Test garbage collection
	err = GarbageCollect(node)
	require.NoError(t, err)

	// Clean up
	node.Close()
	os.RemoveAll(testRepoPath)
}

func TestCloseNode(t *testing.T) {
	// Clean up any previous test repo
	os.RemoveAll(testRepoPath)

	// Create a new repo
	err := CreateRepo(testRepoPath, testBindPort, []string{})
	require.NoError(t, err)

	// Create a new node
	ctx := context.Background()
	node, _, err := CreateNode(ctx, testRepoPath)
	require.NoError(t, err)

	// Test closing the node
	err = CloseNode(node)
	require.NoError(t, err)

	// Clean up
	os.RemoveAll(testRepoPath)
}
