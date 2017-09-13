package step

import (
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
	"github.com/michlabs/fbbot"
)

type OrderName struct {
	fbbot.BaseStep
}

func (s OrderName) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	bot.SendText(msg.Sender, text.T("ask_for_name", &msg.Sender))
	return event.Stay
}

func (s OrderName) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	if name := detectCustomerName(msg.Text); name != "" {
		bot.LTMemory.For(msg.Sender.ID).Set("customerName", name)
		return event.Order
	}

	bot.SendText(msg.Sender, text.T("ask_for_name", &msg.Sender))
	return event.Stay
}

func detectCustomerName(msg string) string {
	//TODO extract customer name or return not found. For now, accept all input as name
	return msg
}
