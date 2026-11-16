````
# Workspace Module — Production-Grade Comprehensive Testing & Quality Engineering

## 1. Mission and Role

Act as a **Principal Go Backend Engineer, Senior Software Architect, Staff-Level Test Automation Engineer, and Application Security Engineer**.

You are responsible for bringing the existing AdOrbit Workspace module to a production-ready testing standard.

Work with the engineering discipline expected from someone responsible for a high-traffic, security-sensitive advertising technology backend.

Your responsibility is NOT merely to write unit tests or increase code coverage.

Your responsibility is to establish confidence that every Workspace API, business operation, database interaction, authentication boundary, authorization rule, middleware behavior, and failure scenario works correctly.

Think critically. Challenge assumptions. Find weaknesses. Verify behavior.

### Engineering expectations

You must demonstrate expertise in:

- Advanced Go programming and idiomatic Go practices.
- Clean Architecture and layered backend systems.
- REST API design and HTTP semantics.
- Unit, integration, component, and end-to-end testing.
- PostgreSQL and SQL query behavior.
- Database constraints, indexes, transactions, and isolation.
- Authentication, authorization, and multi-tenant security.
- JWT validation and middleware composition.
- Idempotency, duplicate requests, and concurrency.
- Race-condition detection and synchronization.
- Error classification and error propagation.
- Test isolation, deterministic execution, and reliability.
- CI/CD and automated quality gates.
- Security testing and regression prevention.

**Do not behave like a junior developer who writes tests based on the happy path.**

Behave like an engineer who assumes that production traffic will eventually expose every weakness that has not been tested.

---

# 2. Project Context

This is the Workspace module of **AdOrbit**, an advertising technology platform.

The backend is written in Go and follows a layered architecture.

Current directory:

```text
internal/workspace/
├── handler.go
├── handler_test.go
├── model.go
├── model_test.go
├── repository.go
├── repository_test.go
├── routes.go
├── routes_test.go
├── service.go
├── service_test.go
└── test_helpers_test.go
````

Existing test files are outdated.

Do not assume they are correct, complete, or compatible with the current production implementation.

Inspect and update them according to the latest source code.

## 2.1 Current architecture

```
                    HTTP Request
                         |
                         v
                       Routes
                         |
                         v
                   JWT Middleware
                         |
                         v
                Idempotency Middleware
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

## 2.2 Current Workspace functionality

The module currently supports five API operations:

| Method | Endpoint           | Purpose                                      |
| ------ | ------------------ | -------------------------------------------- |
| POST   | `/workspaces`      | Create a workspace                           |
| GET    | `/workspaces`      | Fetch authenticated user's active workspaces |
| GET    | `/workspaces/{id}` | Fetch workspace by ID                        |
| PATCH  | `/workspaces/{id}` | Update a workspace                           |
| DELETE | `/workspaces/{id}` | Soft-delete a workspace                      |

Supported workspace types:

- `advertiser`
- `publisher`

Supported workspace statuses:

- `active`
- `deleted`

The repository uses PostgreSQL.

Workspace creation, update, and deletion use idempotency middleware through the routes.

Inspect the actual implementation before assuming any behavior.

# 3. Non-Negotiable Engineering Rules

Follow these rules throughout the task.

1. Do not assume existing production code is correct.
2. Do not assume existing tests are correct.
3. Do not write tests merely to increase coverage.
4. Do not invent application behavior.
5. Do not silently change API contracts.
6. Do not weaken validation to make tests pass.
7. Do not remove meaningful assertions.
8. Do not skip difficult test cases because they require additional setup.
9. Do not introduce unnecessary abstractions.
10. Do not introduce unnecessary dependencies.
11. Do not use production databases for testing.
12. Do not use uncontrolled external services.
13. Do not use arbitrary sleeps to coordinate concurrent tests.
14. Do not use `time.Sleep` as a substitute for synchronization.
15. Do not use `t.Parallel()` without verifying test isolation.
16. Do not ignore errors in test setup or cleanup.
17. Do not mark tests as passing without executing them.
18. Do not hide failures by skipping tests.
19. Do not use excessive mocking when a real dependency is more appropriate.
20. Do not modify unrelated modules.

Every assertion must verify a meaningful application guarantee.

When a behavior is not currently defined, investigate the intended contract and document the ambiguity instead of inventing an expected result.

# 4. Mandatory Initial Inspection

Before writing or modifying tests, inspect the complete implementation.

Do not start coding tests immediately.

Read the following files.

## 4.1 Workspace files

1. `handler.go`
2. `service.go`
3. `repository.go`
4. `model.go`
5. `routes.go`
6. All existing Workspace test files.
7. `test_helpers_test.go`

## 4.2 Supporting project components

Inspect the actual implementation of:

1. Workspace database migrations.
2. Database initialization and connection management.
3. Database driver and PostgreSQL configuration.
4. Error classification in `apperrors`.
5. Error types and sentinel errors.
6. HTTP response utilities.
7. JSON decoding utilities.
8. DTO validation implementation.
9. JWT middleware.
10. Principal context implementation.
11. JWT generation and verification utilities.
12. Idempotency middleware.
13. Idempotency storage and persistence.
14. Relevant database constraints and indexes.
15. Existing project test helpers.
16. Go module dependencies and Go version.
17. Existing CI configuration, if present.

Inspect related modules when their testing patterns are relevant.

## 4.3 Understand actual contracts

Before writing expectations, establish:

- Which errors are classified as conflict errors?
- Which errors are classified as not-found errors?
- Which errors become internal server errors?
- Which errors become service unavailable responses?
- What is the actual HTTP response structure?
- How does JSON decoding handle unknown fields?
- Does the decoder reject trailing JSON?
- How are nullable PostgreSQL fields scanned?
- What does the validation library do with nil pointers?
- What does it do with pointers to empty strings?
- How are JWT principals stored in request context?
- How does idempotency fingerprint requests?
- What happens when two identical requests execute simultaneously?
- What happens when the same key is used with a different payload?
- What happens when the idempotency record expires?
- What database constraints actually exist?
- Does the database enforce workspace ownership or uniqueness?

Never derive expected behavior from a test name or an outdated test implementation.

Derive expectations from the actual production contract.

# 5. Architecture and Dependency Rules

Preserve the existing Clean Architecture structure.

| Layer             | Responsibility                       | Testing strategy                  |
| ----------------- | ------------------------------------ | --------------------------------- |
| Handler           | HTTP request/response handling       | Unit tests with mock Service      |
| Service           | Business logic                       | Unit tests with mock Repository   |
| Repository        | Database operations                  | PostgreSQL integration tests      |
| Routes            | Middleware and endpoint registration | HTTP integration tests            |
| Model/DTO         | Data structures and validation       | Unit tests                        |
| Complete API flow | End-to-end behavior                  | Integration tests where practical |

Architecture rules:

