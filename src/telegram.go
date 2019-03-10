package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var bot *tgbotapi.BotAPI

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(Telegram_bot_token) // KatiuskasBot
	if err != nil {
		log.Panic(err)
	}
}

func tg_private(message *tgbotapi.Message) {
	var userid int
	username := message.Chat.FirstName
	if username == "" {
		username = message.Chat.UserName
	}
	if username == "" {
		username = "chaval"
	}
	if message.Contact != nil {
		tgid := int64(message.Contact.UserID)
		if tgid <= 0 {
			newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s, pero me temo que %s no está en Telegram...", username, message.Contact.FirstName))
			bot.Send(newmsg)
			return
		}
		userid = db_phone_to_userid(message.Contact.PhoneNumber)
		if userid <= 0 {
			var tmp string
			if tgid == message.Chat.ID {
				tmp = "no estás"
			} else {
				tmp = fmt.Sprintf("%s no está", message.Contact.FirstName)
			}
			newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s, pero me temo %s en el listado de socios de Katiuskas.", username, tmp))
			bot.Send(newmsg)
			return
		}
		db_set_phone_tgid(message.Contact.PhoneNumber, tgid)
		if tgid == message.Chat.ID {
			newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s.  Tu número de socio en Katiuskas es el %d.", username, userid))
			newmsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			bot.Send(newmsg)
			return
		} else {
			newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s.  El número de socio en Katiuskas de %s es el %d.", username, message.Contact.FirstName, userid))
			bot.Send(newmsg)
			return
		}
		return
	}
	userid = db_telegram_to_userid(message.Chat.ID)
	if userid > 0 {
		newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hola, %s... lo siento, pero no te entiendo.", username))
		bot.Send(newmsg)
		return
	}
	newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hola, %s.  Para poder saber quién eres, necesito que me envíes tu número de teléfono", username))
	newmsg.ReplyToMessageID = message.MessageID
	newmsg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButtonContact("Enviar teléfono")})
	bot.Send(newmsg)
}

func tg_message(message *tgbotapi.Message) {
	if message.Chat.IsPrivate() {
		tg_private(message)
	}
}

func tgbotHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	Log(r, LOG_DEBUG, fmt.Sprintf("tgbot: %s", bytes))
	if update.Message != nil {
		tg_message(update.Message)
	}
//	if update.Message == nil {
//		log(w, r, LOG_DEBUG, fmt.Sprintf("tgbot: unknown update: %+v", update))
//	} else {
//		log(w, r, LOG_DEBUG, fmt.Sprintf("tgbot: Message: %+v", update.Message))
//		log(w, r, LOG_DEBUG, fmt.Sprintf("tgbot: Chat: %+v", update.Message.Chat))
//	}
	fmt.Fprintln(w, "Hi there!")
}

