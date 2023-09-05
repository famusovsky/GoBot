package main

import (
	"log"
	"os"

	"github.com/famusovsky/GoBot/internal/botlogic"
	"github.com/famusovsky/GoBot/pkg/db"
	"github.com/famusovsky/GoBot/pkg/tgbot"
	_ "github.com/lib/pq"
)

func main() {
	token := os.Getenv("TGBOT_TOKEN")

	getProcessor := botlogic.NewProcessor

	errorLog := log.New(os.Stdout, "ERR\t", log.Ldate|log.Ltime)

	dBase, err := db.OpenViaEnvVars()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dBase.Close()

	tgbot.Run(token, getProcessor, errorLog, dBase)
}
