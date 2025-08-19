import http from 'k6/http';
import { check, sleep } from 'k6';

// Usage:
// K6_BASE_URL=http://localhost:8080 K6_ADMIN_TOKEN=... K6_USER_TOKEN=... k6 run scripts/k6/order_place.js

export const options = {
  vus: Number(__ENV.K6_VUS || 10),
  duration: __ENV.K6_DURATION || '1m',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<400'],
  },
};

const baseURL = __ENV.K6_BASE_URL || 'http://localhost:8080';
const adminToken = __ENV.K6_ADMIN_TOKEN || '';
const userToken = __ENV.K6_USER_TOKEN || '';

export default function () {
  // List products (public)
  const listRes = http.get(`${baseURL}/v1/products`);
  check(listRes, {
    'list products 200': (r) => r.status === 200,
  });

  let chosen = null;
  try {
    const arr = listRes.json();
    if (Array.isArray(arr) && arr.length > 0) {
      // Distribui entre produtos para evitar esgotar estoque de um único item
      const idx = (Number(__ITER) + Number(__VU)) % arr.length;
      chosen = arr[idx];
    }
  } catch (_) {}

  // Place order (private) com Idempotency-Key
  const idemKey = `ord-${__VU}-${Date.now()}`;
  const payload = JSON.stringify({
    items: [
      chosen
        ? { product_id: chosen.id, quantity: 1, price_cents: chosen.price_cents || 100 }
        : { product_id: '00000000-0000-0000-0000-000000000000', quantity: 1, price_cents: 100 },
    ],
  });

  const headers = {
    'Content-Type': 'application/json',
    Authorization: userToken ? `Bearer ${userToken}` : '',
    'Idempotency-Key': idemKey,
  };
  const orderRes = http.post(`${baseURL}/v1/orders`, payload, { headers });
  check(orderRes, {
    'place order 201/200': (r) => r.status === 201 || r.status === 200,
  });

  // Retry same request (should deduplicate)
  const orderRes2 = http.post(`${baseURL}/v1/orders`, payload, { headers });
  check(orderRes2, {
    'idempotent same response': (r) => r.status === orderRes.status,
  });

  // Admin-only list users to measure RBAC path (optional)
  if (adminToken) {
    const resUsers = http.get(`${baseURL}/v1/users`, { headers: { Authorization: `Bearer ${adminToken}` } });
    check(resUsers, { 'list users 200|403': (r) => r.status === 200 || r.status === 403 });
  }

  sleep(1);
}


