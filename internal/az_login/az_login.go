package azlogin

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/janik6n/azlogin/internal/configuration"
)

type Tenant struct {
	TenantName string
}

func RunCommand(c configuration.Configuration) (string, error) {
	fmt.Println("Running az login...")
	fmt.Println("Tenants:")
	fmt.Println(c.GetAzTenantNames())

	var tenantChoices = []string{}
	if len(c.Features.AzLogin.Tenants) > 0 {
		tenantChoices = c.GetAzTenantNames()
	} else {
		return "", errors.New("no az tenants configured")
	}

	var tenant Tenant
	var options = huh.NewOptions(tenantChoices...)

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	tenantSelectionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(options...).
				Title("Choose tenant").
				Description("Which tenant to login to?").
				Value(&tenant.TenantName),
		),
	).WithAccessible(accessible)

	err := tenantSelectionForm.Run()
	if err != nil {
		return "", err
	}

	fmt.Printf("\nSelected tenant:\n")
	fmt.Printf("  Tenant Name: %s\n", tenant.TenantName)
	t, err := c.FindAzTenantByName(tenant.TenantName)
	if err != nil {
		return "", err
	}
	fmt.Printf("  Tenant Id: %s\n", t.TenantId)

	response, err := AzLoginFlow(t)
	if err != nil {
		return "", err
	}

	return response, nil
}

func AzLoginFlow(t configuration.Tenant) (string, error) {
	loginCommand := fmt.Sprintf("az login --tenant %s", t.TenantId)
	fmt.Println(loginCommand)

	// Login to az cli, pass the command output directly to stdout & stderr
	args := strings.Split(loginCommand, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Pretty print response
	var sb strings.Builder
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}
	fmt.Fprintf(&sb,
		"%s\n\n✨ Tenant: %s 💫",
		lipgloss.NewStyle().Bold(true).Render("Az login completed."),
		keyword(t.TenantName),
	)

	return "\n" + lipgloss.NewStyle().
		Width(100).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()), nil
}
