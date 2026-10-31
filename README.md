# Terraform Provider DQL

Terraform functions that compile DQL queries into each backend's native query at plan time, with the [DQL compiler](https://github.com/litebitus/dql-compiler).

## Examples

| Function | Example | Documentation |
|---|---|---|
| `promql` | [`examples/functions/promql/`](examples/functions/promql) | [docs/functions/promql.md](docs/functions/promql.md) |
| `opensearch` | [`examples/functions/opensearch/`](examples/functions/opensearch) | [docs/functions/opensearch.md](docs/functions/opensearch.md) |

## Features

- `provider::dql::promql` compiles a query into a PromQL expression.
- `provider::dql::opensearch` compiles a query into an OpenSearch search request body.
- A query the compiler refuses fails the plan with the compiler's message, on the argument that is wrong.
- The grammar version defaults to the latest the provider's compiler supports, so the provider pin decides it; a query can state its own.

Each function's arguments are on the [Terraform Registry](https://registry.terraform.io/providers/litebitus/dql/latest/docs); the design is in [docs/design.md](docs/design.md).

## Usage

Provider functions need Terraform 1.8 or later.

```hcl
terraform {
  required_providers {
    dql = {
      source = "litebitus/dql"
    }
  }
}

output "mock_request_rate" {
  value = provider::dql::promql("__name__:mock_requests_total and env:mock | rate(5m) | sum() by job")
}
```

## Development

Go and Terraform build and test the repo; GoReleaser builds the snapshot archives.

| Target | Does |
|---|---|
| `make build` | runs the tests, then builds the provider into `bin/` |
| `make test` | `go vet`, a `gofmt` check, the Go tests including plans through Terraform, and `tfplugindocs validate` |
| `make docs` | generates the Registry pages in `docs/` from the functions and `examples/` |
| `make install` | builds the provider and prints the `dev_overrides` block that points local plans at it |
| `make snapshot` | runs the tests, then builds every platform's archive, unsigned, into `dist/` |

### Releases

The `Release` workflow, dispatched by hand, runs the tests, computes the next version, stamps `CHANGELOG.md`, tags it, and publishes the signed archives the Terraform Registry reads. It signs with the GPG key in the `GPG_PRIVATE_KEY` and `PASSPHRASE` secrets, whose public key is registered with the Registry.

The `Snapshot` workflow, dispatched by hand, publishes a signed pre-release such as `v0.1.1-snapshot.20261030175512`. A configuration installs it by pinning that exact version:

```hcl
dql = {
  source  = "litebitus/dql"
  version = "0.1.1-snapshot.20261030175512"
}
```
