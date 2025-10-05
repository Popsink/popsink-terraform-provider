# Create teams associated with the environment
resource "popsink_team" "data_team" {
  name        = "Data Engineering Team 2024"
  description = "Team responsible for data ingestion and processing"
  env_id      = popsink_env.production.id
}

resource "popsink_team" "analytics_team" {
  name        = "Analytics Team"
  description = "Team focused on business intelligence and reporting"
  env_id      = popsink_env.production.id
}
