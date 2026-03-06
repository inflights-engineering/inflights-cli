# inflights flight assign-processor

Assign a data processor to a flight.

## Usage

```bash
inflights flight assign-processor <flightId> <processorId>
```

## Example

```bash
inflights flight assign-processor FL-1042 proc-anna
# → Processor proc-anna assigned to FL-1042
```

## API

```
PATCH /v1/flights/<flightId>/assign
Body: { "processorId": "<processorId>" }
```

## Roles

Customer, Admin.
