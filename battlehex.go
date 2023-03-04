package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

	"cloud.google.com/go/compute/metadata"
	"github.com/golang-jwt/jwt"

    "battle-hex-go/frontend"
)

// app holds the Cloud IAP certificates and audience field for this app, which
// are needed to verify authentication headers set by Cloud IAP.
type app struct {
	certs map[string]string
	aud   string
}

func main() {
	a, err := newApp()
	if err != nil {
		log.Fatal(err)
	}

    fs := http.FileServer(http.Dir("./static"))

	http.HandleFunc("/", a.index)
    http.HandleFunc("/battlehex_vs_js_ai_v1.1", frontend.BattleHexJSHandler)
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


// newApp creates a new app, returning an error if either the Cloud IAP
// certificates or the app's audience field cannot be obtained.
func newApp() (*app, error) {
	certs, err := certificates()
	if err != nil {
		return nil, err
	}

	aud, err := audience()
	if err != nil {
		return nil, err
	}

	a := &app{
		certs: certs,
		aud:   aud,
	}
	return a, nil
}

// [END getting_started_auth_setup]

// [START getting_started_auth_front_controller]

// index responds to requests with our greeting.
func (a *app) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	assertion := r.Header.Get("X-Goog-IAP-JWT-Assertion")
	if assertion == "" {
		fmt.Fprintln(w, "No Cloud IAP header found.")
		return
	}
	email, _, err := validateAssertion(assertion, a.certs, a.aud)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, "Could not validate assertion. Check app logs.")
		return
	}

	fmt.Fprintf(w, "Hello %s\n", email)
}

// [END getting_started_auth_front_controller]

// [START getting_started_auth_validate]

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

// [END getting_started_auth_validate]

// [START getting_started_auth_audience]

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

// [END getting_started_auth_audience]

// [START getting_started_auth_certs]

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