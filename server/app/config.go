package app

import (
	"server/pb"
	"server/util"
)

func (a *App) loadConfig() error {

	a.Config = &pb.Config{}

	err := util.NewFile(`config.json`).ReadJSON(a.Config)
	if err != nil {
		return err
	}

	return nil
}
