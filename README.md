# Scheme Runtime Go

A type-safe, generics-based Go runtime library for the [Scheme](https://github.com/uthereal) ORM framework. Provides a fluent query builder, PostgreSQL SQL compiler, and database execution layer with first-class support for Google API Improvement Proposals (AIPs).

## Features

- **Type-Safe Query Builder** — Fluent API powered by Go generics (`QueryBuilder[Model, Mutator]`)
- **PostgreSQL Compiler** — Generates parameterized SQL with `$N` bind variables and `pgx.Identifier` escaping
- **Full CRUD** — `Get`, `First`, `Insert`, `Update`, `Delete`, `Upsert` with `RETURNING` variants
- **Bulk Mutations** — Efficient multi-row operations using PostgreSQL `UNNEST(...)` with typed array casts
- **AIP-160 Filtering** — Parse Google API filter strings into composable WHERE conditions
- **AIP-132 Ordering** — Parse `order_by` strings into ORDER BY clauses
- **AIP-158 Pagination** — Offset-based pagination with page token encoding/decoding
- **Relationship Eager Loading** — `HasOne`, `HasMany`, `BelongsTo`, `BelongsToMany`
- **Rich Column Types** — Numeric, String, Array, JSON, Geo, Timestamp, UUID, Duration, Decimal, and nullable variants
- **Composable WHERE Conditions** — Basic, Between, Column, Date, Exists, In, JSON, Nested, Null, Raw, Subquery

## Installation

```bash
go get github.com/uthereal/scheme-runtime-go
```

## Architecture

```
pkg/contract/    → Interfaces and type definitions
pkg/orm/         → Query builder and fluent API
pkg/grammar/     → SQL compilation engine (PostgreSQL)
pkg/testutil/    → Test infrastructure (ephemeral Postgres via testcontainers)
examples/        → Generated-code examples and integration tests
```

### Quick Example

```go
// Create a query builder (typically from generated code)
qb := orm.NewQueryBuilder[User, UserMutator](
    db, grammar.NewPostgresGrammar(), userTable, hydrateUser, dehydrateUser,
)

// Fluent queries
users, err := qb.
    Where(UserColumns.Status.Eq("active")).
    OrderBy(UserColumns.CreatedAt.Desc()).
    Limit(10).
    Get(ctx)

// AIP-160 filter strings
qb, err := qb.WhereAip160(`status = "active" AND age >= 18`, fieldMap)

// Bulk insert with RETURNING
users, err := qb.InsertReturningMany(ctx, mutators)

// Eager loading
users, err := qb.With(UserRelations.Posts, UserRelations.Profile).Get(ctx)
```

## Testing

Unit and integration tests use ephemeral PostgreSQL containers via [testcontainers-go](https://github.com/testcontainers/testcontainers-go). Docker must be running.

```bash
go test ./...
```

## Dependencies

| Dependency | Purpose |
|---|---|
| [`pgx/v5`](https://github.com/jackc/pgx) | PostgreSQL driver |
| [`samber/lo`](https://github.com/samber/lo) | Functional utilities |
| [`chromium/luci`](https://chromium.googlesource.com/infra/luci/luci-go) | AIP-160/132 parsing |
| [`einride/aip`](https://github.com/einride/aip-go) | AIP-158 pagination |
| [`testcontainers-go`](https://github.com/testcontainers/testcontainers-go) | Ephemeral Postgres for tests |

## License

[MIT](LICENSE)
