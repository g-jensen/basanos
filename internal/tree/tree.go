package tree

import (
	"path/filepath"

	"basanos/internal/fs"
	"basanos/internal/spec"
)

type SpecTree struct {
	Path     string
	Context  *spec.Context
	Children []*SpecTree
}

func LoadContext(filesystem fs.FileSystem, dirPath string) (*spec.Context, error) {
	data, err := filesystem.ReadFile(filepath.Join(dirPath, "context.yaml"))
	if err != nil {
		return nil, err
	}
	return spec.ParseContext(data)
}

func LoadSpecTree(filesystem fs.FileSystem, rootPath string) (*SpecTree, error) {
	ctx, err := LoadContext(filesystem, rootPath)
	if err != nil {
		return nil, err
	}

	tree := &SpecTree{Path: rootPath, Context: ctx}

	entries, err := filesystem.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		childPath := filepath.Join(rootPath, entry.Name())
		contextFile := filepath.Join(childPath, "context.yaml")
		if _, err := filesystem.Stat(contextFile); err != nil {
			continue
		}
		child, err := LoadSpecTree(filesystem, childPath)
		if err != nil {
			return nil, err
		}
		tree.Children = append(tree.Children, child)
	}

	return tree, nil
}
