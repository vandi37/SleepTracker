# SleepTracker API Reference

Base URL: `http://localhost:8080` (or your configured `PORT`)

All JSON responses use `application/json`. Protected routes require:

```
Authorization: Bearer <access_token>
```

---

## Public endpoints

### Health check

```
GET /ping
```

**Response:** `200 OK` — body: `pong`

---

### Register

```
POST /register
Content-Type: application/json
```

**Request body:**

```json
{
  "username": "string",
  "password": "string",
  "nickname": "string",
  "birth": "YYYY-MM-DD"
}
```

- `username` — 3–40 chars, alphanumeric and underscore
- `nickname` — required
- `birth` — date in `YYYY-MM-DD` format

**Response:** `200 OK` — `UserWithToken` (id, access, refresh, expires)

---

### Login

```
POST /login
Content-Type: application/json
```

**Request body:**

```json
{
  "username": "string",
  "password": "string"
}
```

**Response:** `200 OK` — `UserWithToken` (id, access, refresh, expires)

---

### Refresh tokens

```
POST /refresh
Content-Type: application/json
```

**Request body:**

```json
{
  "token": "refresh_token"
}
```

**Response:** `200 OK` — `UserWithToken` with new access and refresh tokens

---

## Protected endpoints

### Users

#### Get current user

```
GET /users/
Authorization: Bearer <access_token>
```

**Response:** `200 OK` — `User` (id, username, nickname, birth, created_at)

---

#### Get user by ID

```
GET /users/:id
Authorization: Bearer <access_token>
```

**Response:** `200 OK` — `User`

---

#### Update profile

```
PUT /users/
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request body:**

```json
{
  "nickname": "string",
  "birth": "YYYY-MM-DD"
}
```

**Response:** `204 No Content`

---

#### Change password

```
PATCH /users/password
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request body:**

```json
{
  "password": "string"
}
```

**Response:** `204 No Content`

---

#### Delete account

```
DELETE /users/
Authorization: Bearer <access_token>
```

**Response:** `204 No Content`

---

### Friends

#### Send friend request

```
POST /friends/request
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request body:**

```json
{
  "id": 123
}
```

`id` — target user ID

**Response:** `200 OK` — created `Friend` object

---

#### Accept friend request

```
POST /friends/accept
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request body:**

```json
{
  "id": 456
}
```

`id` — friendship ID (from `/friends/` list)

**Response:** `200 OK` — updated `Friend` object

---

#### List friendships

```
GET /friends/?limit=10&offset=0
Authorization: Bearer <access_token>
```

**Query params:**

- `limit` — min 1
- `offset` — min 0

**Response:** `200 OK` — `{ "friendships": [ Friend, ... ] }`

---

#### Delete friendship

```
DELETE /friends/:id
Authorization: Bearer <access_token>
```

`id` — friendship ID

**Response:** `204 No Content`

---

#### Friend's sleep history (by week)

```
GET /friends/:id/history/:page
Authorization: Bearer <access_token>
```

`id` — friendship ID  
`page` — week index (0 = current week)

**Response:** `200 OK` — `{ "sleeps": [ Sleep, ... ] }`  
Only available for accepted friendships.

---

#### Friend's scores (by year)

```
GET /friends/:id/table/:page
Authorization: Bearer <access_token>
```

`id` — friendship ID  
`page` — year index

**Response:** `200 OK` — `{ "scores": [ SleepScore, ... ] }`  
Only available for accepted friendships.

---

### Sleep

#### Add sleep record

```
POST /sleep/
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request body:**

```json
{
  "sleep_time": 1320,
  "wake_time": 420,
  "score": 85,
  "enter_date": "2025-03-08"
}
```

- `sleep_time`, `wake_time` — minutes since midnight (0–2159), optional (both null or both set; sleep must be < wake)
- `score` — 0–100
- `enter_date` — `YYYY-MM-DD`, unique per user per day

**Response:** `200 OK` — created `Sleep` object

---

#### Update sleep record

```
PUT /sleep/:id
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request body:** same fields as add

**Response:** `200 OK` — updated `Sleep` object

---

#### Delete sleep record

```
DELETE /sleep/:id
Authorization: Bearer <access_token>
```

**Response:** `204 No Content`

---

#### Own sleep history (by week)

```
GET /history/:page
Authorization: Bearer <access_token>
```

`page` — week index

**Response:** `200 OK` — `{ "sleeps": [ Sleep, ... ] }`

---

#### Own scores (by year)

```
GET /table/:page
Authorization: Bearer <access_token>
```

`page` — year index

**Response:** `200 OK` — `{ "scores": [ SleepScore, ... ] }`

---

## Error format

Errors return JSON in the form:

```json
{
  "status": 400,
  "message": "description"
}
```

Common status codes: `400` (bad request), `401` (unauthorized), `403` (forbidden), `404` (not found), `409` (conflict), `500` (internal error).
