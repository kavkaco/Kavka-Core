# Fix Delete Error Check Bug

**File:** `database/repo_mongo/message_repository.go:260`
**Priority:** P0 (Critical)
**Impact:** `Delete` silently fails when message doesn't exist

## Problem

Line 260:
```go
if err != nil && result.ModifiedCount != 1 {
```

Uses `&&` instead of `||`. When `err == nil` but `ModifiedCount == 0` (message not found), the condition is `false` and the function returns `nil` instead of an error.

## Fix

Change `&&` to `||`:
```go
if err != nil || result.ModifiedCount != 1 {
```

## Additional Issue

The `Create` method (line 126-137) returns `nil` instead of `err` on line 133:
```go
if err != nil {
    return nil  // should be "return err"
}
```

Fix this too.

## Validation

- Deleting a non-existent message returns `ErrNotModified`
- Deleting an existing message succeeds
- Message store creation propagates errors correctly
