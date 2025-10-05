# Popsink Terraform Provider Examples

This directory contains example Terraform configurations for the Popsink provider, organized in a modular structure for better maintainability and readability.

## File Structure

The configuration has been split into separate files by resource type:

### `providers.tf`
- Terraform configuration block with required providers
- Popsink provider configuration
- Environment variables setup instructions

### `environments.tf`
- Environment resource definitions
- Retention policy configurations
- Kafka connection settings

### `teams.tf`
- Team resource definitions
- Team descriptions and associations with environments
- Team-environment relationships

### `pipelines.tf`
- Pipeline resource definitions
- Data ingestion and analytics pipeline configurations
- Source and target configurations
- Pipeline states (live/draft)

### `outputs.tf`
- Output values for all created resources
- Resource ID exports for reference in other configurations

### `main.tf.backup`
- Original monolithic configuration file (kept for reference)

## Usage

1. **Set up environment variables:**
   ```bash
   export POPSINK_TOKEN="your-popsink-token"
   export POPSINK_BASE_URL="https://onprem.ppsk.uk/api"  # optional
   ```

2. **Initialize Terraform:**
   ```bash
   terraform init
   ```

3. **Plan the deployment:**
   ```bash
   terraform plan
   ```

4. **Apply the configuration:**
   ```bash
   terraform apply
   ```

## Resources Created

This configuration creates:

- **1 Environment**: Production environment with retention policy
- **2 Teams**: Data Engineering Team and Analytics Team
- **2 Pipelines**: 
  - Data ingestion pipeline (live state)
  - Analytics reports pipeline (draft state)

## Dependencies

The resources are created in the following dependency order:

1. Environment (`popsink_env`)
2. Teams (`popsink_team`) - depends on environment
3. Pipelines (`popsink_pipeline`) - depends on teams

## Customization

You can customize the configuration by:

- Modifying environment settings in `environments.tf`
- Adding/removing teams in `teams.tf`
- Configuring different pipeline sources/targets in `pipelines.tf`
- Adding additional outputs in `outputs.tf`

## Benefits of Modular Structure

- **Maintainability**: Easier to find and modify specific resource types
- **Readability**: Clear separation of concerns
- **Reusability**: Individual files can be reused in other configurations
- **Collaboration**: Team members can work on different resource types simultaneously
- **Version Control**: More granular change tracking and easier code reviews