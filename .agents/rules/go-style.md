# Rule: Go Style Guide & Engineering Standards

Apply these conventions when writing, refactoring, or reviewing Go code,
tests, and protocol buffers in this repository.

---

## Source Code Representation

### Ordering

- **Declarations:** Strict declaration ordering within a Go file MUST be:
    1. `package`
    2. `import`
    3. `interface` definitions
    4. `type` definitions (structs, aliases)
    5. `const` declarations
    6. `var` declarations
    7. Methods and functions

### Formatting

- **Declaration Grouping:** Do NOT use block wrappers for declarations
  (e.g., `const (...)`, `var (...)`). Declare each constant or variable
  individually on its own line using an explicit `const` or `var` keyword.
- **80-Character Limit:** Every `.go` file is strictly formatted to a
  maximum of **80 characters per line**. This applies to logic, long
  `fmt.Errorf` chains, and anonymous closures.
    - **String Literals:** Long string literals MUST NOT be concatenated
      with `+` across multiple lines. Keep strings on a single line even if
      they exceed 80 characters to preserve literal formatting.
- **Goimports:** **DO NOT** run `goimports`. The `goimports` engine does
  not enforce the 80-character limit or project formatting styles.
  Manage imports and formatting manually.

---

## Lexical Elements

### Identifiers

- **Import Aliases:** Avoid aliasing imports unless resolving package
  name collisions or standard `go.chromium.org` AIP packages.
    - **Generated Protos:** NEVER alias generated protobuf packages.
      Use default generated package names directly.
    - When collisions require aliases, use `{location}{package}`
      (e.g., `internalaip160`, `gatewayaip132`).
- **Variables / Functions:** Use `camelCase` for unexported and
  `PascalCase` for exported members.
- **Map Variable Naming:** Map variables MUST be named explicitly as
  `map{KeyName}To{ValueName}` (e.g., `mapUserIDToAccountID`,
  `mapTranslationHashToDbID`).
- **Parameters:** Provide explicit types for every parameter
  (e.g., `(min int, max int)`, never `(min, max int)`).
- **Callee Validation:** Do NOT perform `nil` checks on pointer
  parameters (e.g., `if p == nil { return ... }`). Non-nil safety is the
  caller's responsibility.
- **Receivers:** Use short single-letter or two-letter abbreviations
  (e.g., `c *Client`, `o *OrderData`). Standard receiver abbreviations in this
  codebase MUST remain consistent (`g` for `PostgresGrammar`, `qb` for
  `QueryBuilder`).
- **Timestamps:** Use the `_at` suffix for timestamp fields and columns
  (e.g., `operations_ready_at`, `created_at`).
- **Contexts:** When deriving a context, prefix with `ctx` followed by a
  descriptive descriptor (e.g., `ctxLogger`, `ctxTarget`, `ctxInsert`).

---

## Types

### Struct Types (Models & DTOs)

- **Custom Unmarshaling:** Use the `type Alias T` pattern inside
  `UnmarshalJSON` to prevent infinite recursion during post-processing.

### Go Version & Generics

- **Language Target:** Target Go 1.27+ with generics throughout.
- **Type-Safe Columns:**
  - The `Column[Model]` interface provides `ColumnName()`, `PostgresCast()`,
    `ToTypedSlice()`, and `IsArray()` for type-aware SQL generation.
  - Strongly-typed column wrappers (e.g., `NumericColumn[Model, int64]`)
    produce correctly typed `contract.Where` and `contract.Order` nodes
    via fluent methods (`.Eq()`, `.Gt()`, `.In()`, `.IsNull()`).

---

## Declarations and Scope

### Documentation

- **Godocs:** All functions, methods, exported types, and package-private
  variables MUST have Godoc comments consisting of complete sentences
  starting with the element name.

---

## Statements

### Assignment Statements

- **Variable Assignments:** No inline `if` variables. Inline variable
  assignments inside `if` statements (`if err := fn(); err != nil`) are
  banned. Hoist variables above the `if` block.

### If Statements

- **Control Flow & Guard Clauses:** Avoid deeply nested conditional
  logic. Use `if condition { return/continue/break }` guard clauses at
  the start of functions/loops to flatten logical flow.

### Defer Statements

- **Cleanup:** Use `WarnContext` if a logger is available to log resource
  closing errors. Otherwise, explicitly assign to the blank identifier
  (e.g., `_ = f.Close()`).

---

## Errors and Logging

