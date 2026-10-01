# JSON API conventions (AIP) — reference

The shapes behind the ten rules in `SKILL.md`, from the AIPs (<https://google.aip.dev>) and the proto3 JSON mapping (<https://protobuf.dev/programming-guides/json/>), seen 2026-09-30. Every example is a response a Go service here can emit with `encoding/json` and lowerCamelCase struct tags.

## Property naming (AIP-140)

- `lowerCamelCase` on the wire, the JSON form of the proto field name (proto spells it in lower snake form). A server accepts both forms when parsing; it emits camelCase only.
- Letters and digits only; no prepositions (`entryCount`, not `numberOfEntries`); no language keywords (`type`, `class`, `default`, `import` are fine in JSON but hurt generated clients: prefer `kind`, `category`).
- Booleans are adjectives or states without `is`: `enabled`, `deleted`, `public`.
- Quantities carry the unit as a suffix: `delayMs`, `sizeBytes`, `distanceKm`. Counts end in `Count`; sizes of a page or a collection are `pageSize`, `totalSize`.
- Arrays are plural; everything else singular.
- Map keys are data, not property names: a map's keys follow the data's own convention and the spec says so.

```json
{
  "name": "projects/demo/orders/123",
  "displayName": "Pizza",
  "enabled": true,
  "retryCount": 2,
  "timeoutSeconds": 30,
  "labels": {"env": "prod", "team": "kitchen"},
  "items": [{"sku": "A1", "quantity": 2}]
}
```

## Reserved and standard field names

| Field | Meaning | Behaviour |
|---|---|---|
| `name` | the resource's identifier path, `collection/id[/sub/id]` | output-only or identifier |
| `displayName` | the human label | optional |
| `createTime`, `updateTime`, `deleteTime`, `expireTime` | RFC 3339 timestamps | output-only (expireTime may be input) |
| `etag` | opaque version for optimistic concurrency | output-only; echoed back on update |
| `labels` | `map<string,string>` user metadata | optional |
| `state` | the lifecycle enum | output-only |
| `items` | the list response's elements | — |
| `nextPageToken`, `totalSize` | pagination | — |
| `error` | the only top-level wrapper, on failure | — |

A new field that would collide with one of these takes another name; the resource is never wrapped in `data` or versioned with `apiVersion` to avoid the collision.

## Values and types

- Strings, numbers, booleans, objects, arrays. `int32` and `double` are JSON numbers; **`int64` and `uint64` are JSON strings** (`"totalCents": "1250"`) because JavaScript clients lose precision above 2^53; `bytes` are standard base64 with padding.
- An unset optional field is omitted. `null` appears only where the spec names it as a state distinct from absent; a client treats `null` and absent the same unless told otherwise.
- Empty arrays are `[]`, never omitted, so a client can tell "none" from "not returned".
- Enums are strings: `UPPER_SNAKE_CASE`, first value `<ENUM>_UNSPECIFIED` which a server never emits for a known state; an unknown value on input is `INVALID_ARGUMENT`. A field whose set grows every month is a documented `string`, not an enum.

```json
{"state": "PAID", "paymentMethod": "CARD", "totalCents": "1250", "attachments": []}
```

## Time (AIP-142)

- Timestamps: RFC 3339, UTC, `Z`, 0, 3, 6 or 9 fractional digits — `"createTime": "2026-09-30T12:00:00Z"`, `"2026-09-30T12:00:00.123456Z"`. Field names end in `Time`. Go's `time.Time` marshals this way when in UTC; convert with `.UTC()` before encoding.
- Durations: `"timeout": "3.5s"` (seconds with up to nine fractional digits and the `s` suffix, the `google.protobuf.Duration` JSON form); the field name says what it measures (`timeout`, `ttl`), and a plain integer with a unit suffix (`timeoutSeconds`) is the alternative when a number is more useful to the client.
- Civil dates: `"dueDate": "2026-10-15"` (`YYYY-MM-DD`) or `{"year": 2026, "month": 10, "day": 15}` when the parts are addressed separately; field names end in `Date`.
- A day-of-week or month on its own is an enum (`MONDAY`, `OCTOBER`), never a number.

## Pagination (AIP-158)

Request (query parameters on a `GET` list):

| Parameter | Type | Rule |
|---|---|---|
| `pageSize` | int32 | 0 or absent → server default; above the maximum → coerced to the maximum; negative → `INVALID_ARGUMENT` |
| `pageToken` | string | opaque, from a previous response; a token issued for a different filter or order is `INVALID_ARGUMENT` |
| `filter`, `orderBy` | string | optional, their grammar is the API's own (AIP-160, 132) |

Response: the elements first, then the token.

```json
{
  "items": [{"name": "orders/1"}, {"name": "orders/2"}],
  "nextPageToken": "CgoxMjM0NTY3ODkw",
  "totalSize": 57
}
```

- `nextPageToken` is empty or absent at the end of the collection. That is the ONLY end-of-collection signal: a short page is not one, and a client keeps calling while a token comes back.
- `totalSize` is optional and may be an estimate; say which in the spec.
- Offset pagination (`startIndex`, `pageIndex`, `totalPages`, `nextLink`) is not used: it breaks under concurrent writes and leaks the storage shape.

## Errors (AIP-193)

The body of every non-2xx response:

```json
{
  "error": {
    "code": 400,
    "status": "INVALID_ARGUMENT",
    "message": "pageSize must be between 0 and 1000",
    "details": [
      {
        "@type": "type.googleapis.com/google.rpc.ErrorInfo",
        "reason": "PAGE_SIZE_OUT_OF_RANGE",
        "domain": "orders.example.com",
        "metadata": {"pageSize": "5000", "maximum": "1000"}
      },
      {
        "@type": "type.googleapis.com/google.rpc.BadRequest",
        "fieldViolations": [{"field": "pageSize", "description": "must be between 0 and 1000"}]
      }
    ]
  }
}
```

| Field | Meaning |
|---|---|
| `error.code` | the HTTP status code, as a number |
| `error.status` | the `google.rpc.Code` name: `INVALID_ARGUMENT` 400, `UNAUTHENTICATED` 401, `PERMISSION_DENIED` 403, `NOT_FOUND` 404, `ALREADY_EXISTS` 409, `ABORTED` 409, `FAILED_PRECONDITION` 400, `RESOURCE_EXHAUSTED` 429, `INTERNAL` 500, `UNAVAILABLE` 503, `DEADLINE_EXCEEDED` 504 |
| `error.message` | one developer-facing sentence: what was wrong and, when safe, how to fix it; never a stack trace, never a secret |
| `error.details[]` | typed objects, each type at most once: `ErrorInfo` (required: `reason` in `UPPER_SNAKE_CASE` ≤ 63 chars, `domain` = the service's name, `metadata` string map), `BadRequest` (`fieldViolations[]` with `field` as a path like `items[2].quantity`), `LocalizedMessage` (`locale`, `message` for the end user), `RetryInfo` (`retryDelay`), `Help` (`links[]`), `PreconditionFailure`, `QuotaFailure`, `ResourceInfo`, `RequestInfo` (`requestId`) |

A 5xx carries `status: "INTERNAL"` and a message that says nothing about the cause; the cause is in the server log with the request's trace (the `golang` skill's logging section). A 404 for a resource the caller may not even be allowed to see is a 403 or a 404 by the API's stated policy, consistently.

## Field behaviour (AIP-203)

Every field in a request body is annotated in the spec, and the handler enforces it:

| Behaviour | On input | Example |
|---|---|---|
| required | missing → `INVALID_ARGUMENT` with a `fieldViolations` entry | `displayName` on create |
| optional | may be absent; the default is documented | `labels` |
| output-only | ignored when sent (not an error) | `name`, `createTime`, `state`, `etag` |
| input-only | accepted, never returned | a `password` field |
| immutable | settable on create, `INVALID_ARGUMENT` on update | `name`'s id segment, a `region` |
| identifier | the resource's `name` | `name` |

## A list and a get, end to end

`GET /v1/orders?pageSize=2&pageToken=…` → the pagination response above. `GET /v1/orders/123` →

```json
{
  "name": "orders/123",
  "displayName": "Pizza",
  "state": "PAID",
  "totalCents": "1250",
  "items": [{"sku": "A1", "quantity": 2}],
  "createTime": "2026-09-30T12:00:00Z",
  "updateTime": "2026-09-30T12:05:30.250Z",
  "etag": "W/\"3\""
}
```

`PATCH /v1/orders/123` with `{"displayName": "Large pizza", "etag": "W/\"3\""}` returns the updated resource or `{"error": {"code": 409, "status": "ABORTED", …}}` when the etag no longer matches.
