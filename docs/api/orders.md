# Orders

Base path: `/v1/orders`

## Create order
POST `/v1/orders`
```bash
curl -s -X POST http://localhost:8080/v1/orders \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"USER_ID","items":[{"product_id":"PRODUCT_ID","quantity":1}],"payment_method":"mock"}'
```

## Get order by ID
GET `/v1/orders/:id`
```bash
curl -s -H 'Authorization: Bearer <JWT>' http://localhost:8080/v1/orders/ORDER_ID
```

## List orders for a user
GET `/v1/users/:id/orders`
```bash
curl -s -H 'Authorization: Bearer <JWT>' http://localhost:8080/v1/users/USER_ID/orders
```

## Update status (admin)
PUT `/v1/orders/:id/status`
```bash
curl -s -X PUT http://localhost:8080/v1/orders/ORDER_ID/status \
  -H 'Authorization: Bearer <JWT>' \
  -H 'Content-Type: application/json' \
  -d '{"status":"shipped"}'
```
