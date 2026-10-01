# Go samples — API auth on Cloud Run

Every block is a complete file that compiles (`go vet`), each in its own package. The
libraries: `google.golang.org/api/idtoken` (IAP assertions, service-to-service tokens, the
scheduler route) and `firebase.google.com/go/v4` (public apps). Nothing here parses a token
by hand.

## IAP middleware

```go
// Package iapauth verifies the Identity-Aware Proxy assertion on every request.
// Audience: /projects/PROJECT_NUMBER/locations/REGION/services/SERVICE (the project NUMBER).
package iapauth

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/api/idtoken"
)

// Issuer is the one issuer IAP uses; idtoken.Validate does not check it, so this package does.
const Issuer = "https://cloud.google.com/iap"

type userKey struct{}

// User is who IAP authenticated: the stable subject and the email the team recognises.
type User struct {
	Subject string
	Email   string
}

// New returns the middleware. It refuses an empty audience at construction time:
// idtoken.Validate skips the aud check when the audience is "", which would accept an
// assertion minted for any other IAP resource.
func New(audience string) (func(http.Handler) http.Handler, error) {
	if audience == "" {
		return nil, errors.New("iapauth: IAP_AUDIENCE is empty; refusing to start")
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertion := r.Header.Get("x-goog-iap-jwt-assertion")
			if assertion == "" { // the request did not come through IAP
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			payload, err := idtoken.Validate(r.Context(), assertion, audience)
			if err != nil || payload.Issuer != Issuer {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			email, _ := payload.Claims["email"].(string)
			u := User{Subject: payload.Subject, Email: email}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey{}, u)))
		})
	}, nil
}

// FromContext returns the authenticated user, or false before the middleware ran.
func FromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userKey{}).(User)
	return u, ok
}
```

## Scheduler route on a public service

```go
// Package tick serves a Cloud Scheduler route on a PUBLIC Cloud Run service, where Cloud Run
// no longer verifies the OIDC token: the handler checks audience, issuer, caller email and
// email_verified itself. Any service account anywhere can mint a token for a public URL,
// so the audience alone proves nothing about the caller.
package tick

import (
	"errors"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

// New returns the /tick handler, refusing an empty audience or scheduler SA at construction
// time so a misconfigured service fails to start instead of answering 500 per request.
// audience is the service URL (what the Scheduler job was created with); schedulerSA is the
// job's service-account email.
func New(audience, schedulerSA string, run func() error) (http.Handler, error) {
	if audience == "" || schedulerSA == "" {
		return nil, errors.New("tick: audience and scheduler service account must be set")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		payload, err := idtoken.Validate(r.Context(), token, audience)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		// Validate never compares iss: a Google-signed ID token's issuer is accounts.google.com.
		if payload.Issuer != "https://accounts.google.com" && payload.Issuer != "accounts.google.com" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		email, _ := payload.Claims["email"].(string)
		verified, _ := payload.Claims["email_verified"].(bool)
		if email != schedulerSA || !verified {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if err := run(); err != nil {
			http.Error(w, "tick failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}), nil
}
```

## Service to service

```go
// Package callee shows a Cloud Run service calling another IAM-gated service: the client
// mints an OIDC token for the callee's URL from the instance's own service account (no key
// file); the caller's SA holds roles/run.invoker on the callee.
package callee

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/api/idtoken"
)

// Client wraps an *http.Client that attaches a fresh identity token for audience.
type Client struct {
	http *http.Client
	base string
}

// New builds the client once at startup; audience is the callee's URL.
func New(ctx context.Context, audience string) (*Client, error) {
	if audience == "" {
		return nil, fmt.Errorf("callee: empty audience")
	}
	c, err := idtoken.NewClient(ctx, audience)
	if err != nil {
		return nil, fmt.Errorf("callee: id token client: %w", err)
	}
	return &Client{http: c, base: audience}, nil
}

// Health calls GET /healthz on the callee and returns its status code.
func (c *Client) Health(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.base+"/healthz", nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("callee: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}
```

## Firebase

```go
// Package fbauth verifies Firebase Auth ID tokens for a public app. The SPA sends
// Authorization: Bearer <ID token>; the Admin SDK checks signature, aud, iss and exp.
package fbauth

import (
	"context"
	"net/http"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type uidKey struct{}

// New initialises the Admin SDK with Application Default Credentials (the Cloud Run
// service account; locally `gcloud auth application-default login`).
func New(ctx context.Context) (*auth.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	return app.Auth(ctx)
}

// Middleware rejects a request without a valid ID token and stores the UID in the context.
// checkRevoked adds a network call per request: true on sensitive routes only.
func Middleware(client *auth.Client, checkRevoked bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if raw == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			var tok *auth.Token
			var err error
			if checkRevoked {
				tok, err = client.VerifyIDTokenAndCheckRevoked(r.Context(), raw)
			} else {
				tok, err = client.VerifyIDToken(r.Context(), raw)
			}
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), uidKey{}, tok.UID)))
		})
	}
}

// UID returns the verified user id, or "" before the middleware ran.
func UID(ctx context.Context) string {
	uid, _ := ctx.Value(uidKey{}).(string)
	return uid
}
```
