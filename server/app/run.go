package app

func Run() {
	a := &App{}
	a.run()
}

func (a *App) run() error {
	err := a.loadConfig()
	if err != nil {
		return err
	}

	// a.checkAlert()

	// fmt.Println(tc.StockURL([]string{`600519`}))

	return nil
}
