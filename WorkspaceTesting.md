# Workspace Module — Advanced Testing Implementation

## 1. Role and Engineering Expectations

Act as a **Senior Go Backend Engineer, Software Architect, and Test Automation Engineer** with extensive experience in:

- Production-grade Go backend development.
- Clean Architecture and layered application design.
- REST API testing and HTTP middleware.
- Unit testing and integration testing.
- PostgreSQL database testing.
- SQL query behavior and transaction management.
- Authentication and authorization testing.
- Idempotency and concurrency testing.
- Test isolation, reliability, and race-condition detection.
- CI/CD testing practices.

You are responsible for implementing comprehensive automated tests for the existing Workspace module.

Think and work like an engineer responsible for maintaining a production backend system.

Do not write superficial tests just to increase test coverage.

Your objective is to verify that the Workspace module behaves correctly under normal conditions, invalid inputs, failures, edge cases, and concurrent requests.

**Important:** Do not assume that the existing implementation is correct. Inspect it critically and identify potential bugs, missing safeguards, and architectural weaknesses.

---

# 2. Project Context

This is the Workspace module of AdOrbit, an advertising technology platform.

The backend is written in Go and follows a layered architecture.

The Workspace module currently contains:

```text
internal/workspace/
    handler.go
    handler_test.go

    service.go
    service_test.go

    repository.go
    repository_test.go

    model.go
    model_test.go

    routes.go
    routes_test.go
```

The module is responsible for creating and fetching advertiser and publisher workspaces.

### Current architecture

```text
HTTP Request
     |
     v
Routes
     |
     v
JWT Middleware
     |
     v
Idempotency Middleware (POST only)
     |
     v
Handler
     |
     v
Service Interface
     |
     v
Service Implementation
     |
     v
Repository Interface
     |
     v
Repository Implementation
     |
     v
PostgreSQL
```

### Architecture rules

- Handler is responsible for HTTP-level concerns.
- Service is responsible for business logic.
- Repository is responsible for database operations.
- Handler depends on the Service interface.
- Service depends on the Repository interface.
- Interfaces are defined at the consumer side.
- Repository implementation must not leak into the handler.
- Service implementation must not leak into the handler.
- Business logic should not be duplicated across layers.

Preserve this architecture.

Do not introduce unnecessary abstractions.

---

# 3. Mandatory Initial Inspection

Before writing any tests, inspect the complete existing implementation.

Read:

1. `handler.go`
2. `service.go`
3. `repository.go`
4. `model.go`
5. `routes.go`
6. Workspace-related database migrations.
7. Error classification implementation.
8. Response encoding and decoding utilities.
9. JWT middleware.
10. Idempotency middleware.
11. Existing testing utilities and project configuration.

Also inspect related modules if they provide useful patterns for:

- Database initialization.
- Test database setup.
- JWT generation.
- Middleware testing.
- Error assertions.
- SQL transaction handling.

### Do not skip this step.

Understand the actual implementation before defining expected results.

For example:

- Check actual validation tags instead of assuming validation behavior.
- Check actual SQL constraints instead of inventing database errors.
- Check actual idempotency middleware behavior instead of assuming duplicate requests return a particular status.
- Check the response structure before asserting response JSON.
- Check the error classifier before asserting HTTP status codes.

If an existing test expectation conflicts with the actual intended application contract, investigate it and document the discrepancy.

Do not silently change production behavior just to satisfy a test.

---

# 4. Testing Architecture

Implement different types of tests for different layers.

| Layer      | Testing Type           | Database                      |
| ---------- | ---------------------- | ----------------------------- |
| Handler    | Unit tests             | Mock Service                  |
| Service    | Unit tests             | Mock Repository               |
| Repository | Integration tests      | Real PostgreSQL test database |
| Routes     | HTTP integration tests | Isolated dependencies         |
| DTO        | Unit tests             | Not required                  |

Use Go's standard testing package wherever practical.

Preferred tools:

- `testing`
- `net/http/httptest`
- `context`
- `encoding/json`
- `database/sql`
- `errors`
- Existing project dependencies.

