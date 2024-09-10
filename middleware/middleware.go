package middleware

import (
    "net/http"
    "strings"
)

func BasicAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/probe-bsign") {
            next.ServeHTTP(w, r)
            return
        }

        // Accepter les requêtes avec AWS4-HMAC-SHA256 dans l'en-tête Authorization
        if strings.Contains(r.Header.Get("Authorization"), "AWS4-HMAC-SHA256") {
            next.ServeHTTP(w, r)
            return
        }

        // Appliquer l'authentification basique pour les autres routes
        user, pass, ok := r.BasicAuth()
        if !ok || user != "accessuser" || pass != "accesspassword" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
