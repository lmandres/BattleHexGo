package main

import (
    //"encoding/json"
    "fmt"
    //"io/ioutil"
    "log"
    "net/http"
    "os"

    "battle-hex-go/frontend"

    "github.com/gorilla/sessions"
    //"github.com/coreos/go-oidc/v3/oidc"
    //"golang.org/x/net/context"
    //"golang.org/x/oauth2"
)

var store = sessions.NewCookieStore([]byte("development_session_key"))

func main() {

    fs := http.FileServer(http.Dir("./static"))
    http.HandleFunc("/battlehex_vs_js_ai_v1.1", frontend.BattleHexJSHandler)
    http.HandleFunc("/set_session", setSessionHandler)
    http.HandleFunc("/get_session", getSessionHandler)
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    port :=  os.Getenv("PORT")
    if port == "" {
        port = "8080"
        log.Printf("Defaulting to port %s", port)
    }

    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
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
    session.Values[42] = 43
    // Save it before we write to the response/return from the handler.
    err = session.Save(r, w)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "session.Values[\"foo\"] = %s<br />\n", session.Values["foo"])
    fmt.Fprintf(w, "session.Values[42] = %d<br />\n", session.Values[42])
}

func getSessionHandler(w http.ResponseWriter, r *http.Request) {

    session, err := store.Get(r, "session-name")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }


    fmt.Fprintf(w, "session.Values[\"foo\"] = %s<br />\n", session.Values["foo"])
    fmt.Fprintf(w, "session.Values[42] = %d<br />\n", session.Values[42])
}
