# Go Project Naming & Coding Convention

## General Naming Rules
1. **Structs and Interfaces**
   - Use PascalCase (UpperCamelCase) for exported types, camelCase for unexported.
     - Example: `User`, `FriendshipService`
   - Do **not** prefix interfaces with `I` (Go convention).
     - Example: `Repository`, `Service`

2. **Variables and Parameters**
   - Use camelCase for local variables and parameters.
     - Example: `userID`, `friendUserID`, `retryCount`
   - Use PascalCase for exported fields.
     - Example: `UserID`, `CreatedAt`

3. **Functions and Methods**
   - Use PascalCase for exported, camelCase for unexported.
     - Example: `Create`, `getFriendshipStatus`
   - Use verbs for function/method names.
   - Use `ctx` as the first parameter for context.
     - Example: `func (s *FriendshipService) Create(ctx context.Context, userID, friendUserID string)`
   - Receiver name: dùng viết tắt tên struct, ví dụ: `func (s *Service)`

4. **Constants**
   - Use PascalCase for exported constants.
     - Example: `FriendshipStatusAccepted`, `ErrNotFound`
   - Prefix error variables with `Err`.
   - Group related constants using `const (...)`.
     - Example:
       ```go
       const (
           FriendshipStatusPending   = "pending"
           FriendshipStatusAccepted = "accepted"
           FriendshipStatusRejected = "rejected"
       )
       ```

5. **Database Fields**
   - Use snake_case for DB columns/fields.
     - Example: `user_id`, `friend_user_id`

6. **Packages**
   - Use all lowercase, no underscores, no plurals.
     - Example: `service`, `adapter`, `constant`

7. **File Names**
   - Use lowercase, words separated by underscores.
     - Example: `user.go`, `friendship_service.go`, `user_test.go`

## Examples

### Structs
```go
type Friendship struct {
    UserID       string
    FriendUserID string
    CreatedAt    time.Time
}
```

### Functions
```go
func (s *FriendshipService) Create(ctx context.Context, userID, friendUserID string) error {
    // ...
}
```

### Variables
```go
var retryCount int
var isActive bool
var friendIDs []string
```

### Constants
```go
const (
    FriendshipStatusPending   = "pending"
    FriendshipStatusAccepted = "accepted"
    FriendshipStatusRejected = "rejected"
)

var (
    ErrInvalidInput = errors.New("invalid input")
    ErrNotFound     = errors.New("not found")
)
```

### Packages & Files
- `service`, `adapter`, `constant`
- `user.go`, `friendship_service.go`, `user_test.go`

## Logging
- Use descriptive, structured log messages.
- Include context fields where possible.
```go
log.WithFields(log.Fields{
    "user_id": userID,
    "friend_id": friendID,
}).Info("friendship created")
```

## Error Handling
- Prefix error variables with `Err`.
- Use descriptive error messages.
- Group error definitions logically.

## Testing
- Use PascalCase for test function names.
- Use meaningful names for mocks.
```go
func TestFriendshipService_Create(t *testing.T) {
    // ...
}
```

## Comments
- Use complete sentences for comments.
- Document exported functions and types.
```go
// FriendshipService provides methods for managing friendships.
type FriendshipService struct {
    // ...
}
```

## Imports
- Group imports: standard, third-party, project-specific. Separate by blank lines.
- Use aliases only if necessary.
```go
import (
    "fmt"
    "time"

    "github.com/some/package"

    "project/internal/repository"
)
```
- Remove unused imports.
- Use full path for project imports.

## Consistency
- Use linters and formatters to enforce style.
- Ensure consistent naming and organization across all layers.

---

This document provides naming and organization conventions for Go projects. Adhering to these rules ensures readability, maintainability, and consistency across the codebase.
