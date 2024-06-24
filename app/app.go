package app

import (
	"context"
	"sync"

	"bharvest.io/beramon/utils"
)

func NewBaseApp(cfg *Config) *BaseApp {
	return &BaseApp{
		cfg: cfg,
	}
}

func (app *BaseApp) Run(ctx context.Context) {
	appCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(app.cfg.EL.JsonRPCList))

	for _, jsonrpc := range app.cfg.EL.JsonRPCList {
		go func(jsonrpc string) {
			// Check EL
			defer wg.Done()

			err := app.checkEL(appCtx, jsonrpc)
			if err != nil {
				utils.SendTg(err.Error())
				utils.Error(err, true)
				return
			}
		}(jsonrpc)
	}

	wg.Wait()

	return
}
