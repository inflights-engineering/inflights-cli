# Authentication

## `inflights login`

Opens a browser to `inflights.com/login/:token` where `:token` is a locally generated UUID. The CLI polls or listens for the token exchange to complete, then stores a bearer token in `~/.inflights/credentials`.

```bash
inflights login
# → Opening browser… https://inflights.com/login/a3f1c9e2-...
# → Authenticated as xavier@inflights.com (customer)
```

**Token storage:** `~/.inflights/credentials` (JSON, `0600` permissions).

All subsequent requests send `Authorization: Bearer <token>` in headers.

## `inflights logout`

Clears the local token file.

## `inflights whoami`

Prints the current authenticated user, email, and role(s).

```bash
inflights whoami
# → xavier@inflights.com  roles: customer, pilot
```

## API pattern

Every authenticated command makes HTTPS requests like:

```
curl -H "Authorization: Bearer ${token}" https://api.inflights.com/v1/...
```

See [[API Mapping]] for the full route table.
