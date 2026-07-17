// Package util 工具箱
package util

import (
	"os"
	"server/config"
)

type File struct {
	Base   string
	Static string
}

func NewFile(base string) *File {
	return &File{
		Base:   base,
		Static: config.StaticDir + base,
	}
}

func (f *File) String() string {
	return f.Base
}

func (f *File) Read() ([]byte, error) {
	ab, err := os.ReadFile(f.Static)
	if err != nil {
		return nil, err
	}
	return ab, nil
}

func (f *File) Write(data []byte) error {
	return nil
}
