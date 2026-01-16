package tree

import (
	"testing"

	memfs "basanos/internal/testutil/fs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadContext(t *testing.T) {
	mfs := memfs.NewMemoryFS()
	mfs.AddDir("/spec")
	mfs.AddFile("/spec/context.yaml", []byte(`name: "Test Context"`))

	ctx, err := LoadContext(mfs, "/spec")

	require.NoError(t, err)
	assert.Equal(t, "Test Context", ctx.Name)
}

func TestLoadSpecTree(t *testing.T) {
	mfs := memfs.NewMemoryFS()
	mfs.AddDir("/spec")
	mfs.AddDir("/spec/child")
	mfs.AddFile("/spec/context.yaml", []byte(`name: "Root"`))
	mfs.AddFile("/spec/child/context.yaml", []byte(`name: "Child"`))

	tree, err := LoadSpecTree(mfs, "/spec")

	require.NoError(t, err)
	assert.Equal(t, "Root", tree.Context.Name)
	require.Len(t, tree.Children, 1)
	assert.Equal(t, "Child", tree.Children[0].Context.Name)
}

func TestLoadSpecTree_TracksPath(t *testing.T) {
	mfs := memfs.NewMemoryFS()
	mfs.AddDir("/spec")
	mfs.AddDir("/spec/child")
	mfs.AddFile("/spec/context.yaml", []byte(`name: "Root"`))
	mfs.AddFile("/spec/child/context.yaml", []byte(`name: "Child"`))

	tree, err := LoadSpecTree(mfs, "/spec")

	require.NoError(t, err)
	assert.Equal(t, "/spec", tree.Path)
	assert.Equal(t, "/spec/child", tree.Children[0].Path)
}
