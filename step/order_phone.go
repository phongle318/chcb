package step

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type OrderPhone struct {
	fbbot.BaseStep
}

func (s OrderPhone) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	bot.SendText(msg.Sender, text.T("ask_for_phone", &msg.Sender))
	return event.Stay
}

func (s OrderPhone) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	if phone := detectCustomerPhone(msg.Text); phone != "" {
		bot.LTMemory.For(msg.Sender.ID).Set("phone", phone)
		return event.Order
	}
	return s.Enter(bot, msg)
}

func detectCustomerPhone(msg string) string {
	r, err := regexp.Compile(`[0-9][0-9\s\-]+`)
	if err != nil {
		log.Error("Error in phone regexp: ", err)
		return ""
	}
	return r.FindString(msg)
}
