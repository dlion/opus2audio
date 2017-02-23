package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	URL          = "api2.online-convert.com/jobs"
	BOTAPI       = "TELEGRAM_BOT_API_KEY_HERE"
	DEBUG        = true
	APICONVERTER = "ONLINE-CONVERTER_API_KEY_HERE"
)

func main() {
	fmt.Printf("--- OPUS2AUDIO ---\n" +
		"Author: Domenico (DLion) Luciani\n" +
		"Site: https://domenicoluciani.com\n" +
		"License: MIT\n")

	bot, err := tgbotapi.NewBotAPI(BOTAPI)
	if err != nil {
		log.Panic(err)
	}

	if DEBUG {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panicln(err)
	}

	var mJob JsonResponsePost
	mFiles := make(JsonResponseGet, 1)

	for update := range updates {
		//If no messages -> skip
		if update.Message == nil {
			continue
		}
		//Only Documet type accepted
		if update.Message.Document == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "I need to receive only Document files")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			continue
		}

		source, err := bot.GetFileDirectURL(update.Message.Document.FileID)
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error on getting file's direct URL")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		if DEBUG {
			log.Println("Source: " + source)
		}

		payload := strings.NewReader("{\"input\":[{\"type\":\"remote\",\"source\":\"" + source + "\"}],\"conversion\":[{\"category\":\"audio\",\"target\":\"ogg\"}]}")

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://"+URL, payload)
		req.Header.Add("x-oc-api-key", APICONVERTER)
		req.Header.Add("content-type", "application/json")
		req.Header.Add("cache-control", "no-cache")

		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error on the POST request to online-converter")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		res, err := client.Do(req)
		defer res.Body.Close()

		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error doing the POST request")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error on the reading the body's response POST request")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		err = json.Unmarshal(body, &mJob)
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error decoding the POST response's json")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		//Warn the user
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Request taken, please wait...")
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)

		//Wait 5 seconds
		time.Sleep(5000 * time.Millisecond)

		req, err = http.NewRequest("GET", "http://"+string(URL)+"/"+string(mJob.ID)+"/output", nil)
		req.Header.Add("x-oc-api-key", APICONVERTER)
		req.Header.Add("content-type", "application/json")
		req.Header.Add("cache-control", "no-cache")
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error on the GET request to online-converter")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		if DEBUG {
			log.Println("URL GET: https://" + string(URL) + "/" + string(mJob.ID) + "/output")
		}

		res, err = client.Do(req)
		defer res.Body.Close()
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error doing the GET request")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		body, err = ioutil.ReadAll(res.Body)
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error on the reading the GET response's body")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		err = json.Unmarshal(body, &mFiles)
		if err != nil && DEBUG {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error decoding the GET response's json")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		if DEBUG {
			log.Printf("URI: %s\nSIZE: %d\n", mFiles[0].URI, mFiles[0].Size)
		}

		out, err := os.Create("tmpFiles/" + mFiles[0].ID + ".ogg")
		defer out.Close()
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error on create a new file on the local disk")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Panicln(err)
		}

		//Download file
		resp, err := http.Get(mFiles[0].URI)
		defer resp.Body.Close()
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting the result from the URI")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
			log.Println(err)
		}

		io.Copy(out, resp.Body)

		//Send file
		result := tgbotapi.NewAudioUpload(update.Message.Chat.ID, "tmpFiles/"+mFiles[0].ID+".ogg")
		result.Duration = 20 //20 seconds
		result.BaseFile.BaseChat.ReplyToMessageID = update.Message.MessageID
		bot.Send(result)

		//Remove tmpFile
		if _, err := os.Stat("tmpFiles/" + mFiles[0].ID + ".ogg"); err == nil {
			err = os.Remove("tmpFiles/" + mFiles[0].ID + ".ogg")
			if err != nil && DEBUG {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Error getting the result from the URI")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				log.Panicln("Error on deleting tmpFile")
			}
		}
	}
}
