# Ordering a Flight

End-to-end workflow from "I need drone data" to a confirmed flight mission.

## Steps

```
Customer                          Platform                    Pilot
   │                                 │                           │
   ├─ inflights order ─────────────►│                           │
   │   area.geojson --service 3     │  Flight created            │
   │                                 │                           │
   │                                 │  (inflights prices it)    │
   │                                 │                           │
   ├─ inflights quotes ────────────►│                           │
   │   (sees quote with price)       │                           │
   │                                 │                           │
   ├─ inflights quote confirm ──────►│                           │
   │   QT-2091                       │  → Proposal PR-401        │
   │                                 │    sent to pilot           │
   │                                 │                           │
   │                                 │  inflights proposal show ─┤
   │                                 │                           │
   │                                 │  inflights proposal       │
   │                                 │    accept PR-401 ─────────┤
   │                                 │                           │
   │                                 │  → Flight FL-1055 active  │
   │                                 │                           │
```

## Use case example

> "Hey, I need drone images of this warehouse roof."

```bash
inflights order warehouse.geojson --service 2 --description "Flat roof, 2000 sqm"
```

Inflights reviews and prices the quote. Once priced, the customer confirms it:

```bash
inflights quote confirm QT-2091
```

A proposal is then sent to an assigned pilot. The pilot reviews and accepts it:

```bash
inflights proposal accept PR-401
# → Proposal PR-401 accepted.
```

If the pilot rejects the proposal, inflights follows up to negotiate or reassign to another pilot.
