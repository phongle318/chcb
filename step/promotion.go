package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/extservice"
	"github.com/phongle318/chcb/text"
)

type PromotionData interface {
	GetTitle() string
	GetUrl() string
	GetImage() string
	GetDescription() string
}

type Promotion struct {
	fbbot.BaseStep
}

func (s Promotion) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	promotions, err := extservice.GetHomepagePromotions()
	if err == nil {
		if len(promotions) > 0 {
			bot.SendText(msg.Sender, text.T("promotion", &msg.Sender))
			promotionData := make([]PromotionData, len(promotions))
			for i, d := range promotions {
				promotionData[i] = d
			}
			bot.Send(msg.Sender, CreatePromotionCarousel(promotionData))
		} else {
			var callToAction = fbbot.NewButtonMessage()
			callToAction.Text = text.T("promotion_not_found", &msg.Sender)
			callToAction.AddPostbackButton(text.TitleAskForProduct, text.PayloadAskForProduct)
			callToAction.AddPostbackButton(text.TitleAskForSupport, text.PayloadAskForSupport)
			bot.Send(msg.Sender, callToAction)
		}
		return event.Stay
	} else {
		return event.HasError
	}
}

func (s Promotion) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	return event.GoSilence
}

func CreatePromotionCarousel(promotions []PromotionData) *fbbot.GenericMessage {
	genericMessage := fbbot.NewGenericMessage()
	var bubbles []fbbot.Bubble
	for _, p := range promotions {
		bubble := fbbot.Bubble{
			Title:    p.GetTitle(),
			SubTitle: p.GetDescription(),
			ItemURL:  p.GetUrl(),
			ImageURL: text.GetResizedImgURL(p.GetImage()),
		}
		var buttons []fbbot.Button
		buttons = append(buttons, fbbot.Button{
			Title: "Xem chi tiết",
			Type:  "web_url",
			URL:   p.GetUrl(),
		})
		buttons = append(buttons, fbbot.Button{
			Type: "element_share",
		})
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
