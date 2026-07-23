package app

import (
	"project/pb"
	"project/util"
)

func (app *App) loadConfig() error {

	app.Config = &pb.Config{}

	err := util.NewFile(`config.json`).ReadJSON(app.Config)
	if err != nil {
		return err
	}

	return nil
}
