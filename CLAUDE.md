# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build & Install

```sh
make build    # Build for Linux (amd64)
make install  # Build and install to ~/.terraform.d/plugins/
make release  # Multi-platform release builds (Darwin, Linux, Windows, etc.)
```

### Testing

```sh
make test     # Unit tests (go clean -testcache && go test ./...)
make testacc  # Acceptance tests against a live Couchbase instance
```

Acceptance tests require environment variables: `CB_ADDRESS`, `CB_CLIENT_PORT`, `CB_NODE_PORT`, `CB_USERNAME`, `CB_PASSWORD` and `TF_ACC=1`.

To run a single test:

```sh
go test ./couchbase/... -run TestAccBucketManager -v
```

### Lint

```sh
make lint  # Runs super-linter via Docker
```

### Local Couchbase (Docker)

```sh
make cbnetup  # Create Docker network
make cbup     # Start Couchbase via docker-compose
make cbinit   # Initialize cluster
make cbdown   # Stop Couchbase
```

## Architecture

The provider is in the `couchbase/` package. All resources follow the Terraform Plugin SDK v2 pattern.

**Entry points:**

- `main.go` — calls `plugin.Serve` with the provider factory
- `couchbase/provider.go` — defines all 7 resources, provider schema (address, username, password, TLS, SASL config), and builds the `gocb.Cluster` connection
- `couchbase/connection.go` — Couchbase cluster connection initialization and validation

**Resources** (one file each):

| Resource                        | File                           |
| ------------------------------- | ------------------------------ |
| `couchbase_bucket_manager`      | `resourceBucket.go`            |
| `couchbase_security_user`       | `resourceSecurityUser.go`      |
| `couchbase_security_group`      | `resourceSecurityGroup.go`     |
| `couchbase_query_index`         | `resourceQueryIndex.go`        |
| `couchbase_primary_query_index` | `resourcePrimaryQueryIndex.go` |
| `couchbase_bucket_scope`        | `resourceScope.go`             |
| `couchbase_bucket_collection`   | `resourceCollection.go`        |

**Supporting files:**

- `constants.go` — schema key constants shared across resources
- `validationBucket.go` — bucket-specific validation logic
- `validateProvider.go` — SASL mechanism and TLS certificate validation
- `helper.go` — diagnostic error helpers

Each resource file implements the standard `schema.Resource` with `Create`, `Read`, `Update`, `Delete` functions using the `gocb/v2` Couchbase client.

## Important Notes

- **Query index parallelism**: When applying query index resources, use `terraform apply -parallelism=1`. Couchbase rejects concurrent index creation requests on the same bucket.
- **Go version**: 1.26.2 (see `go.mod`)
- **Provider namespace**: `budisky.com/couchbase/couchbase` (local dev install path)
