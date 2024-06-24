package el

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"bharvest.io/beramon/utils"
)

func jsonRPCQuery(ctx context.Context, host string, method string, params []string) (*JSONRPCResponse, error) {
	payload, err := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		host,
		bytes.NewBuffer(payload),
	)
	req.Header.Set("Content-Type", "application/json")

	msg := fmt.Sprintf("Querying %s", method)
	utils.Info(msg)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := JSONRPCResponse{}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetSyncStatus(ctx context.Context) (bool, error) {
	resp, err := jsonRPCQuery(ctx, c.host, "eth_syncing", []string{})
	if err != nil {
		return true, err
	}

	if string(resp.Result) != "false" {
		return true, nil
	}
	return false, nil
}

func (c *Client) GetLatestBlock(ctx context.Context) (uint64, error) {
	resp, err := jsonRPCQuery(ctx, c.host, "eth_blockNumber", []string{})
	if err != nil {
		return 0, err
	}

	height := string(resp.Result[3:len(resp.Result)-1])

	result, err := strconv.ParseUint(height, 16, 64)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (c *Client) GetPeerCnt(ctx context.Context) (uint64, error) {
	resp, err := jsonRPCQuery(ctx, c.host, "net_peerCount", []string{})
	if err != nil {
		return 0, err
	}

	peers := string(resp.Result[3:len(resp.Result)-1])

	result, err := strconv.ParseUint(peers, 16, 64)
	if err != nil {
		panic(err)
	}

	return result, nil
}

func (c *Client) GetTxQueuedCnt(ctx context.Context) (int, error) {
	resp, err := jsonRPCQuery(ctx, c.host, "txpool_content", []string{})
	if err != nil {
		return 0, err
	}

	result := TxpoolContentResponse{}
	err = json.Unmarshal(resp.Result, &result)
	if err != nil {
		return 0, err
	}

	return len(result.Queued), nil
}
