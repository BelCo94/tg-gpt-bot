package bot

import (
	"fmt"

	"github.com/BelCo94/tg-gpt-bot/openai"
	"github.com/BelCo94/tg-gpt-bot/storage"
	"github.com/glebarez/sqlite"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Bot struct {
	tgbot   *tgbotapi.BotAPI
	storage *storage.Storage
	config  Config
}

func NewBot(config Config) (*Bot, error) {
	db, err := connectDB(config.DB)
	if err != nil {
		return nil, err
	}
	storage_ := &storage.Storage{
		DB: db,
	}
	storage_.InitModels()

	tgbot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		tgbot:   tgbot,
		storage: storage_,
		config:  config,
	}
	err = bot.updateInitialUsers()
	if err != nil {
		return nil, err
	}
	return bot, nil
}

func (bot *Bot) updateInitialUsers() error {
	me, err := bot.tgbot.GetMe()
	if err != nil {
		return err
	}
	err = bot.storage.CreateOrUpdateUser(
		&storage.User{
			ID:         int64(me.ID),
			Name:       me.UserName,
			IsDisabled: false,
		},
	)
	if err != nil {
		return err
	}
	adminChat, err := bot.tgbot.GetChat(tgbotapi.ChatConfig{
		ChatID: bot.config.Base.AdminID,
	})
	if err != nil {
		return err
	}
	err = bot.storage.CreateOrUpdateUser(
		&storage.User{
			ID:         adminChat.ID,
			Name:       adminChat.UserName,
			IsDisabled: false,
		},
	)
	if err != nil {
		return err
	}
	err = bot.storage.CreateOrUpdateChat(
		&storage.Chat{
			ID:         adminChat.ID,
			Name:       adminChat.UserName,
			IsDisabled: false,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func connectDB(config DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (bot *Bot) Run() {
	updater := tgbotapi.NewUpdate(0)
	updater.Timeout = 5

	updates, err := bot.tgbot.GetUpdatesChan(updater)
	if err != nil {
		fmt.Printf("failed to connect to update channel: %v", err)
		return
	}
	fmt.Printf("running %s ...\n", bot.tgbot.Self.UserName)
	for update := range updates {
		go bot.handleUpdate(update)
	}
}

func (bot *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message != nil {
		err := bot.messageHandler(update.Message)
		if err != nil {
			fmt.Printf("Fail to handle a message: %v", err)
		}
	}
}

func (bot *Bot) getMessageRole(msg *storage.Message) string {
	if msg.UserID == int64(bot.tgbot.Self.ID) {
		return "assistant"
	}
	return "user"
}

func (bot *Bot) askOpenAI(messages []*storage.Message) error {
	if len(messages) == 0 {
		return fmt.Errorf("empty message list")
	}
	conversation := make([]openai.MessageT, 0, bot.config.Base.ChatHistorySize)
	for _, msg := range messages {
		conversation = append(conversation, openai.MessageT{
			Role:    bot.getMessageRole(msg),
			Content: msg.Text,
		})
	}
	paylaod, err := openai.PreparePayload(conversation, bot.config.OpenAI)
	if err != nil {
		return err
	}
	answer, err := openai.AskOpenAI(paylaod, bot.config.OpenAI)
	if err != nil {
		return err
	}
	text := prepareTextForMarkdownV2(answer.Content)
	newMsg := tgbotapi.NewMessage(messages[0].ChatID, text)
	newMsg.ReplyToMessageID = int(messages[0].ID)
	newMsg.ParseMode = "MarkdownV2"

	sentMessage, err := bot.tgbot.Send(newMsg)
	if err != nil {
		return err
	}
	storedMessage := &storage.Message{
		ID:        int64(sentMessage.MessageID),
		ChatID:    sentMessage.Chat.ID,
		UserID:    int64(sentMessage.From.ID),
		Text:      sentMessage.Text,
		ReplyToID: int64(sentMessage.ReplyToMessage.MessageID),
	}
	return bot.storage.CreateMessage(storedMessage)
}
