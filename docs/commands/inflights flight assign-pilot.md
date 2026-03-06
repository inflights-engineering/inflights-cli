# inflights flight assign-pilot

Assign a pilot to a flight.

## Usage

```bash
inflights flight assign-pilot <flightId> <pilotId>
```

## Example

```bash
inflights flight assign-pilot FL-1042 pilot-xavier
# → Pilot pilot-xavier assigned to FL-1042
```

## API

```
PATCH /v1/flights/<flightId>/assign
Body: { "pilotId": "<pilotId>" }
```

## Roles

Customer, Admin.
