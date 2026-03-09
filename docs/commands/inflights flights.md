# inflights flights

List flights linked to the current user (as customer, pilot, or processor).

## Usage

```bash
inflights flights [options]
```

## Options

| Flag                      | Description                | Example                 |
| ------------------------- | -------------------------- | ----------------------- |
| `--location <location>`   | Filter by location         | `--location brussels`   |
| `--status <status>`       | Filter by flight status    | `--status scheduled`    |
| `--number <number>`       | Search by flight number    | `--number FL234`        |
| `--reference <reference>` | Search by flight reference | `--reference Basic-fit` |


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
