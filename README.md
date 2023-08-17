# tg-gpt-bot

The simple telegram bot for the ChatGPT. Just for example and experiments.

Run:

```bash
tg-gpt-bot /path/to/conf/conf.toml
```


Config file example (`conf.toml`): 
```toml
[base]
USERS_WHITELIST_PATH = "/path/to/conf/users.txt"
CHAT_HISTORY_SIZE = 10

[telegram]
BOT_TOKEN = "yourtelegrambottoken"
ADMIN_ID = 11111111

[openai]
TOKEN = "youropenaitoken"
MODEL = "gpt-3.5-turbo"
SYSTEM_PROMPT = "Always start a message with ....."

```

Parameters:

| Parameter                       | Type   | Description                             |
| -                               | -      | -                                       |
| base / USERS_WHITELIST_PATH     | string | Path to the file with whitelisted users |
| base / CHAT_HISTORY_SIZE        | int    | Number of the last messages from the user's chat to be sent as context to the model | 
| telegram / BOT_TOKEN            | string | Telegram bot token |
| telegram / ADMIN_ID             | int    | Telegram ID of the admin user (for now it's just an extra whitelisted id) |
| openai / TOKEN                  | string | OpenAI token |
| openai / MODEL                  | string | The name of the OpenAI model to be used |
| openai / SYSTEM_PROMPT          | string | Text to be sent as system prompt to the model (aka context) |
