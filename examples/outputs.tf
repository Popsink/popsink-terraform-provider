output "environment_id" {
  value       = popsink_env.production.id
  description = "The ID of the production environment"
}

output "data_team_id" {
  value       = popsink_team.data_team.id
  description = "The ID of the data engineering team"
}

output "analytics_team_id" {
  value       = popsink_team.analytics_team.id
  description = "The ID of the analytics team"
}

output "data_pipeline_id" {
  value       = popsink_pipeline.data_ingestion.id
  description = "The ID of the data ingestion pipeline"
}

output "analytics_pipeline_id" {
  value       = popsink_pipeline.analytics_reports.id
  description = "The ID of the analytics pipeline"
}
