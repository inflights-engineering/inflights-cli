# Flight Lifecycle

A flight moves through these statuses:

```
scheduled → in-progress → images-uploaded → processing → completed → delivered
                                                 │
                                          (also: cancelled, on-hold)
```

## Status transitions

All transitions are automatic unless noted otherwise. The system updates the status when the corresponding action is performed.

| From → To                         | Trigger                                          | Status update                                         |
| --------------------------------- | ------------------------------------------------ | ----------------------------------------------------- |
| (new) → `scheduled`               | Quote confirmed or order created                 | automatic                                             |
| `scheduled` → `in-progress`       | Pilot starts the mission                         | automatic                                             |
| `in-progress` → `images-uploaded` | Pilot uploads images (`inflights upload images`) | automatic                                             |
| `images-uploaded` → `processing`  | Processor begins work                            | automatic                                             |
| `processing` → `completed`        | Processor marks as done                          | **manual** — `inflights flight status <id> completed` |
| `completed` → `delivered`         | Deliverables sent to customer                    | automatic                                             |
| any → `cancelled`                 | Admin or customer cancels the flight             | **manual** — `inflights flight status <id> cancelled` |


## Assigning people

Before a flight can progress, it needs a pilot:

```bash
inflights flight assign-pilot FL-1055 pilot-xavier
```