- Handler must depend on the `Service` interface.
- Service must depend on the `Repository` interface.
- Repository implementation must not leak into Handler.
- Service implementation must not leak into Handler.
- Business logic must not be duplicated across layers.
- Keep mocks small, deterministic, and explicit.
- Do not change interfaces unnecessarily.
- Do not introduce a mocking framework without a clear technical reason.

Test code may use additional helpers, but production architecture must remain clean.

# 6. Handler Testing — Complete HTTP Contract

File: `handler_test.go`

Handler tests must verify HTTP behavior independently of PostgreSQL.

Create a handwritten mock implementing the exact existing `Service` interface.

The mock must support:

- Configurable method responses.
- Configurable errors.
- Method call tracking.
- Call counts.
- Capturing received arguments.
- Capturing authenticated UserID.
- Capturing DTO values.
- Deterministic behavior.
- Safe access to shared state if concurrent tests use the mock.

Use the actual principal context mechanism.

Do not bypass authentication checks by modifying production code or introducing fake handler behavior.

## 6.1 POST /workspaces — Create

Test every meaningful request and response path.

### A. Successful requests

1. Valid advertiser workspace.
2. Valid publisher workspace.
3. All required fields populated.
4. All optional fields omitted.
5. Optional fields populated.
6. Optional fields represented by valid nullable values.
7. Service returns a fully populated workspace.
8. Service returns a workspace with nullable fields set to nil.
9. Service receives the correct authenticated UserID.
10. Service receives the correct values for every DTO field.
11. Service is called exactly once.
12. Response contains the correct workspace data.
13. Response has the correct content type.
14. Response uses the actual success status defined by the handler.

### B. Authentication and request decoding

15. Missing principal.
16. Invalid or unavailable principal context.
17. Empty request body.
18. Whitespace-only request body.
19. Invalid JSON syntax.
20. Truncated JSON.
21. Malformed JSON object.
22. JSON array instead of an object.
23. JSON string instead of an object.
24. JSON number instead of an object.
25. `null` JSON payload.
26. Incorrect JSON field types.
27. Unknown JSON fields.
28. Duplicate JSON keys.
29. Valid JSON followed by trailing garbage.
30. Two consecutive JSON objects.
31. Extremely large request body, according to the actual request-size policy.
32. Unsupported content type, if the handler enforces one.

Verify actual decoder behavior rather than assuming unknown fields are rejected.

### C. DTO validation

33. Missing name.
34. Empty name.
35. Name exceeding maximum length.
36. Missing workspace type.
37. Unsupported workspace type.
38. Missing country.
39. Country with invalid length.
40. Missing timezone.
41. Timezone exceeding maximum length.
42. Missing currency.
43. Currency with invalid length.
44. Industry exceeding maximum length.
45. Invalid website.
46. Invalid description according to the actual intended contract.
47. Optional field set to nil.
48. Optional field set to an empty string.
49. Optional field set to a valid string.
50. Optional field set to an invalid string.

### D. Service errors

51. Conflict error.
52. Not-found error, if relevant to the service contract.
53. Validation error returned by service.
54. Unauthorized error, if applicable.
55. Forbidden error, if applicable.
56. Internal error.
57. Service unavailable error.
58. Context cancellation error.
59. Context deadline exceeded.
60. Unexpected unknown error.

Verify the actual error-to-HTTP mapping.

### E. Defensive behavior

61. Service returns nil workspace and nil error.
62. Service returns workspace and error simultaneously.
63. Service is not called when authentication fails.
64. Service is not called when decoding fails.
65. Service is not called when validation fails.
66. No accidental dependency calls after an earlier failure.
67. No panic for unexpected service responses.
68. Correct response encoding for nullable fields.

For nil workspace responses, inspect the current contract and decide whether defensive handling is required. Do not blindly expect HTTP 201 with a null body.

## 6.2 GET /workspaces — Fetch

### A. Successful requests

1. User has one workspace.
2. User has multiple workspaces.
3. User has no workspaces.
4. User has workspaces with nullable fields.
5. User has fully populated workspaces.
6. Service returns an empty non-nil slice.
7. Service returns a nil slice.
8. Service returns workspaces in a specific order.
9. Handler preserves the service-returned order.
10. Correct response JSON structure.
11. Correct content type.
12. Correct success status.
13. Correct authenticated UserID passed to service.
14. Service called exactly once.

### B. Authentication failures

15. Missing principal.
16. Invalid principal context.
17. Service must not be called when principal is missing.

### C. Service failures

18. Conflict error, if applicable.
19. Not-found error, if applicable.
20. Internal error.
21. Service unavailable.
22. Context cancellation.
23. Deadline exceeded.
24. Unexpected error.

### D. Response behavior

25. Empty slice serialization.
26. Nil slice serialization.
27. Nullable field serialization.
28. No accidental omission of required fields.
29. No unexpected data transformation.
30. No panic on unexpected service results.

Verify whether empty results serialize as `[]` or `null` based on the actual contract.

## 6.3 GET /workspaces/ — FetchByID

This endpoint must receive its own complete test suite.

Do not assume that testing `Fetch` covers `FetchByID`.

### A. Successful requests

1. Existing active workspace.
2. Existing deleted workspace, if the API contract allows fetching deleted records.
3. Correct workspace ID passed to service.
4. All workspace fields returned correctly.
5. Nullable fields returned correctly.
6. Correct response status.
7. Correct response JSON structure.
8. Correct content type.
9. Service called exactly once.

### B. Invalid identifiers

10. Empty path identifier, where route-level behavior permits it.
11. Malformed UUID.
12. Incorrect identifier format.
13. Extremely long identifier.
14. Identifier containing unexpected characters.

Do not assume malformed UUIDs are rejected at the handler layer if validation happens in the database.

### C. Error cases

15. Workspace not found.
16. Forbidden access.
17. Conflict, if applicable.
18. Internal server error.
19. Service unavailable.
20. Context cancellation.
21. Deadline exceeded.
22. Unexpected error.

### D. Defensive behavior

23. Service returns nil workspace and nil error.
24. Service returns workspace and error.
25. Service is not called when authentication fails.
26. Correct error response is written.
27. No panic for unexpected service responses.

## 6.4 PATCH /workspaces/ — Update

Treat Update as a complete API, not as a minor variation of Create.

### A. Successful update scenarios

1. Update workspace name only.
2. Update country only.
3. Update timezone only.
4. Update currency only.
5. Update industry only.
6. Update website only.
7. Update description only.
8. Update multiple fields.
9. Update all supported fields.
10. Update name and regenerate slug.
11. Update a non-name field without regenerating slug.
12. Preserve unchanged fields.
13. Correct workspace ID passed to service.
14. Correct DTO passed to service.
15. Correct response body.
16. Correct success status.
17. Nullable fields handled correctly.
18. Service called exactly once.

### B. Partial update semantics

19. Field omitted from JSON.
20. Field explicitly set to null.
21. Field explicitly set to an empty string.
22. Field set to a valid value.
23. Field set to an invalid value.
24. Empty JSON object.
25. All fields omitted.
26. Only whitespace values.
27. Multiple fields with a mix of nil and populated pointers.

