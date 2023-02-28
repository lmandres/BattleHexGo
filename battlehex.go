package main

import (
    "crypto/rand"
    "encoding/base64"
    //"encoding/json"
    //"fmt"
    "io"
    //"io/ioutil"
    "log"
    "net/http"
    "os"
    "time"

    "battle-hex-go/helper"

    //"github.com/coreos/go-oidc/v3/oidc"
    //"golang.org/x/net/context"
    //"golang.org/x/oauth2"
)

type Secret struct {
    Web Web `json:"web"`
}

type Web struct {
    ClientID string `json:"client_id"`
    ProjectID string `json:"project_id"`
    AuthURI string `json:"auth_uri"`
    TokenURI string `json:"token_uri"`
    AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
    ClientSecret string `json:"client_secret"`
    RedirectURIs []string `json:"redirect_uris"`
}

func main() {

    fs := http.FileServer(http.Dir("./static"))
    http.HandleFunc("/battlehex_vs_js_ai_v1.1", helper.BattleHexJSHandler)
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

func randString(nByte int) (string, error) {
    b := make([]byte, nByte)
    if _, err := io.ReadFull(rand.Reader, b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
    c := &http.Cookie{
        Name:     name,
        Value:    value,
        MaxAge:   int(time.Hour.Seconds()),
        Secure:   r.TLS != nil,
        HttpOnly: true,
    }
    http.SetCookie(w, c)
}
