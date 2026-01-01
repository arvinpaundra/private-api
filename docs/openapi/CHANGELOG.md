# OpenAPI Specification Changelog

## 2024-01-XX - Major Response Format Update

### Fixed

#### 1. RegisterRequest Schema

**Before:**

```yaml
RegisterRequest:
  properties:
    name: string
    email: string
    username: string
    password: string
```

**After:**

```yaml
RegisterRequest:
  properties:
    email: string (required, format: email, maxLength: 50)
    password: string (required, minLength: 8)
    fullname: string (required, minLength: 3, maxLength: 100)
```

**Reasoning:** The actual API implementation (`UserRegisterCommand` in `domain/auth/service/user_register.go`) only accepts `email`, `password`, and `fullname` fields. The `name` and `username` fields did not exist in the actual implementation.

---

#### 2. Response Format Structure

**Before:**

```yaml
{ 'status': 'success|error', 'message': 'some message', 'data': { ... } }
```

**After:**

```yaml
{
  "meta": {
    "code": 200,
    "message": "some message",
    "pagination": {...}  // optional, only for paginated responses
  },
  "data": {...},
  "errors": {...}  // optional, only for validation errors
}
```

**Reasoning:** The actual API uses `format.Response` structure defined in `core/format/response.go` which has:

- `meta` object with `code` (int), `message` (string), and optional `pagination`
- `data` field for response payload
- `errors` field for validation errors

---

#### 3. Response Schemas Updated

**New Schemas Created:**

- `Meta`: Contains code, message, and optional pagination
- `SuccessResponse`: Standard success response with meta and data
- `ErrorResponse`: Standard error response with meta and data
- `ValidationErrorResponse`: Error response with meta, data, and errors object

**Old Schemas Removed:**

- `Error`: Replaced by `ErrorResponse`
- `ValidationError`: Replaced by `ValidationErrorResponse`

---

#### 4. Response Examples Updated

All endpoint response examples have been updated to match the actual API response format:

**Register Endpoint (POST /v1/auth/register):**

```yaml
{
  'meta': { 'code': 201, 'message': 'user registered successfully' },
  'data': null,
}
```

**Login Endpoint (POST /v1/auth/login):**

```yaml
{
  'meta': { 'code': 200, 'message': 'login successful' },
  'data':
    {
      'access_token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
      'refresh_token': 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...',
      'token_type': 'Bearer',
      'expires_in': 900,
    },
}
```

**Validation Error (400 Bad Request):**

```yaml
{
  'meta': { 'code': 400, 'message': 'invalid request body' },
  'data': null,
  'errors': { 'email': 'email is required' },
}
```

**Error Response (401 Unauthorized, 403 Forbidden, 404 Not Found, 409 Conflict, 500 Internal Server Error):**

```yaml
{ 'meta': { 'code': 401, 'message': 'unauthorized' }, 'data': null }
```

---

#### 5. Pagination Structure

**Format:**

```yaml
Pagination:
  page: integer # Current page number
  per_page: integer # Items per page
  total: integer # Total number of items
  total_pages: integer # Total number of pages
```

**Note:** Pagination appears in `meta.pagination` for paginated list responses.

---

### Implementation Notes

The changes align the OpenAPI specification with the actual Go API implementation:

1. **Response Helpers Used:**

   - `format.SuccessOK(message, data, pagination...)` - 200 OK
   - `format.SuccessCreated(message, data)` - 201 Created
   - `format.BadRequest(message, errors)` - 400 Bad Request
   - `format.Unauthorized(message)` - 401 Unauthorized
   - `format.Forbidden(message)` - 403 Forbidden
   - `format.NotFound(message)` - 404 Not Found
   - `format.Conflict(message)` - 409 Conflict
   - `format.InternalServerError()` - 500 Internal Server Error

2. **Response Structure (`core/format/response.go`):**

   ```go
   type Meta struct {
       Code       int         `json:"code"`
       Message    string      `json:"message"`
       Pagination *Pagination `json:"pagination,omitempty"`
   }

   type Response struct {
       Meta   Meta            `json:"meta"`
       Data   any             `json:"data"`
       Errors validator.Error `json:"errors,omitempty"`
   }
   ```

3. **Pagination Structure (`core/format/pagination.go`):**
   ```go
   type Pagination struct {
       Page       int `json:"page"`
       PerPage    int `json:"per_page"`
       Total      int `json:"total"`
       TotalPages int `json:"total_pages"`
   }
   ```

---

### Validation Status

✅ OpenAPI specification passes validation:

```bash
npx @apidevtools/swagger-cli validate docs/openapi/openapi.yaml
# Output: docs/openapi/openapi.yaml is valid
```

---

### Next Steps

To test the updated OpenAPI specification:

1. **View in Swagger UI:**

   ```bash
   npx swagger-ui-watcher docs/openapi/openapi.yaml
   ```

2. **Generate API Client:**

   ```bash
   # For TypeScript
   npx @openapitools/openapi-generator-cli generate \
     -i docs/openapi/openapi.yaml \
     -g typescript-axios \
     -o generated/typescript-client

   # For Go
   npx @openapitools/openapi-generator-cli generate \
     -i docs/openapi/openapi.yaml \
     -g go \
     -o generated/go-client
   ```

3. **Run API Tests:**
   Ensure all integration tests pass with the corrected response format expectations.