Verify the actual intended distinction between omitted, null, and empty values.

An empty PATCH request must not silently generate invalid SQL.

If the current implementation does not define this behavior correctly, identify and fix it with a regression test.

### C. Validation

28. Name exceeds maximum length.
29. Invalid country length.
30. Invalid timezone length.
31. Invalid currency length.
32. Invalid industry length.
33. Invalid website URL.
34. Invalid description according to the intended contract.
35. Incorrect JSON types.
36. Malformed JSON.
37. Trailing JSON data.
38. Unknown fields.
39. Empty request body.
40. Invalid JSON object structure.

### D. Error cases

41. Workspace not found.
42. Duplicate slug conflict.
43. Other unique constraint violation.
44. Database conflict.
45. Internal error.
46. Service unavailable.
47. Context cancellation.
48. Deadline exceeded.
49. Unexpected error.

### E. Defensive behavior

50. Service returns nil workspace and nil error.
51. Service returns workspace and error.
52. Service is not called when authentication fails.
53. Service is not called when decoding fails.
54. Service is not called when validation fails.
55. No panic for unexpected responses.

## 6.5 DELETE /workspaces/ — Delete

The current implementation uses soft deletion.

Test the actual soft-delete contract.

### A. Successful deletion

1. Existing active workspace.
2. Correct workspace ID passed to service.
3. Correct success status.
4. Correct success response.
5. Service called exactly once.
6. No workspace object unexpectedly included in response.
7. Correct content type.

### B. Error cases

8. Workspace not found.
9. Forbidden access.
10. Internal server error.
11. Service unavailable.
12. Context cancellation.
13. Deadline exceeded.
14. Unexpected error.

### C. Defensive behavior

15. Missing principal.
16. Service not called when authentication fails.
17. Empty workspace ID where route behavior permits it.
18. Malformed workspace ID.
19. Unexpected service response is handled correctly.
20. No panic on failure.

### D. Soft-delete contract

Verify that successful deletion means the record is marked as deleted rather than physically removed.

Repository integration tests must verify the database state.

# 7. Service Unit Tests — Business Logic Verification

File: `service_test.go`

Create a mock implementing the exact `Repository` interface.

Track:

- Method call counts.
- Context propagation.
- Received DTO values.
- Received UserID.
- Workspace ID.
- Generated slug.
- Returned workspace.
- Returned errors.

Test every service method independently.

## 7.1 Create Service

### Slug generation

1. Basic workspace name.
2. Name containing spaces.
3. Name containing multiple consecutive spaces.
4. Name with uppercase characters.
5. Name with leading spaces.
6. Name with trailing spaces.
7. Name containing punctuation.
8. Name containing special characters.
9. Name containing numbers.
10. Name containing Unicode characters.
11. Name containing accented characters.
12. Name containing non-Latin characters.
13. Name containing only punctuation.
14. Empty name.
15. Whitespace-only name.
16. Name containing mixed punctuation and spaces.
17. Name containing repeated separators.

Use the actual `gosimple/slug` library behavior.

Do not guess expected slugs.

### Repository interaction

18. Correct DTO passed to repository.
19. Correct UserID passed to repository.
20. Generated slug passed to repository.
21. Repository called exactly once.
22. Original DTO remains unchanged.
23. Nil optional fields remain nil.
24. Populated optional fields remain unchanged.
25. Context is propagated.
26. Repository returns a valid workspace.
27. Repository returns nil workspace and nil error.
28. Repository returns workspace and error.
29. Repository returns a classified conflict error.
30. Repository returns an unknown error.
31. Error classification preserves expected semantics.

## 7.2 Fetch Service

1. Repository returns one workspace.
2. Repository returns multiple workspaces.
3. Repository returns an empty slice.
4. Repository returns a nil slice.
5. Repository returns nullable workspace fields.
6. Correct UserID passed to repository.
7. Repository called exactly once.
8. Context propagated.
9. Workspace order preserved.
10. Workspace values remain unchanged.
11. Repository returns an error.
12. Repository returns a classified error.
13. Repository returns an unexpected error.
14. Error classification verified.

## 7.3 FetchByID Service

1. Existing workspace returned.
2. Correct workspace ID passed.
3. Repository called exactly once.
4. Context propagated.
5. Repository returns nil workspace and nil error.
6. Repository returns workspace and error.
7. Repository returns not-found error.
8. Repository returns unexpected error.
9. Error classification verified.
10. No accidental modification of returned workspace.

## 7.4 Update Service

### Slug behavior

1. Name is nil: slug must remain nil.
2. Name is populated: slug must be generated.
3. Name contains spaces.
4. Name contains uppercase characters.
5. Name contains punctuation.
6. Name contains Unicode characters.
7. Name contains leading/trailing spaces.
8. Name contains numbers.
9. Name contains only special characters.
10. Empty name pointer.
11. Whitespace-only name pointer.

### Repository interaction

12. Correct DTO passed.
13. Correct workspace ID passed.
14. Correct generated slug passed.
15. Non-name update passes nil slug.
16. Repository called exactly once.
17. Context propagated.
18. Repository returns updated workspace.
19. Repository returns nil workspace and nil error.
20. Repository returns workspace and error.
21. Repository returns not-found error.
22. Repository returns duplicate slug error.
23. Repository returns internal error.
24. Repository returns unexpected error.
25. Error classification verified.
26. Original DTO remains unchanged.

## 7.5 Delete Service

1. Successful deletion.
2. Correct workspace ID passed.
3. Repository called exactly once.
4. Context propagated.
5. Repository returns not-found error.
6. Repository returns conflict error, if applicable.
7. Repository returns internal error.
8. Repository returns unavailable error.
9. Repository returns unexpected error.
10. Nil repository error produces nil service error.
11. Error classification verified.

### Service testing rule

Do not add business rules that do not exist.

However, if an existing behavior is unsafe or incomplete, document it and propose a minimal correction.

# 8. DTO and Model Testing

File: `model_test.go`

Test `CreateWorkspaceDTO` thoroughly.

Use the same validator and validation mechanism used by the actual application.

Do not test a different validator configuration.

## 8.1 CreateWorkspaceDTO

### Required fields

1. All required fields valid.
2. Missing name.
3. Missing type.
4. Missing country.
5. Missing timezone.
6. Missing currency.
7. Empty name.
8. Empty type.
9. Empty country.
10. Empty timezone.
11. Empty currency.
12. Nil required pointers where possible.

### Name validation

13. Minimum valid name.
14. Name at maximum length.
15. Name one character above maximum.
16. Empty string.
17. Whitespace-only string.
18. Leading whitespace.
19. Trailing whitespace.
20. Unicode name.
21. Special characters.
22. Numbers.
23. Mixed letters and numbers.

### Workspace type

24. `advertiser`.
25. `publisher`.
26. Unsupported type.
27. Uppercase type.
28. Mixed-case type.
29. Empty type.
30. Nil type if the underlying representation allows it.

### Country

