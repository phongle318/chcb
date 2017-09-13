package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type HasError struct {
	fbbot.BaseStep
}

func (s HasError) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	bot.SendText(msg.Sender, text.T("error", &msg.Sender))
	return event.GoSilence
}
