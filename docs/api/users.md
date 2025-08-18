# Users

Base paths: `/v1/auth` and `/v1/users`

## Register (public)
POST `/v1/auth/register`
```bash
curl -s -X POST http://localhost:8080/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","name":"User","password":"secret"}'
```

## Login (public)
POST `/v1/auth/login`
```bash
curl -s -X POST http://localhost:8080/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"user@example.com","password":"secret"}'
```
Response includes `{"token":"<JWT>"}`.

## Get user by ID (private)
GET `/v1/users/:id`
```bash
curl -H 'Authorization: Bearer <JWT>' http://localhost:8080/v1/users/USER_ID
```

## List users (private)
GET `/v1/users`
```bash
curl -H 'Authorization: Bearer <JWT>' http://localhost:8080/v1/users
```

## Update my profile (private)
PUT `/v1/users/me`
```bash
curl -s -X PUT http://localhost:8080/v1/users/me \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"New Name"}'
```
