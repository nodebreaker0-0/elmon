package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"bharvest.io/elmon/app"
	"bharvest.io/elmon/server"
	"bharvest.io/elmon/store"
	"bharvest.io/elmon/utils"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	ctx := context.Background()

	f, err := os.ReadFile("config.toml")
	if err != nil {
		utils.Error(err, true)
		return
	}
	cfg := app.Config{}
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		utils.Error(err, true)
		return
	}

	tgTitle := fmt.Sprintf("🤖 Monad elmon 🤖")
	utils.SetTg(cfg.Tg.Enable, tgTitle, cfg.Tg.Token, cfg.Tg.ChatID)

	// Init JsonRPCList & memory store
	cfg.EL.JsonRPCList = strings.Split(cfg.EL.JsonRPCs, ",")
	for i, rpc := range cfg.EL.JsonRPCList {
		r := strings.TrimSpace(rpc)
		cfg.EL.JsonRPCList[i] = r
		store.GlobalState.ELs[r] = &store.ELType{
			Status: true,
		}
	}

	go server.Run(cfg.General.APIListenPort)

	baseapp := app.NewBaseApp(&cfg)
	for {
		baseapp.Run(ctx)
		time.Sleep(time.Duration(cfg.General.Period) * time.Minute)
	}
}
