# inflights quote assign

Assign a quote to a user for review (e.g. a pilot to price it out).

## Usage

```bash
inflights quote assign <quoteId> <userId>
```

## Example

```bash
inflights quote assign QT-2091 pilot-xavier
# → QT-2091 assigned to pilot-xavier for review
```

## API

```
PATCH /v1/quotes/<quoteId>/assign
Body: { "reviewerId": "<userId>" }
```

## Roles

Admin, Customer.
