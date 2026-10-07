Build a simple CRM booking system in this (empty) git repository.

Constraints
- Go, standard library only (no third-party modules). Main package at ./cmd/crm.
- Listens on $PORT (default 8080). Persists all data to the JSON file at $DATA_FILE (default ./data.json); data must survive a restart.
- JSON field names are exactly as written below. Timestamps are RFC 3339.

API
- POST /api/customers  {"name","email","phone"} -> 201 with the customer incl. "id" (string).
  400 if name is empty or email is not a valid address; 409 if the email (case-insensitive) already exists.
- GET  /api/customers -> 200, JSON array.   GET /api/customers/{id} -> 200, or 404.
- POST /api/bookings {"customerId","start","end","note"} -> 201 with {"id","customerId","start","end","note","status":"confirmed"}.
  400 if end <= start, a time does not parse, or the customer does not exist.
  409 if it overlaps any other confirmed booking (there is one shared resource; touching end==start is not an overlap).
- GET  /api/bookings[?customerId=..&from=..&to=..] -> 200, JSON array sorted by start; from/to keep bookings that overlap [from,to).
- POST /api/bookings/{id}/cancel -> 200 with the booking, "status":"cancelled"; 404 if unknown. A cancelled booking blocks nothing.
- GET / -> a minimal HTML page to list/add customers and bookings and cancel a booking.

Include automated tests. `go build ./...`, `go vet ./...` and `go test ./...` must pass when you are done, and a README must say how to run it.

You are running unattended: no human will answer questions. Wherever you would normally ask, take your own recommended option, note it, and continue to the end.
