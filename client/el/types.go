package el

import (
	"encoding/json"
)

type Client struct {
	host string
}

type (
	JSONRPCRequest struct {
		JSONRPC string   `json:"jsonrpc"`
		Method  string   `json:"method"`
		Params  []string `json:"params"`
		ID      int      `json:"id"`
	}

	JSONRPCResponse struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      int             `json:"id"`
		Result  json.RawMessage `json:"result"`
		Error   *RPCError       `json:"error"`
	}
	RPCError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
)

/*
type TxpoolContentResponse struct {
	Pending map[string]any `json:"pending"`
	Queued map[string]any `json:"queued"`
}
*/
