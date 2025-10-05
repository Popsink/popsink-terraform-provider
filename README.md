# Terraform Provider for Popsink

[![Build Status](https://github.com/Popsink/popsink-terraform-provider/workflows/test/badge.svg)](https://github.com/Popsink/popsink-terraform-provider/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Popsink/popsink-terraform-provider)](https://golang.org/)

The Popsink Terraform Provider allows you to manage Popsink data pipelines using Infrastructure as Code.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25

## Using the Provider

### Installation

Add the provider to your Terraform configuration:

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
  # export POPSINK_BASE_URL="your-base-url" or use base_url
  base_url = "your-base-url"

  # export POPSINK_TOKEN="your-api-token" or use token
  token    = var.popsink_token
}
```

### Quick Example

```hcl
# Create an environment
resource "popsink_env" "production" {
  name          = "production"
}

# Create a team
resource "popsink_team" "analytics" {
  name        = "Analytics Team"
  description = "Analytics and reporting"
  env_id      = popsink_env.production.id
}

# Create a pipeline
resource "popsink_pipeline" "example" {
  name    = "my-pipeline"
  team_id = popsink_team.analytics.id
  state   = "draft"

  json_configuration = jsonencode({
    source_name = "postgres-source"
    source_type = "postgres"
    # ... additional configuration
  })
}
```

## Documentation

- **Resources**: See [docs/resources/](./docs/resources/) for detailed documentation on each resource
  - [popsink_env](./docs/resources/env.md)
  - [popsink_team](./docs/resources/team.md)
  - [popsink_pipeline](./docs/resources/pipeline.md)

- **Examples**: See [examples/](./examples/) for complete working configurations

## Development

To build and test the provider locally:

```bash
# Build the provider
make build

# Run tests
make test

# Install locally for testing
make install
```