31. Valid two-character country.
32. One-character country.
33. Three-character country.
34. Empty string.
35. Lowercase country.
36. Uppercase country.
37. Numbers.
38. Special characters.
39. Nil pointer.

### Timezone

40. Valid timezone.
41. Empty timezone.
42. Timezone at maximum length.
43. Timezone above maximum length.
44. Nil pointer.
45. Invalid timezone identifier if the application has semantic timezone validation.

Do not expect semantic timezone validation unless it is actually implemented.

### Currency

46. Valid three-character currency.
47. Two-character currency.
48. Four-character currency.
49. Empty string.
50. Lowercase currency.
51. Uppercase currency.
52. Numeric characters.
53. Special characters.
54. Nil pointer.

### Industry

55. Nil industry.
56. Empty industry.
57. Valid industry.
58. Maximum length.
59. One character above maximum.
60. Unicode industry.
61. Special characters.

### Website

62. Nil website.
63. Empty website.
64. Valid HTTP URL.
65. Valid HTTPS URL.
66. URL without scheme.
67. Malformed URL.
68. URL with query parameters.
69. URL with fragment.
70. URL with port.
71. URL with subdomain.
72. URL with whitespace.
73. Non-HTTP scheme.
74. Boundary and unusual URL values.

Verify actual `url` validation semantics.

### Description

75. Nil description.
76. Empty description.
77. Normal text description.
78. Description containing spaces.
79. Description containing punctuation.
80. Description containing Unicode.
81. Description at the intended maximum length.
82. Description exceeding the intended maximum length.
83. URL-like description.
84. Text that is not a URL.

Important: Inspect the actual `Description` validation tag. If it currently uses URL validation but the field is intended to contain descriptive text, identify this as a production bug. Do not write tests that incorrectly enforce URL-only descriptions.

### JSON serialization

85. Correct JSON field names.
86. Nil optional fields.
87. Populated optional fields.
88. Correct workspace type serialization.
89. Correct workspace status serialization.
90. Workspace timestamps serialize correctly.
91. Nullable fields serialize according to the actual contract.

## 8.2 UpdateWorkspaceDTO

Unlike the outdated testing plan, Update DTO validation is IN SCOPE.

Test it thoroughly.

Cover:

- All fields nil.
- Each field updated individually.
- Multiple fields updated together.
- All fields populated.
- Empty strings.
- Nil pointers.
- Invalid values.
- Maximum allowed lengths.
- Values exceeding maximum lengths.
- Invalid URLs.
- Unicode values.
- Whitespace-only values.
- Pointer-to-empty versus nil pointer behavior.
- JSON omitted fields.
- JSON null fields.
- JSON empty strings.
- Incorrect JSON field types.
- Unknown fields.

Also test the actual semantics of an empty Update DTO.

If the API must reject an empty PATCH request, verify that this behavior is enforced by the actual validation or business logic.

`omitempty` on individual validation tags does not automatically prove that an empty DTO is invalid.

Do not assume otherwise.

# 9. Repository Integration Tests — Real PostgreSQL

File: `repository_test.go`

Repository tests must verify actual PostgreSQL behavior.

Use a dedicated isolated PostgreSQL test database.

Do not use:

- Development database.
- Production database.
- Shared application database.
- Hardcoded credentials.
- Uncontrolled external database instances.

## 9.1 Test database strategy

Inspect the existing project database setup first.

Reuse the existing initialization and migration mechanisms where practical.

Requirements:

1. Dedicated test database.
2. Explicit test database configuration.
3. Actual project migrations applied.
4. Deterministic test data.
5. Proper cleanup.
6. No test-order dependencies.
7. No uncontrolled shared state.
8. Safe execution with repeated test runs.
9. Safe behavior when database configuration is missing.
10. Clear failure messages when PostgreSQL is unavailable.

Do not silently fall back to a production or development database.

Do not silently skip all repository tests.

If integration tests require an environment variable, document it clearly.

For example, use a dedicated variable such as:

```
TEST_DATABASE_URL
```

Only use the actual project configuration convention if one already exists.

## 9.2 Create Repository

### Successful insertion

1. Create advertiser workspace.
2. Create publisher workspace.
3. Correct name persisted.
4. Correct slug persisted.
5. Correct type persisted.
6. Correct country persisted.
7. Correct timezone persisted.
8. Correct currency persisted.
9. Correct created_by persisted.
10. Default status is active.
11. UUID generated correctly.
12. Created timestamp populated.
13. Updated timestamp populated.
14. Returned workspace matches persisted data.
15. Nullable fields are nil when omitted.
16. Nullable fields are populated when provided.
17. All optional fields populated together.
18. Unique ID behavior verified.

### Database constraints

19. Duplicate slug.
20. Invalid workspace type.
21. Invalid country according to actual schema.
22. Invalid currency according to actual schema.
23. Invalid creator UUID.
24. Missing required database fields.
25. Invalid data violating actual CHECK constraints.
26. Other actual unique constraint violations.
27. Foreign-key violations, if applicable.

Do not invent constraints that do not exist.

### Failure behavior

28. Database unavailable.
29. Connection failure.
30. Context cancelled before execution.
31. Context deadline exceeded.
32. Query execution failure.
33. Scan failure, if reproducible.
34. Transaction rollback and cleanup where applicable.

### Persistence integrity

35. Failed insertion does not create a record.
36. Duplicate insertion does not create an additional record.
37. Nullable values are preserved.
38. No unexpected data mutation.
39. Database-generated values are returned correctly.

## 9.3 Fetch Repository

The current implementation filters by:

- `created_by`
- `status = 'active'`

It orders results by newest creation time first.

Verify this exact contract.

### User isolation

1. User with no workspaces.
2. User with one active workspace.
3. User with multiple active workspaces.
4. User with only deleted workspaces.
5. User with active and deleted workspaces.
6. Different users have separate workspaces.
7. User A cannot retrieve User B's workspaces.
8. Identical workspace names belonging to different users.
9. Empty result behavior.

### Ordering

10. Newest workspace appears first.
11. Oldest workspace appears last.
12. Deterministic timestamp-based ordering.
13. Multiple records with identical timestamps.

If deterministic ordering is required when timestamps match, identify whether the query needs a secondary ordering key. Add a regression test if the contract requires it.

### Data integrity

14. All columns scanned correctly.
15. Nullable industry scanned.
16. Nullable website scanned.
17. Nullable description scanned.
18. All nullable fields populated.
19. All nullable fields nil.
20. Workspace status scanned correctly.
21. UUID values scanned correctly.
22. Timestamp values scanned correctly.

### Error behavior

23. Database query error.
24. Context cancellation.
25. Context deadline exceeded.
26. Row scanning error using SQL mocking where appropriate.
27. Rows iteration error using SQL mocking where appropriate.
28. `rows.Err()` propagation.
29. Rows are closed correctly.

Use SQL mocking only where a real PostgreSQL integration test cannot reliably reproduce the condition.

Do not modify production repository code solely to simulate artificial failures.

