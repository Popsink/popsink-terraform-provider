# Terraform Provider for Popsink

[![Build Status](https://github.com/Popsink/popsink-terraform-provider/workflows/test/badge.svg)](https://github.com/Popsink/popsink-terraform-provider/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Popsink/popsink-terraform-provider)](https://golang.org/)
[![License](https://img.shields.io/github/license/Popsink/popsink-terraform-provider)](LICENSE)
[![GitHub release](https://img.shields.io/github/v/release/Popsink/popsink-terraform-provider)](https://github.com/Popsink/popsink-terraform-provider/releases)

The Popsink Terraform Provider allows you to manage Popsink data pipelines using Infrastructure as Code.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Authentication](#authentication)
- [Configuration](#configuration)
- [Resources](#resources)
- [Example Configurations](#example-configurations)
- [Development](#development)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for building the provider)

## Installation

### Using the Provider from Terraform Registry

Once published, you can use the provider by adding it to your Terraform configuration:

```hcl
terraform {
  required_providers {
    popsink = {
      source  = "popsink/popsink"
      version = "~> 1.0"
    }
  }
}

provider "popsink" {
  # Configuration options
}
```

## Quick Start

### Usage Example

For a complete working example, see [`examples/main.tf`](./examples/main.tf) which demonstrates:
- Creating an environment with retention policy
- Setting up two teams associated with the environment  
- Creating one pipeline per team

```

## Building from Source

If you want to build the provider from source:

```bash
git clone https://github.com/Popsink/popsink-terraform-provider
cd popsink-terraform-provider
go build -o terraform-provider-popsink
```

## Authentication

The provider requires an API token to authenticate with the Popsink API. You can provide the token in two ways:

### Environment Variable (Recommended)

```bash
export POPSINK_TOKEN="your-api-token"
```

### Provider Configuration

```hcl
provider "popsink" {
  token = "your-api-token"
}
```

## Configuration

The provider supports the following configuration options:

| Argument   | Description                                      | Required | Default                        | Environment Variable |
|------------|--------------------------------------------------|----------|--------------------------------|---------------------|
| `base_url` | The base URL for the Popsink API                | No       | `https://onprem.ppsk.uk/api`  | `POPSINK_BASE_URL`  |
| `token`    | The API token for authentication                | Yes      | -                              | `POPSINK_TOKEN`     |

## Resources

### `popsink_env`

Manages a Popsink environment resource.

#### Example Usage

```hcl
resource "popsink_env" "production" {
  name          = "production"
  use_retention = true

  retention_configuration = jsonencode({
    retention_ms   = 604800000 # 7 days in milliseconds
    segment_ms     = 86400000  # 1 day in milliseconds
    cleanup_policy = "delete"
  })
}
```

#### Argument Reference

- `name` (String, Required) - The name of the environment.
- `use_retention` (Boolean, Optional) - Whether message retention is enabled. Defaults to `false`.
- `retention_configuration` (String, Optional) - Retention policy configuration as a JSON string.

#### Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` (String) - The unique identifier of the environment.

#### Import

Environments can be imported using their ID:

```bash
terraform import popsink_env.example 550e8400-e29b-41d4-a716-446655440000
```

### `popsink_team`

Manages a Popsink team resource.

#### Example Usage

```hcl
resource "popsink_team" "analytics_team" {
  name        = "Analytics Team"
  description = "Team focused on analytics and reporting"
  env_id      = popsink_env.production.id
}
```

#### Argument Reference

- `name` (String, Required) - The name of the team.
- `description` (String, Required) - Short description of the team.
- `env_id` (String, Optional) - Environment ID the team is associated with.

#### Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` (String) - The unique identifier of the team.

#### Import

Teams can be imported using their ID:

```bash
terraform import popsink_team.example 550e8400-e29b-41d4-a716-446655440000
```

### `popsink_pipeline`

Manages a Popsink data pipeline.

#### Example Usage

```hcl
resource "popsink_pipeline" "example" {
  name    = "my-data-pipeline"
  team_id = popsink_team.analytics_team.id
  state   = "draft"

  json_configuration = jsonencode({
    source_name   = "postgres-source"
    source_type   = "postgres"
    source_config = {
      host     = "localhost"
      port     = 5432
      database = "mydb"
    }
    target_name   = "kafka-target"
    target_type   = "kafka"
    target_config = {
      bootstrap_servers = "localhost:9092"
      topic            = "my-topic"
    }
    smt_name   = "basic-transform"
    smt_config = []
    draft_step = "config"
  })
}
```

#### Argument Reference

- `name` (String, Required) - The name of the pipeline.
- `team_id` (String, Required) - The UUID of the team that owns the pipeline.
- `state` (String, Required) - The state of the pipeline. Valid values: `draft`, `paused`, `live`, `error`, `building`.
- `json_configuration` (String, Required) - The complete configuration of the pipeline as a JSON string.

#### Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `id` (String) - The unique identifier of the pipeline.
- `team_name` (String) - The name of the team that owns the pipeline.

#### Import

Pipelines can be imported using their ID:

```bash
terraform import popsink_pipeline.example 550e8400-e29b-41d4-a716-446655440000
```

## Example Configurations

See the [examples](./examples) directory for complete example configurations:

- [`main.tf`](./examples/main.tf) - Complete example with environment, teams, and pipelines
- [`env.tf`](./examples/env.tf) - Environment and resource creation examples
- [`team.tf`](./examples/team.tf) - Team management examples
- [`pipeline.tf`](./examples/pipeline.tf) - Pipeline configuration examples

## Development

For detailed development instructions, please see [CONTRIBUTING.md](CONTRIBUTING.md).

### Quick Start for Developers

#### Prerequisites

- Go 1.23 or later
- Terraform 1.0 or later
- golangci-lint (for linting)

#### Building the Provider

```bash
make build
```

#### Running Tests

```bash
# Run unit tests
make test

# Run tests with coverage
make coverage

# Run acceptance tests (requires valid API credentials)
make testacc
```

#### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all checks
make check
```

#### Installing Locally for Testing

```bash
# Build and install the provider locally
make install
```

This will install the provider to `~/.terraform.d/plugins/` for local testing.

### Available Make Targets

Run `make help` to see all available targets:

```bash
make help
```

## Contributing

We welcome contributions from the community! Here's how you can help:

- **Report bugs**: Open an [issue](https://github.com/Popsink/popsink-terraform-provider/issues/new?template=bug_report.yml) using our bug report template
- **Request features**: Submit a [feature request](https://github.com/Popsink/popsink-terraform-provider/issues/new?template=feature_request.yml)
- **Submit pull requests**: See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines

Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

### Development Resources

- [Contributing Guide](CONTRIBUTING.md) - Detailed contribution guidelines
- [Code of Conduct](CODE_OF_CONDUCT.md) - Community standards and expectations
- [Security Policy](SECURITY.md) - How to report security vulnerabilities
- [Changelog](CHANGELOG.md) - Release notes and version history

## Security

If you discover a security vulnerability, please review our [Security Policy](SECURITY.md) for responsible disclosure guidelines. **Do not** create a public issue for security vulnerabilities.

## Support

- **Documentation**: Check the [docs](./docs) directory for detailed resource documentation
- **Issues**: Search [existing issues](https://github.com/Popsink/popsink-terraform-provider/issues) or create a new one
- **Examples**: See the [examples](./examples) directory for complete working examples

## License

This provider is distributed under the Mozilla Public License 2.0. See [LICENSE](./LICENSE) for more information.

## Acknowledgments

This provider is maintained by the Popsink team and the community. Thank you to all [contributors](https://github.com/Popsink/popsink-terraform-provider/graphs/contributors)!
