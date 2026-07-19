package util

import (
	"os"
	"project/config"
)

func Init() error {
	return os.MkdirAll(TmpDir.Static, config.DirFileMode)
}
