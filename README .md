# Telegram Command Bot

Use Telegram Bot as remote controller to assign some works for local device by Golang

## How To Use
1. Write command as .sh file to folder, `cmd`
2. Edit button labels and the paths connection in `./models/totalListRaw.json`
3. Overwrite the settings in `conf/config.go`
4. Install imports dependencies
```sh
go mod download
```
## How To Run
```
go run main.go
```
## How To Deploy
```sh
go build && ./telegram-cmd_bot 
```

## Dependencies

- [go-telegram-bot-api/telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)

## License

MIT License