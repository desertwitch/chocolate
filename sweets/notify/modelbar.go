package notify

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
)

func NewNotificationBar(c *chocolate.Chocolate, useFlavour bool) *NotifyModel {
	nmodel := &NotifyModel{useFlavour: useFlavour}
	nmodel.initDefaultNotificationTypes()

	notificaitonBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:                   nmodel,
			UpdateHandlerFct:        notifyModelUpdateHandler,
			FlavourCustomizeHandler: notifyModelFlavourCustomizeHandler,
		},
		chocolate.WithBarXScaler(chocolate.DYNAMIC, 0),
		chocolate.WithBarYScaler(chocolate.DYNAMIC, 0),
		chocolate.WithBarID("overlay"),
	)
	notificaitonBar.Hide(true)
	c.AddOverlayRoot(notificaitonBar)

	c.RegisterUpdateFor(notifyMsg{}, func(b chocolate.ChocolateBar) func(tea.Msg) (tea.Cmd, bool) {
		return func(msg tea.Msg) (tea.Cmd, bool) {
			return b.HandleUpdate(msg), true
		}
	}(notificaitonBar))
	c.RegisterUpdateFor(destroyMsg(""), func(b chocolate.ChocolateBar) func(tea.Msg) (tea.Cmd, bool) {
		return func(msg tea.Msg) (tea.Cmd, bool) {
			return b.HandleUpdate(msg), true
		}
	}(notificaitonBar))

	return nmodel
}
