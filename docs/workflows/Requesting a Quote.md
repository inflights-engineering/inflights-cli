# Requesting a Quote

End-to-end workflow from "I need drone data" to a confirmed flight mission.

## Steps

```
Customer                          Platform                    Pilot
   │                                 │                           │
   ├─ inflights quote request ──────►│                           │
   │   --geo site.geojson            │                           │
   │   --product "Roof Mapping"      │  QT-2091 created          │
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

> "Hey, I need drone images of this warehouse roof. Can you get a quote from @inflights?"

```bash
inflights quote request --geo warehouse.geojson --product "Roof Mapping with CAD" --notes "Flat roof, 2000 sqm, need by end of month"
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
