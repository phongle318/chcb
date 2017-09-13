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

type ProductVariant struct {
	fbbot.BaseStep
}

func (s ProductVariant) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	sku := session.Get("productSku")
	productDetail, err := extservice.GetProductDetails(sku)
	if err != nil {
		log.Error(err)
		return event.HasError
	}
	variants := productDetail.Variant
	if len(variants) <= 1 { // Do not set variant id if there's less than 2 variants
		session.Set("productVariantID", "0")
		return event.Order
	}

	variantOptions := new(fbbot.QuickRepliesMessage)
	variantOptions.Text = text.T("product_variants", &msg.Sender)
	variantOptions.Items = make([]fbbot.QuickRepliesItem, len(variants))
	for i, v := range variants {
		variantOptions.Items[i] = fbbot.NewQuickRepliesText(
			v.ColorName,
			fmt.Sprintf(`{"variant_id": "%d", "variant_name": "%s"}`, v.ID, v.ColorName),
		)
		if v.ColorHex != "" {
			colorURL := text.GetColorImgURL(v.ColorHex)
			if colorURL != "" {
				variantOptions.Items[i].ImageURL = colorURL
			}
		}

	}
	bot.Send(msg.Sender, variantOptions)
	return event.Stay
}

func (s ProductVariant) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	session := bot.STMemory.For(msg.Sender.ID)
	if text.IsJSON(msg.Quickreply.Payload) {
		var payload = make(map[string]string)
		err := json.Unmarshal([]byte(msg.Quickreply.Payload), &payload)
		if err != nil {
			log.Error(err)
			return event.HasError
		}
		session.Set("productVariantID", payload["variant_id"])
		session.Set("productVariantName", payload["variant_name"])
		return event.Order
	}

	return event.GoSilence
}
