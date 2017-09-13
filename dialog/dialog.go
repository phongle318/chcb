package dialog

import (
	"github.com/michlabs/fbbot"
	"github.com/phongle318/chcb/event"
	"github.com/phongle318/chcb/step"
)

var fptshop *fbbot.Dialog

func NewDialog() *fbbot.Dialog {
	d := fbbot.NewDialog()

	var welcome step.Welcome
	var goodbye step.Goodbye
	var order step.Order
	var product step.Product
	var hasError step.HasError

	d.SetBeginStep(welcome)
	d.SetEndStep(goodbye)

	d.AddTransition(event.AskForProduct, product)
	d.AddTransition(event.Order, order)
	d.AddTransition(event.Goodbye, goodbye)
	d.AddTransition(event.HasError, hasError)

	d.PreHandleMessageHook = PreHandleMessageHook

	fptshop = d
	return fptshop
}
