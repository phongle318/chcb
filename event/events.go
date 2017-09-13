package event

import (
	"github.com/michlabs/fbbot"
)

const Stay fbbot.Event = ""
const GoToWelcome fbbot.Event = "go to welcome"
const Order fbbot.Event = "order"
const AskForProduct fbbot.Event = "ask for product"
const HasError fbbot.Event = "has error"
const GoSilence fbbot.Event = "go silent"
const Goodbye fbbot.Event = "goodbye"
