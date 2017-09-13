package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type Welcome struct {
	fbbot.BaseStep
}

func (s Welcome) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	sender := msg.Sender
	// Try to get sender name from long term memory
	// senderName := bot.LTMemory.For(msg.Sender.ID).Get("customerName")
	// log.Debugf("Starting chat with %s(%s)", senderName, sender.ID)
	bot.STMemory.Delete(sender.ID)
	bot.TypingOn(sender)
	bot.SendImage(sender, "https://fptshop.fpt.ai/diem_my.gif")
	var callToAction = fbbot.NewButtonMessage()
	callToAction.Text = text.T("prompt_actions", &sender)
	callToAction.AddPostbackButton(text.TitleAskForProduct, text.PayloadAskForProduct)
	callToAction.AddPostbackButton(text.TitleAskForPromotion, text.PayloadAskForPromotion)
	callToAction.AddPostbackButton(text.TitleAskForSupport, text.PayloadAskForSupport)
	bot.Send(sender, callToAction)
	return event.Stay
}

func (s Welcome) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	switch msg.Text {
	case text.PayloadAskForProduct:
		return event.AskForProduct
	// case text.PayloadAskForPromotion:
	// 	return event.GoToPromotion
	case text.PayloadAskForSupport:
		return event.GoSilence
	default:
		return event.GoSilence
	}
}
