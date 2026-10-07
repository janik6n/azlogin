package account

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Account struct {
	SubscriptionName    string `json:"name"`
	SubscriptionID      string `json:"id"`
	TenantID            string `json:"tenantId"`
	TenantDisplayName   string `json:"tenantDisplayName"`
	TenantDefaultDomain string `json:"tenantDefaultDomain"`
	User                struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"user"`
}

// Show returns the active Azure CLI account.
func Show() (Account, error) {
	var a Account
	out, err := exec.Command("az", "account", "show", "--output", "json").Output()
	if err != nil {
		return a, fmt.Errorf("could not get current account, are you logged in? %w", err)
	}
	if err := json.Unmarshal(out, &a); err != nil {
		return a, err
	}
	return a, nil
}

// TenantGUID fails if Azure CLI has no valid session for the tenant, so it doubles as a session check.
func TenantGUID(tenantId string) (string, error) {
	out, err := exec.Command("az", "account", "get-access-token", "--tenant", tenantId, "--query", "tenant", "--output", "tsv").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// SubscriptionIDsInTenant returns the subscriptions Azure CLI knows of in the tenant.
func SubscriptionIDsInTenant(tenantGUID string) ([]string, error) {
	out, err := exec.Command("az", "account", "list", "--output", "json").Output()
	if err != nil {
		return nil, err
	}
	var subs []Account
	if err := json.Unmarshal(out, &subs); err != nil {
		return nil, err
	}
	var ids []string
	for _, s := range subs {
		if strings.EqualFold(s.TenantID, tenantGUID) {
			ids = append(ids, s.SubscriptionID)
		}
	}
	return ids, nil
}

func SetSubscription(subscriptionID string) error {
	cmd := exec.Command("az", "account", "set", "--subscription", subscriptionID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Render(title string, a Account) string {
	keyword := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}

	tenant := a.TenantID
	if name := a.TenantDisplayName; name != "" {
		tenant = fmt.Sprintf("%s (%s)", name, a.TenantID)
	} else if domain := a.TenantDefaultDomain; domain != "" {
		tenant = fmt.Sprintf("%s (%s)", domain, a.TenantID)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb,
		"%s\n\n✨ Subscription: %s 💫\n📌 Subscription ID: %s\n🏢 Tenant: %s\n👤 User: %s (%s)",
		lipgloss.NewStyle().Bold(true).Render(title),
		keyword(a.SubscriptionName),
		a.SubscriptionID,
		tenant,
		a.User.Name,
		a.User.Type,
	)

	return "\n" + lipgloss.NewStyle().
		Width(100).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String())
}