## 9.4 FetchByID Repository

This method requires independent coverage.

### Successful retrieval

1. Existing active workspace.
2. Existing deleted workspace.
3. Correct ID returned.
4. All fields returned.
5. Nullable fields returned.
6. Timestamps returned correctly.
7. Correct workspace type.
8. Correct status.
9. Record matches database state.

### Failure cases

10. Nonexistent workspace ID.
11. Malformed UUID.
12. Empty ID where applicable.
13. Context cancellation.
14. Deadline exceeded.
15. Database unavailable.
16. Query error.
17. Scan error using an appropriate SQL mock.

### Security and data isolation

18. Verify the repository query's ownership behavior.
19. Verify whether the current method allows fetching another user's workspace by ID.
20. Identify the authorization responsibility of this layer.

A successful lookup by ID alone does not prove that the authenticated user owns the workspace.

Do not confuse authentication with authorization.

## 9.5 Update Repository

This is a high-priority area because the current implementation dynamically builds SQL.

Test every update field independently.

### Individual field updates

1. Update name.
2. Update slug.
3. Update country.
4. Update timezone.
5. Update currency.
6. Update industry.
7. Update website.
8. Update description.

### Multiple field updates

9. Update two fields.
10. Update three fields.
11. Update all supported fields.
12. Update name and slug together.
13. Update multiple nullable fields.
14. Update a mixture of populated and nil pointers.

### Partial update behavior

15. Omitted fields remain unchanged.
16. Updated fields persist correctly.
17. Explicit empty string behavior.
18. Nil pointer behavior.
19. Empty DTO behavior.
20. No malformed SQL for empty updates.
21. No accidental overwriting of unchanged fields.

### Slug and uniqueness

22. Change name and slug.
23. Change name to another workspace's slug.
24. Update to an already existing slug.
25. Update to the same existing slug.
26. Update a non-name field without changing slug.
27. Verify unique constraint errors.
28. Verify failed update does not partially modify the record.

### Not-found and database failures

29. Update nonexistent workspace.
30. Update malformed UUID.
31. Database unavailable.
32. Context cancellation.
33. Deadline exceeded.
34. Query error.
35. Scan error using SQL mocking.

### Persistence integrity

36. Verify returned workspace.
37. Verify persisted workspace.
38. Verify updated_at behavior.
39. Verify created_at remains unchanged.
40. Verify status remains unchanged.
41. Verify created_by remains unchanged.
42. Verify unrelated fields remain unchanged.

### SQL correctness

43. Correct placeholder ordering.
44. Correct number of SQL arguments.
45. No SQL injection through user-provided values.
46. Correct handling of special characters.
47. Correct handling of Unicode.
48. No accidental SQL syntax errors.

Use actual database tests for SQL behavior.

Do not test SQL injection merely by asserting that a string contains a particular substring. Verify that user-provided values are treated as parameters.

## 9.6 Delete Repository

The current implementation performs a soft delete:

```
UPDATE workspaces
SET status = 'deleted',
    updated_at = NOW()
WHERE id = $1
```

Test the actual database behavior.

### Successful deletion

1. Delete existing active workspace.
2. Status changes to deleted.
3. Record remains in the database.
4. ID remains unchanged.
5. created_at remains unchanged.
6. updated_at changes.
7. Other workspace fields remain unchanged.
8. Deleting one workspace does not affect another workspace.

### Failure and edge cases

9. Delete nonexistent workspace.
10. Delete malformed UUID.
11. Delete already deleted workspace.
12. Context cancellation.
13. Deadline exceeded.
14. Database unavailable.
15. Query execution error.
16. RowsAffected error using SQL mocking if needed.
17. Zero rows affected.
18. Multiple rows affected should be investigated as an invariant violation if the schema permits it.

Important production review:

The current repository implementation discards the `sql.Result` returned by `ExecContext`.

Investigate whether it must check `RowsAffected()` to distinguish:

- A successful soft delete.
- A workspace that does not exist.
- A workspace already deleted.

Do not silently accept incorrect not-found behavior.

Add regression tests for the chosen contract.

# 10. Routes and Middleware Integration Tests

File: `routes_test.go`

These tests must verify the complete route chain.

Use:

- Actual `RegisterRoutes`.
- Actual JWT middleware.
- Actual idempotency middleware.
- `httptest`.
- Mock Service.
- Isolated idempotency storage.
- Realistic authenticated principals.

Do not replace the route chain with direct handler calls.

Do not mock JWT middleware in tests intended to verify JWT behavior.

## 10.1 Authentication tests

Test every endpoint.

| Endpoint            | Missing JWT | Malformed JWT | Expired JWT | Valid JWT |
| ------------------- | ----------- | ------------- | ----------- | --------- |
| POST /workspaces    | Reject      | Reject        | Reject      | Allow     |
| GET /workspaces     | Reject      | Reject        | Reject      | Allow     |
| GET /workspaces/    | Reject      | Reject        | Reject      | Allow     |
| PATCH /workspaces/  | Reject      | Reject        | Reject      | Allow     |
| DELETE /workspaces/ | Reject      | Reject        | Reject      | Allow     |

Also test:

1. Invalid signature.
2. Wrong signing key.
3. Missing required claims.
4. Expired token.
5. Not-yet-valid token if supported.
6. Unsupported signing algorithm.
7. Malformed authorization header.
8. Missing Bearer prefix.
9. Empty token.
10. Correct principal reaches the handler.
11. Authenticated identity is not replaced by request data.
12. Rejected authentication never reaches the service.

Use the actual JWT configuration and token-generation helpers.

## 10.2 Route registration

Test:

1. POST route registration.
2. GET collection route registration.
3. GET by ID route registration.
4. PATCH by ID route registration.
5. DELETE by ID route registration.
6. Unsupported methods.
7. Incorrect paths.
8. Trailing slash behavior.
9. Path parameter extraction.
10. Correct method dispatch.
11. No accidental route conflicts.
12. No service invocation for unsupported methods.

Verify actual Go `ServeMux` behavior and the project's Go version.

## 10.3 POST idempotency

The current route uses:

```
workspace.create
workspace:create:<principal.UserID>
```

Inspect the middleware before asserting expected behavior.

Test:

1. Missing Idempotency-Key.
2. Empty Idempotency-Key.
3. Malformed key according to actual rules.
4. Valid key.
5. Same key and identical payload.
6. Same key and different payload.
7. Same key used by different users.
8. Different keys used by the same user.
9. Repeated request after completion.
10. Request while the original request is still processing.
11. Concurrent identical requests.
12. Concurrent requests using different keys.
13. Stored response replay.
14. Error response handling.
15. Idempotency storage failure.
16. Expired idempotency record.
17. Retry after expiration.
18. Request fingerprint mismatch.
19. Verify only one workspace is created for a duplicate request where guaranteed by middleware semantics.
20. Verify user-specific scope isolation.

## 10.4 PATCH idempotency

The current route also applies idempotency middleware to PATCH.

Test:

