# Opus2Audio

Be able to listen Whatsapp's voice audio files on Telegram

## What is it?

Opus2Audio is a Telegram's bot that allows you to be able to listen Whatsapp's voice audio files on Telegram

## How it work?

Opus2Audio use online-convert's APIs to converts OPUS files, so to use it you need an APIKEY from http://www.online-convert.com/.

Obviously you need an APIKEY from Telegram to use your own bot.

Pay attention: **The free account allows to you to convert 30files per day**.

## How to use it?

First of all you need to create your own bot on Telegram, follows these steps: https://core.telegram.org/bots

Then you have to create an account on online-convert to retrieve your APIKEY, follows this doc: http://apiv2.online-convert.com/

Then fill the constants in main.go file


I suggest to you to use glide as dependencies manager, so you can use it (`glide install`) or you can use `go get`.

That's all, now you can run or build it without any problems.

## License

MIT

## Author

Domenico Luciani aka DLion
https://domenicoluciani.com