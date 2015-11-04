package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/thorduri/pushover"
)

const (
	NotifyFrequency = 60
)

func notify(po pushover.Pushover, status string, last *time.Time) {
	if time.Since(*last).Minutes() < NotifyFrequency {
		return
	}

	*last = time.Now()

	err := po.Message(status)

	if err != nil {
		log.Println(err)
	}
}

func main() {
	token := os.Getenv("PUSHOVER_TOKEN")
	user := os.Getenv("PUSHOVER_USER")

	if len(token) == 0 || len(user) == 0 {
		fmt.Println("Environment variable `PUSHOVER_TOKEN` or `PUSHOVER_USER` was not defined")
		return
	}

	po, err := pushover.NewPushover(token, user)

	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		// 1440 (24 * 60)
		for _ = range time.Tick(5 * time.Second) {
			var last time.Time

			res, err := http.Get("http://www.atualizarboleto.com.br/")

			if err != nil {
				log.Println(err)
			}

			if res.StatusCode == http.StatusOK {
				notify(*po, res.Status, &last)
			}

			log.Printf("%s %s\n", res.Proto, res.Status)
		}
	}()

	select {}
}
