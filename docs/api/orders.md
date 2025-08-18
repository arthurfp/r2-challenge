# Orders API

Base paths: `/v1/orders`, `/v1/users/{id}/orders`

## Models (domain)
```json
{
  "id": "string",
  "user_id": "string",
  "status": "created|paid|shipped|...",
  "total_cents": 1234,
  "items": [
    { "product_id": "string", "quantity": 1, "price_cents": 999 }
  ]
}
```

## Endpoints

### Place order (private)
POST `/v1/orders`
- Body: `items[{product_id, quantity, price_cents}]` (user comes from JWT)
- Success: 201 `Order`
- Errors: 400 validation, 401, 500

Example:
```bash
curl -s -X POST http://localhost:8080/v1/orders \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"items":[{"product_id":"P1","quantity":1,"price_cents":199990}]}'
```

### Get by ID (private)
GET `/v1/orders/{id}`
- Success: 200 `Order`
- Errors: 400, 404

### List my orders (private)
GET `/v1/users/{id}/orders`
- Path: `{id}` must match authenticated user (enforced at handler level)
- Query: `limit`, `offset`
- Success: 200 `[Order]`
- Errors: 401, 500

### Update status (admin)
PUT `/v1/orders/{id}/status`
- Body: `{ "status": "shipped" }` (example)
- Success: 200 `Order`
- Errors: 400, 401/403, 500

## Error handling (patterns)
- Consistent `{ "error": "..." }` body across 4xx/5xx
- Business errors return appropriate HTTP status (404 not found, 401/403 auth)

## Notes
- Payment capture is mocked but persisted to `payments` table
- Notifications are mocked (no-op sender)
