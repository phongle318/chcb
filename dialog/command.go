package dialog

import (
	log "github.com/Sirupsen/logrus"
	"github.com/michlabs/fbbot"
)

func NewCommander() *fbbot.Commander {
	commander := fbbot.NewCommander()

	commander.Add("restart", func(bot *fbbot.Bot, echoMsg *fbbot.Message, params string) {
		senderID := echoMsg.Page.ID
		fptshop.Reset(senderID)
		bot.STMemory.Delete(senderID)
		bot.SendText(fbbot.User{ID: senderID}, "Dialog restarted")
		log.Debugf("debug: dialog with %s was restarted", senderID)
	})

	return commander
}
