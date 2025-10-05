resource "popsink_env" "production" {
  name          = "production-demo"
  use_retention = true

  retention_configuration = jsonencode({
    bootstrap_server  = "kafka.example.com:9092"
    security_protocol = "SASL_SSL"
    sasl_mechanism    = "SCRAM-SHA-256"
    sasl_username     = "kafka_user"
    sasl_password     = "kafka_password"
  })
}
