# inflights flight show

Show details of a specific flight.

## Usage

```bash
inflights flight <flight-id or public-uid>
```

Accepts either a numeric ID or a public UID (e.g. `FL-1042`).

## Output

```
ID:          42
UID:         FL-1042
Status:      flight_flown
Product:     Roof Mapping with CAD
Scheduled:   2026-03-10
Area (ha):   5.2
Price:       750.00€
Description: Flat roof inspection
Reference:   REF-123
Pilot:       Xavier Dupont
Customer:    Acme Corp

Deliverable ID  Name
3               Orthomosaic - 2D
7               CAD Model Roof - 3D
```

The deliverables table shows available deliverable type IDs, which can be used with `inflights upload data --deliverable <id>`.

## Options

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |

## API

```
GET /v1/flights/:id
```

## Roles

Any role linked to the flight.
