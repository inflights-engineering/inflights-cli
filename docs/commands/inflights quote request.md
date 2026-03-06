# inflights quote request

Request a new quote for drone services over a given area.

## Usage

```bash
inflights quote request [options]
```

## Options

| Flag | Description | Required |
|------|-------------|----------|
| `--geo <geofence>` | Area of interest (WKT string or path to .geojson) | yes |
| `--insight <type>` | Desired insight (orthomosaic, ndvi, 3d-model, thermal, etc.) | yes |
| `--notes <text>` | Additional context for the quote | no |
| `--urgency <level>` | `normal` (default), `rush` | no |

## Example

```bash
inflights quote request --geo site.geojson --insight orthomosaic --notes "50-acre vineyard"
# → Quote QT-2091 created. Pending review.
```

## API

```
POST /v1/quotes
Body: { "geofence": "...", "insight": "...", "notes": "...", "urgency": "..." }
```

## Roles

Customer.
