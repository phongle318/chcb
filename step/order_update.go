package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type OrderUpdate struct {
	fbbot.BaseStep
}

func (s OrderUpdate) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	var callToAction = new(fbbot.QuickRepliesMessage)
	callToAction.Text = text.T("order_update", &msg.Sender)
	callToAction.Items = []fbbot.QuickRepliesItem{
		fbbot.NewQuickRepliesText(text.TitleChangeName, text.PayloadChangeName),
		fbbot.NewQuickRepliesText(text.TitleChangePhone, text.PayloadChangePhone),
		fbbot.NewQuickRepliesText(text.TitleChangeProduct, text.PayloadAskForProduct),
	}
	if session.Get("productVariantName") != "" {
		callToAction.Items = append(callToAction.Items,
			fbbot.NewQuickRepliesText(text.TitleChangeColor, text.PayloadChangeColor),
		)
	}
	if session.Get("pickupStore") != "" {
		callToAction.Items = append(callToAction.Items,
			fbbot.NewQuickRepliesText(text.TitleChangeShop, text.PayloadChangeShop),
		)
	}
	bot.Send(msg.Sender, callToAction)
	return event.Stay
}

func (s OrderUpdate) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	switch msg.Quickreply.Payload {
	// case text.PayloadChangeName:
	// 	bot.LTMemory.For(msg.Sender.ID).Delete("customerName")
	// 	return event.AskForName
	// case text.PayloadChangePhone:
	// 	bot.LTMemory.For(msg.Sender.ID).Delete("phone")
	// 	return event.AskForPhoneNumber
	case text.PayloadAskForProduct:
		bot.STMemory.For(msg.Sender.ID).Delete("productSku")
		bot.STMemory.For(msg.Sender.ID).Delete("productVariantID")
		bot.STMemory.For(msg.Sender.ID).Delete("productVariantName")
		return event.AskForProduct
	// case text.PayloadChangeColor:
	// 	bot.STMemory.For(msg.Sender.ID).Delete("productVariantID")
	// 	bot.STMemory.For(msg.Sender.ID).Delete("productVariantName")
	// 	return event.ProductVariant
	// case text.PayloadChangeShop:
	// 	bot.STMemory.For(msg.Sender.ID).Delete("pickupStore")
	// 	return event.PickupStore
	default:
		return event.GoSilence
	}
}