- **Explicit Checks:** Always check errors immediately after calls. All
  map lookups returning `ok` booleans must be checked.
- **Error Wrapping:** Native Go errors must NOT be complete sentences.
  Start with a lowercase letter, omit ending punctuation, and use `->` as
  the wrapping delimiter:
  `fmt.Errorf("failed to map base type -> %w", err)`
  Use `errors.New("...")` for static error strings.
- **Logging:** Use `slog` for structured logging (`ErrorContext`,
  `WarnContext`). Log messages MUST be complete sentences starting with a
  capital letter and ending with a period (e.g.,
  `logger.Error("The schema name is invalid.")`).
- **Log OR Return:** Never both log and return the same error.
    - **Exceptions:**
        1. Critical local state snapshots too large/sensitive for errors.
        2. Asynchronous error handoffs to channels or goroutines.
- **Error Handling Guide:**
  | Action | When to use it |
  | --- | --- |
  | Return `err` | Default behavior. Let caller decide. |
  | Wrap & Return | Adding context (e.g., "processing record X"). |
  | Log & Stop | Top level of application. |
  | Log & Continue | Non-critical errors (e.g., cache miss). |
- **Structured Attributes:** Pass contextual variables using strongly
  typed `slog` attributes (e.g., `slog.String("user_id", id)`). When an
  error is logged, `slog.Any("error", err)` MUST be the first attribute.

---

## Concurrency and Context Management

- **Parallelism:** Use `golang.org/x/sync/errgroup` for concurrent tasks
  that return errors.
- **Context Propagation:** Always propagate `context.Context`. Use
  `context.WithoutCancel` for cleanup/logging that must finish post-cancel.
- **Strict Timeouts:** All database queries, external API calls, and
  blocking I/O operations MUST use `context.WithTimeout`. Never reuse raw
  parent/global request contexts directly for execution calls.

---

## System Considerations (Database & API)

### API Design (AIP Compliance)

- **Filtering (AIP-160):** Use `pkg/orm/aip160` via `qb.WhereAip160(...)`.
  Filter strings are parsed via `go.chromium.org/luci` into an AST and
  walked with an iterative post-order traversal using a stack to build
  `contract.Where` trees via typed `Aip160Field` definitions.
- **Ordering (AIP-132):** Use `qb.OrderByAip132(...)`. Order strings are
  parsed via `go.chromium.org/luci` and mapped to `contract.Order` nodes
  through a column-name lookup map.
- **Pagination (AIP-158):** Use `qb.PaginateAip158(...)`. Offset-based
  pagination uses `go.einride.tech/aip` for page token encoding/decoding.
- **Query Builder:** Always use `QueryBuilder[Model, Mutator]`, typed
  columns, and fluent methods (`qb.Where(...)`, `column.Lt(...)`,
  `column.IsNull()`) rather than raw SQL strings.

### SQL Safety & Mutations

- **Identifier Sanitization:** All table, schema, and column names MUST
  always be sanitized via `pgx.Identifier{name}.Sanitize()`.
- **Positional Parameters:** All user values MUST be bound as positional
  parameters (`$1`, `$2`, ...) tracked by the `Parameterized` helper
  struct. Never interpolate raw values directly into SQL strings.
- **Bulk Mutations:** Bulk mutations (INSERT, UPDATE, UPSERT) MUST use
  PostgreSQL `UNNEST(...)` with typed array casts for efficient multi-row
  operations.

---

## Testing Standards

- **Unit Testing:** Unit tests sit alongside source files (`*_test.go`).
- **Table-Driven & Closures:** Use table-driven tests (`[]struct{...}`)
  or `t.Run("name", func(t *testing.T) { ... })` subtest closures.
- **Assertions:** Use `github.com/stretchr/testify/assert` and `require`.
  Use `require` for setup/fatal checks, `assert` for properties.
- **Testing Packages:**
  - **White Box:** Use `package mypkg` in `{{file}}_test.go` for internals.
  - **Black Box:** Use `package mypkg_test` to test public surface and
    prevent import cycles.
- **Organization & Lifecycle:** Use `TestMain` in `main_test.go` for shared
  container lifecycle management and template DB preparation.
- **Integration:** Use `testcontainers-go` with `PostgresContainer` from
  `pkg/testutil/`. Template DB cloning ensures each test gets an isolated
  database with the schema pre-applied, rather than table truncation.
- **Timestamp Fidelity:** Use `.Truncate(time.Microsecond)` on timestamps
  to ensure stable equality comparisons across DB/JSON round-trips.
