# inflights order

Order a new flight by providing a GeoJSON file defining the area.

## Usage

```bash
inflights order <geojson-file> [options]
```

## Options

| Flag | Description | Required |
|------|-------------|----------|
| `--service <id>` | Service ID (see `inflights services`) | yes |
| `--description <text>` | Description or notes for the flight | no |

## Area format

The CLI automatically normalizes common GeoJSON formats into the structure the API expects:

| Input format | Handled |
|-------------|---------|
| Feature with GeometryCollection of Polygons | yes (native) |
| Feature with a single Polygon | yes (auto-wrapped) |
| Feature with a MultiPolygon | yes (split into polygons) |
| FeatureCollection with polygon features | yes (merged) |
| Bare Polygon geometry | yes (auto-wrapped) |
| Bare MultiPolygon geometry | yes (split) |
| Bare GeometryCollection | yes (auto-wrapped) |

Coordinates are `[longitude, latitude]`. Each polygon must be closed (first and last point match).

## Example

```bash
inflights order area.geojson --service 3
# → Flight FL-1060 created.
# → ID:      42
# → UID:     FL-1060
# → Status:  needs_flight_proposal

inflights order area.geojson --service 3 --description "Roof inspection"
```

## API

```
POST /v1/flights
Body: { product_id, areas: <normalized GeoJSON>, description_user, skip_obtain: false }
```

## Roles

Customer.
