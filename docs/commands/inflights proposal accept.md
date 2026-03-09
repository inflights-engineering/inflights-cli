# inflights proposal accept

Accept a flight proposal. Once accepted, the proposal becomes an active flight mission.

## Usage

```bash
inflights proposal accept <proposalId>
```

## Example

```bash
inflights proposal accept PR-401
# → Proposal PR-401 accepted. Flight FL-1070 created.
```

## API

```
POST /proposals/:id/accept
```

## Roles

Pilot.
