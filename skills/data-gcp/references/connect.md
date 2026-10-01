# Connecting a Go service to Cloud SQL Postgres

Two complete files that compile (`go vet`). The first is the default path (the `/cloudsql`
socket, a password from Secret Manager); the second is the connector with IAM database
authentication (no password; needs Go 1.26 or later and `cloud.google.com/go/cloudsqlconn`).
Both return a `*pgxpool.Pool` the caller closes; neither logs the DSN.

## Socket (default)

```go
// Package db opens a pgx pool over the Cloud SQL Unix socket Cloud Run mounts with
// --add-cloudsql-instances. The DSN comes from configuration (Secret Manager through
// --set-secrets), of the form
//   host=/cloudsql/PROJECT:REGION:INSTANCE user=app dbname=app password=…
// and locally postgres://app:app@localhost:5432/app?sslmode=disable (the compose database).
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Open connects and pings. The DSN is never logged: it carries the password.
func Open(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("db: parse dsn: %w", err) // pgx's message does not echo the password
	}
	cfg.MaxConns = 4 // Cloud Run instances are small; concurrency × instances must stay under the tier's max_connections
	cfg.MaxConnLifetime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("db: connect: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}
	return pool, nil
}
```

## Connector with IAM authentication

```go
// Package dbiam opens a pgx pool through the Cloud SQL Go connector with IAM database
// authentication: the Cloud Run service account is the database user, there is no password,
// and the connector handles the TLS to the instance. The DSN names no host:
//   user=APP-run@PROJECT.iam dbname=app
// (the service-account email without .gserviceaccount.com). Go 1.26 or later.
package dbiam

import (
	"context"
	"fmt"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool holds the pool and the dialer that must be closed with it.
type Pool struct {
	*pgxpool.Pool
	dialer *cloudsqlconn.Dialer
}

// Open dials PROJECT:REGION:INSTANCE for every connection in the pool. WithLazyRefresh
// refreshes certificates on demand, the mode for scale-to-zero; WithIAMAuthN uses the
// instance identity as the Postgres login.
func Open(ctx context.Context, dsn, instance string) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("dbiam: parse dsn: %w", err)
	}
	d, err := cloudsqlconn.NewDialer(ctx, cloudsqlconn.WithIAMAuthN(), cloudsqlconn.WithLazyRefresh())
	if err != nil {
		return nil, fmt.Errorf("dbiam: dialer: %w", err)
	}
	cfg.ConnConfig.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
		return d.Dial(ctx, instance)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		_ = d.Close()
		return nil, fmt.Errorf("dbiam: connect: %w", err)
	}
	return &Pool{Pool: pool, dialer: d}, nil
}

// Close closes the pool, then the dialer.
func (p *Pool) Close() error {
	p.Pool.Close()
	return p.dialer.Close()
}
```

## Which one

| Project wants | Use |
|---|---|
| The simplest working thing; a password in Secret Manager is acceptable | Socket |
| No database password anywhere; IAM is the only credential | Connector + `WithIAMAuthN` |
| A private-IP instance reached over Direct VPC egress | Connector (the socket needs the instance to accept the Cloud Run connection path) |
| Local development | the compose database (`docker` skill); the Auth Proxy for production data only |
