# tg-gpt-bot

The simple telegram bot for the ChatGPT. Just for example and experiments.

Run:

```bash
tg-gpt-bot --config /path/to/conf/conf.toml
```


Config file example (`conf/conf.toml` by default): 
```toml
[base]
ADMIN_ID = 11111111
CHAT_HISTORY_SIZE = 10

[telegram]
TOKEN = "yourtelegrambottoken"

[openai]
TOKEN = "youropenaitoken"
MODEL = "gpt-3.5-turbo"
SYSTEM_PROMPT = "Always start a message with ....."

```

Parameters:

| Parameter                       | Type   | Description                             |
| -                               | -      | -                                       |
| base / ADMIN_ID             | int    | Telegram ID of the admin user |
| base / CHAT_HISTORY_SIZE        | int    | Number of the last messages from the user's chat to be sent as context to the model | 
| telegram / BOT_TOKEN            | string | Telegram bot token |
| openai / TOKEN                  | string | OpenAI token |
| openai / MODEL                  | string | The name of the OpenAI model to be used |
| openai / SYSTEM_PROMPT          | string | Text to be sent as system prompt to the model (aka context) |
