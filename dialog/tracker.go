package dialog

import (
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/phongle318/chcb/db"
	"github.com/phongle318/chcb/text"
	"github.com/michlabs/fbbot"
)

type ActivityTracker struct{}

func (h *ActivityTracker) HandleEcho(bot *fbbot.Bot, echoMsg *fbbot.Message) {
	// Do not handle echo from bot itself
	if echoMsg.AppID > 0 {
		return
	}
	// Do not count commands as last echo
	if strings.HasPrefix(echoMsg.Text, "/") {
		return
	}
	// In Echo message, Page.ID is recipient ID
	bot.STMemory.For(echoMsg.Page.ID).Set("lastEcho", text.FromTime(time.Now()))
}

func (h *ActivityTracker) HandleMessage(bot *fbbot.Bot, msg *fbbot.Message) {
	bot.STMemory.For(msg.Sender.ID).Set("lastMessage", text.FromTime(time.Now()))
}

func (h *ActivityTracker) HandlePostback(bot *fbbot.Bot, msg *fbbot.Postback) {
	bot.STMemory.For(msg.Sender.ID).Set("lastPostback", text.FromTime(time.Now()))
}

func (h *ActivityTracker) HandleRead(bot *fbbot.Bot, msg *fbbot.Read) {
	readAt := time.Unix(int64(msg.Watermark)/1000, 0)
	err := db.UpdatePromotionRead(msg.Sender.ID, readAt)
	if err != nil {
		log.Error(err)
	}
}
