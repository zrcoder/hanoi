package style

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/zrcoder/tgame/pkg/style/color"
)

var (
	Success = lipgloss.NewStyle().Foreground(color.Green)
	Error   = lipgloss.NewStyle().Foreground(color.Red)
)
