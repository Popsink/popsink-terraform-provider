# popsink_pipeline Resource

Manages a Popsink pipeline resource. A pipeline defines the data flow from a source to a target with optional transformations.

## Example Usage

### Basic Kafka to Oracle Pipeline

```hcl
resource "popsink_pipeline" "example" {
  name    = "my-data-pipeline"
  team_id = popsink_team.my_team.id
  state   = "draft"

  json_configuration = jsonencode({
    source_name = "kafka-source"
    source_type = "KAFKA_SOURCE"
    source_config = {
      bootstrap_servers = "kafka.example.com:9092"
      topic            = "input-topic"
      consumer_group   = "my-consumer-group"
    }
    target_name = "oracle-target"
    target_type = "ORACLE_TARGET"
    target_config = {
      host     = "oracle.example.com"
      port     = 1521
      database = "ORCL"
      user     = "myuser"
      password = "mypassword"
    }
    smt_name   = "my-transform"
    smt_config = []
    draft_step = "config"
  })
}
```

### Pipeline with Transformations

```hcl
resource "popsink_pipeline" "transform_example" {
  name    = "transform-pipeline"
  team_id = popsink_team.my_team.id
  state   = "draft"

  json_configuration = jsonencode({
    source_name = "source-connector"
    source_type = "KAFKA_SOURCE"
    source_config = {
      bootstrap_servers = "kafka.example.com:9092"
      topic            = "raw-data"
      consumer_group   = "transform-group"
    }
    target_name = "target-connector"
    target_type = "ORACLE_TARGET"
    target_config = {
      host     = "oracle.example.com"
      port     = 1521
      database = "PROD"
    }
    smt_name = "data-transformation"
    smt_config = [
      {
        function_type = "mapper"
        function_config = [
          {
            table_name = "users"
            fields = [
              {
                key      = "id"
                path     = "user.id"
                nullable = false
              },
              {
                key  = "name"
                path = "user.full_name"
              }
            ]
          }
        ]
      }
    ]
    draft_step = "config"
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the pipeline.
* `team_id` - (Required) The UUID of the team that owns the pipeline.
* `state` - (Required) The state of the pipeline. Must be one of:
  * `draft` - Pipeline is in draft mode
  * `paused` - Pipeline is paused
  * `live` - Pipeline is running
  * `error` - Pipeline has errors
  * `building` - Pipeline is being built
* `json_configuration` - (Required) The complete configuration of the pipeline as a JSON string.

### JSON Configuration Structure

The `json_configuration` must be a valid JSON string containing:

* `source_name` - (Required) Name of the source connector
* `source_type` - (Optional) Type of the source connector. Valid values: `JOB_SMT`, `KAFKA_SOURCE`, `ORACLE_TARGET`
* `source_config` - (Required) Configuration object for the source connector
* `target_name` - (Required) Name of the target connector
* `target_type` - (Optional) Type of the target connector. Valid values: `JOB_SMT`, `KAFKA_SOURCE`, `ORACLE_TARGET`
* `target_config` - (Required) Configuration object for the target connector
* `smt_name` - (Required) Name of the SMT (Simple Message Transform)
* `smt_config` - (Required) Array of transformation configurations
* `draft_step` - (Required) Current draft step (e.g., "config", "review")

### Source/Target Configuration Examples

#### Kafka Source Configuration

```json
{
  "bootstrap_servers": "kafka.example.com:9092",
  "topic": "my-topic",
  "consumer_group": "my-group",
  "security_protocol": "SASL_SSL",
  "sasl_mechanism": "SCRAM-SHA-256",
  "sasl_username": "user",
  "sasl_password": "password"
}
```

#### Oracle Target Configuration

```json
{
  "host": "oracle.example.com",
  "port": 1521,
  "database": "ORCL",
  "user": "username",
  "password": "password",
  "server_name": "XE",
  "server_id": "oraclesrv01"
}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the pipeline.
* `team_name` - The name of the team that owns the pipeline.

## Import

Pipelines can be imported using the pipeline ID:

```shell
terraform import popsink_pipeline.example 12345678-1234-1234-1234-123456789abc
```

## Validation

The provider performs the following validations:

- **State**: Must be one of: `draft`, `paused`, `live`, `error`, `building`
- **JSON Configuration**: Must be valid JSON
- **Connector Types**: If `source_type` or `target_type` are specified, they must be one of: `JOB_SMT`, `KAFKA_SOURCE`, `ORACLE_TARGET`

## Notes

- **State Changes**: When changing the `state` from `draft` to `live`, ensure the pipeline configuration is complete and valid.
- **Configuration Format**: The `json_configuration` is stored as a JSON string in Terraform state. Ensure valid JSON when specifying this field.
- **Transformations**: The configuration supports complex transformation pipelines with multiple SMT steps.