You may introduce a testing library if it provides meaningful benefits, but do not add dependencies unnecessarily.

Prefer small handwritten mocks over complicated mocking frameworks.

---

# 5. Handler Tests

File: `handler_test.go`

Handler tests must verify HTTP behavior independently of PostgreSQL.

Create a mock implementation of the existing `Service` interface.

The mock should support:

- Configurable Create response.
- Configurable Fetch response.
- Configurable errors.
- Tracking method calls.
- Tracking received UserID.
- Tracking received DTO.
- Verifying whether methods were called.

## A. Create Handler

Implement tests for:

| #   | Test Case                         | Expected Result           |
| --- | --------------------------------- | ------------------------- |
| 1   | Valid workspace creation          | 201 Created               |
| 2   | Missing authenticated principal   | 401 Unauthorized          |
| 3   | Invalid JSON payload              | 400 Bad Request           |
| 4   | Empty request body                | 400 Bad Request           |
| 5   | Missing workspace name            | 422 Validation Error      |
| 6   | Invalid workspace type            | 422 Validation Error      |
| 7   | Invalid country length            | 422 Validation Error      |
| 8   | Invalid currency length           | 422 Validation Error      |
| 9   | Invalid website URL               | 422 Validation Error      |
| 10  | Service returns conflict error    | 409 Conflict              |
| 11  | Service returns internal error    | 500 Internal Server Error |
| 12  | Service returns unavailable error | 503 Service Unavailable   |

Additional advanced test scenarios:

13. Service returns nil workspace without error.
14. Service returns workspace with nullable fields set to nil.
15. Service returns workspace with optional fields populated.
16. Request contains unknown JSON fields.
17. Request contains malformed JSON.
18. Request contains JSON with incorrect field types.
19. Request contains trailing invalid JSON data.
20. Service receives the correct authenticated UserID.
21. Service receives every expected DTO field.
22. Service is not called when authentication fails.
23. Service is not called when JSON decoding fails.
24. Service is not called when validation fails.
25. Verify service method call count.

For nil workspace responses, first inspect the current response contract and decide whether the handler should return an error or a successful response. Do not invent behavior.

## B. Fetch Handler

Implement tests for:

| #   | Test Case                                     | Expected Result             |
| --- | --------------------------------------------- | --------------------------- |
| 1   | User has multiple workspaces                  | 200 OK with workspace array |
| 2   | User has no workspaces                        | 200 OK with empty array     |
| 3   | Missing authenticated principal               | 401 Unauthorized            |
| 4   | Service returns internal error                | 500 Internal Server Error   |
| 5   | Service returns unavailable error             | 503 Service Unavailable     |
| 6   | Verify authenticated UserID passed to service | Correct UserID              |

Additional advanced test scenarios:

7. Service returns a single workspace.
8. Service returns workspaces containing nullable fields.
9. Service returns an unexpected error.
10. Verify returned workspace order is preserved.
11. Verify service is not called when principal is missing.
12. Verify service method call count.
13. Verify response JSON structure.
14. Verify empty result serialization behavior.

### Handler testing requirements

- Use `httptest.NewRequest`.
- Use `httptest.NewRecorder`.
- Add authenticated principals using the actual principal mechanism.
- Do not bypass handler authentication checks.
- Use the actual response utilities.
- Verify HTTP status and response body.
- Verify dependency interactions.

Do not use a real database in handler unit tests.

---

# 6. Service Tests

File: `service_test.go`

Service tests must verify business logic independently of the database.

Create a mock implementation of the existing `Repository` interface.

The mock must support:

- Configurable Create response.
- Configurable Fetch response.
- Configurable errors.
- Tracking received DTO.
- Tracking received UserID.
- Tracking generated slug.
- Tracking method calls.

## A. Create Service

Implement the following tests:

| #   | Test Case                           | Expected Result              |
| --- | ----------------------------------- | ---------------------------- |
| 1   | Valid workspace creation            | Workspace returned           |
| 2   | Workspace name converted into slug  | Correct slug                 |
| 3   | Name containing spaces              | Hyphens generated            |
| 4   | Name containing special characters  | Correct slug                 |
| 5   | Empty workspace name                | Empty slug behavior verified |
| 6   | Repository returns error            | Classified error returned    |
| 7   | Correct UserID passed to repository | Correct UserID               |
| 8   | All DTO fields passed correctly     | Correct arguments            |

