package step

import (
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/config"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/text"
)

var silenceWaitTime float64 // in minutes

func init() {
	silenceWaitTime, _ = strconv.ParseFloat(config.Env.SilenceWaitTime, 32)
}

type Silence struct {
	fbbot.BaseStep
}

func (s Silence) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	bot.STMemory.For(msg.Sender.ID).Set("lastEcho", text.FromTime(time.Now()))
	lastEcho := text.ToTime(bot.STMemory.For(msg.Sender.ID).Get("lastEcho"))
	log.Debugf("Bot entered silence, lastEcho = %s", lastEcho.Format("15:04:05"))

	if config.GetStaffId() != "" {
		bot.SendText(fbbot.User{ID: config.GetStaffId()}, text.T("notify_staff", &msg.Sender))
	}
	pleaseWait := new(fbbot.QuickRepliesMessage)
	pleaseWait.Text = text.T("staff_support", &msg.Sender)
	pleaseWait.Items = []fbbot.QuickRepliesItem{
		fbbot.NewQuickRepliesText("OK", "ok"),
	}
	bot.Send(msg.Sender, pleaseWait)
	return event.Stay
}

func (s Silence) Process(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	lastEcho := bot.STMemory.For(msg.Sender.ID).Get("lastEcho")
	if TimeExpired(lastEcho, silenceWaitTime) {
		log.Debugf("Bot was waken up, lastEcho = %s", text.ToTime(lastEcho).Format("15:04:05"))
		return event.GoToWelcome
	} else {
		return event.Stay
	}
}

func TimeExpired(lastTimeStr string, duration float64) bool {
	currentTime := time.Now()
	if lastTimeStr == "" {
		return false
	}
	lastTime := text.ToTime(lastTimeStr)
	if float64(currentTime.Sub(lastTime).Minutes()) < duration {
		return false
	}

	return true
}
