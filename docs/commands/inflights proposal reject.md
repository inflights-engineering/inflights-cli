# inflights proposal reject

Reject a flight proposal with a reason. Inflights may come back to negotiate.

## Usage

```bash
inflights proposal reject <proposalId> --reason <text>
```

## Options

| Flag               | Description          | Required |
| ------------------ | -------------------- | -------- |
| `--reason <text>`  | Reason for rejection | yes      |

## Example

```bash
inflights proposal reject PR-401 --reason "Schedule conflict, unavailable on that date"
# → Proposal PR-401 rejected. Inflights will follow up.
```

## API

```
POST /proposals/:id/reject
Body: { "reason": "..." }
```

## Roles

Pilot.