1. Missing key.
2. Empty key.
3. Valid key.
4. Same key and same payload.
5. Same key and different payload.
6. Different users using the same key.
7. Concurrent identical updates.
8. Concurrent requests with different keys.
9. Replay of completed response.
10. In-progress duplicate request.
11. Expired key.
12. Storage failure.
13. Fingerprint mismatch.
14. Verify duplicate processing behavior.
15. Verify workspace state remains consistent.

Do not assume PATCH is naturally safe to retry. Test the middleware contract.

## 10.5 DELETE idempotency

The current route applies idempotency middleware to DELETE.

Test:

1. Missing key.
2. Empty key.
3. Valid key.
4. Repeated request with same key.
5. Same key with a different request.
6. Different users using the same key.
7. Concurrent identical delete requests.
8. In-progress duplicate request.
9. Completed response replay.
10. Expired key.
11. Storage failure.
12. Fingerprint mismatch.
13. Workspace state after repeated requests.
14. Behavior when workspace is already deleted.

Verify actual behavior instead of assuming every repeated DELETE must return the same status.

## 10.6 Idempotency scope isolation

Verify the actual scope strings:

```
workspace:create:<userID>
workspace:update:<userID>
workspace:delete:<userID>
```

Test that:

- Different users do not share idempotency scope.
- Different operations do not accidentally share scope.
- The same user can use independent keys for independent operations.
- Scope is derived from authenticated identity.
- User-controlled request fields cannot override the scope.

## 10.7 Complete route-chain behavior

For each API, verify:

1. Authentication occurs before handler execution.
2. Idempotency is applied only to the intended methods.
3. Handler receives the correct context.
4. Correct UserID reaches the handler.
5. Correct path ID reaches the handler.
6. Rejected requests do not call the service.
7. Correct response status.
8. Correct response body.
9. Correct headers.
10. Correct middleware order.

# 11. Authorization and Multi-Tenant Security Testing

This section is mandatory and has high priority.

AdOrbit is a multi-user advertising platform.

Authentication alone is not sufficient.

An authenticated user must not be able to access or modify another user's workspace merely by knowing or guessing its ID.

Inspect the complete ownership enforcement path.

The current interfaces shown in the implementation use:

```
FetchByID(ctx, workspaceID)
Update(ctx, dto, workspaceID, slug)
Delete(ctx, workspaceID)
```

These methods do not directly receive a UserID.

Investigate where ownership is enforced.

Do not assume the JWT middleware automatically provides resource-level authorization.

## 11.1 Required security scenarios

Test:

1. User A creates a workspace.
2. User A fetches their own workspace.
3. User B attempts to fetch User A's workspace.
4. User B attempts to update User A's workspace.
5. User B attempts to delete User A's workspace.
6. User B attempts to modify User A's workspace name.
7. User B attempts to change User A's workspace slug.
8. User A can still access their own workspace.
9. Deleted workspace access follows the intended policy.
10. Unauthenticated users cannot access any workspace.
11. User ID supplied in JSON cannot override authenticated identity.
12. User ID supplied in URL or query parameters cannot override authenticated identity.
13. Cross-user idempotency keys do not bypass authorization.

Test for insecure direct object reference (IDOR).

## 11.2 Required response behavior

Inspect the intended security contract.

Depending on the API design, unauthorized resource access may return:

- `403 Forbidden`, or
- `404 Not Found` to avoid revealing resource existence.

Do not arbitrarily choose one.

Use a consistent policy and test it.

If ownership is currently missing:

1. Identify the production vulnerability.
2. Explain the attack scenario.
3. Propose the smallest correct fix.
4. Implement the fix if it is necessary for secure Workspace behavior.
5. Add regression tests for every affected endpoint.
6. Update affected interfaces and mocks only when justified.
7. Avoid unrelated architectural refactoring.

A workspace module must not be considered production-ready while cross-user resource access remains possible.

# 12. Advanced Database and Concurrency Testing

Do not limit testing to sequential requests.

Investigate concurrency-sensitive behavior.

## 12.1 Repository concurrency

Test relevant scenarios:

1. Concurrent workspace creation.
2. Concurrent updates to the same workspace.
3. Concurrent updates to different workspaces.
4. Concurrent deletion of the same workspace.
5. Concurrent fetch and update.
6. Concurrent fetch and delete.
7. Duplicate slug creation under concurrency.
8. Database constraint enforcement under concurrent writes.
9. Context cancellation during concurrent operations.
10. Data consistency after concurrent operations.

Use synchronization primitives and deterministic coordination.

Do not rely on arbitrary sleeps.

## 12.2 Idempotency concurrency

Test:

1. Same user, same operation, same key, same payload.
2. Same user, same operation, same key, different payload.
3. Different users, same key.
4. Same user, different operations, same key.
5. Many simultaneous requests with the same key.
6. Concurrent requests while the original request is processing.
7. Retry after the original request fails.
8. Retry after an idempotency storage error.
9. Expiration and cleanup races, if supported.

Assert the exact guarantees documented by the middleware.

## 12.3 Race detection

Ensure test helpers and mocks do not introduce data races.

Run:

```
go test -race ./...
```

Investigate all race detector reports.

Do not suppress race reports without a technically justified explanation.

# 13. Test Helpers and Infrastructure

File: `test_helpers_test.go`

Inspect the existing helpers before replacing them.

Build reusable helpers only where they improve consistency.

Potential helpers include:

- Test principal creation.
- JWT generation.
- HTTP request construction.
- JSON response decoding.
- Error response assertions.
- Workspace fixtures.
- DTO fixtures.
- PostgreSQL setup.
- Migration execution.
- Test transaction setup.
- Database cleanup.
- Idempotency storage setup.
- Context cancellation.
- Deterministic time control.
- Mock Service.
- Mock Repository.

Requirements:

1. Helpers must fail clearly.
2. Setup errors must not be ignored.
3. Cleanup must use `t.Cleanup`.
4. Helpers must not hide assertions.
5. Avoid global mutable fixtures.
6. Avoid hidden shared database state.
7. Keep test setup explicit.
8. Keep helpers small and reusable.
9. Do not build an unnecessary testing framework.
10. Make failure messages identify the relevant operation and input.

Do not make tests harder to understand merely to reduce repeated lines.

# 14. Error Classification and HTTP Error Contract

Inspect `apperrors` and response utilities.

Create a clear mapping between:

- Domain errors.
- Repository errors.
- Classified service errors.
- Handler errors.
- HTTP status codes.

Test every meaningful error class.

At minimum, investigate:

1. Not found.
2. Conflict.
3. Validation error.
4. Unauthorized.
5. Forbidden.
6. Internal server error.
7. Service unavailable.
8. Context cancellation.
9. Deadline exceeded.
10. Unknown error.

Use:

```
errors.Is
errors.As
```

where appropriate.

Do not compare raw error strings when typed or sentinel errors are available.

Verify:

