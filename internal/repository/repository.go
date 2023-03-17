package repository

import "path/filepath"

type repo struct {
	basePath string
}

var _ Interface = repo{}

type Interface interface {
	GetPath(...string) string
}

func New(basePath string) Interface {
	return repo{
		basePath: basePath,
	}
}

func (r repo) GetPath(rel ...string) string {
	return filepath.Join(append([]string{r.basePath}, rel...)...)
}
