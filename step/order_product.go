package step

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/db"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

type OrderProduct struct {
	fbbot.BaseStep
}

func (s OrderProduct) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	return event.Stay
}

func (s OrderProduct) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	if text.IsJSON(msg.Text) {
		var payload = make(map[string]string)
		err := json.Unmarshal([]byte(msg.Text), &payload)
		if err != nil {
			log.Error(err)
			return event.HasError
		}
		switch payload["action"] {
		// case "promotion":
		// 	session.Set("productSku", payload["product_sku"])
		// 	return event.ProductPromotion
		case "order":
			session.Set("productSku", payload["product_sku"])
			session.Set("productName", payload["product_name"])
			session.Set("productID", payload["product_id"])
			return event.Order
		case "notify":
			err := db.SubscribeProductNotification(
				db.ProductNotification{
					ProductID:   payload["product_id"],
					ProductName: payload["product_name"],
					Url:         payload["url"],
					ImageUrl:    payload["image_url"],
				},
				db.Subscriber{
					ProductID:  payload["product_id"],
					SenderID:   msg.Sender.ID,
					SenderName: msg.Sender.FullName(),
				})
			if err != nil {
				log.Error(err)
				return event.HasError
			}
			bot.SendText(msg.Sender, text.T("notification_subscribed", &msg.Sender))
			return event.Goodbye

		default:
			log.Debug("Unknown action: ", payload["action"])
			return event.GoSilence
		}
	} else {
		numberOfSearch := bot.STMemory.For(msg.Sender.ID).Get("numberOfSearch")
		if numberOfSearch == "x" {
			return event.GoSilence
		}
		// text.PayloadAskForSupport is handled by hook already
		bot.STMemory.For(msg.Sender.ID).Set("searchProduct", msg.Text)
		bot.STMemory.For(msg.Sender.ID).Set("numberOfSearch", "x"+numberOfSearch)
		return event.AskForProduct
	}
}
