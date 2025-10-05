# popsink_env Resource

Manages a Popsink environment resource.

An environment in Popsink represents an isolated workspace where teams can manage their data pipelines. Environments can optionally be configured with message retention policies for data lifecycle management.

## Example Usage

### Basic Environment

```hcl
resource "popsink_env" "development" {
  name = "development"
}
```

### Environment with Retention

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

### Complete Example with Team Association

```hcl
# Create an environment
resource "popsink_env" "staging" {
  name          = "staging"
  use_retention = true

  retention_configuration = jsonencode({
    retention_ms   = 86400000 # 1 day in milliseconds
    cleanup_policy = "delete"
  })
}

# Associate a team with the environment
resource "popsink_team" "staging_team" {
  name        = "Staging Team"
  description = "Team managing staging pipelines"
  env_id      = popsink_env.staging.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the environment. This should be a unique, descriptive identifier for the environment.

* `use_retention` - (Optional) Whether message retention is enabled for this environment. Defaults to `false`. When set to `true`, you can optionally provide a `retention_configuration`.

* `retention_configuration` - (Optional) Retention policy configuration as a JSON string. This is only used when `use_retention` is `true`. The configuration should be a valid JSON object containing broker-specific retention settings.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the environment (UUID format).

## Retention Configuration

When `use_retention` is set to `true`, you can provide a `retention_configuration` as a JSON string. The exact structure of this configuration depends on your broker setup, but common fields include:

* `retention_ms` - Retention time in milliseconds
* `segment_ms` - Log segment size in milliseconds  
* `cleanup_policy` - Cleanup policy (e.g., "delete", "compact")

Example retention configuration:

```json
{
  "retention_ms": 604800000,
  "segment_ms": 86400000,
  "cleanup_policy": "delete"
}
```

## Import

Environments can be imported using their UUID:

```shell
terraform import popsink_env.example 550e8400-e29b-41d4-a716-446655440000
```

## Important Notes

* **Retention Configuration**: The retention configuration is stored as a JSON string in Terraform state. Make sure to use valid JSON when specifying this field.

* **Environment Names**: Environment names should be unique within your Popsink instance.