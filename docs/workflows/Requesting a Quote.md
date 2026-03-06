# Requesting a Quote

End-to-end workflow from "I need drone data" to a confirmed order.

## Steps

```
Customer                          Platform                    Reviewer
   │                                 │                           │
   ├─ inflights quote request ──────►│                           │
   │   --geo site.geojson            │                           │
   │   --insight orthomosaic         │  QT-2091 created          │
   │                                 │                           │
   │                                 │◄── inflights quote assign │
   │                                 │    QT-2091 pilot-xavier   │
   │                                 │                           │
   │                                 │    (reviewer prices it)   │
   │                                 │                           │
   ├─ inflights quotes ────────────►│                           │
   │   (sees QT-2091 with price)     │                           │
   │                                 │                           │
   ├─ inflights quote confirm ──────►│                           │
   │   QT-2091                       │  → Flight FL-1055 created │
   │                                 │                           │
```

## Use case example

> "Hey, I need drone images of this vineyard. Can you get a quote from @inflights?"

```bash
inflights quote request --geo vineyard.geojson --insight orthomosaic --notes "50 acres, need by end of month"
```

The quote gets assigned to a pilot for review. Once priced, the customer confirms it, and a flight is automatically scheduled.
