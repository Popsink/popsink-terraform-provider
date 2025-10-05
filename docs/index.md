---
page_title: "Popsink Provider"
subcategory: ""
description: |-
  The Popsink provider is used to interact with Popsink API resources to manage data pipelines.
---

# Popsink Provider

The Popsink provider allows you to manage Popsink data pipelines, teams, and environments using Terraform.

Use the navigation to the left to read about the available resources.

## Schema

### Required

- `base_url` (String) The base URL for the Popsink API. Can also be set via the `POPSINK_BASE_URL` environment variable.
- `token` (String, Sensitive) The API token for authenticating with the Popsink API. Can also be set via the `POPSINK_TOKEN` environment variable.

## Authentication

The provider requires both a base URL and an API token to authenticate with the Popsink API.

### Environment Variables

The recommended way to configure the provider is using environment variables:

```bash
export POPSINK_BASE_URL="your-base-url"
export POPSINK_TOKEN="your-api-token"
```

Then in your Terraform configuration:

```hcl
provider "popsink" {
  # Configuration will be read from environment variables
}
```

### Provider Configuration

You can also configure the provider directly in your Terraform configuration:

```hcl
provider "popsink" {
  base_url = "your-base-url"
  token    = var.popsink_token
}
```

**Note**: It is not recommended to hardcode the API token in your configuration. Use environment variables or Terraform variables instead.

## Resources

The following resources are available:

- [popsink_env](resources/env.md) - Manage Popsink environments
- [popsink_team](resources/team.md) - Manage Popsink teams
- [popsink_pipeline](resources/pipeline.md) - Manage Popsink data pipelines
