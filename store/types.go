package store

type GlobalStateType struct {
	// Host -> ELType
	ELs map[string]*ELType `json:"EL"`
}

type ELType struct {
	Status         bool   `json:"status"`
	Sync           bool   `json:"sync"`
	CurrentHeight  uint64 `json:"current_height"`
	Peers          uint64 `json:"peers"`
	TxpoolQueued   int    `json:"txpool_queued"`
}
