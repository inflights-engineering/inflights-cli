# Ordering a Flight

End-to-end workflow from "I need drone data" to a confirmed flight mission.

## Steps

```
Customer                          Platform                    Pilot
   │                                 │                           │
   ├─ inflights order ─────────────►│                           │
   │   SV-01 --geo site.geojson     │                           │
   │   --location 50.85,4.38        │  Quote QT-2091 created    │
   │                                 │                           │
   │                                 │  (inflights prices it)    │
   │                                 │                           │
   ├─ inflights quotes ────────────►│                           │
   │   (sees QT-2091 with price)     │                           │
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
   │                                 │  → Flight FL-1055 created │
   │                                 │                           │
```

## Use case example

> "Hey, I need drone images of this warehouse roof."

```bash
inflights order SV-02 --geo warehouse.geojson --location 50.85,4.38 --notes "Flat roof, 2000 sqm"
```

Inflights reviews and prices the quote. Once priced, the customer confirms it:

```bash
inflights quote confirm QT-2091
```

A proposal is then sent to an assigned pilot. The pilot reviews and accepts it:

```bash
inflights proposal accept PR-401
# → Flight FL-1055 created.
```

If the pilot rejects the proposal, inflights follows up to negotiate or reassign to another pilot.
