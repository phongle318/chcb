package step

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
)

type Goodbye struct {
	fbbot.BaseStep
}

func (s Goodbye) Enter(bot *fbbot.Bot, msg *fbbot.Message) fbbot.Event {
	return event.Stay
}
