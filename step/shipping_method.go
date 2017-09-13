package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type ShippingMethod struct {
	fbbot.BaseStep
}

func (s ShippingMethod) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	shippingOptions := fbbot.NewButtonMessage()
	shippingOptions.Text = text.T("shipping_method", &msg.Sender)
	shippingOptions.Buttons = []fbbot.Button{
		fbbot.NewPostbackButton(text.TitlePickupStore, text.PayloadPickupStore),
		fbbot.NewPostbackButton(text.TitleHomeDelivery, text.PayloadHomeDelivery),
	}
	bot.Send(msg.Sender, shippingOptions)
	return event.Stay
}

func (s ShippingMethod) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	switch msg.Text {
	case text.PayloadHomeDelivery, text.PayloadPickupStore:
		session.Set("shippingMethod", msg.Text)
		return event.Order
	default:
		return event.GoSilence
	}
}
