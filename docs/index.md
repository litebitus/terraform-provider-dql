---
page_title: "dql Provider"
description: |-
  Compiles DQL queries into each backend's native query.
---

# dql Provider

Compiles DQL queries into each backend's native query at plan time. The provider has no configuration, resources or data sources: it is its functions.

Provider functions need Terraform 1.8 or later.

## Functions

| Function | Compiles a query into |
|---|---|
| [`promql`](functions/promql) | a PromQL expression |
| [`opensearch`](functions/opensearch) | an OpenSearch search request body |

A query is read in the latest grammar version the provider supports unless it states its own. The language is documented with the [DQL compiler](https://github.com/litebitus/dql-compiler/blob/main/docs/design.md).

## Example Usage

```terraform
terraform {
  required_providers {
    dql = {
      source = "litebitus/dql"
    }
  }
}
```
