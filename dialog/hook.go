package dialog

import (
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/config"
	"github.com/phongle318/chcb/db"
	"github.com/phongle318/chcb/step"
	"github.com/phongle318/chcb/text"
)

func PreHandleMessageHook(bot *fbbot.Bot, msg *fbbot.Message) bool {
	sender := msg.Sender
	// Try to get sender name from long term memory
	senderName := bot.LTMemory.For(sender.ID).Get("customerName")
	gender := bot.LTMemory.For(sender.ID).Get("gender")

	// Note: Supposedly new sender, but cannot differentiate from update sender name
	// Anyway, the db.NewSender() won't throw error in this case since it uses ON DUPLICATE KEY UPDATE
	if senderName == "" || gender == "" {
		senderName = sender.FullName()
		gender = sender.Gender()
		bot.LTMemory.For(sender.ID).Set("customerName", senderName)
		bot.LTMemory.For(sender.ID).Set("gender", gender)
		err := db.NewSender(db.Sender{ID: sender.ID, FullName: senderName, Gender: gender})
		if err != nil {
			log.Error("New sender error: ", err)
		}
	}
	// msg.Sender.SetFullName(senderName)
	// msg.Sender.SetGender(gender)

	return AdHocPostback(bot, msg) || CheckTiming(bot, msg)
}

func CheckTiming(bot *fbbot.Bot, msg *fbbot.Message) bool {
	// Prevent bot interfering when humans are chatting
	silenceWaitTime, _ := strconv.ParseFloat(config.Env.SilenceWaitTime, 32)
	lastEcho := bot.STMemory.For(msg.Sender.ID).Get("lastEcho")
	if lastEcho != "" && !step.TimeExpired(lastEcho, silenceWaitTime) {
		log.Debugf("Bot do nothing since staff is chatting or in silence, lastEcho = %s", text.ToTime(lastEcho).Format("15:04:05"))
		return true
	}

	// Let dialog be expired after certain time, otherwise human won't remember what state of dialog he is in
	dialogExpiredTime, _ := strconv.ParseFloat(config.Env.DialogExpiredTime, 32)
	lastActiveTime := bot.STMemory.For(msg.Sender.ID).Get("lastMessage")
	if step.TimeExpired(lastActiveTime, dialogExpiredTime) {
		log.Debug("Dialog is expired, move to Welcome state")
		// Dialog is expired, move to Welcome state
		fptshop.Move(msg, step.Welcome{})
		return true
	}

	return false
}

func AdHocPostback(bot *fbbot.Bot, msg *fbbot.Message) bool {
	switch msg.Text {
	case text.PayloadAskForProduct:
		bot.STMemory.Delete(msg.Sender.ID)
		fptshop.Move(msg, step.Product{})
		return true
	case text.PayloadAskForPromotion:
		bot.STMemory.Delete(msg.Sender.ID)
		fptshop.Move(msg, step.Promotion{})
		return true
	case text.PayloadAskForSupport:
		fptshop.Move(msg, step.Silence{})
		return true
	}

	return false
}
