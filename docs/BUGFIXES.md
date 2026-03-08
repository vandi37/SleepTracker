# Bug Fixes

## Backend bugs fixed

### 1. SQL syntax error in `user_repo/repo.go` (line 98)

**Bug:** Extra comma in SQL query caused login to fail.

```go
// Before (broken)
`select id, password_hash, from users where username = $1`

// After (fixed)
`select id, password_hash from users where username = $1`
```

### 2. Missing query parameters in `friend_repo/repo.go` (GetSecond)

**Bug:** The query uses `$1` and `$2` but only `id` was passed. This caused friend sleep history and scores to fail.

```go
// Before (broken)
tx.QueryRowContext(ctx, `... where id = $1 and (user1_id = $2 or user2_id = $2)`, id)

// After (fixed)
tx.QueryRowContext(ctx, `... where id = $1 and (user1_id = $2 or user2_id = $2)`, id, user_id)
```

### 3. Wrong error response in `handler/handler.go` (GetSelf, GetUser)

**Bug:** `ctx.AbortWithError()` was used, which does not send a JSON body to the client. The client received an empty response on errors.

```go
// Before (broken)
ctx.AbortWithError(err.Code(), err.JsonError())

// After (fixed)
ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
```

### 4. Incorrect token expiry time in `service/service.go` (tokens)

**Bug:** `expires` was set to a time in the past (using `Add(-duration)`), so the frontend received an expiry timestamp that had already passed.

```go
// Before (broken)
expires := time.Now().Add(-s.AccessJwt.GetExpiration())

// After (fixed)
expires := time.Now().Add(s.AccessJwt.GetExpiration())
```

### 5. Variable shadowing in `user_repo/repo.go` (Update)

**Bug:** When `err` was a non-pq Error, the type assertion `err, ok := err.(*pq.Error)` set `err` to `nil`, so `err != nil` was false and the real error was silently dropped.

```go
// Before (broken)
if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
    return UsernameTakenError(username)
} else if err != nil {
    // ... - when err was non-pq, we'd skip this
}

// After (fixed)
if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
    return UsernameTakenError(username)
}
if err != nil {
    // ... - now correctly handles all errors
}
```

---

## Frontend bugs fixed

### 1. API client — 204 No Content

**Bug:** `res.json()` on a 204 response throws (empty body). Update/delete endpoints return 204.

**Fix:** Check `res.status === 204` and return without parsing.

### 2. Router — user not loaded on auth routes

**Bug:** `user` stayed null until a component fetched it, so FriendsView `otherUser()` used uid 0 and showed the wrong person.

**Fix:** Call `ensureUser()` in router `beforeEach` for auth routes.

### 3. FriendsView — Accept button for wrong user

**Bug:** Accept was shown for every pending request; only the recipient (user2) can accept.

**Fix:** `v-if="!f.is_accepted && f.user2.id === user?.id"`.

### 4. ProfileView — loading state

**Bug:** Form appeared before user data loaded, causing empty/flash state.

**Fix:** Add `profileLoading` and show "Loading…" until fetch completes.

### 5. ProfileView — update missing username

**Bug:** Backend `UpdateUser` requires `username`; sending only nickname/birth caused validation errors.

**Fix:** Include `username` in the update body.
