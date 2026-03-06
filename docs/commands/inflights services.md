# inflights services

List available services (optionally filtered by pilot).

## Usage

```bash
inflights services [pilotId]
```

## Example

```bash
inflights services
```

```
ID       NAME                  PILOT           PRICE FROM
SV-01    Orthomosaic Survey    pilot-xavier    $500
SV-02    Thermal Inspection    pilot-xavier    $750
SV-03    3D Mapping            pilot-anna      $900
```

```bash
inflights services pilot-xavier
# → shows only pilot-xavier's services
```

## API

```
GET /v1/services?pilotId=...
```

## Roles

All roles.
