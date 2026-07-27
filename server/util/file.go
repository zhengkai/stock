// Package util 工具箱
package util

import (
	"crypto/sha256"
	"fmt"
	"os"
	"project/config"
	"project/zj"

	"golang.org/x/sys/unix"
)

const xattrHashKey = `user.sha256hash`

var (
	TmpDir = NewFile(`tmp`)
)

type File struct {
	Base   string
	Static string
	hash   *[sha256.Size]byte
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

func (f *File) WriteWithHash(data []byte) error {

	h := sha256.Sum256(data)
	if f.checkHash(h) {
		zj.J(`hash same, skip write`, f.Static)
		return nil
	}
	zj.F(`hash %x`, h[:8])

	return writeBinWithHash(f.Static, data, h)
}

func (f *File) checkHash(h [sha256.Size]byte) bool {
	ph := f.getHash()
	if ph == nil {
		return false
	}
	return *ph == h
}

func (f *File) getHash() *[sha256.Size]byte {
	if f.hash != nil {
		return f.hash
	}

	buf := make([]byte, sha256.Size)
	size, err := unix.Getxattr(f.Static, xattrHashKey, buf)
	if size != sha256.Size {
		// err = fmt.Errorf(`invalid hash size: %d`, size)
		return nil
	}
	if err != nil {
		return nil
	}
	f.hasDir = true
	hash := [sha256.Size]byte{}
	copy(hash[:], buf)
	f.hash = &hash
	return f.hash
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

func (f *File) Append(s string) error {
	h, err := os.OpenFile(f.Static, os.O_WRONLY|os.O_CREATE|os.O_APPEND, config.FileMode)
	if err != nil {
		return err
	}
	defer h.Close()

	_, err = h.WriteString(s)
	return err
}

func (f *File) AppendF(format string, arg ...any) error {
	h, err := os.OpenFile(f.Static, os.O_WRONLY|os.O_CREATE|os.O_APPEND, config.FileMode)
	if err != nil {
		return err
	}
	defer h.Close()

	_, err = fmt.Fprintf(h, format, arg...)
	return err
}

func (f *File) Remove() error {
	return os.Remove(f.Static)
}

func writeBinWithHash(file string, content []byte, hash [sha256.Size]byte) (err error) {

	f, err := TmpFile()
	if err != nil {
		return
	}

	f.Chmod(config.FileMode)
	tmpName := f.Name()

	if _, err = f.Write(content); err != nil {
		zj.W(`write bin fail`, file, len(content))
		os.Remove(tmpName)
		f.Close()
		return
	}

	fd := int(f.Fd())
	unix.Fsetxattr(fd, xattrHashKey, hash[:], 0) // 失败也无所谓，就当没 hash

	return os.Rename(tmpName, file)
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
			zj.W(`write bin fail`, file, len(ab))
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