Additional advanced scenarios:

9. Name containing uppercase characters.
10. Name containing leading and trailing spaces.
11. Name containing multiple consecutive spaces.
12. Name containing Unicode characters.
13. Name containing punctuation.
14. Name containing numbers.
15. Name containing only special characters.
16. Repository returns nil workspace and nil error.
17. Repository returns a workspace and an error, if the implementation permits this.
18. Verify repository Create is called exactly once.
19. Verify the original DTO is not unexpectedly modified.
20. Verify nullable DTO fields are passed correctly.
21. Verify repository errors retain the expected classification.
22. Verify context is propagated to the repository.

Use realistic examples such as:

```text
My Workspace       -> my-workspace
AdOrbit Publisher  -> adorbit-publisher
My Workspace 2026  -> my-workspace-2026
```

Check the actual slug library behavior before defining expected values.

Do not write assumptions about slug generation.

## B. Fetch Service

Implement:

| #   | Test Case                              | Expected Result           |
| --- | -------------------------------------- | ------------------------- |
| 1   | Repository returns multiple workspaces | All returned              |
| 2   | Repository returns empty slice         | Empty slice returned      |
| 3   | Repository returns error               | Classified error returned |
| 4   | Correct UserID passed to repository    | Correct UserID            |

Additional advanced scenarios:

5. Repository returns nil slice.
6. Repository returns a populated slice.
7. Verify repository method call count.
8. Verify workspace values remain unchanged.
9. Verify context propagation.
10. Verify error classification.
11. Verify returned ordering is preserved.

Do not add business rules that do not exist in the current service implementation.

---

# 7. Repository Integration Tests

File: `repository_test.go`

Repository tests must verify actual PostgreSQL interactions.

Use a dedicated PostgreSQL test database.

Do not use the development or production database.

### Test database requirements

- Use isolated database configuration.
- Apply the actual project migrations.
- Use deterministic test data.
- Clean up test records.
- Avoid shared mutable fixtures.
- Ensure tests can run repeatedly.
- Ensure tests do not depend on execution order.

Use the existing database initialization patterns wherever possible.

Do not create an entirely different database architecture for tests.

## A. Create Repository

Implement:

| #   | Test Case                        | Expected Result           |
| --- | -------------------------------- | ------------------------- |
| 1   | Valid advertiser workspace       | Successfully inserted     |
| 2   | Valid publisher workspace        | Successfully inserted     |
| 3   | Verify generated UUID            | Valid UUID returned       |
| 4   | Verify generated slug            | Correct slug stored       |
| 5   | Verify default status            | `active`                  |
| 6   | Verify created_at and updated_at | Timestamps populated      |
| 7   | Duplicate slug                   | Unique constraint error   |
| 8   | Invalid workspace type           | CHECK constraint error    |
| 9   | Invalid country length           | Database constraint error |
| 10  | Invalid currency length          | Database constraint error |
| 11  | Invalid created_by UUID          | Database error            |
| 12  | Database unavailable             | Error returned            |
| 13  | Context cancellation             | Query cancelled           |

Additional advanced scenarios:

14. Nullable industry field.
15. Nullable website field.
16. Nullable description field.
17. All nullable fields populated.
18. All nullable fields set to nil.
19. Verify returned record matches inserted data.
20. Verify timestamps are valid.
21. Verify generated UUID is unique.
22. Verify duplicate slug does not insert another record.
23. Verify database errors are returned without being swallowed.
24. Verify context deadline exceeded behavior.
25. Verify transaction rollback and cleanup behavior where applicable.

### Database constraints

Inspect the actual migrations before implementing constraint tests.

Do not assume that a constraint exists.

If a requested constraint is absent, document the missing constraint rather than inventing a test that expects it.

## B. Fetch Repository

Implement:

