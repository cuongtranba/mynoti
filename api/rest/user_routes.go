package rest

import (
    "net/http"
)

// UserRoutes sets up the user routes.
func UserRoutes() {
    http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("User endpoint"))
    })
}
