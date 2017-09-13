package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type ProductNotFound struct {
	fbbot.BaseStep
}

func (s ProductNotFound) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	var callToAction = fbbot.NewButtonMessage()
	callToAction.Text = text.T("product_not_found", &msg.Sender)
	callToAction.AddPostbackButton(text.TitleAskForPromotion, text.PayloadAskForPromotion)
	callToAction.AddPostbackButton(text.TitleAskForSupport, text.PayloadAskForSupport)
	bot.Send(msg.Sender, callToAction)

	return event.Stay
}

func (s ProductNotFound) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	switch msg.Text {
	// case text.PayloadAskForSupport:
	// return event.GoToPromotion
	default:
		return event.GoSilence
	}
}
