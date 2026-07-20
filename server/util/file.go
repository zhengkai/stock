// Package util 工具箱
package util

import (
	"fmt"
	"os"
	"project/config"
)

var (
	TmpDir = NewFile(`tmp`)
)

type File struct {
	Base   string
	Static string
	hasDir bool
}

func NewFile(base string) *File {
	return &File{
		Base:   base,
		Static: config.StaticDir + base,
	}
}

func NewFileF(format string, a ...any) *File {
	return NewFile(fmt.Sprintf(format, a...))
}

func (f *File) String() string {
	return f.Base
}

func (f *File) Read() ([]byte, error) {
	ab, err := os.ReadFile(f.Static)
	if err != nil {
		return nil, err
	}
	f.hasDir = true
	return ab, nil
}

func (f *File) Write(data []byte) error {
	f.Mkdir()
	return writeBin(f.Static, data)
}

func (f *File) IsExists() bool {
	_, err := os.Stat(f.Static)
	if err == nil {
		return true
	}
	exists := !os.IsNotExist(err)
	if exists {
		f.hasDir = true
	}
	return exists
}

func (f *File) Remove() error {
	return os.Remove(f.Static)
}

func writeBin(file string, li ...[]byte) (err error) {

	f, err := TmpFile()
	if err != nil {
		return
	}

	f.Chmod(config.FileMode)
	tmpName := f.Name()

	for _, ab := range li {
		if _, err = f.Write(ab); err != nil {
			fmt.Println(`write bin fail`, file, len(ab))
			os.Remove(tmpName)
			f.Close()
			return
		}
	}
	f.Close()

	return os.Rename(tmpName, file)
}

func TmpFile() (*os.File, error) {
	return os.CreateTemp(TmpDir.Static, `tmp-go-*`)
}
