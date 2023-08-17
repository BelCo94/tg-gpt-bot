package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/BelCo94/tg-gpt-bot/openai"
	"github.com/BurntSushi/toml"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type TelegramConfig struct {
	Token   string `toml:"BOT_TOKEN"`
	AdminID int64  `toml:"ADMIN_ID"`
}

type BaseConfig struct {
	Whitelist       string `toml:"USERS_WHITELIST_PATH"`
	ChatHistorySize int64  `toml:"CHAT_HISTORY_SIZE"`
}

type Config struct {
	Base     BaseConfig          `toml:"base"`
	Telegram TelegramConfig      `toml:"telegram"`
	OpenAI   openai.OpenAIConfig `toml:"openai"`
}

type ChatStorageT struct {
	Lock     *sync.Mutex
	Messages []openai.MessageT
}

var Storage = make(map[int64]*ChatStorageT)
var StorageLock = sync.Mutex{}

var users = make(map[int64]struct{})
var config Config

func main() {
	configFile := flag.String("config", "conf/conf.toml", "config file")
	flag.Parse()

	fmt.Println(*configFile)

	if _, err := toml.DecodeFile(*configFile, &config); err != nil {
		log.Fatalf("can't parse config file: %v", err)
	}

	if err := parseUserFile(config.Base.Whitelist); err != nil {
		log.Fatalf("failed to parse user file: %v", err)
	}

	bot, err := gotgbot.NewBot(config.Telegram.Token, &gotgbot.BotOpts{
		Client: http.Client{},
	})
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	updater := ext.NewUpdater(&ext.UpdaterOpts{
		Dispatcher: ext.NewDispatcher(&ext.DispatcherOpts{
			Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
				log.Println("an error occured while handling update:", err.Error())
				return ext.DispatcherActionNoop
			},
		}),
	})
	dispatcher := updater.Dispatcher
	dispatcher.AddHandler(handlers.NewMessage(privateMsgFilter, privateMsgHandler))
	err = updater.StartPolling(bot, nil)
	if err != nil {
		log.Fatalf("failed to start polling: %v", err.Error())
	}
	log.Printf("%s has been started...\n", bot.Username)

	updater.Idle()
}

func privateMsgFilter(msg *gotgbot.Message) bool {
	_, allowed := users[msg.From.Id]
	return (msg.From.Id == config.Telegram.AdminID || allowed) && msg.From.Id == msg.Chat.Id
}

func privateMsgHandler(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage

	StorageLock.Lock()
	if _, ok := Storage[msg.Chat.Id]; !ok {
		Storage[msg.Chat.Id] = &ChatStorageT{
			Lock:     &sync.Mutex{},
			Messages: make([]openai.MessageT, 0, config.Base.ChatHistorySize+1),
		}
	}
	chatStorage := Storage[msg.Chat.Id]
	StorageLock.Unlock()
	chatStorage.Lock.Lock()
	defer chatStorage.Lock.Unlock()
	addMessageToStorage(chatStorage, openai.MessageT{
		Role:    "user",
		Content: msg.Text,
	})
	paylaod, err := openai.PreparePayload(chatStorage.Messages, config.OpenAI)
	if err != nil {
		return err
	}
	answer, err := openai.AskOpenAI(paylaod, config.OpenAI)
	if err != nil {
		return err
	}
	addMessageToStorage(chatStorage, *answer)
	sentMsg, err := msg.Reply(bot, answer.Content, &gotgbot.SendMessageOpts{
		ParseMode: "Markdown",
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	fmt.Printf("sent %d to %d\n", sentMsg.MessageId, sentMsg.Chat.Id)
	return nil
}

func addMessageToStorage(storage *ChatStorageT, msg openai.MessageT) {
	if len(storage.Messages) == int(config.Base.ChatHistorySize) {
		storage.Messages = storage.Messages[1:]
	}
	storage.Messages = append(storage.Messages, msg)
}

func parseUserFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		userID, err := strconv.ParseInt(line, 10, 64)
		if err == nil {
			users[userID] = struct{}{}
		}
	}
	return nil
}
