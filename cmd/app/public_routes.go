package main

const (
    GET  = "GET"
    POST = "POST"
    PUT  = "PUT"
    DELETE  = "DELETE"
)

type routeKey struct {
    Method string
    Path   string
}

// publicRoutes defines routes that do not require JWT auth.
var publicRoutes = map[routeKey]struct{}{
    {Method: POST, Path: "/v1/auth/register"}: {},
    {Method: POST, Path: "/v1/auth/login"}:   {},
    {Method: GET, Path: "/v1/products"}:      {},
    {Method: GET, Path: "/v1/products/:id"}:  {},
    {Method: GET, Path: "/v1/orders/:id"}:    {},
}


