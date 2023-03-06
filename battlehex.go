package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "cloud.google.com/go/compute/metadata"
    "cloud.google.com/go/firestore"
    "github.com/golang-jwt/jwt"
    firestoregorilla "github.com/GoogleCloudPlatform/firestore-gorilla-sessions"

    "battle-hex-go/frontend"
)

const sessionName string = "session-name"

func main() {

    projectID, err := metadata.ProjectID()
    if err != nil {
        log.Fatalf("metadata.ProjectID: %v", err)
    }

    session, err := newSession(projectID)
    if err != nil {
        log.Fatalf("newSessionStore: %v", err)
    }

    setSessionHandler := getSetSessionHandlerFunc(session)
    getSessionHandler := getGetSessionHandlerFunc(session)

    certs, aud, err := getAuthItems()
    if err != nil {
        log.Fatalf("getAuthItems: %v", err)
    }
    authIndexHandler := getAuthIndexHandlerFunc(certs, aud)

    http.HandleFunc("/", authIndexHandler)
    http.HandleFunc("/battlehex_vs_js_ai_v1.1", frontend.BattleHexJSHandler)
    http.HandleFunc("/set_session", setSessionHandler)
    http.HandleFunc("/get_session", getSessionHandler)

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

func newSession(projectID string) (*firestoregorilla.Store, error) {

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

func getAuthItems() (map[string]string, string, error) {

    certs, err := certificates()
    if err != nil {
        return nil, "", err
    }

    aud, err := audience()
    if err != nil {
        return certs, "", err
    }

    return certs, aud, nil
}

func getAuthIndexHandlerFunc(certs map[string]string, aud string) func(http.ResponseWriter, *http.Request) {

    // function responds to requests with our greeting.
    return func (w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }

        assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
        if assertion == "" {
           fmt.Fprintln(w, "No Cloud IAP header found.")
           return
        }
        email, userID, err := validateAssertion(assertion, certs, aud)
        if err != nil {
            log.Println(err)
            fmt.Fprintln(w, "Could not validate assertion. Check app logs.")
            return
        }

        fmt.Fprintf(w, "Hello %s with userID = %s\n", email, userID)
    }
}

func getSetSessionHandlerFunc(store *firestoregorilla.Store) func(http.ResponseWriter, *http.Request) {

    return func (w http.ResponseWriter, r *http.Request) {

        // Get a session. Get() always returns a session, even if empty.
        session, err := store.Get(r, sessionName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }

        // Set some session values.
        session.Values["foo"] = "bar"
        fooString := r.URL.Query().Get("foo")
        if fooString != "" {
            session.Values["foo"] = fooString
        }

        // Save it before we write to the response/return from the handler.
        err = session.Save(r, w)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "session.Values[\"foo\"] = %s<br />\n", session.Values["foo"])
    }
}

func getGetSessionHandlerFunc(store *firestoregorilla.Store) func(http.ResponseWriter, *http.Request) {

    return func (w http.ResponseWriter, r *http.Request) {
        // Get a session. Get() always returns a session, even if empty.
        session, err := store.Get(r, sessionName)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        fmt.Fprintf(w, "session.Values[\"foo\"] = %s<br />\n", session.Values["foo"])
    }
}

// validateAssertion validates assertion was signed by Google and returns the
// associated email and userID.
func validateAssertion(assertion string, certs map[string]string, aud string) (email string, userID string, err error) {
    token, err := jwt.Parse(assertion, func(token *jwt.Token) (interface{}, error) {
        keyID := token.Header["kid"].(string)

        _, ok := token.Method.(*jwt.SigningMethodECDSA)
        if !ok {
            return nil, fmt.Errorf("unexpected signing method: %q", token.Header["alg"])
        }

        cert := certs[keyID]
        return jwt.ParseECPublicKeyFromPEM([]byte(cert))
    })

    if err != nil {
        return "", "", err
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", "", fmt.Errorf("could not extract claims (%T): %+v", token.Claims, token.Claims)
    }

    if claims["aud"].(string) != aud {
        return "", "", fmt.Errorf("mismatched audience. aud field %q does not match %q", claims["aud"], aud)
    }
    return claims["email"].(string), claims["sub"].(string), nil
}

// audience returns the expected audience value for this service.
func audience() (string, error) {
    projectNumber, err := metadata.NumericProjectID()
    if err != nil {
        return "", fmt.Errorf("metadata.NumericProjectID: %v", err)
    }

    projectID, err := metadata.ProjectID()
    if err != nil {
        return "", fmt.Errorf("metadata.ProjectID: %v", err)
    }

    return "/projects/" + projectNumber + "/apps/" + projectID, nil
}

// certificates returns Cloud IAP's cryptographic public keys.
func certificates() (map[string]string, error) {
    const url = "https://www.gstatic.com/iap/verify/public_key"
    client := http.Client{
        Timeout: 5 * time.Second,
    }
    resp, err := client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("Get: %v", err)
    }

    var certs map[string]string
    dec := json.NewDecoder(resp.Body)
    if err := dec.Decode(&certs); err != nil {
        return nil, fmt.Errorf("Decode: %v", err)
    }

    return certs, nil
}
