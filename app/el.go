package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"bharvest.io/elmon/client/el"
	"bharvest.io/elmon/store"
	"bharvest.io/elmon/utils"
)

func (app *BaseApp) checkEL(ctx context.Context, jsonrpc string) error {
	utils.Info(fmt.Sprintf("Start check EL for %s", jsonrpc))

	// Init EL status
	store.GlobalState.ELs[jsonrpc].Status = true

	now := time.Now()
	appCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(2)

	client, err := el.New(jsonrpc)
	if err != nil {
		return err
	}

	// Check sync status
	go func() {
		defer wg.Done()

		isSyncing, err := client.GetSyncStatus(appCtx)
		if err != nil {
			utils.SendTg(err.Error())
			utils.SendPd(err.Error())
			utils.SendSl(err.Error())
			utils.Error(err, true)
			return
		}
		store.GlobalState.ELs[jsonrpc].Sync = isSyncing

		if isSyncing {
			store.GlobalState.ELs[jsonrpc].Status = false

			msg := "Monad EL Node is syncing"
			utils.SendTg(msg)
			utils.SendPd(msg)
			utils.SendSl(msg)
			utils.Error(errors.New(msg), false)

			return
		}
	}()

	// Check Latest Block
	go func() {
		defer wg.Done()

		height, err := client.GetLatestBlock(appCtx)
		if err != nil {
			utils.SendTg(err.Error())
			utils.SendPd(err.Error())
			utils.SendSl(err.Error())
			utils.Error(err, true)
			return
		}

		if store.GlobalState.ELs[jsonrpc].CurrentHeight == height {
			store.GlobalState.ELs[jsonrpc].Status = false

			msg := "Monad Height is not increasing"
			utils.SendTg(msg)
			utils.SendPd(msg)
			utils.SendSl(msg)
			utils.Error(errors.New(msg), false)
		}
		store.GlobalState.ELs[jsonrpc].CurrentHeight = height
	}()

	/* Check Peer Count
	go func() {
		defer wg.Done()

		peers, err := client.GetPeerCnt(appCtx)
		if err != nil {
			utils.SendTg(err.Error())
			utils.SendPd(err.Error())
			utils.SendSl(err.Error())
			utils.Error(err, true)
			return
		}
		store.GlobalState.ELs[jsonrpc].Peers = peers

		if peers < app.cfg.EL.PeerThreshold {
			store.GlobalState.ELs[jsonrpc].Status = false

			msg := fmt.Sprintf("EL Node has low peers: %d", peers)
			utils.SendTg(msg)
			utils.SendPd(msg)
			utils.SendSl(msg)
			utils.Error(errors.New(msg), false)

			return
		}
	}()

	// Check Txpool Queued
	go func() {
		defer wg.Done()

		cnt, err := client.GetTxQueuedCnt(appCtx)
		if err != nil {
			utils.SendTg(err.Error())
			utils.SendPd(err.Error())
			utils.SendSl(err.Error())
			utils.Error(err, true)
			return
		}
		store.GlobalState.ELs[jsonrpc].TxpoolQueued = cnt

		if cnt >= int(app.cfg.EL.TxpoolQueuedThreshold) {
			store.GlobalState.ELs[jsonrpc].Status = false

			msg := fmt.Sprintf("Txpool Queued is too high: %d", cnt)
			utils.SendTg(msg)
			utils.SendPd(msg)
			utils.SendSl(msg)
			utils.Error(errors.New(msg), false)

			return
		}
	}()
	*/
	wg.Wait()

	utils.Debug(fmt.Sprintf("Finish check EL for %s and Elapsed Time: %s", jsonrpc, time.Since(now)))

	return nil
}
