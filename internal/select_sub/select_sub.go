package selectsub

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/charmbracelet/huh"

	"github.com/janik6n/azlogin/internal/account"
	"github.com/janik6n/azlogin/internal/configuration"
	"github.com/janik6n/azlogin/internal/logger"
)

func RunCommand(tenantId string, c configuration.Configuration) (string, error) {
	funcName := "select_sub - RunCommand"

	logline := fmt.Sprintf("Running sub selection for tenantId: %s", tenantId)
	logger.LogInfo(logline, funcName, c)

	fmt.Printf("\nFetching Subscriptions for tenantId: %s\n", tenantId)

	tenantGUID, err := account.TenantGUID(tenantId)
	if err != nil {
		return "", errors.New("No valid Azure CLI session for tenant, " + err.Error())
	}

	// List Azure subscriptions
	cred, err := azidentity.NewAzureCLICredential(&azidentity.AzureCLICredentialOptions{TenantID: tenantGUID})
	if err != nil {
		return "", errors.New("Could not get Azure CLI credentials, " + err.Error())
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
			if subscription.TenantID == nil || !strings.EqualFold(*subscription.TenantID, tenantGUID) {
				continue
			}
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
					Description("Which Subscription to select? Press / to search.").
					Filtering(true).
					Height(15).
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
	subscriptionID := parts[1]

	logger.LogInfo("az account set --subscription "+subscriptionID, funcName, configuration.Configuration{})
	if err := account.SetSubscription(subscriptionID); err != nil {
		return "", err
	}

	a, err := account.Show()
	if err != nil {
		return "", err
	}

	return account.Render("Subscription selected", a), nil
}
