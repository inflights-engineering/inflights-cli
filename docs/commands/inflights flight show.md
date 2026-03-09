# inflights flight show

Show details of a specific flight.

## Usage

```bash
inflights flight show <flight publicUid>
```

## Output

```
Flight:      FL-1042
Status:      scheduled
Date:        2026-03-10
Geofence:    POLYGON((-80.2 25.7, ...))
Customer:    acme-corp
Pilot:       xavier
Processor:   —
Insight:     orthomosaic
Images:      not uploaded
Processed:   —
```

## API

```
GET /v1/flights/<flightId>
```

## Roles

Any role linked to the flight.
