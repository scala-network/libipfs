package ipfs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRepo(t *testing.T) {
	dir := t.TempDir()

	err := CreateRepo(dir, 4001, []string{})
	require.NoError(t, err)

	// Calling CreateRepo again on the same initialized path should be idempotent.
	err = CreateRepo(dir, 4001, []string{})
	require.NoError(t, err)

	// Invalid port.
	err = CreateRepo(t.TempDir(), 0, []string{})
	require.Error(t, err)
}

func TestCreateNode(t *testing.T) {
	dir := t.TempDir()
	err := CreateRepo(dir, 4001, []string{})
	require.NoError(t, err)

	ctx := context.Background()
	node, api, err := CreateNode(ctx, dir)
	require.NoError(t, err)
	require.NotNil(t, node)
	require.NotNil(t, api)
	defer node.Close()

	// Invalid repo path should fail.
	_, _, err = CreateNode(ctx, t.TempDir()+"/does-not-exist")
	require.Error(t, err)
}

func TestGarbageCollect(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, CreateRepo(dir, 4001, []string{}))

	ctx := context.Background()
	node, _, err := CreateNode(ctx, dir)
	require.NoError(t, err)
	defer node.Close()

	require.NoError(t, GarbageCollect(node))
}

func TestCloseNode(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, CreateRepo(dir, 4001, []string{}))

	ctx := context.Background()
	node, _, err := CreateNode(ctx, dir)
	require.NoError(t, err)

	require.NoError(t, CloseNode(node))
}
