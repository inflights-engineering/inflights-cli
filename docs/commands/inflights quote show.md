# inflights quote show

Show details of a specific quote.

## Usage

```bash
inflights quote show <quoteId>
```

## Output

```
Quote:     QT-2091
Status:    pending
Geofence:  POLYGON((-80.2 25.7, ...))
Insight:   orthomosaic
Customer:  acme-corp
Reviewer:  pilot-xavier
Price:     —
Notes:     50-acre vineyard
Created:   2026-03-01
```

## API

```
GET /v1/quotes/<quoteId>
```
