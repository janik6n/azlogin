# Azlogin

Azure CLI login helper. You have multipe tenants to login to, but cannot remember the tenant IDs? Azlogin to the rescue!


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
      - tenant_name: "alpha.onmicrosoft.com"
        tenant_id: "12345-12345"
      - tenant_name: "bravo.onmicrosoft.com"
        tenant_id: "23456-23456"
```
Configuration file location:
- If the `environment` is `DEV`: `./configuration.yaml`
- If the `environment` is `PROD`: `$HOME/azlogin/configuration.yaml`

## Logging

The app optionally logs to a file `./azlogin.log` in `DEV` and to `$HOME/azlogin/azlogin.log` in `PROD`.

## How to run in dev

Before installing:

```bash
ENVIRONMENT=DEV go run cmd/azlogin/main.go
```

## Build and install

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

# macOS & ZHS: To install, run in project root. Runs the build too.
./install_azlogin.sh
```

Install will install the binary to location defined in `.zshrc`:
```bash
export PATH="$HOME/go/bin:$PATH"
```

## Automation

New GitHub Release is created with [GoReleaser](https://goreleaser.com/ci/actions/) and GitHub Actions when a new tag is pushed to the repository.


## Run binary release

With the binary releases all you need is the configuration file described above and a binary for your platform.

The binaries in each release are built with [GoReleaser](https://goreleaser.com/ci/actions/) and GitHub Actions. They are not signed or notarized, so warnings may arise when you run the app. Run at your own risk.


## References

### Packages used

- https://github.com/charmbracelet/huh for terminal UI (License: MIT)
- https://github.com/charmbracelet/huh/spinner for terminal UI (License: MIT)
- https://github.com/charmbracelet/lipgloss for terminal UI (License: MIT)
- https://github.com/go-yaml/yaml for reading and writing YAML files (License: MIT & Apache 2.0)
- https://github.com/natefinch/lumberjack for logging (License: MIT)


## Changelog

[CHANGELOG](CHANGELOG.md)

## License

[MIT](LICENSE)
