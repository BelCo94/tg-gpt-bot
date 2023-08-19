package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BelCo94/tg-gpt-bot/bot"
	"github.com/BurntSushi/toml"
)

func main() {
	configFile := flag.String("config", "conf/conf.toml", "config file")
	flag.Parse()

	fmt.Printf("Using config: %s\n", *configFile)

	var config bot.Config

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatalf("can't parse config file: %v", err)
	}

	b, err := bot.NewBot(config)
	if err != nil {
		log.Fatalf("failed to create a bot: %v", err)
	}

	b.Run()
}
