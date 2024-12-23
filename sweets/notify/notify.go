package notify

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	LEVEL_INFO  = "Info"
	LEVEL_WARN  = "Warn"
	LEVEL_ERROR = "Error"
)

const (
	COLOR_INFO  = "10"
	COLOR_WARN  = "11"
	COLOR_ERROR = "9"
)

type notifyMsg struct {
	level    string
	msg      string
	duration time.Duration
}

type notification struct {
	uid      string
	msg      string
	color    string
	duration time.Duration
}

type NotifyDefinition struct {
	Level string
	Color string
}

func NewNotifyMsg(level, message string, duration time.Duration) tea.Cmd {
	return func() tea.Msg {
		return notifyMsg{
			level:    level,
			msg:      message,
			duration: duration,
		}
	}
}
