package logout

import (
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/janik6n/azlogin/internal/configuration"
	"github.com/janik6n/azlogin/internal/logger"
)

const (
	currentAccount = "Current account"
	allAccounts    = "All accounts"
)

func RunCommand(c configuration.Configuration) (string, error) {
	funcName := "logout - RunCommand"

	var scope string
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions(currentAccount, allAccounts)...).
				Title("Logout").
				Description("Which accounts to log out?").
				Value(&scope),
		),
	).WithAccessible(accessible)

	if err := form.Run(); err != nil {
		return "", err
	}

	args := []string{"logout"}
	if scope == allAccounts {
		args = []string{"account", "clear"}
	}
	logger.LogInfo("Selected logout scope: "+scope, funcName, c)

	cmd := exec.Command("az", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return "\n" + lipgloss.NewStyle().
		Width(100).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(lipgloss.NewStyle().Bold(true).Render("Logged out: "+scope)), nil
}
