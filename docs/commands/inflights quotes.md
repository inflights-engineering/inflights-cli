# inflights quotes

List quotes assigned to the current user.

## Usage

```bash
inflights quotes [options]
```

## Options

| Flag | Description |
|------|-------------|
| `--status <status>` | Filter: `pending`, `accepted` |
| `--json` | Output as JSON |

## Example

```bash
inflights quotes --status pending
```

```
NUMBER    FLIGHT     STATUS    TYPE       PRICE
QT-2091   FL-1042    pending   quote      500.00€
QT-2088   FL-1038    pending   estimate   —
```

## API

```
GET /v1/quotes?status=...
```

## Roles

All roles (scoped to quotes linked to the user).
