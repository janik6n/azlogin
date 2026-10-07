package azlogin

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/janik6n/azlogin/internal/about"
	"github.com/janik6n/azlogin/internal/account"
	azlogin "github.com/janik6n/azlogin/internal/az_login"
	"github.com/janik6n/azlogin/internal/configuration"
	errorhandler "github.com/janik6n/azlogin/internal/error_handler"
	"github.com/janik6n/azlogin/internal/logger"
	"github.com/janik6n/azlogin/internal/logout"
	selectsub "github.com/janik6n/azlogin/internal/select_sub"
	"github.com/janik6n/azlogin/internal/utils"
)

type Flow struct {
	SelectedFlow string
}

func RunCLI() {
	funcName := "azlogin - RunCLI"
	// Init the azloginilable commands
	var flowChoices = []string{}

	// Read configuration
	config, err := configuration.ReadConfiguration()
	if err != nil {
		errorhandler.HandleFatalWithoutLogger("Error loading configuration", err, configuration.Configuration{})
	}
	// prinf configuration if print_config is enabled
	if config.General.PrintConfig {
		fmt.Print(config.Print())
	}

	// Setup logger
	err = logger.SetupLogger(*config)
	if err != nil {
		errorhandler.HandleFatalWithoutLogger("Error setting up logger", err, *config)
	} else {
		logger.LogInfo("----- Begin main -----", funcName, *config)
		logger.LogInfo("Logger set up successfully", funcName, *config)
		logger.LogInfo("Configuration loaded successfully:\n"+config.Print(), funcName, *config)
	}

	if _, err := exec.LookPath("az"); err != nil {
		errorhandler.HandleFatal(
			"Azure CLI is required, install it from https://learn.microsoft.com/en-us/cli/azure/install-azure-cli",
			errors.New("Azure CLI (az) was not found in PATH"),
			*config,
		)
	}

	if len(config.Features.AzLogin.Tenants) > 0 {
		flowChoices = append(flowChoices, "Login to Azure cli")
	}

	var flow Flow

	flowChoices = append(flowChoices, "Select subscription", "Show current login", "Logout", "About")

	var options = huh.NewOptions(flowChoices...)

	// Log available flows
	logger.LogInfo("Available flows: "+utils.SliceOfStringsToString(flowChoices), funcName, *config)

	// Should we run in accessible mode?
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	flowSelectionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(options...).
				Title("Choose command").
				Description("What do you want to do?").
				Value(&flow.SelectedFlow),
		),
	).WithAccessible(accessible)

	err = flowSelectionForm.Run()
	if err != nil {
		errorhandler.HandleFatal("running commandSelectionForm on main", err, *config)
	}

	switch flow.SelectedFlow {
	case "Login to Azure cli":
		logger.LogInfo("Selected flow: Login to Azure cli", funcName, *config)
		azLoginResponse, tenantId, err := azlogin.RunCommand(*config)
		if err != nil {
			errorhandler.HandleError("running Login to Azure cli flow", err, *config)
		}
		// would log garbage because of formatting
		// logger.LogInfo("AWS Login response: "+awsLoginResponse, *config)
		fmt.Println(azLoginResponse)

		// if Features.AzLogin.SelectSubscription is enabled
		if config.Features.AzLogin.SelectSubscription {
			// run select subscription flow
			logger.LogInfo("Select subscription", funcName, *config)
			selectSubscriptionResponse, err := selectsub.RunCommand(tenantId, *config)
			if err != nil {
				errorhandler.HandleError("running Select subscription flow", err, *config)
			}

			fmt.Println(selectSubscriptionResponse)
		}
	case "Select subscription":
		logger.LogInfo("Selected flow: Select subscription", funcName, *config)
		tenantId := ""
		// az account show succeeds from cache even when the session has expired, so verify the token too.
		if current, err := account.Show(); err == nil {
			if _, err := account.TenantGUID(current.TenantID); err == nil {
				tenantId = current.TenantID
			}
		}
		if tenantId == "" {
			logger.LogInfo("No active session, running Login to Azure cli flow", funcName, *config)
			fmt.Println("No active Azure CLI session, login first.")
			azLoginResponse, loggedInTenantId, err := azlogin.RunCommand(*config)
			if err != nil {
				errorhandler.HandleError("running Login to Azure cli flow", err, *config)
			}
			fmt.Println(azLoginResponse)
			tenantId = loggedInTenantId
		}
		selectSubscriptionResponse, err := selectsub.RunCommand(tenantId, *config)
		if err != nil {
			errorhandler.HandleError("running Select subscription flow", err, *config)
		}
		fmt.Println(selectSubscriptionResponse)
	case "Show current login":
		logger.LogInfo("Selected flow: Show current login", funcName, *config)
		current, err := account.Show()
		if err != nil {
			errorhandler.HandleError("getting current login", err, *config)
		}
		fmt.Println(account.Render("Current login", current))
	case "Logout":
		logger.LogInfo("Selected flow: Logout", funcName, *config)
		logoutResponse, err := logout.RunCommand(*config)
		if err != nil {
			errorhandler.HandleError("running Logout flow", err, *config)
		}
		fmt.Println(logoutResponse)
	case "About":
		logger.LogInfo("Selected flow: About", funcName, *config)
		about, err := about.ShowAbout()
		if err != nil {
			errorhandler.HandleError("running About flow", err, *config)
		}

		fmt.Println(about)
	}

	logger.LogInfo("----- Fin. -----", funcName, *config)
}
