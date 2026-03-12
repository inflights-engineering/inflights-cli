# Data Pipeline

How drone imagery flows from capture to customer delivery.

```
Pilot                    Processor                 Customer
  │                         │                         │
  ├─ flies mission          │                         │
  ├─ upload images ────────►│                         │
  │   FL-1042 ./photos/     │                         │
  │                         ├─ download ──► raw imgs  │
  │                         │   FL-1042 -o ./raw/     │
  │                         │                         │
  │                         ├─ (process locally)      │
  │                         │                         │
  │                         ├─ upload data ──────────►│
  │                         │   FL-1042 output.tif    │
  │                         │   --deliverable 3       │
  │                         │                         │
  │                         │                    download
  │                         │                    FL-1042
```

## Commands involved

1. **Pilot** uploads raw imagery:
   `inflights upload images FL-1042 ./flight-photos/`

2. **Processor** downloads raw data:
   `inflights download FL-1042 -o ./raw/`

3. **Processor** uploads deliverables:
   `inflights upload data FL-1042 ./orthomosaic.tif --deliverable 3`

4. **Customer** downloads final deliverables:
   `inflights download FL-1042 -o ./deliverables/`
