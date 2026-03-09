# inflights quote confirm

Confirm (accept) a quote, turning it into an active order.

## Usage

```bash
inflights quote confirm <quoteId>
```

## Example

```bash
inflights quote confirm QT-2091
# → QT-2091 confirmed. Proposal PR-401 sent to pilot.
```

Confirming a quote creates a proposal for an assigned pilot. Once the pilot accepts, a flight is created.

## API

```
POST /v1/quotes/<quoteId>/confirm
```

## Roles

Customer.
