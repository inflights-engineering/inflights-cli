# inflights quotes

List quotes assigned to the current user.

## Usage

```bash
inflights quotes [options]
```

## Options

| Flag | Description |
|------|-------------|
| `--status <status>` | Filter: `pending`, `reviewed`, `confirmed`, `expired` |
| `--format <fmt>` | `table` (default), `json`, `csv` |

## Example

```bash
inflights quotes --status pending
```

```
ID       AREA              INSIGHT       STATUS    CUSTOMER
QT-2091  50-acre vineyard  orthomosaic   pending   acme-corp
QT-2088  solar farm B      thermal       pending   greenfield-inc
```

## API

```
GET /v1/quotes?assignee=me&status=...
```

## Roles

All roles (scoped to quotes linked to the user).
