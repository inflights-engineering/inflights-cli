# Flight Lifecycle

A flight moves through these statuses:

```
scheduled → in-progress → images-uploaded → processing → completed → delivered
                                                 │
                                          (also: cancelled, on-hold)
```

## Who triggers each transition

| From → To | Who | Command |
|-----------|-----|---------|
| (new) → `scheduled` | System (on quote confirm or order) | automatic |
| `scheduled` → `in-progress` | Pilot | `inflights flight status <id> in-progress` |
| `in-progress` → `images-uploaded` | Pilot | `inflights upload images <id> <path>` (auto) |
| `images-uploaded` → `processing` | Processor | `inflights flight status <id> processing` |
| `processing` → `completed` | Processor | `inflights upload data <id> <path>` (auto) |
| `completed` → `delivered` | System / Admin | automatic or manual |
| any → `cancelled` | Customer / Admin | `inflights flight status <id> cancelled` |
| any → `on-hold` | Any linked role | `inflights flight status <id> on-hold` |

## Assigning people

Before a flight can progress, it needs a pilot and processor:

```bash
inflights flight assign-pilot FL-1055 pilot-xavier
inflights flight assign-processor FL-1055 proc-anna
```
