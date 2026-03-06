# inflights flight add

Register an external flight (one made outside the platform) on your profile.

## Usage

```bash
inflights flight add [options]
```

## Options

| Flag | Description | Required |
|------|-------------|----------|
| `--date <date>` | Flight date (ISO 8601) | yes |
| `--geo <geofence>` | Geofence (WKT or file path to .geojson) | yes |
| `--insight <type>` | Insight type (e.g. orthomosaic, ndvi, 3d-model) | no |
| `--notes <text>` | Free-form notes | no |

## Example

```bash
inflights flight add --date 2026-03-15 --geo area.geojson --insight orthomosaic
# → Created flight FL-1055
```

## API

```
POST /v1/flights
```

## Roles

Pilot, Customer.


