package configuration

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"

	"github.com/janik6n/azlogin/internal/utils"
)

type Configuration struct {
	General  General  `yaml:"general"`
	Features Features `yaml:"features"`
}

type General struct {
	Environment  string `yaml:"environment"`
	Logging      bool   `yaml:"logging"`
	LoggingLevel string `yaml:"logging_level"`
	PrintConfig  bool   `yaml:"print_config"`
}

type Features struct {
	AzLogin AzLogin `yaml:"azlogin"`
}

type AzLogin struct {
	Tenants []Tenant `yaml:"tenants"`
}

type Tenant struct {
	TenantName string `yaml:"tenant_name"`
	TenantId   string `yaml:"tenant_id"`
}

func ReadConfiguration() (*Configuration, error) {
	// Check environment variable ENVIRONMENT. If it is DEV,
	// read DEV configuration file. Otherwise, read PROD configuration file.

	// Get user's home directory
	userHomeDirectory, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user's home directory, %v", err)
	}

	// Set the congifuration file location based on the environment
	configurationFilename := userHomeDirectory + "/azlogin/configuration.yaml"
	if os.Getenv("ENVIRONMENT") == "DEV" {
		configurationFilename = "./configuration.yaml"
	}

	f, err := os.Open(configurationFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var configuration Configuration

	// decode the yaml file using KnownFields
	dec := yaml.NewDecoder(f)
	// Catches unknown fields in yaml, does not catch missing keys
	dec.KnownFields(true)
	if err := dec.Decode(&configuration); err != nil {
		return nil, err
	}

	// Validate the struct
	if err := configuration.Validate(); err != nil {
		return nil, err
	}

	return &configuration, nil
}

// Validate the configuration struct
func (c Configuration) Validate() error {
	validationErrors := []string{}
	if c.General.Environment == "" {
		validationErrors = append(validationErrors, "general.environment is required")
	}
	if c.General.Logging && c.General.LoggingLevel == "" {
		validationErrors = append(validationErrors, "general.logging_level is required")
	}
	if c.General.Logging &&
		c.General.LoggingLevel != "INFO" &&
		c.General.LoggingLevel != "ERROR" &&
		c.General.LoggingLevel != "FATAL" {
		validationErrors = append(
			validationErrors,
			"general.logging_level must be one of INFO, ERROR, FATAL. Got: "+c.General.LoggingLevel,
		)
	}

	// Validate AzLogin
	if len(c.Features.AzLogin.Tenants) > 0 {
		for _, tenant := range c.Features.AzLogin.Tenants {
			if tenant.TenantName == "" {
				validationErrors = append(validationErrors, "features.azlogin.tenants.tenant_name is required")
			}
			if tenant.TenantId == "" {
				validationErrors = append(validationErrors, "features.azlogin.tenants.tenant_id is required")
			}
		}
	}

	// Compose error message
	if len(validationErrors) > 0 {
		return errors.New("Validation errors: " + utils.SliceOfStringsToString(validationErrors))
	}

	return nil
}

// Print the configuration struct
func (c Configuration) Print() string {

	composed := "Environment: " + c.General.Environment + "\n" +
		"Logging: " + strconv.FormatBool(c.General.Logging) + "\n" +
		"Logging Level: " + c.General.LoggingLevel + "\n"

	if len(c.Features.AzLogin.Tenants) > 0 {
		composed += "\nAzLogin is enabled. Tenants:\n"
		for _, tenant := range c.Features.AzLogin.Tenants {
			composed += "  Tenant Name: " + tenant.TenantName + "\n" +
				"    Tenant ID: " + tenant.TenantId + "\n"
		}
	}

	// return c as string
	return composed
}

// function returns Features.AzLogin.Tenants TenantNames as list of strings
func (c Configuration) GetAzTenantNames() []string {
	var tenantNames []string
	for _, tenant := range c.Features.AzLogin.Tenants {
		tenantNames = append(tenantNames, tenant.TenantName)
	}
	return tenantNames
}

// function finds a tenant by name and returns it
func (c Configuration) FindAzTenantByName(n string) (Tenant, error) {
	for _, tenant := range c.Features.AzLogin.Tenants {
		if tenant.TenantName == n {
			return tenant, nil
		}
	}
	return Tenant{}, errors.New("tenant not found")
}
