package app

import (
	"bharvest.io/elmon/utils"
)

type Config struct {
	General struct {
		APIListenPort int  `toml:"api_listen_port"`
		Period        uint `toml:"period"`
	} `toml:"general"`
	Tg utils.TgConfig `toml:"tg"`
	EL struct {
		JsonRPCs string `toml:"json_rpcs"`
		//PeerThreshold         uint64 `toml:"peer_threshold"`
		//TxpoolQueuedThreshold uint64 `toml:"txpool_queued_threshold"`
		JsonRPCList []string
	} `toml:"el"`
}

type BaseApp struct {
	cfg *Config
}
