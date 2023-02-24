package repository

import "path/filepath"

type repo struct {
	basePath string
}

type GetPathInterface interface {
	GetPath(...string) string
}

func New(basePath string) interface{} {
	return repo{
		basePath: basePath,
	}
}

func (r repo) GetPath(rel ...string) string {
	return filepath.Join(append([]string{r.basePath}, rel...)...)
}
