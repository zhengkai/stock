// Package itick 连接 itick api
package itick

import (
	"fmt"
	"regexp"
	"project/config"
)

var (
	regexpToken = regexp.MustCompile(`^[a-f0-9]{64}$`)
)

func Init() error {

	s := config.ItickToken

	if !regexpToken.MatchString(s) {
		if s == `` {
			return fmt.Errorf(`no env %s`, config.EnvItickToken)
		}
		return fmt.Errorf(`invalid env %s: %s`, config.EnvItickToken, s)
	}

	theClient = &Client{
		token: s,
	}
	return nil
}
