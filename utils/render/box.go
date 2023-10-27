package render

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func Box(content string, color string) {

	width, _, _ := term.GetSize(0)
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(color)).
		Padding(0, 1, 0, 1).
		Width(width - 2)

	fmt.Println(style.Render(content))
}
