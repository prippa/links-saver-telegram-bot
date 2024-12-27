package main

import (
	"flag"
	tgClient "links-saver-telegram-bot/clients/telegram"
	event_consumer "links-saver-telegram-bot/consumer/event-consumer"
	"links-saver-telegram-bot/events/telegram"
	"links-saver-telegram-bot/storage/files"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {
	token := mustToken()
	client := tgClient.New(tgBotHost, token)
	eventsProcessor := telegram.New(client, files.New(storagePath))

	log.Print("Telegram bot started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("Error starting consumer: %s", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"Token to acess to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("Token is required")
	}

	return *token
}
