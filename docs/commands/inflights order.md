# inflights order

Order a specific service.

## Usage

```bash
inflights order <serviceId> [options]
```

## Options

| Flag                   | Description                                 |
| ---------------------- | ------------------------------------------- |
| `--geo <file>`         | GeoJSON file defining the area of interest  |
| `--location <lat,lng>` | Location for the flight (e.g. `50.85,4.38`) |
| `--date <date>`        | Preferred flight date                       |
| `--notes <text>`       | Additional instructions                     |

## Area format

The `--geo` flag accepts a GeoJSON file containing a Feature with a GeometryCollection of Polygons:

```json
{
  "type": "Feature",
  "geometry": {
    "type": "GeometryCollection",
    "geometries": [
      {
        "type": "Polygon",
        "coordinates": [
          [
            [4.3517, 50.8503],
            [4.3527, 50.8503],
            [4.3527, 50.8513],
            [4.3517, 50.8513],
            [4.3517, 50.8503]
          ]
        ]
      }
    ]
  },
  "properties": {}
}
```

Coordinates are `[longitude, latitude]`. Each polygon must be closed (first and last point must match) with at least 4 points.

## Example

```bash
inflights order SV-01 --geo site.geojson --location 50.85,4.38 --date 2026-04-01
# → Order ORD-310 created. Flight FL-1060 scheduled.
```

## API

```
POST /flights
Body:
{
  "product_id": 0,
  "areas": {
    "type": "Feature",
    "properties": {},
    "geometry": {
      "type": "GeometryCollection",
      "geometries": [
        {
          "type": "Polygon",
          "coordinates": [[[lng, lat], ...]]
        }
      ]
    }
  },
  "location": {
    "latitude": 50.85,
    "longitude": 4.38
  }
}
```

## Roles

Customer.
