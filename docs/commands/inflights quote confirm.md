# inflights quote confirm

Confirm (accept) a quote, turning it into an active order.

## Usage

```bash
inflights quote confirm <quoteId>
```

## Example

```bash
inflights quote confirm QT-2091
# → QT-2091 confirmed. Flight FL-1055 created.
```

Confirming a quote should automatically create a flight record.

## API

```
POST /v1/quotes/<quoteId>/confirm
```

## Roles

Customer.