| #   | Test Case                              | Expected Result                  |
| --- | -------------------------------------- | -------------------------------- |
| 1   | User has one workspace                 | One returned                     |
| 2   | User has multiple workspaces           | All returned                     |
| 3   | User has no workspaces                 | Empty slice                      |
| 4   | Multiple users have workspaces         | Only requested user's workspaces |
| 5   | User has active and deleted workspaces | Only active returned             |
| 6   | Verify newest-first ordering           | DESC ordering maintained         |
| 7   | Database query error                   | Error returned                   |
| 8   | Row scanning error                     | Error returned                   |
| 9   | Rows iteration error                   | Error returned                   |
| 10  | Context cancellation                   | Query cancelled                  |

Additional advanced scenarios:

11. User has multiple active workspaces.
12. User has only deleted workspaces.
13. Different users have identical workspace names.
14. Nullable columns are scanned correctly.
15. Verify newest-first ordering with deterministic timestamps.
16. Verify user isolation.
17. Verify deleted records are excluded.
18. Verify empty results return the expected slice behavior.
19. Verify query errors are propagated.
20. Verify `rows.Err()` is handled.
21. Verify context deadline exceeded behavior.

### Important SQL testing note

Some database failures are difficult to reproduce reliably with real PostgreSQL.

For row scanning errors and rows iteration errors, use an appropriate SQL mocking approach if necessary.

Do not add artificial complexity to production repository code merely to simulate impossible database states.

Keep integration tests focused on actual PostgreSQL behavior.

---

# 8. Routes and Middleware Integration Tests

File: `routes_test.go`

These tests must verify the complete HTTP route chain.

Use:

- Actual route registration.
- Actual JWT middleware.
- Actual idempotency middleware.
- `httptest`.
- A mock Service implementation.
- Isolated idempotency storage where required.

Do not test only the handler directly in this file.

## A. Route Tests

| #   | Test Case                              | Expected Result                     |
| --- | -------------------------------------- | ----------------------------------- |
| 1   | POST without JWT                       | 401 Unauthorized                    |
| 2   | GET without JWT                        | 401 Unauthorized                    |
| 3   | POST with valid JWT                    | Request reaches Create handler      |
| 4   | GET with valid JWT                     | Request reaches Fetch handler       |
| 5   | POST without Idempotency-Key           | Rejected                            |
| 6   | POST with valid Idempotency-Key        | Accepted                            |
| 7   | POST with same idempotency key         | Duplicate request handled correctly |
| 8   | POST with different users and same key | Separate idempotency scopes         |
| 9   | GET without Idempotency-Key            | Works successfully                  |
| 10  | Unsupported HTTP method                | 405 Method Not Allowed              |

Additional advanced scenarios:

11. POST with expired JWT.
12. POST with malformed JWT.
13. GET with expired JWT.
14. GET with malformed JWT.
15. POST with empty Idempotency-Key.
16. POST with malformed Idempotency-Key.
17. POST with repeated key and identical payload.
18. POST with repeated key and different payload.
19. Different authenticated users using the same key.
20. Verify JWT identity reaches the handler.
21. Verify rejected requests never reach the handler.
22. Verify GET does not invoke idempotency middleware.
23. Verify unsupported methods do not invoke the service.

### Idempotency testing requirements

Inspect the existing idempotency middleware implementation.

Understand:

- Idempotency key scope.
- Request fingerprint behavior.
- Duplicate request behavior.
- Stored response behavior.
- Expiration behavior.
- Concurrent request behavior.
- Error handling.

Do not assume that every duplicate request must return the same status.

Assert the behavior actually guaranteed by the middleware contract.

If the middleware supports concurrency safety, add tests that submit simultaneous requests with the same idempotency key.

Verify that the same request does not create duplicate workspace records.

---

# 9. DTO Validation Tests

File: `model_test.go`

Test `CreateWorkspaceDTO` independently.

Use the same validation library and mechanism used by the application.

## A. CreateWorkspaceDTO

Implement:

