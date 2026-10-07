# CHANGELOG

All notable changes to this template will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.5] - 2026-10-07

### Added

### Changed

### Deprecated

### Fixed

- Subscription selection lists only subscriptions belonging to the tenant logged into. Tenants configured by domain name are resolved to their tenant ID.

### Removed

### Security

- Subscription selection uses only Azure CLI credentials (`AzureCLICredential`) scoped to the selected tenant, instead of `DefaultAzureCredential`, which could pick up environment or managed identity credentials.

### Internal

- `github.com/google/uuid` is now a direct dependency.

## [1.1.4] - 2026-10-06

### Added

### Changed

### Deprecated

### Fixed

### Removed

### Security

- GitHub Actions pinned to commit hashes.

### Internal

- Dependency version bumps.
- GitHub Actions Action version bumps.

## [1.1.3] - 2026-02-14

### Added

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal

- Dependency version bumps.
- GitHub Actions Action version bumps.

## [1.1.2] - 2025-09-18

### Added

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal

- Dependency version bumps.

## [1.1.1] - 2025-04-05

### Added

### Changed

### Deprecated

### Fixed

### Removed

### Security

- Fixed security issues in transitive dependencies:
  - Fixed golang.org/x/net: HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net
  - Fixed github.com/golang-jwt/jwt/v5: allows excessive memory allocation during header parsing

### Internal

- Transitive dependency version bumps.

## [1.1.0] - 2025-04-05

### Added

- Optionally let user choose Azure Subscription after successful login.
- Improved logging.

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal

## [1.0.6] - 2025-04-03

### Added

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal

- Dependency version bumps.

## [1.0.5] - 2024-11-30

### Added

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal

- Dependency version bumps.

## [1.0.4] - 2024-08-27

### Added

- Improved docs.

### Changed

### Deprecated

### Fixed

- Removed extra debug prints.

### Removed

### Security

### Internal

## [1.0.3] - 2024-08-26

### Added

- Improved docs.

### Changed

- Improved error messages for configuration loading.

### Deprecated

### Fixed

### Removed

### Security

### Internal

## [1.0.2] - 2024-08-26

### Added

### Changed

### Deprecated

### Fixed

- The configuration should now load correctly on Windows.

### Removed

### Security

### Internal

## [1.0.1] - 2024-08-25

### Added

- Updated docs.

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal


## [1.0.0] - 2024-08-25

Initial release of Azlogin.

### Added

- az login flow.
- About flow.

### Changed

### Deprecated

### Fixed

### Removed

### Security

### Internal
