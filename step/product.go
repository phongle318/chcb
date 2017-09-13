package step

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/extservice"
	"github.com/phongle318/chcb/text"
)

type Product struct {
	fbbot.BaseStep
}

func (s Product) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	searchProduct := bot.STMemory.For(msg.Sender.ID).Get("searchProduct")
	if searchProduct != "" {
		return s.Process(bot, msg)
	} else {
		bot.SendText(msg.Sender, text.T("ask_for_product", &msg.Sender))
		return event.Stay
	}
}

func (s Product) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	searchTerm := bot.STMemory.For(msg.Sender.ID).Get("searchProduct")
	bot.STMemory.For(msg.Sender.ID).Delete("searchProduct")
	if searchTerm == "" {
		searchTerm = extractProductName(msg.Text)
		log.Debug("Extracted product name: ", searchTerm)
	}
	bot.TypingOn(msg.Sender)
	products, err := extservice.SearchForProducts(searchTerm)
	if err == nil && len(products) > 0 {
		bot.SendText(msg.Sender, text.T("show_product", &msg.Sender))
		bot.Send(msg.Sender, createProductsCarousel(products))
		return event.GoSilence
		// 	return event.OrderProduct
	} else {
		log.Error(err)
	}

	// session := bot.STMemory.For(msg.Sender.ID)
	// if session.Get("product_not_found") != "" {
	// 	return event.ProductNotFound
	// }
	// session.Set("product_not_found", "1")
	bot.SendText(msg.Sender, text.T("ask_for_product_again", &msg.Sender))
	return event.Stay
}

func createProductsCarousel(products []extservice.Product) *fbbot.GenericMessage {
	genericMessage := fbbot.NewGenericMessage()
	var bubbles []fbbot.Bubble
	for _, p := range products {
		productUrl := fmt.Sprintf("https://fptshop.com.vn/%s/%s", p.Type, p.UrlKey)
		price := text.ToInt(p.Price)
		bubble := fbbot.Bubble{
			Title: p.Name,
			// Add full-width space character (\u3000) to push Price to the right, so that it not be hidden by carousel arrow
			SubTitle: fmt.Sprintf("\u3000\u3000\u3000Giá: %s", text.FormatPriceVND(price)),
			ItemURL:  productUrl,
			ImageURL: text.GetResizedImgURL(p.ImageUrl),
		}
		var buttons []fbbot.Button
		buttons = append(buttons, fbbot.Button{
			Title:   text.TitleAskForPromotion,
			Type:    "postback",
			Payload: text.MarshalToString(map[string]interface{}{"action": "promotion", "product_sku": p.SKU}),
		})
		if price > 0 {
			buttons = append(buttons, fbbot.Button{
				Title: text.TitlePlaceOrder,
				Type:  "postback",
				Payload: text.MarshalToString(map[string]interface{}{
					"action":       "order",
					"product_id":   strconv.Itoa(p.ID),
					"product_name": p.Name,
					"product_sku":  p.SKU}),
			})
		} else {
			bubble.SubTitle = "\u3000\u3000\u3000Hàng sắp về"
			buttons = append(buttons, fbbot.Button{
				Title: text.TitleInStockNotification,
				Type:  "postback",
				Payload: text.MarshalToString(map[string]interface{}{
					"action":       "notify",
					"product_id":   strconv.Itoa(p.ID),
					"product_name": p.Name,
					"url":          productUrl,
					"image_url":    p.ImageUrl}),
			})
		}
		buttons = append(buttons, fbbot.Button{
			Title:   text.TitleAskForSupport,
			Type:    "postback",
			Payload: text.PayloadAskForSupport,
		})

		bubble.Buttons = buttons
		bubbles = append(bubbles, bubble)
	}
	genericMessage.Bubbles = bubbles
	return genericMessage
}

func extractProductName(msg string) string {
	// Disable FPT.AI for now since it extract "Iphone 7" -> "7"
	return msg
}
