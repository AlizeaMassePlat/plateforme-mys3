package middleware

import (
    "net/http"
    "strings"
    "log"
    "bytes"
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

func LogRequestMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Requête reçue 1: %s %s", r.Method, r.RequestURI)

        if len(r.URL.Query()) > 0 {
            log.Printf("Query Params: %v", r.URL.Query())
        }

        next.ServeHTTP(w, r)
    })
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
    responseBody bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
    lrw.responseBody.Write(b) 
    return lrw.ResponseWriter.Write(b) 
}

func LogResponseMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(lrw, r)
        
        // Log la réponse
        log.Printf("Response status: %d", lrw.statusCode)
        log.Printf("Response body: %s", lrw.responseBody.String())
    })
}

// CorsMiddleware permet de configurer les en-têtes CORS
func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Autorise les requêtes provenant de localhost:3000 
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        
        // Autorise les méthodes HTTP spécifiques
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        
        // Autorise les en-têtes spécifiques
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        // Si c'est une requête OPTIONS, on renvoie simplement les en-têtes CORS
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        // Appelle le prochain handler dans la chaîne
        next.ServeHTTP(w, r)
    })
}
