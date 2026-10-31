output "mock_error_count" {
  value = jsondecode(provider::dql::opensearch(
    "app:mock and level:error | count() by app",
    { time_field = "@timestamp", group_limit = 20 },
  ))
}
