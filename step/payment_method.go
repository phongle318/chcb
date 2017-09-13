package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type PaymentMethod struct {
	fbbot.BaseStep
}

func (s PaymentMethod) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	wholePay := fbbot.NewQuickRepliesText(text.TitleWholePay, text.PayloadWholePay)
	installmentPay := fbbot.NewQuickRepliesText(text.TitleInstallmentPay, text.PayloadInstallmentPay)
	paymentOptions := new(fbbot.QuickRepliesMessage)
	paymentOptions.Text = text.T("payment_method", &msg.Sender)
	paymentOptions.Items = []fbbot.QuickRepliesItem{wholePay, installmentPay}
	bot.Send(msg.Sender, paymentOptions)
	return event.Stay
}

func (s PaymentMethod) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	switch msg.Quickreply.Payload {
	case text.PayloadWholePay, text.PayloadInstallmentPay:
		bot.STMemory.For(msg.Sender.ID).Set("paymentMethod", msg.Quickreply.Payload)
		return event.Order
	default:
		return event.GoSilence
	}
}
