# inflights flights

List flights linked to the current user (as customer, pilot, or processor).

## Usage

```bash
inflights flights [options]
```

## Options

| Flag | Description |
|------|-------------|
| `--status <status>` | Filter by flight status |
| `--public-uid <uid>` | Filter by public UID |
| `--json` | Output as JSON |

### Valid statuses

| Status | Description |
|--------|-------------|
| `needs_flight_proposal` | Needs pilot proposals |
| `proposal_pending` | Proposals awaiting action |
| `price_not_final` | Awaiting quote from Inflights |
| `quote_sent` | Quote sent to client |
| `pilot_found` | Pilot scheduling flight |
| `flight_scheduled` | Flight date set |
| `flight_flown` | Flight completed, awaiting upload |
| `raw_data_uploaded` | Data uploaded, awaiting processing |
| `insights_generated` | Processing complete |
| `done` | Invoice created |

## Example

```bash
inflights flights --status flight_flown
```

```
ID    UID       STATUS        PRODUCT                  SCHEDULED
42    FL-1042   flight_flown  Roof Mapping with CAD    2026-03-10
38    FL-1038   flight_flown  Terrain Mapping          2026-03-12
```

## API

```
GET /v1/flights?status=...&public_uid=...
```

## Roles

All roles see flights they are linked to.
