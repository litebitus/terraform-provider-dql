output "mock_request_rate" {
  value = provider::dql::promql("__name__:mock_requests_total and env:mock | rate(5m) | sum() by job")
}

output "mock_request_rate_grammar_1_0" {
  value = provider::dql::promql("__name__:mock_requests_total and env:mock | rate(5m) | sum() by job", "1.0")
}
