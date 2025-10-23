# gFlyDev HTTP

This is `gflydev-http`, a Go library that provides HTTP utilities and abstractions for building REST APIs with the [gFlyDev framework](https://github.com/gflydev/core). It's a library package (not an executable), designed to be imported by other projects.

## Features

- **Generic Response Structures** - Type-safe response models for lists, success, and error responses
- **Request Processing Helpers** - Automated request validation, parsing, and sanitization pipeline
- **XSS Protection** - Built-in sanitization for all string inputs
- **Input Validation** - Integration with gFlyDev validation system
- **Pagination Support** - Ready-to-use filtering and pagination structures
- **Type Safety** - Generic functions for type-safe data transformations

## Core Dependencies

- `github.com/gflydev/core` - Core framework providing `core.Api`, `core.Ctx`, and `core.Data`
- `github.com/gflydev/validation` - Input validation utilities
- `github.com/gflydev/utils` - General utility functions including `fn.TransformList`
- Built on top of `valyala/fasthttp` (indirect dependency)

## Architecture

### Request Processing Pattern

The library implements a two-phase request processing pattern (Validate -> Handle):

1. **Validation Phase**: Extract, parse, sanitize, and validate request data
2. **Handle Phase**: Process the validated data and return response

### Core Components

**Response Structures** (`generic_response.go`):
- `List[T]` - Generic paginated list responses with metadata
- `Success` - Generic success responses with optional data
- `Error` - Error responses with code, message, and validation data
- `Meta` - Pagination metadata (page, per_page, total)

**Request Processing** (`request_helpers.go`):
Three generic helper functions that handle the full request processing pipeline:

- `ProcessPathID(c)` - Extracts and validates path ID parameter, stores in context as `DataPathID`
- `ProcessFilter(c)` - Parses query params (page, per_page, keyword, order_by), validates, stores as `DataFilter`
- `ProcessData[T AddData](c)` - Parses body, sanitizes, validates, stores as `DataRequest` (for CREATE)
- `ProcessUpdateData[T UpdateData](c)` - Same as ProcessData but also extracts path ID and calls `SetID()` (for UPDATE)

**HTTP Helpers** (`http_helpers.go`):
Low-level utilities used by the Process* functions:
- `PathID(c, idName)` - Extracts integer ID from path parameter
- `Parse[T](c, structData)` - Parses request body into struct
- `FilterData(c)` - Constructs Filter DTO from query parameters
- `Validate(structData)` - Performs validation.Check with error formatting

**Security** (`secure.go`):
- `SanitizeStruct(target)` - Recursively sanitizes all string fields in structs to prevent XSS
- `SanitizeString(input)` - Removes script tags, trims, unescapes HTML, removes null bytes
- Called automatically by `ProcessData` and `ProcessUpdateData`

**Transformers** (`generic_transformer.go`):
- `ToListResponse[T, R](records, transformerFn)` - Transforms lists of models to response DTOs

**Context Data Keys** (`constants.go`):
- `DataPathID` - Stores extracted path ID
- `DataRequest` - Stores parsed and validated request body
- `DataFilter` - Stores filter/pagination parameters

## Installation

```bash
go get github.com/gflydev/http
```

## Quick Start

### Basic API Implementation

```go
package api

import (
    "github.com/gflydev/core"
    "github.com/gflydev/http"
)

// Define your request structure
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// Implement the API
type CreateUserApi struct {
    core.Api
}

func (h *CreateUserApi) Validate(c *core.Ctx) error {
    // Automatically parses, sanitizes, and validates the request
    return http.ProcessData[CreateUserRequest](c)
}

func (h *CreateUserApi) Handle(c *core.Ctx) error {
    // Retrieve the validated data
    req := c.Data(http.DataRequest).(CreateUserRequest)

    // Your business logic here
    // user := createUser(req.Name, req.Email)

    return c.JSON(http.Success{
        Message: "User created successfully",
    })
}
```

### List API with Pagination

```go
type ListUsersApi struct {
    core.Api
}

func (h *ListUsersApi) Validate(c *core.Ctx) error {
    // Automatically parses pagination parameters: page, per_page, keyword, order_by
    return http.ProcessFilter(c)
}

func (h *ListUsersApi) Handle(c *core.Ctx) error {
    filter := c.Data(http.DataFilter).(http.Filter)

    // Fetch users with pagination
    users, total := fetchUsers(filter.Page, filter.PerPage, filter.Keyword)

    // Transform to response DTOs
    userResponses := http.ToListResponse(users, func(u User) UserResponse {
        return UserResponse{
            ID:    u.ID,
            Name:  u.Name,
            Email: u.Email,
        }
    })

    return c.JSON(http.List[UserResponse]{
        Meta: http.Meta{
            Page:    filter.Page,
            PerPage: filter.PerPage,
            Total:   total,
        },
        Data: userResponses,
    })
}
```

### Update API with Path ID

```go
type UpdateUserRequest struct {
    ID    int    `json:"id"`
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// Must implement SetID with pointer receiver
func (r *UpdateUserRequest) SetID(id int) {
    r.ID = id
}

type UpdateUserApi struct {
    core.Api
}

func (h *UpdateUserApi) Validate(c *core.Ctx) error {
    // Extracts ID from path, parses body, sanitizes, validates, and calls SetID
    return http.ProcessUpdateData[*UpdateUserRequest](c)
}

func (h *UpdateUserApi) Handle(c *core.Ctx) error {
    req := c.Data(http.DataRequest).(*UpdateUserRequest)

    // Update user with req.ID, req.Name, req.Email

    return c.JSON(http.Success{
        Message: "User updated successfully",
    })
}
```

### Delete API

```go
type DeleteUserApi struct {
    core.Api
}

func (h *DeleteUserApi) Validate(c *core.Ctx) error {
    // Extracts and validates ID from path parameter
    return http.ProcessPathID(c)
}

func (h *DeleteUserApi) Handle(c *core.Ctx) error {
    id := c.Data(http.DataPathID).(int)

    // Delete user by ID

    return c.JSON(http.Success{
        Message: "User deleted successfully",
    })
}
```

## API Reference

### Request Processing Functions

#### `ProcessData[T AddData](c *core.Ctx) error`
Processes CREATE requests by parsing the request body, sanitizing string fields, validating, and storing in context.

**Stores in context:** `http.DataRequest`

#### `ProcessUpdateData[T UpdateData](c *core.Ctx) error`
Processes UPDATE requests by extracting path ID, parsing body, sanitizing, setting ID on the request struct, validating, and storing in context.

**Requirements:** Type `T` must implement `SetID(int)` method with pointer receiver
**Stores in context:** `http.DataRequest`

#### `ProcessFilter(c *core.Ctx) error`
Processes list/filter requests by parsing query parameters (page, per_page, keyword, order_by), validating, and storing in context.

**Stores in context:** `http.DataFilter`

#### `ProcessPathID(c *core.Ctx) error`
Extracts and validates an integer ID from the path parameter, storing it in context.

**Stores in context:** `http.DataPathID`

### Response Structures

#### `Success`
```go
type Success struct {
    Message string    `json:"message"`
    Data    core.Data `json:"data"`
}
```

#### `Error`
```go
type Error struct {
    Code    string    `json:"code"`
    Message string    `json:"message"`
    Data    core.Data `json:"data"` // Validation errors
}
```

#### `List[T]`
```go
type List[T any] struct {
    Meta Meta `json:"meta"`
    Data []T  `json:"data"`
}

type Meta struct {
    Page    int `json:"page,omitempty"`
    PerPage int `json:"per_page,omitempty"`
    Total   int `json:"total"`
}
```

### Filter DTO

```go
type Filter struct {
    Page    int    `json:"page"`      // Current page number (default: 1)
    PerPage int    `json:"per_page"`  // Items per page (default: 10)
    Keyword string `json:"keyword"`   // Search keyword
    OrderBy string `json:"order_by"`  // Sort field (prefix with '-' for DESC)
}
```

### Helper Functions

#### `PathID(c *core.Ctx, idName ...string) (int, *Error)`
Extracts integer ID from path parameter. Default parameter name is "id".

#### `Parse[T any](c *core.Ctx, structData *T) *Error`
Parses request body into the provided struct.

#### `Validate(structData any, msgForTagFunc ...validation.MsgForTagFunc) *Error`
Validates struct using gFlyDev validation rules.

#### `SanitizeStruct(target any)`
Recursively sanitizes all string fields to prevent XSS attacks by removing script tags, trimming whitespace, and removing null bytes.

#### `ToListResponse[T any, R any](records []T, transformerFn func(T) R) []R`
Transforms a list of models to response DTOs using the provided transformer function.

## Security Features

### Automatic XSS Protection

All string fields are automatically sanitized when using `ProcessData` or `ProcessUpdateData`:
- Removes `<script>` tags
- Trims whitespace
- Unescapes HTML entities
- Removes null bytes

### Manual Sanitization

```go
// Sanitize entire struct
http.SanitizeStruct(&myStruct)

// Sanitize single string
clean := http.SanitizeString(userInput)
```

## Context Data Keys

The library uses these constant keys to store data in `core.Ctx.Data`:

- `http.DataPathID` - Extracted path ID parameter (int)
- `http.DataRequest` - Parsed and validated request body (interface{})
- `http.DataFilter` - Filter/pagination parameters (http.Filter)

## Requirements

- Go 1.24.0 or higher
- github.com/gflydev/core v1.17.11+
- github.com/gflydev/validation v1.2.1+

## License

[Add your license information here]

## Contributing

[Add contributing guidelines if applicable]