- Correct error identity.
- Correct error classification.
- Wrapped error behavior.
- Original error context where required.
- Correct HTTP status.
- Correct response structure.
- No internal database details leaked to clients.
- No sensitive information leaked in error responses.

# 15. Production Code Review During Testing

Testing must also be used to identify weaknesses in the actual Workspace implementation.

Review the following areas carefully.

## 15.1 Handler

- Nil dependency handling.
- Nil workspace responses.
- JSON decoding behavior.
- Trailing JSON data.
- Unknown JSON fields.
- Validation execution.
- Error response consistency.
- Correct HTTP status codes.
- Correct authenticated UserID.
- Correct path ID handling.
- Context propagation.
- Panic prevention.

## 15.2 Service

- Slug generation.
- Nil pointer handling.
- Error classification.
- Context propagation.
- Correct repository arguments.
- Unnecessary transformations.
- Business logic consistency.

## 15.3 Repository

- SQL injection protection.
- Placeholder ordering.
- Dynamic SQL construction.
- Empty PATCH handling.
- RowsAffected checks.
- Not-found behavior.
- Soft-delete correctness.
- Nullable field scanning.
- `rows.Err()` handling.
- `rows.Close()`.
- Context propagation.
- Timestamp updates.
- Unique constraint handling.
- Transaction consistency.

## 15.4 Routes

- JWT middleware ordering.
- Idempotency middleware ordering.
- Scope isolation.
- Principal extraction.
- Path registration.
- Unsupported methods.
- Consistent protection across all endpoints.

## 15.5 Security

- IDOR.
- User isolation.
- Ownership enforcement.
- Authentication versus authorization.
- Information disclosure.
- Replay behavior.
- Concurrent duplicate requests.
- Trust boundaries around client-provided identity.

## 15.6 Production bug handling

When you discover a bug:

1. Identify the exact file and function.
2. Explain the root cause.
3. Explain the production impact.
4. Create a regression test.
5. Propose the smallest correct fix.
6. Implement it when required for correctness and within scope.
7. Re-run all affected tests.
8. Verify that the fix does not break other APIs.

Do not silently change production behavior.

Do not refactor unrelated code.

Do not change API contracts without a clear justification.

# 16. Test Quality Standards

## 16.1 Table-driven testing

Use table-driven tests where cases share the same structure.

Example:

```
tests := []struct {
    name string
    // input
    // expected output
    // expected error
}{
    {
        name: "valid advertiser workspace",
    },
    {
        name: "invalid workspace type",
    },
}
```

Do not force table-driven tests when individual setup is clearer.

## 16.2 Test naming

Use descriptive names.

Examples:

```
func TestCreateHandler_Success(t *testing.T)

func TestCreateHandler_MissingPrincipal(t *testing.T)

func TestFetchByIDHandler_NotFound(t *testing.T)

func TestUpdateHandler_EmptyPayload(t *testing.T)

func TestDeleteHandler_ServiceUnavailable(t *testing.T)

func TestCreateService_GeneratesSlug(t *testing.T)

func TestUpdateService_NameChanged_RegeneratesSlug(t *testing.T)

func TestCreateRepository_DuplicateSlug(t *testing.T)

func TestUpdateRepository_EmptyDTO(t *testing.T)

func TestWorkspaceRoutes_PATCHWithoutJWT(t *testing.T)

func TestWorkspaceRoutes_Delete_CrossUserAccess(t *testing.T)
```

## 16.3 Assertion quality

Where applicable, assert:

- Exact status codes.
- Exact response structure.
- Expected response fields.
- Expected error classification.
- Correct dependency call count.
- Correct DTO.
- Correct authenticated UserID.
- Correct workspace ID.
- Correct generated slug.
- Correct persisted database values.
- Correct nullable values.
- Correct timestamps.
- Correct ordering.
- Correct authorization behavior.
- Correct middleware execution.
- Correct context propagation.

A test that only checks `err != nil` is insufficient when the actual error type or behavior matters.

## 16.4 Test independence

Every test must be independently executable.

Requirements:

- No test-order dependency.
- No reliance on previous test data.
- No shared mutable fixtures.
- No production credentials.
- No uncontrolled external APIs.
- No arbitrary sleeps.
- No uncontrolled randomness.
- No ignored cleanup errors.
- No hidden test state.

## 16.5 Flaky test prevention

Avoid:

- Timing-dependent assertions.
- Unbounded goroutines.
- Shared database records.
- Shared idempotency keys.
- Global mutable configuration.
- Non-deterministic ordering assumptions.
- Arbitrary polling without deadlines.

Use deterministic synchronization and explicit timeouts.

# 17. Test Execution and Verification

Do not stop after writing tests.

Execute them.

Fix failures.

Re-run them.

## 17.1 Formatting

```
gofmt -w internal/workspace
```

Inspect the resulting diff.

Do not format unrelated production files unnecessarily.

## 17.2 Workspace tests

```
go test ./internal/workspace/...
```

## 17.3 Full project tests

```
go test ./...
```

## 17.4 Verbose tests

```
go test -v ./internal/workspace/...
```

## 17.5 Race detection

```
go test -race ./internal/workspace/...
```

Also run project-wide race detection when practical:

```
go test -race ./...
```

## 17.6 Coverage

```
go test -cover ./internal/workspace/...
```

## 17.7 Coverage profile

```
go test -coverprofile=coverage.out ./internal/workspace/...
```

## 17.8 Function-level coverage

```
go tool cover -func=coverage.out
```

## 17.9 HTML coverage report

```
go tool cover -html=coverage.out
```

## 17.10 Repeated execution

```
go test -count=10 ./internal/workspace/...
```

Run repeated execution to detect flaky tests.

## 17.11 Static analysis

Where supported by the project, run:

```
go vet ./...
```

Also inspect:

- Compiler warnings and failures.
- Unused helpers.
- Test-only race conditions.
- Unhandled errors.
- Suspicious assertions.
- Unnecessary test dependencies.

## 17.12 Database-dependent execution

Run repository integration tests against the isolated PostgreSQL test database.

Do not claim repository tests passed if they were skipped or were not executed.

If the database is unavailable:

1. Report the exact missing requirement.
2. Do not report integration tests as passing.
3. Continue with tests that can run independently.
4. Provide the exact command needed to execute the database-dependent tests.

# 18. Coverage Expectations

Coverage targets are quality indicators, not the objective.

Suggested targets:

| Layer          | Target |
| -------------- | ------ |
| Handler        | 95%+   |
| Service        | 95%+   |
| Repository     | 90%+   |
| DTO validation | 95%+   |
| Routes         | 95%+   |

Prioritize meaningful branch and behavior coverage.

Pay special attention to:

- Error paths.
- Security boundaries.
- Ownership enforcement.
- Invalid input.
- Nil values.
- Database failures.
- Partial updates.
- Soft deletion.
- Idempotency.
- Concurrent requests.
- Context cancellation.

Do not manipulate coverage by testing trivial lines repeatedly.

Do not claim 100% coverage is necessary or sufficient for production readiness.

Explain any uncovered code that is difficult, unsafe, or inappropriate to test.

