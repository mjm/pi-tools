package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/mjm/pi-tools/homebase/bot/telegram"
	"github.com/mjm/pi-tools/observability"
	"github.com/mjm/pi-tools/pkg/signal"
)

func main() {
	flag.Parse()

	stopObs, err := observability.Start("homebase-bot-srv")
	if err != nil {
		log.Panicf("Error setting up observability: %v", err)
	}
	defer stopObs()

	c, err := telegram.New(telegram.Config{
		Token: os.Getenv("TELEGRAM_TOKEN"),
	})
	if err != nil {
		log.Panicf("Error creating Telegram client: %v", err)
	}

	ch := make(chan telegram.UpdateOrError, 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c.WatchUpdates(ctx, ch, telegram.GetUpdatesRequest{
		Timeout: 30,
	})

	go func() {
		for updateOrErr := range ch {
			if updateOrErr.Err != nil {
				log.Printf("Error getting updates: %v", updateOrErr.Err)
			} else {
				update := updateOrErr.Update
				if update.Message != nil {
					log.Printf("Received message: %#v", update.Message)
				} else {
					log.Printf("Received update: %+v", updateOrErr.Update)
				}
			}
		}
	}()

	signal.Wait()
}