| #   | Test Case                         | Expected Result        |
| --- | --------------------------------- | ---------------------- |
| 1   | All valid fields                  | Valid                  |
| 2   | Missing name                      | Invalid                |
| 3   | Name exceeding 150 characters     | Invalid                |
| 4   | Advertiser type                   | Valid                  |
| 5   | Publisher type                    | Valid                  |
| 6   | Unsupported type                  | Invalid                |
| 7   | Country with 2 characters         | Valid                  |
| 8   | Country with incorrect length     | Invalid                |
| 9   | Currency with 3 characters        | Valid                  |
| 10  | Currency with incorrect length    | Invalid                |
| 11  | Timezone exceeding 100 characters | Invalid                |
| 12  | Industry exceeding 100 characters | Invalid                |
| 13  | Valid website URL                 | Valid                  |
| 14  | Invalid website URL               | Invalid                |
| 15  | Nil optional fields               | Valid                  |
| 16  | Empty optional fields             | Verify actual behavior |

Additional advanced scenarios:

17. Name containing only whitespace.
18. Name containing leading/trailing whitespace.
19. Name at maximum allowed length.
20. Name one character above maximum length.
21. Country with lowercase characters.
22. Currency with lowercase characters.
23. Empty required country.
24. Empty required currency.
25. Empty required timezone.
26. Timezone at maximum allowed length.
27. Industry at maximum allowed length.
28. Website with HTTP scheme.
29. Website with HTTPS scheme.
30. Website with missing scheme.
31. Website with malformed URL.
32. Nil optional pointers.
33. Empty optional strings.
34. Populated optional strings.
35. Boundary values for every length-constrained field.

### Validation rules

Inspect the actual validation tags.

Do not change validation rules to make incorrect tests pass.

For pointer fields, explicitly test:

- Nil pointer.
- Pointer to empty string.
- Pointer to valid string.
- Pointer to invalid string.

Nil and empty values must not be treated as equivalent unless the validation implementation guarantees that behavior.

## B. UpdateWorkspaceDTO

Do not implement tests for `UpdateWorkspaceDTO` in this task.

It is currently outside the scope of this implementation.

---

# 10. Advanced Testing Standards

The following requirements apply to every test file.

## A. Table-Driven Tests

Use table-driven tests when several scenarios share the same test structure.

Example:

```go
tests := []struct {
    name string
    // input
    // expected output
    // expected error
}{
    {
        name: "valid workspace",
    },
    {
        name: "missing workspace name",
    },
}
```

Do not force table-driven tests where individual setup is clearer.

## B. Test Naming

Use descriptive test names.

Examples:

```go
func TestCreateHandler_Success(t *testing.T)

func TestCreateHandler_MissingPrincipal(t *testing.T)

func TestCreateService_GeneratesSlug(t *testing.T)

func TestFetchService_RepositoryError(t *testing.T)

func TestCreateRepository_DuplicateSlug(t *testing.T)

func TestWorkspaceRoutes_POSTWithoutJWT(t *testing.T)
```

## C. Assertions

Tests must verify actual behavior, not just the existence of an error.

Where applicable, verify:

- Exact HTTP status.
- Response JSON structure.
- Response message.
- Returned workspace fields.
- Error classification.
- Dependency call count.
- Correct UserID.
- Correct DTO.
- Correct slug.
- Correct SQL persistence.
- Database isolation.
- Middleware execution.
- Context propagation.

## D. Mock Design

Mocks must implement the existing interfaces:

- `Service`
- `Repository`

Do not change interface definitions unnecessarily.

Mocks should be deterministic.

Avoid global mutable mock state.

Avoid sharing mocks between parallel tests unless they are concurrency-safe.

## E. Test Isolation

Every test must be independently executable.

Requirements:

- No test-order dependency.
- No dependency on external production services.
- No hardcoded production credentials.
- No uncontrolled shared database state.
- Proper cleanup.
- No reliance on real external APIs.
- No random values without deterministic control.
- No unnecessary sleeps or timing-based assertions.

## F. Error Testing

Use `errors.Is` and `errors.As` where appropriate.

Verify that classified errors retain the expected semantics.

Do not compare error strings when a typed or sentinel error is available.

## G. Concurrency Testing

Where relevant, test:

- Concurrent idempotency requests.
- Concurrent repository operations.
- Shared state safety.
- Mock safety.

