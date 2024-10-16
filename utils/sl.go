package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

type SL struct {
	enable     bool
	webhookURL string
}
type SlackMessage struct {
	Text string `json:"text"`
}

var sl SL
var slQueue chan func()

func SetSl(enable bool, webhookURL string) {
	if !enable {
		return
	}

	// Set SL
	// It is singleton
	sl = SL{
		enable,
		webhookURL,
	}
	slQueue = make(chan func())

	// For thread safe
	go func() {
		for sl := range slQueue {
			sl()
		}
	}()
}

func slenqueue(sl func()) {
	slQueue <- sl
}

func SendSl(msg string) {
	if !sl.enable {
		return
	}

	msg = fmt.Sprintf("%s\n%s", "monad-v", msg)

	slm := SlackMessage{
		msg,
	}

	SlBytes, err := json.Marshal(slm)
	if err != nil {
		log.Error().Err(err).Msg("\n" + string(debug.Stack()))
		return
	}
	buff := bytes.NewBuffer(SlBytes)

	req, err := http.NewRequest(
		"POST",
		sl.webhookURL,
		buff,
	)
	if err != nil {
		log.Error().Err(err).Msg("\n" + string(debug.Stack()))
		return
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}

	sl := func() {
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("\n" + string(debug.Stack()))
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		// Body 출력
		fmt.Println(string(body))
		defer resp.Body.Close()
		if !(200 <= resp.StatusCode && resp.StatusCode <= 299) {
			log.Error().Msg("Fail to seding msg from sl module")
			return
		}
	}
	slenqueue(sl)
}
