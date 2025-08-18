# Products

Base path: `/v1/products`

## List
GET `/v1/products`
```bash
curl -s http://localhost:8080/v1/products
```

## Get by ID
GET `/v1/products/:id`
```bash
curl -s http://localhost:8080/v1/products/PRODUCT_ID
```

## Create (admin)
POST `/v1/products`
Headers: `Authorization: Bearer <JWT>`
```bash
curl -s -X POST http://localhost:8080/v1/products \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Laptop","category":"electronics","description":"...","price":1999.90,"inventory":10}'
```

## Update (admin)
PUT `/v1/products/:id`
```bash
curl -s -X PUT http://localhost:8080/v1/products/PRODUCT_ID \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"name":"Laptop Pro","price":2499.90,"inventory":8}'
```

## Delete (admin)
DELETE `/v1/products/:id`
```bash
curl -s -X DELETE http://localhost:8080/v1/products/PRODUCT_ID \
  -H 'Authorization: Bearer <JWT>'
```
