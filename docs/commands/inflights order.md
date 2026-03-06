# inflights order

Order a specific service, optionally from a specific pilot.

## Usage

```bash
inflights order <serviceId> [pilotId] [options]
```

## Options

| Flag | Description |
|------|-------------|
| `--geo <geofence>` | Area of interest (if not already set on the service) |
| `--date <date>` | Preferred flight date |
| `--notes <text>` | Instructions for the pilot |

## Example

```bash
inflights order SV-01 pilot-xavier --geo site.geojson --date 2026-04-01
# → Order ORD-310 created. Flight FL-1060 scheduled.
```

## API

```
POST /v1/orders
Body: { "serviceId": "...", "pilotId": "...", "geofence": "...", ... }
```

## Roles

Customer.