Do not add `t.Parallel()` indiscriminately.

Only parallelize tests that are genuinely isolated and safe.

---

# 11. Production Code Review During Testing

While implementing tests, inspect production code for test-relevant bugs.

Pay particular attention to:

1. Nil pointer handling.
2. Incorrect error propagation.
3. Missing context propagation.
4. Incorrect HTTP status mapping.
5. SQL resource cleanup.
6. Missing `rows.Err()` checks.
7. Incorrect nullable field scanning.
8. User isolation.
9. Incorrect middleware ordering.
10. Idempotency scope.
11. Race conditions.
12. Unnecessary business logic in handlers.

If a bug is discovered:

- Explain it.
- Identify the affected production file.
- Explain its impact.
- Propose the smallest correct fix.
- Implement the fix only when necessary for correctness and within the task scope.
- Add a regression test.

Do not refactor unrelated modules.

Do not change API contracts without explicit justification.

---

# 12. Test Execution

After implementation, execute the tests.

### Run all project tests

```bash
go test ./...
```

### Verbose output

```bash
go test -v ./...
```

### Race detection

```bash
go test -race ./...
```

### Coverage

```bash
go test -cover ./...
```

### Coverage profile

```bash
go test -coverprofile=coverage.out ./...
```

### Function-level coverage

```bash
go tool cover -func=coverage.out
```

### HTML coverage report

```bash
go tool cover -html=coverage.out
```

### Repeat tests to detect flaky behavior

```bash
go test -count=10 ./internal/workspace/...
```

Run database-dependent tests using the isolated test database configuration.

Do not report success for commands that were not actually executed.

If a command fails, investigate and fix the underlying issue rather than hiding or skipping the failure.

---

# 13. Coverage Expectations

Aim for strong meaningful coverage.

Suggested targets:

| Layer          | Target |
| -------------- | -----: |
| Handler        |   90%+ |
| Service        |   95%+ |
| Repository     |   85%+ |
| DTO Validation |   95%+ |
| Routes         |   90%+ |

These are targets, not reasons to write meaningless tests.

Do not manipulate coverage by testing trivial lines repeatedly.

Prioritize:

- Business-critical paths.
- Error paths.
- Security boundaries.
- User isolation.
- Database consistency.
- Middleware behavior.

Explain uncovered code that is difficult or inappropriate to test.

---

# 14. Final Deliverables

After completing the work, provide a clear implementation report.

Include:

## A. Files Created or Updated

List every test file created or modified.

## B. Test Coverage Summary

Provide a table:

| Layer      | Tests Implemented | Passed | Failed | Coverage |
| ---------- | ----------------: | -----: | -----: | -------: |
| Handler    |                   |        |        |          |
| Service    |                   |        |        |          |
| Repository |                   |        |        |          |
| Routes     |                   |        |        |          |
| DTO        |                   |        |        |          |

Use actual execution results.

## C. Production Code Changes

List every production code change and explain why it was necessary.

If no production code was changed, explicitly state that.

## D. Missing Tests

List any test scenarios that could not be implemented and explain why.

## E. Known Risks

Identify any remaining weaknesses in the Workspace module.

## F. Execution Results

Provide the actual output or a concise summary of:

- `go test ./...`
- `go test -race ./...`
- `go test -cover ./...`
- Repeated test execution.

---

# 15. Final Instructions

Work incrementally.

1. Inspect the complete existing implementation.
2. Understand the architecture and contracts.
3. Prepare the test strategy.
4. Implement handler tests.
5. Implement service tests.
6. Implement repository integration tests.
7. Implement route and middleware integration tests.
8. Implement DTO validation tests.
9. Review production code for test-relevant bugs.
10. Execute all tests.
11. Fix failures.
12. Run race detection and coverage.
13. Provide the final report.

**Do not stop after creating test files. Execute them and verify their results.**

**Do not mark tests as passing without actually running them.**

**Do not skip difficult edge cases merely because they require additional setup.**

**Do not modify unrelated production modules.**

Your final output must represent the quality expected from a senior Go backend engineer working on a production-oriented backend system.
