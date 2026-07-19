package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	EnvItickToken = `STOCK_ITICK_TOKEN`
)

func init() {

	Dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	list := map[string]*string{
		`STOCK_WEB`:               &Web,
		`STOCK_DIR`:               &StaticDir,
		`STOCK_VAPID_PUBLIC_KEY`:  &VapidPublicKey,
		`STOCK_VAPID_PRIVATE_KEY`: &VapidPrivateKey,
		EnvItickToken:             &ItickToken,
	}
	for k, v := range list {
		s := os.Getenv(k)
		if len(s) > 1 {
			*v = s
		}
	}

	StaticDir = strings.TrimRight(StaticDir, `/`) + `/`
}
