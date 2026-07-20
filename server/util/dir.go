package util

import (
	"os"
	"path/filepath"
	"project/config"
	"strings"
)

func (f *File) ReadDir(fn FnSelector) ([]string, error) {
	el, err := os.ReadDir(f.Static)
	if err != nil {
		return nil, err
	}

	var fl []string
	for _, e := range el {
		path := e.Name()
		if e.IsDir() {
			continue
		}
		if fn == nil || fn(path) {
			path = filepath.Join(f.Base, path)
			fl = append(fl, path)
		}
	}
	return fl, nil
}

func (f *File) Walk(fn FnSelector) ([]string, error) {

	var fl []string

	err := filepath.WalkDir(f.Static, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if fn == nil || fn(path) {
			path = strings.TrimPrefix(path, config.StaticDir)
			fl = append(fl, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fl, nil
}

func (f *File) Mkdir() error {
	if f.hasDir {
		return nil
	}
	err := os.MkdirAll(filepath.Dir(f.Static), config.DirFileMode)
	if err == nil {
		f.hasDir = true
	}
	return err
}
