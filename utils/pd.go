package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

type PD struct {
	enable  bool
	Pdevent PDEvent
}
type PDEvent struct {
	RoutingKey  string  `json:"routing_key"`
	EventAction string  `json:"event_action"`
	Payload     Payload `json:"payload"`
}

type Payload struct {
	Summary  string `json:"summary"`
	Severity string `json:"severity"`
	Source   string `json:"source"`
}

var pd PD
var pdQueue chan func()

func SetPd(enable bool, routing_key string, event_action string, severity string, source string) {
	if !enable {
		return
	}

	// Set PD
	// It is singleton
	pd = PD{
		enable,
		PDEvent{
			routing_key,
			event_action,
			Payload{
				"",
				severity,
				source,
			},
		},
	}
	pdQueue = make(chan func())

	// For thread safe
	go func() {
		for pd := range pdQueue {
			pd()
		}
	}()
}

func pdenqueue(pd func()) {
	pdQueue <- pd
}

func SendPd(msg string) {
	if !pd.enable {
		return
	}

	url := "https://events.pagerduty.com/v2/enqueue"
	msg = fmt.Sprintf("%s\n%s", pd.Pdevent.Payload.Summary, msg)

	Pdevent := PDEvent{
		pd.Pdevent.RoutingKey,
		pd.Pdevent.EventAction,
		Payload{
			msg,
			pd.Pdevent.Payload.Severity,
			pd.Pdevent.Payload.Source,
		},
	}

	PdeventBytes, err := json.Marshal(Pdevent)
	if err != nil {
		log.Error().Err(err).Msg("\n" + string(debug.Stack()))
		return
	}
	buff := bytes.NewBuffer(PdeventBytes)

	req, err := http.NewRequest(
		"POST",
		url,
		buff,
	)
	if err != nil {
		log.Error().Err(err).Msg("\n" + string(debug.Stack()))
		return
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}

	pd := func() {
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("\n" + string(debug.Stack()))
			return
		}
		defer resp.Body.Close()
		/*
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading body:", err)
				return
			}

			// Body 출력
			fmt.Println(string(body))
		*/
		if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
			log.Error().Msg("Fail to seding msg from pd module")
			return
		}
	}
	pdenqueue(pd)
}
