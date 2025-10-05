# Create pipelines for each team

resource "popsink_pipeline" "data_ingestion" {
  name    = "user-data-ingestion"
  team_id = popsink_team.data_team.id
  state   = "live"

  json_configuration = jsonencode({
    source_name = "postgres-users"
    source_type = "KAFKA_SOURCE"
    source_config = {
      host     = "db.example.com"
      port     = 5432
      database = "users"
      table    = "events"
    }
    target_name = "oracle-events"
    target_type = "ORACLE_TARGET"
    target_config = {
      host        = "oracle.example.com"
      port        = 1521
      database    = "ORCL"
      user        = "oracle_user"
      password    = "oracle_password"
      server_name = "XE"
      server_id   = "oraclesrv01"
    }
    smt_name   = "basic-transform"
    smt_config = []
    draft_step = "config"
  })
}

resource "popsink_pipeline" "analytics_reports" {
  name    = "weekly-analytics"
  team_id = popsink_team.analytics_team.id
  state   = "draft"

  json_configuration = jsonencode({
    source_name = "kafka-analytics"
    source_type = "KAFKA_SOURCE"
    source_config = {
      bootstrap_servers = "kafka.example.com:9092"
      topic             = "analytics-events"
      consumer_group    = "analytics-group"
    }
    target_name = "s3-reports"
    target_type = "ORACLE_TARGET"
    target_config = {
      bucket = "analytics-reports"
      prefix = "weekly/"
      region = "us-east-1"
    }
    smt_name   = "aggregation-transform"
    smt_config = []
    draft_step = "config"
  })
}
