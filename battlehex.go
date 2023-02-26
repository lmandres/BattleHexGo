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

    /*
    var secret Secret

    ctx := context.Background()

    // Open our jsonFile
    jsonFile, err := os.Open("client_secret.json")
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
    }
    // defer the closing of our jsonFile so that we can parse it later on
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)
    json.Unmarshal(byteValue, &secret)

    provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
    if err != nil {
        log.Fatal(err)
    }
    oidcConfig := &oidc.Config{
        ClientID: secret.Web.ClientID,
    }
    verifier := provider.Verifier(oidcConfig)

    config := oauth2.Config{
        ClientID:     secret.Web.ClientID,
        ClientSecret: secret.Web.ClientSecret,
        Endpoint:     provider.Endpoint(),
        RedirectURL:  secret.Web.RedirectURIs[0],
        Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
    }

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {

        state, err := randString(16)
        if err != nil {
            http.Error(w, "Internal error", http.StatusInternalServerError)
            return
        }
        nonce, err := randString(16)
        if err != nil {
            http.Error(w, "Internal error", http.StatusInternalServerError)
            return
        }
        setCallbackCookie(w, r, "state", state)
        setCallbackCookie(w, r, "nonce", nonce)

        http.Redirect(w, r, config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
    })
    http.HandleFunc("/auth/google/callback", func (w http.ResponseWriter, r *http.Request) {

        state, err := r.Cookie("state")
        if err != nil {
            http.Error(w, "state not found", http.StatusBadRequest)
            return
        }
        if r.URL.Query().Get("state") != state.Value {
            http.Error(w, "state did not match", http.StatusBadRequest)
            return
        }

        oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
        if err != nil {
            http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
            return
        }
        rawIDToken, ok := oauth2Token.Extra("id_token").(string)
        if !ok {
            http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
            return
        }
        idToken, err := verifier.Verify(ctx, rawIDToken)
        if err != nil {
            http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
            return
        }

        nonce, err := r.Cookie("nonce")
        if err != nil {
            http.Error(w, "nonce not found", http.StatusBadRequest)
            return
        }
        if idToken.Nonce != nonce.Value {
            http.Error(w, "nonce did not match", http.StatusBadRequest)
            return
        }

        oauth2Token.AccessToken = "*REDACTED*"

        resp := struct {
            OAuth2Token   *oauth2.Token
            IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
        }{oauth2Token, new(json.RawMessage)}

        if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        data, err := json.MarshalIndent(resp, "", "    ")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Write(data)
    })
    */

    http.HandleFunc("/battlehex_vs_js_ai_v1.1", BattleHexJSHandler)

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
