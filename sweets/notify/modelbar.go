package notify

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mfulz/chocolate"
)

func NewNotificationBar(c *chocolate.Chocolate, useFlavour bool, opts ...chocolate.BaseBarOption) *NotifyModel {
	nmodel := &NotifyModel{useFlavour: useFlavour}
	nmodel.initDefaultNotificationTypes()

	notificaitonBar := chocolate.NewModelBar(
		&chocolate.BarModel{
			Model:                   nmodel,
			UpdateHandlerFct:        notifyModelUpdateHandler,
			FlavourCustomizeHandler: notifyModelFlavourCustomizeHandler,
		},
		chocolate.WithBarYScaler(chocolate.DYNAMIC, 0),
		chocolate.WithBarXScaler(chocolate.DYNAMIC, 0),
		chocolate.WithBarXPlacer(chocolate.CENTER, 0),
		chocolate.WithBarYPlacer(chocolate.START, 0),
		chocolate.WithBarID("overlay"),
	)
	notificaitonBar.Hide(true)
	for _, opt := range opts {
		opt(notificaitonBar)
	}

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
