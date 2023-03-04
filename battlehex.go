package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "cloud.google.com/go/firestore"
    firestoregorilla "github.com/GoogleCloudPlatform/firestore-gorilla-sessions"

    "battle-hex-go/frontend"
)

var store *firestoregorilla.Store

func main() {

    http.HandleFunc("/battlehex_vs_js_ai_v1.1", frontend.BattleHexJSHandler)
    http.HandleFunc("/set_session", setSessionHandler)
    http.HandleFunc("/get_session", getSessionHandler)

    port :=  os.Getenv("PORT")
    if port == "" {
        port = "8080"
        log.Printf("Defaulting to port %s", port)
    }

    projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
    if projectID == "" {
        log.Fatal("GOOGLE_CLOUD_PROJECT must be set")
    }

    var err error
    store, err = newSessionStore(projectID)
    if err != nil {
        log.Fatalf("newSessionStore: %v", err)
    }

    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}

func newSessionStore(projectID string) (*firestoregorilla.Store, error) {

    ctx := context.Background()
    client, err := firestore.NewClient(ctx, projectID)
    if err != nil {
        log.Fatalf("firestore.NewClient: %v", err)
    }

    store, err := firestoregorilla.New(ctx, client)
    if err != nil {
        log.Fatalf("firestoregorilla.New: %v", err)
    }

    return store, nil
}

func setSessionHandler(w http.ResponseWriter, r *http.Request) {
    // Get a session. Get() always returns a session, even if empty.
    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set some session values.
    session.Values["foo"] = "bar"
    // Save it before we write to the response/return from the handler.
    err = session.Save(r, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "session.Values[\"foo\"] = %s<br />\n", session.Values["foo"])
}

func getSessionHandler(w http.ResponseWriter, r *http.Request) {

    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


    fmt.Fprintf(w, "session.Values[\"foo\"] = %s<br />\n", session.Values["foo"])
}
