package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

func SetupCron(handler func(), time string) {
	c := cron.New()

	_, err := c.AddFunc(time, func() {
		handler()
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
