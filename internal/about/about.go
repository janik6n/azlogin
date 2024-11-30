package about

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func GetVersion() string {
	return "v1.0.5"
}

func ShowAbout() (string, error) {

	// Pretty print response
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n\n- version: %s\n- full CHANGELOG at %s",
		lipgloss.NewStyle().Bold(true).Render("About Azlogin"),
		keyword(GetVersion()),
		lipgloss.NewStyle().Underline(true).Render("https://github.com/janik6n/azlogin/blob/main/CHANGELOG.md"),
	)

	return "\n" + lipgloss.NewStyle().
		Width(90).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()), nil
}
