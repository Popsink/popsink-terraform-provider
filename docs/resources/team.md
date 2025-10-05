# popsink_team Resource

Manages a Popsink team resource. Teams are used to organize users and pipelines within an environment.

## Example Usage

### Basic Team

```hcl
resource "popsink_team" "example" {
  name        = "Data Engineering Team"
  description = "Team responsible for data ingestion and processing"
}
```

### Team Associated with Environment

```hcl
resource "popsink_team" "example" {
  name        = "Analytics Team"
  description = "Team focused on business intelligence and reporting"
  env_id      = popsink_env.production.id
}
```

### Multiple Teams in Different Environments

```hcl
resource "popsink_env" "production" {
  name          = "production"
  use_retention = true
}

resource "popsink_env" "staging" {
  name          = "staging"
  use_retention = false
}

resource "popsink_team" "prod_data_team" {
  name        = "Production Data Team"
  description = "Manages production data pipelines"
  env_id      = popsink_env.production.id
}

resource "popsink_team" "staging_data_team" {
  name        = "Staging Data Team"
  description = "Manages staging data pipelines"
  env_id      = popsink_env.staging.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the team. Must be unique within the organization.
* `description` - (Required) Short description of the team and its purpose.
* `env_id` - (Optional) The UUID of the environment the team is associated with. If not specified, the team will not be tied to a specific environment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier (UUID) of the team.

## Import

Teams can be imported using the team ID (UUID):

```shell
terraform import popsink_team.example 12345678-1234-1234-1234-123456789abc
```

### Finding Team IDs

To find the ID of an existing team, you can:

1. Use the Popsink web interface to view team details
2. Use the Popsink API to list teams:
   ```bash
   curl -H "Authorization: Bearer $POPSINK_TOKEN" \
        "https://your-popsink-instance.com/api/teams/"
   ```
3. Check existing Terraform state if the team was previously managed

### Import Example

```shell
# Import an existing team
terraform import popsink_team.data_engineering a1b2c3d4-e5f6-7890-abcd-ef1234567890

# After import, create the corresponding resource configuration
cat >> main.tf << EOF
resource "popsink_team" "data_engineering" {
  name        = "Data Engineering Team"
  description = "Team managing data pipelines"
  env_id      = "env-uuid-here"  # Optional, based on existing team
}
EOF

# Plan to see if configuration matches imported state
terraform plan
```

## Usage Notes

- Team names must be unique within your Popsink organization
- Teams can be created without being associated with a specific environment
- Once a team is created, pipelines can be assigned to it
- Deleting a team will also delete all associated pipelines
- Teams can be moved between environments by updating the `env_id`

## Example with Pipeline

```hcl
resource "popsink_env" "production" {
  name          = "production"
  use_retention = true
}

resource "popsink_team" "data_team" {
  name        = "Data Engineering"
  description = "Manages all data pipelines"
  env_id      = popsink_env.production.id
}

resource "popsink_pipeline" "user_data" {
  name    = "user-data-pipeline"
  team_id = popsink_team.data_team.id
  state   = "draft"
  
  json_configuration = jsonencode({
    source_name = "kafka-users"
    source_type = "KAFKA_SOURCE"
    # ... rest of configuration
  })
}
```
