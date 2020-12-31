# TaigaBot
Simple discord bot written in go.

## Configuration
To configure this bot you need to create `config.toml` with `[bot]` section in it:
```toml
[bot]
token = "<Your token here>"
logs_id = "<ID of channel for error logging>"
db_uri = "<URI to connect mongo db>"
```
Or write `prod = true` for pick values from environment variables:
- `TOKEN` - is your discord bot token
- `LOGS_ID` - is ID of the channel for error logs and announcements
- `DB_URI` - is connect uri of your mongo datebase

## Commands
List of all commands can be found by executing `!help` command, or `!help moderation` for moderation commands.

## Contribution
All contributions are welcome.