# 19. Required Incremental Workflow

Work in the following sequence.

Do not skip stages.

## Phase 1 — Inspection

1. Inspect all Workspace production files.
2. Inspect all existing Workspace tests.
3. Inspect migrations.
4. Inspect validation.
5. Inspect error classification.
6. Inspect JWT middleware.
7. Inspect idempotency middleware.
8. Inspect database test infrastructure.
9. Inspect response utilities.
10. Identify architectural and security risks.

Deliver a concise inspection report before implementing tests.

## Phase 2 — Test Strategy

Create a test matrix covering all five APIs.

Identify:

- Unit test cases.
- Integration test cases.
- Security test cases.
- Middleware test cases.
- Database test cases.
- Concurrency test cases.
- Existing bugs.
- Missing test infrastructure.

Identify dependencies between tests.

## Phase 3 — Test Infrastructure

1. Improve `test_helpers_test.go`.
2. Implement deterministic fixtures.
3. Implement mock Service.
4. Implement mock Repository.
5. Configure isolated PostgreSQL tests.
6. Configure actual migration execution.
7. Configure JWT helpers.
8. Configure idempotency test storage.

Avoid unnecessary abstractions.

## Phase 4 — DTO Tests

Implement complete Create and Update DTO validation tests.

Verify actual validation behavior.

## Phase 5 — Service Tests

Implement complete tests for:

- Create.
- Fetch.
- FetchByID.
- Update.
- Delete.

## Phase 6 — Handler Tests

Implement complete HTTP unit tests for:

- Create.
- Fetch.
- FetchByID.
- Update.
- Delete.

## Phase 7 — Repository Tests

Implement PostgreSQL integration tests for:

- Create.
- Fetch.
- FetchByID.
- Update.
- Delete.

## Phase 8 — Routes and Middleware Tests

Implement complete route-chain tests for all five APIs.

Include:

- Authentication.
- Authorization.
- Idempotency.
- Request routing.
- Error propagation.
- User isolation.
- Duplicate requests.
- Concurrency.

## Phase 9 — Security Review

Verify cross-user access protections.

Verify that no user can access, update, or delete another user's workspace.

Fix security issues that are necessary for the Workspace module to be safe.

## Phase 10 — Production Review

Inspect production code for all issues discovered during testing.

Add regression tests for every production bug fixed.

## Phase 11 — Execution

Run:

1. Workspace tests.
2. Full project tests.
3. Race detection.
4. Coverage.
5. Repeated tests.
6. Static analysis.

Fix failures and repeat relevant commands.

## Phase 12 — Final Review

Review the final diff.

Verify:

- No unrelated files changed.
- No production credentials introduced.
- No test depends on production services.
- No tests silently skipped.
- No ignored errors.
- No race conditions.
- No flaky behavior.
- No false coverage claims.
- No unresolved critical security issue.

# 20. Final Deliverables

After implementation, provide a detailed engineering report.

## A. Files Created or Updated

List every file created, modified, or deleted.

Separate:

- Production files.
- Test files.
- Test infrastructure.
- Configuration files.

Explain why each file changed.

## B. Test Inventory

Provide the actual number of tests implemented.

Break them down by:

- Handler.
- Service.
- Repository.
- Routes.
- DTO.
- Security.
- Concurrency.

Do not count subtests as independent top-level tests unless clearly identified.

## C. Execution Summary

Use actual command results.

| Command                                      | Result | Details |
| -------------------------------------------- | ------ | ------- |
| `go test ./internal/workspace/...`           |        |         |
| `go test ./...`                              |        |         |
| `go test -race ./internal/workspace/...`     |        |         |
| `go test -race ./...`                        |        |         |
| `go test -cover ./internal/workspace/...`    |        |         |
| `go test -count=10 ./internal/workspace/...` |        |         |
| `go vet ./...`                               |        |         |

Do not fill in results for commands that were not executed.

## D. Coverage Summary

| Layer      | Tests Implemented | Passed | Failed | Coverage |
| ---------- | ----------------- | ------ | ------ | -------- |
| Handler    |                   |        |        |          |
| Service    |                   |        |        |          |
| Repository |                   |        |        |          |
| Routes     |                   |        |        |          |
| DTO        |                   |        |        |          |
| Overall    |                   |        |        |          |

Use real coverage output.

## E. Production Code Changes

For every production change, explain:

- File.
- Function.
- Original problem.
- Root cause.
- Fix.
- Regression test.

If no production code changed, explicitly state that.

## F. Security Findings

Report:

- Authentication issues.
- Authorization issues.
- Cross-user access risks.
- IDOR findings.
- Idempotency weaknesses.
- Information disclosure.
- Remaining security concerns.

Do not mark the module production-ready if a critical ownership or authorization vulnerability remains unresolved.

## G. Missing or Unimplemented Tests

List every test scenario that could not be implemented.

For each one, explain:

- What could not be tested.
- Why it could not be tested.
- What infrastructure or change is required.
- Whether the omission creates a production risk.

Never hide incomplete work.

## H. Remaining Risks

List any unresolved issues in:

- Handler.
- Service.
- Repository.
- Routes.
- Database.
- Middleware.
- Security.
- Concurrency.
- Test infrastructure.

## I. Final Engineering Assessment

Provide an honest assessment of the Workspace module.

Distinguish clearly between:

- Tests implemented.
- Tests executed.
- Tests passed.
- Tests skipped.
- Coverage achieved.
- Production issues fixed.
- Production issues remaining.

Do not call something production-ready simply because the test suite passes.

# 21. Final Non-Negotiable Instructions

1. Work like a Principal/Staff-level Go engineer.
2. Inspect before implementing.
3. Test all five APIs independently.
4. Cover every meaningful success path.
5. Cover every meaningful failure path.
6. Cover invalid input and boundary conditions.
7. Cover nil pointers and nullable values.
8. Cover database constraints and SQL behavior.
9. Cover context cancellation and timeouts.
10. Cover authentication and authorization.
11. Cover cross-user access.
12. Cover idempotency for POST, PATCH, and DELETE.
13. Cover concurrent requests where the implementation promises concurrency safety.
14. Cover error classification.
15. Cover response serialization.
16. Cover database persistence.
17. Cover soft-delete behavior.
18. Cover regression scenarios for discovered bugs.
19. Run the actual tests.
20. Run race detection.
21. Run coverage.
22. Run repeated tests.
23. Report actual results.
24. Never fabricate successful execution.
25. Never leave a critical security issue undisclosed.

Do not optimize for the number of test files or lines of test code. Optimize for confidence in the correctness, security, reliability, and maintainability of the Workspace module.

The objective is comprehensive, production-oriented test coverage of every API and every meaningful execution path—not superficial coverage percentages.

If an important behavior is missing from the current implementation, identify it, explain its impact, and address it with the smallest justified production change and a regression test.

Do not stop at writing tests. Finish the implementation, execute it, investigate failures, fix them, and deliver an evidence-based engineering report.

```

```
