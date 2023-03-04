package main

import (
    "fmt"
    "io"
    "log"
    "math/rand"
    "net/http"
    "os"

    "battle-hex-go/frontend"
)

func main() {

    fs := http.FileServer(http.Dir("./static"))
    http.HandleFunc("/battlehex_vs_js_ai_v1.1", frontend.BattleHexJSHandler)

    port :=  os.Getenv("PORT")
    if port == "" {
        port = "8080"
        log.Printf("Defaulting to port %s", port)
    }

    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
    }

    session, err := newSession(projectID)
    if err != nil {
        log.Fatal("newSession: %v", err)
    }

    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}
