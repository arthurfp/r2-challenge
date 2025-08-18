# Users API

Base paths: `/v1/auth`, `/v1/users`

## Model (domain)
```json
{
  "id": "string",
  "email": "string",
  "name": "string",
  "role": "user|admin"
}
```

## Auth

### Register (public)
POST `/v1/auth/register`
- Body: `email`, `name`, `password`
- Success: 201 `User`
- Errors: 400 validation, 500

Example:
```bash
curl -s -X POST http://localhost:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","name":"User","password":"secret"}'
```

### Login (public)
POST `/v1/auth/login`
- Body: `email`, `password`
- Success: 200 `{ "access_token": "..." }`
- Errors: 400, 401 invalid credentials

## Users

### Get by ID (private)
GET `/v1/users/{id}`
- Success: 200 `User`
- Errors: 400 invalid id, 404 not found

### List (private)
GET `/v1/users`
- Query: `email`, `name`, `limit`, `offset`
- Success: 200 `[User]`
- Error: 500

### Update my profile (private)
PUT `/v1/users/me`
- Body: `name`, `email`
- Success: 200 `User`
- Errors: 400 validation, 401 unauthorized, 500

## Error handling (patterns)
- Same shapes as Products; 401/403 when JWT missing/invalid or role insufficient
