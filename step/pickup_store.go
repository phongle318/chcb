package step

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/extservice"
	"github.com/phongle318/chcb/text"
)

type PickupStore struct {
	fbbot.BaseStep
}

func (s PickupStore) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	bot.SendText(msg.Sender, text.T("pickup_store", &msg.Sender))
	return event.Stay
}

func (s PickupStore) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	if text.IsJSON(msg.Text) {
		var payload = make(map[string]string)
		err := json.Unmarshal([]byte(msg.Text), &payload)
		if err != nil {
			log.Error(err)
			return event.HasError
		}
		switch payload["action"] {
		case "shop_select":
			session.Set("pickupStore", payload["shop_name"])
			return event.Order
		default:
			return event.GoSilence
		}
	}

	sku := session.Get("productSku")
	bot.TypingOn(msg.Sender)
	shops := extservice.GetNearestShopHasSufficientStock(msg.Text, sku)
	if len(shops) == 0 {
		buttonMsg := fbbot.NewButtonMessage()
		buttonMsg.Text = text.T("store_not_found", &msg.Sender)
		buttonMsg.AddPostbackButton(text.TitleBackOrder, `{"action":"shop_select", "shop_name":"back_order"}`)
		buttonMsg.AddPostbackButton(text.TitleOtherProduct, text.PayloadAskForProduct)
		buttonMsg.AddPostbackButton(text.TitleAskForSupport, text.PayloadAskForSupport)
		bot.Send(msg.Sender, buttonMsg)
		return event.Stay
	} else {
		shopSelection := fbbot.NewGenericMessage()
		var bubbles []fbbot.Bubble
		for i, shop := range shops {
			if i >= 10 { // FB allows 10 elements only
				break
			}
			mapImg := fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?size=573x300&center=%s&zoom=18&markers=%s", shop.Location, shop.Location)
			bubble := fbbot.Bubble{
				Title:    fmt.Sprintf("FPT Shop %s, %s, %s", shop.ShopName[4:], shop.District, shop.Province),
				ImageURL: mapImg,
			}
			var buttons []fbbot.Button
			buttons = append(buttons, fbbot.Button{
				Title:   text.TitleSelect,
				Type:    "postback",
				Payload: text.MarshalToString(map[string]interface{}{"action": "shop_select", "shop_name": shop.ShopName}),
			})
			buttons = append(buttons, fbbot.Button{
				Title:   text.TitleAskForSupport,
				Type:    "postback",
				Payload: text.PayloadAskForSupport,
			})
			bubble.Buttons = buttons
			bubbles = append(bubbles, bubble)
		}
		shopSelection.Text = text.T("please_select_store", &msg.Sender)
		shopSelection.Bubbles = bubbles
		bot.Send(msg.Sender, shopSelection)
		return event.Stay
	}
}
