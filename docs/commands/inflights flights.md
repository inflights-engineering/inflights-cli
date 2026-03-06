# inflights flights

List flights linked to the current user (as customer, pilot, or processor).

## Usage

```bash
inflights flights [options]
```

## Options

| Flag | Description | Example |
|------|-------------|---------|
| `--geo <geofence>` | Filter by geofence (WKT or bounding box) | `--geo "POLYGON((-80 25,-80 26,-79 26,-79 25,-80 25))"` |
| `--daterange <range>` | Filter by date range | `--daterange 2026-01-01:2026-03-01` |
| `--status <status>` | Filter by flight status | `--status scheduled` |
| `--role <role>` | Filter by your role on the flight | `--role pilot` |
| `--format <fmt>` | Output format: `table` (default), `json`, `csv` | `--format json` |

## Example

```bash
inflights flights --status scheduled --format table
```

```
ID       DATE        STATUS      CUSTOMER         PILOT
FL-1042  2026-03-10  scheduled   acme-corp        xavier
FL-1038  2026-03-12  scheduled   greenfield-inc   xavier
```

## API

```
GET /v1/flights?geo=...&from=...&to=...&status=...
Authorization: Bearer <token>
```

## Roles

All roles see flights they are linked to.
