# Changelog

## Unreleased

## v0.1.0 — 2026-11-06

- Functions
  - `promql`: compiles a DQL query into a PromQL expression.
  - `opensearch`: compiles a DQL query into an OpenSearch search request body.
- Development
  - `make build`, which tests first, `make test`, `make docs`, `make snapshot` and `make install`.
  - CI runs the tests on every push and pull request.
  - A dispatched snapshot workflow publishes a signed pre-release that a configuration installs by pinning its exact version.
  - A dispatched release workflow stamps the changelog, tags the version, and publishes the signed archives the Terraform Registry reads.
