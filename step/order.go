package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
)

type Order struct {
	fbbot.BaseStep
}

func (s Order) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	// persist := bot.LTMemory.For(msg.Sender.ID)

	if session.Get("productSku") == "" {
		return event.AskForProduct
	}

	// if session.Get("paymentMethod") == "" {
	// 	return event.PaymentMethod
	// }

	// if session.Get("paymentMethod") == text.PayloadInstallmentPay {
	// 	log.Debug("Customer chose Installment payment method, bot goes silence")
	// 	return event.GoSilence
	// }

	// if session.Get("productVariantID") == "" {
	// 	return event.ProductVariant
	// }

	// if persist.Get("customerName") == "" {
	// 	return event.AskForName
	// }

	// if persist.Get("phone") == "" {
	// 	return event.AskForPhoneNumber
	// }

	// if session.Get("shippingMethod") == "" {
	// 	return event.ShippingMethod
	// }
	// if session.Get("shippingMethod") == text.PayloadPickupStore &&
	// 	session.Get("pickupStore") == "" {
	// 	return event.PickupStore
	// }

	return event.GoSilence
}
