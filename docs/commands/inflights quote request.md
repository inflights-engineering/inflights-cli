# inflights quote request

Request a new quote for drone services over a given area.

## Usage

```bash
inflights quote request [options]
```

## Options

| Flag                | Description                                       | Required |
| ------------------- | ------------------------------------------------- | -------- |
| `--geo <geofence>`  | Area of interest (WKT string or path to .geojson) | yes      |
| `--insight <type>`  | Desired insight (Roof mapping with cad...)        | yes      |
| `--notes <text>`    | Additional context for the quote                  | no       |

## Example

```bash
inflights quote request --geo site.geojson --insight roof-mapping --notes "50-acre vineyard"
# → Quote QT-2091 created. Pending review.
```

## API

```
POST /v1/quotes
Body: { "geofence": "...", "insight": "...", "notes": "..." }
```

## Roles

Customer.
