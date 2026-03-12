# inflights quote reject

Reject a quote or estimate.

## Usage

```bash
inflights quote reject <quote-number or flight-uid>
```

The command auto-detects whether the identifier matches a quote number or a flight public UID (for estimates), and calls the appropriate endpoint.

## Example

```bash
# Reject a quote by quote number
inflights quote reject QT-2091
# → Quote QT-2091 rejected.

# Reject an estimate by flight UID
inflights quote reject FL-1042
# → Estimate for FL-1042 rejected.
```

## API

- Quote: `POST /v1/quotes/:id/reject`
- Estimate: `POST /v1/quotes/reject_estimate` with `{ flight_public_uid: "..." }`

## Roles

Customer.
