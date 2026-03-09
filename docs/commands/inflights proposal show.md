# inflights proposal show

Show details of a specific flight proposal.

## Usage

```bash
inflights proposal show <proposalId>
```

## Output

```
Proposal:  PR-401
Status:    pending
Product:   Roof Mapping with CAD
Area:      Brussels warehouse
Location:  50.85, 4.38
Date:      2026-03-15
Notes:     Flat roof, 2000 sqm
Created:   2026-03-08
```

## API

```
GET /proposals/:id
```

## Roles

Pilot.
