package katintranet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBotAPI = tgbotapi.BotAPI

func (s *server) TelegramInit() error {
	return nil
	var err error
	s.telegramBotAPI, err = tgbotapi.NewBotAPI(s.config["telegram_bot_token"]) // KatiuskasBot
	return err
}

func (s *server) telegramMessage(message *tgbotapi.Message) {
	if !message.Chat.IsPrivate() {
		return
	}
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
			s.telegramBotAPI.Send(newmsg)
			return
		}
		userid = s.DBphoneToUserid(message.Contact.PhoneNumber)
		if userid <= 0 {
			var tmp string
			if tgid == message.Chat.ID {
				tmp = "no estás"
			} else {
				tmp = fmt.Sprintf("%s no está", message.Contact.FirstName)
			}
			newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s, pero me temo que %s en el listado de socios de Katiuskas.", username, tmp))
			s.telegramBotAPI.Send(newmsg)
			return
		}
		s.DBsetPhoneTgid(message.Contact.PhoneNumber, tgid)
		if tgid == message.Chat.ID {
			newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s.  Tu número de socio en Katiuskas es el %d.", username, userid))
			newmsg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
			s.telegramBotAPI.Send(newmsg)
			return
		}
		newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Gracias, %s.  El número de socio en Katiuskas de %s es el %d.", username, message.Contact.FirstName, userid))
		s.telegramBotAPI.Send(newmsg)
		return
	}
	userid = s.DBtelegramToUserid(message.Chat.ID)
	if userid > 0 {
		newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hola, %s... lo siento, pero no te entiendo.", username))
		s.telegramBotAPI.Send(newmsg)
		return
	}
	newmsg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Hola, %s.  Para poder saber quién eres, necesito que me envíes tu número de teléfono", username))
	newmsg.ReplyToMessageID = message.MessageID
	newmsg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButtonContact("Enviar teléfono")})
	s.telegramBotAPI.Send(newmsg)
}

func (s *server) telegramBotHandler(w http.ResponseWriter, r *http.Request) {
	bytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	var update tgbotapi.Update
	json.Unmarshal(bytes, &update)

	s.Log(r, LOG_DEBUG, fmt.Sprintf("tgbot: %s", bytes))
	if update.Message != nil {
		s.telegramMessage(update.Message)
	}
	//	if update.Message == nil {
	//		log(w, r, LOG_DEBUG, fmt.Sprintf("tgbot: unknown update: %+v", update))
	//	} else {
	//		log(w, r, LOG_DEBUG, fmt.Sprintf("tgbot: Message: %+v", update.Message))
	//		log(w, r, LOG_DEBUG, fmt.Sprintf("tgbot: Chat: %+v", update.Message.Chat))
	//	}
	fmt.Fprintln(w, "Hi there!")
}
