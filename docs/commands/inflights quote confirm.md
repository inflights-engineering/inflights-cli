# inflights quote confirm

Confirm (accept) a quote or estimate.

## Usage

```bash
inflights quote confirm <quote-number or flight-uid>
```

The command auto-detects whether the identifier matches a quote number or a flight public UID (for estimates), and calls the appropriate endpoint.

## Example

```bash
# Confirm a quote by quote number
inflights quote confirm QT-2091
# → Quote QT-2091 confirmed.

# Confirm an estimate by flight UID
inflights quote confirm FL-1042
# → Estimate for FL-1042 accepted.
```

Confirming a quote creates a proposal for a pilot. Once the pilot accepts, a flight is created.

## API

- Quote: `POST /v1/quotes/:id/accept`
- Estimate: `POST /v1/quotes/accept_estimate` with `{ flight_public_uid: "..." }`

## Roles

Customer.
