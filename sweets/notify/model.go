package notify

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/mfulz/chocolate"
	"github.com/mfulz/chocolate/flavour"
)

type NotifyModel struct {
	useFlavour         bool
	notifyTypes        map[string]NotifyDefinition
	activeNotification *notification
}

func (m *NotifyModel) RegisterNewNotificationType(definition NotifyDefinition) {
	if m.notifyTypes == nil {
		m.notifyTypes = make(map[string]NotifyDefinition)
	}

	m.notifyTypes[definition.Level] = definition
}

func (m *NotifyModel) initDefaultNotificationTypes() {
	m.RegisterNewNotificationType(NotifyDefinition{
		Level: LEVEL_INFO,
		Color: COLOR_INFO,
	})
	m.RegisterNewNotificationType(NotifyDefinition{
		Level: LEVEL_WARN,
		Color: COLOR_WARN,
	})
	m.RegisterNewNotificationType(NotifyDefinition{
		Level: LEVEL_ERROR,
		Color: COLOR_ERROR,
	})
}

func (m NotifyModel) newNotification(level, msg string, duration time.Duration) *notification {
	ntype, ok := m.notifyTypes[level]
	if !ok {
		return nil
	}

	return &notification{
		uid:      uuid.NewString(),
		msg:      msg,
		color:    ntype.Color,
		duration: duration,
	}
}

func (m NotifyModel) Init() tea.Cmd                            { return nil }
func (m *NotifyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m NotifyModel) View() string {
	if m.activeNotification != nil {
		return m.activeNotification.msg
	}
	return ""
}

var notifyModelUpdateHandler = func(b chocolate.ChocolateBar, m tea.Model) func(tea.Msg) tea.Cmd {
	return func(msg tea.Msg) tea.Cmd {
		model := m.(*NotifyModel)
		switch msg := msg.(type) {
		case destroyMsg:
			if model.activeNotification != nil {
				if model.activeNotification.uid == string(msg) {
					b.Hide(true)
					model.activeNotification = nil
				}
			}

		case notifyMsg:
			model.activeNotification = model.newNotification(msg.level, msg.msg, msg.duration)
			if model.activeNotification != nil {
				b.Hide(false)
				return destroyCmd(model.activeNotification.uid, msg.duration)
			}
		}
		return nil
	}
}

var notifyModelFlavourCustomizeHandler = func(b chocolate.ChocolateBar, m tea.Model, s lipgloss.Style) func() lipgloss.Style {
	return func() lipgloss.Style {
		model := m.(*NotifyModel)
		if model.activeNotification == nil {
			return s
		}

		if model.useFlavour {
			return flavour.GetPresetNoErr(flavour.StylePreset(model.activeNotification.color))
		}
		color := lipgloss.Color(model.activeNotification.color)
		return s.Foreground(color).
			BorderForeground(color)
	}
}

type destroyMsg string

func destroyCmd(uid string, duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return destroyMsg(uid)
	})
}
