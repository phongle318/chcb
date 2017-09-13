package step

import (
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/extservice"
	"github.com/phongle318/chcb/text"
)

type ProductPromotion struct {
	fbbot.BaseStep
}

func (s ProductPromotion) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	sku := session.Get("productSku")
	if sku == "" {
		return event.AskForProduct
	}

	productDetails, err := extservice.GetProductDetails(sku)
	if err != nil {
		log.Error("Product details API: ", err)
		return event.HasError
	}

	product := productDetails.Product
	var callToAction = fbbot.NewButtonMessage()
	if product.Promotion == "" {
		callToAction.Text = fmt.Sprintf(text.T("product_no_promotion", &msg.Sender), product.Name)
	} else {
		callToAction.Text = fmt.Sprintf(text.T("product_promotion", &msg.Sender)+"\n%s", product.Name, text.Sanitize(product.Promotion))
	}
	orderPayload := text.MarshalToString(map[string]interface{}{
		"action":       "order",
		"product_sku":  sku,
		"product_id":   strconv.Itoa(product.ID),
		"product_name": product.Name,
	})
	callToAction.AddPostbackButton(text.TitlePlaceOrder, orderPayload)
	callToAction.AddPostbackButton(text.TitleOtherProduct, text.PayloadAskForProduct)
	callToAction.AddPostbackButton(text.TitleAskForSupport, text.PayloadAskForSupport)
	bot.Send(msg.Sender, callToAction)
	return event.Stay
}

func (s ProductPromotion) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	if text.IsJSON(msg.Text) {
		var payload = make(map[string]string)
		err := json.Unmarshal([]byte(msg.Text), &payload)
		if err != nil {
			log.Error(err)
			return event.HasError
		}
		switch payload["action"] {
		case "order":
			session.Set("productName", payload["product_name"])
			session.Set("productID", payload["product_id"])
			session.Set("productSku", payload["product_sku"])
			return event.Order
		case "promotion":
			session.Set("productSku", payload["product_sku"])
			return s.Enter(bot, msg)
		default:
			log.Debug("Unknown action: ", payload["action"])
		}
	}
	return event.GoSilence
}
