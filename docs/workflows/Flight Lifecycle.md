# Flight Lifecycle

A flight moves through these statuses:

```
needs_flight_proposal → proposal_pending → price_not_final → quote_sent →
pilot_found → flight_scheduled → flight_flown → raw_data_uploaded →
insights_generated → done
```

## Status transitions

All transitions are automatic. The system updates the status when the corresponding action is performed.

| From → To | Trigger |
|-----------|---------|
| (new) → `needs_flight_proposal` | Flight ordered |
| → `proposal_pending` | Proposal sent to pilot |
| → `pilot_found` | Pilot accepts proposal |
| → `flight_scheduled` | Flight date set |
| → `flight_flown` | Pilot completes mission |
| → `raw_data_uploaded` | Pilot uploads images (`inflights upload images`) |
| → `insights_generated` | Processing complete |
| → `done` | Invoice created |

## Assigning people

Pilots are assigned through the proposal system. A proposal is sent to a pilot, who accepts or rejects it:

```bash
inflights proposal accept PR-401
# → Proposal accepted. Flight FL-1055 assigned.
```
