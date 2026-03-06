# inflights flight status

Update the status of a flight.

## Usage

```bash
inflights flight status <flightId> <newStatus>
```

## Valid statuses

`scheduled` → `in-progress` → `images-uploaded` → `processing` → `completed` → `delivered`

Also: `cancelled`, `on-hold`.

## Example

```bash
inflights flight status FL-1042 in-progress
# → FL-1042 status updated to in-progress
```

## API

```
PATCH /v1/flights/<flightId>
Body: { "status": "<newStatus>" }
```

## Roles

Depends on the transition. See [[Flight Lifecycle]].
