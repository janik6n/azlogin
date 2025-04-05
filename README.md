![Go](https://shields.io/badge/Go-00ADD8?logo=Go&logoColor=FFF&style=flat-square)
![Azure](https://shields.io/badge/Azure-0078D4?logo=Azure&logoColor=FFF&style=flat-square)

# Azlogin

Azure CLI login helper. You have multiple tenants to login to, but cannot remember the tenant IDs? Azlogin to the rescue!

Functionally & securitywise there is nothing too special; This is just a wrapper for Azure CLI. All the information stored about your tenants is defined [below](#configuration), and all the information is provided by you.

Running the app with the selected tenant will trigger `az login` flow for the selected tenant, and that's it.

If you have configured `select_subscription: true` after successful login, a Subscription selection will be run.


## Prerequisites

[Azure CLI](https://learn.microsoft.com/en-us/cli/azure/) v2 is expected to be installed and configured.


## Configuration

Configuration is handled in `configuration.yaml`

```yaml
general:
  environment: DEV|PROD
  logging: true|false
  logging_level: INFO|WARNING|ERROR|FATAL
  print_config: true|false
features:
  azlogin:
    tenants:
      - tenant_name: "Alpha"
        tenant_id: "alpha.onmicrosoft.com"
      - tenant_name: "Bravo"
        tenant_id: "23456-23456"
    select_subscription: true|false
```
Configuration file location:
- If the `environment` is `DEV`: `./configuration.yaml`. You also need to set environment variable `ENVIRONMENT=DEV`.
- If the `environment` is `PROD`: `$HOME/azlogin/configuration.yaml` in macOS/Linux and `%USERPROFILE%\azlogin\configuration.yaml` in Windows.

## Logging

The app optionally logs to a file `./azlogin.log` for `DEV` and to `$HOME/azlogin/azlogin.log` in macOS/Linux or `%USERPROFILE%\azlogin/azlogin.log` in Windows for `PROD`.


## Run binary release

With the binary releases all you need is the configuration file described [above](#configuration) and a binary for your platform.

⚠️ The binaries in each release are built with [GoReleaser](https://goreleaser.com/ci/actions/) and GitHub Actions. They are *not signed or notarized*, so warnings may arise when you run the app. Run at your own risk.


## How to run in dev

Before installing:

```bash
ENVIRONMENT=DEV go run cmd/azlogin/main.go
```

### Build and install

Build, and installation, must be done where the `main` package is located.

```bash
# To build
cd cmd/azlogin
# For architecture where developed:
go build -o "../../build/azlogin"
# For macOS ARM64
GOARCH=arm64 GOOS=darwin go build -o "../../build/azlogin"
# For Windows AMD64
GOARCH=amd64 GOOS=windows go build -o "../../build/azlogin.exe"
```
There is a little helper script available for installation on macOS & ZSH. To install, run in project root:
```bash
# Runs the build too.
./install_azlogin.sh
```
The binary will be installed to location defined in `.zshrc`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

## Automation

New GitHub Release is created with [GoReleaser](https://goreleaser.com/ci/actions/) and GitHub Actions when a new tag is pushed to the repository.


## References

- [Azure CLI](https://learn.microsoft.com/en-us/cli/azure/)
- [Package index for Azure SDK libraries for Go | Microsoft Learn](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-library-package-index)


### Packages used

- https://github.com/Azure/azure-sdk-for-go (License: MIT)
- https://github.com/charmbracelet/huh for terminal UI (License: MIT)
- https://github.com/charmbracelet/huh/spinner for terminal UI (License: MIT)
- https://github.com/charmbracelet/lipgloss for terminal UI (License: MIT)
- https://github.com/go-yaml/yaml for reading and writing YAML files (License: MIT & Apache 2.0)
- https://github.com/natefinch/lumberjack for logging (License: MIT)


## Changelog

[CHANGELOG](CHANGELOG.md)


## License

[MIT](LICENSE)
