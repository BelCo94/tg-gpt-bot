package bot

import "github.com/BelCo94/tg-gpt-bot/openai"

type TelegramConfig struct {
	Token string `toml:"TOKEN"`
}

type BaseConfig struct {
	AdminID         int64  `toml:"ADMIN_ID"`
	ChatHistorySize uint16 `toml:"CHAT_HISTORY_SIZE"`
}

type DBConfig struct {
	Path string `toml:"PATH"`
}

type Config struct {
	Base     BaseConfig          `toml:"base"`
	Telegram TelegramConfig      `toml:"telegram"`
	OpenAI   openai.OpenAIConfig `toml:"openai"`
	DB       DBConfig            `toml:"db"`
}
