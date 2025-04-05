package selectsub

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/janik6n/azlogin/internal/configuration"
	"github.com/janik6n/azlogin/internal/logger"
)

func RunCommand(tenantId string, c configuration.Configuration) (string, error) {
	funcName := "select_sub - RunCommand"

	logline := fmt.Sprintf("Running sub selection for tenantId: %s", tenantId)
	logger.LogInfo(logline, funcName, c)

	fmt.Printf("\nFetching Subscriptions for tenantId: %s\n", tenantId)

	// List Azure subscriptions
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", errors.New("No Azure credentials found within DefaultAzureCredential, " + err.Error())
	}
	// Create a context for the operation
	ctx := context.Background()

	// Create a client to interact with Azure subscriptions
	client, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		return "", errors.New("Could not create Azure subscriptions client, " + err.Error())
	}

	// List all subscriptions
	pager := client.NewListPager(nil)

	var subscriptionList = []string{}

	// Iterate through the pager to get subscriptions
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return "", errors.New("Error listing subscriptions, " + err.Error())
		}

		// Print subscription details
		for _, subscription := range page.Value {
			logger.LogInfo(fmt.Sprintf("Subscription ID: %s, Subscription Name: %s", *subscription.SubscriptionID, *subscription.DisplayName), funcName, c)
			subscriptionList = append(subscriptionList, fmt.Sprintf("%s | %s", *subscription.DisplayName, *subscription.SubscriptionID))
		}
	}

	if len(subscriptionList) > 0 {
		var selectedSubscription string
		var options = huh.NewOptions(subscriptionList...)
		accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

		subscriptionSelectionForm := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Options(options...).
					Title("Choose Subscription").
					Description("Which Subscription to select?").
					Value(&selectedSubscription),
			),
		).WithAccessible(accessible)

		err := subscriptionSelectionForm.Run()
		if err != nil {
			return "", err
		}

		logger.LogInfo("Selected subscription: "+selectedSubscription, funcName, c)

		response, err := SelectSubscriptionFlow(selectedSubscription)
		if err != nil {
			return "", err
		}
		return response, nil
	} else {
		return "", errors.New("No subscriptions found for tenantId: " + tenantId)
	}
}

func SelectSubscriptionFlow(s string) (string, error) {
	funcName := "select_sub - SelectSubscriptionFlow"

	// input is like "Subscription Name | Subscription ID"
	parts := strings.Split(s, " | ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid subscription format: %s. Expected name | id", s)
	}
	subscriptionName := parts[0]
	subscriptionID := parts[1]

	// Set the subscription using az cli
	setSubscriptionCommand := fmt.Sprintf("az account set --subscription %s", subscriptionID)
	logger.LogInfo(setSubscriptionCommand, funcName, configuration.Configuration{})

	// Login to az cli, pass the command output directly to stdout & stderr
	args := strings.Split(setSubscriptionCommand, " ")
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
		"%s\n\n✨ Subscription: %s 💫\n\n📌 ID: %s",
		lipgloss.NewStyle().Bold(true).Render("Subscription selected"),
		keyword(subscriptionName),
		subscriptionID,
	)

	return "\n" + lipgloss.NewStyle().
		Width(100).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2).
		Render(sb.String()), nil
}
