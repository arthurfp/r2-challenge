package auth

import (
    "net/http"
    "strings"
    "github.com/labstack/echo/v4"
)

const (
    CtxUserID = "userId"
    CtxRole   = "role"
)

func JWTMiddleware(tm TokenManager, isPublic func(method, path string) bool) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            method := c.Request().Method
            path := c.Path()
            if path == "" {
                if c.Request().URL != nil {
                    path = c.Request().URL.Path
                } else {
                    path = c.Request().RequestURI
                }
            }
            // Always allow public routes and special paths
            if isPublic(method, path) ||
                method == "OPTIONS" ||
                strings.HasPrefix(path, "/swagger") || path == "/swagger.yaml" {
                return next(c)
            }

            authz := c.Request().Header.Get("Authorization")
            if !strings.HasPrefix(authz, "Bearer ") {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
            }

            token := strings.TrimPrefix(authz, "Bearer ")
            claims, err := tm.Verify(token)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
            }

            if sub, ok := claims["sub"].(string); ok {
                c.Set(CtxUserID, sub)
            }
            if role, ok := claims["role"].(string); ok {
                c.Set(CtxRole, role)
            }

            return next(c)
        }
    }
}

func RequireRoles(allowed ...string) echo.MiddlewareFunc {
    allowedSet := map[string]struct{}{}
    for _, r := range allowed {
        allowedSet[r] = struct{}{}
    }

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            role, _ := c.Get(CtxRole).(string)
            if _, ok := allowedSet[role]; !ok {
                return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
            }
            return next(c)
        }
    }
}

// RequireSelfOrRoles allows access if the authenticated user's ID matches the
// given path parameter or if the user's role is one of the allowed roles.
func RequireSelfOrRoles(param string, allowed ...string) echo.MiddlewareFunc {
    allowedSet := map[string]struct{}{}
    for _, r := range allowed {
        allowedSet[r] = struct{}{}
    }

    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            role, _ := c.Get(CtxRole).(string)
            if _, ok := allowedSet[role]; ok {
                return next(c)
            }

            pathID := c.Param(param)
            userID, _ := c.Get(CtxUserID).(string)
            if pathID != "" && userID != "" && pathID == userID {
                return next(c)
            }

            return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
        }
    }
}


