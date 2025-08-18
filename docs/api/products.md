# Products API

Base path: `/v1/products`

## Model (domain)
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "category": "string",
  "price_cents": 1234,
  "inventory": 10
}
```

## Endpoints

### List
GET `/v1/products`
- Query: `category`, `name`, `sort` (e.g., `price_cents`), `order` (`asc|desc`), `limit`, `offset`
- Success: 200 `[Product]`
- Error: 500 `{ "error": "..." }`

Example:
```bash
curl -s 'http://localhost:8080/v1/products?category=books&sort=price_cents&order=asc&limit=10'
```

### Get by ID
GET `/v1/products/{id}`
- Success: 200 `Product`
- Errors: 400 invalid id, 404 not found

Example:
```bash
curl -s http://localhost:8080/v1/products/PRODUCT_ID
```

### Create (admin)
POST `/v1/products`
- Auth: `Authorization: Bearer <JWT>` with `role=admin`
- Body: `name`, `category`, `description?`, `price_cents`, `inventory?`
- Success: 201 `Product`
- Errors: 400 validation, 401/403 auth, 500

Example:
```bash
curl -s -X POST http://localhost:8080/v1/products \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Notebook","category":"electronics","price_cents":199990,"inventory":10}'
```

### Update (admin)
PUT `/v1/products/{id}`
- Body: same as create
- Success: 200 `Product`
- Errors: 400, 404, 401/403

### Delete (admin)
DELETE `/v1/products/{id}`
- Success: 204
- Errors: 400, 404, 401/403

## Error handling (patterns)
- Validation: 400 `{ "error": "<validation message>" }`
- Not found: 404 `{ "error": "not found" }`
- Unauthorized/Forbidden: 401/403 `{ "error": "..." }`
- Internal: 500 `{ "error": "..." }`

## Tips
- Filtering + pagination are optional; defaults are safe
- Sorting accepts domain field names (e.g., `price_cents`)
