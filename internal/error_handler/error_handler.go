package errorhandler

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/janik6n/azlogin/internal/configuration"
	"github.com/janik6n/azlogin/internal/logger"
)

func HandleError(message string, err error, c configuration.Configuration) {
	funcName := "error_handler - HandleError"
	logger.LogError(fmt.Errorf(message, err), funcName, c)

	// Pretty print response
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n\n%s\n\nwhile %s",
		lipgloss.NewStyle().Bold(true).Render("🚫 Uh oh, we got an error!"),
		keyword(err.Error()),
		keyword(message),
	)

	fmt.Println(lipgloss.NewStyle().
		Width(80).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()))

	logger.LogInfo("----- Fin with errors. -----", funcName, c)
	os.Exit(1)
}

func HandleFatal(message string, err error, c configuration.Configuration) {
	handleFatalinternal(message, err, c, true)
}

func HandleFatalWithoutLogger(message string, err error, c configuration.Configuration) {
	handleFatalinternal(message, err, c, false)
}

func handleFatalinternal(message string, err error, c configuration.Configuration, hasLogger bool) {
	funcName := "error_handler - handleFatalinternal"
	if hasLogger {
		logger.LogFatal(fmt.Errorf(message, err), funcName, c)
	}

	// Pretty print response
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n\n%s\n\n%s",
		lipgloss.NewStyle().Bold(true).Render("🚫 Uh oh, we got a fatal error!"),
		keyword(message),
		keyword(err.Error()),
	)

	fmt.Println(lipgloss.NewStyle().
		Width(100).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()))

	if hasLogger {
		logger.LogInfo("----- Fin with fatal errors. -----", funcName, c)
	}
	os.Exit(1)
}
