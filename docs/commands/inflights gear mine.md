# inflights gear mine

List equipment you have added to your pilot profile.

## Usage

```bash
inflights gear mine [--format table|json]
```

## Example

```bash
inflights gear mine
```

```
ID       CATEGORY           NAME                          DRONE TYPE    SENSOR TYPES
EQ-001   drone              DJI Phantom 4 RTK             rotary_wing   rolling_shutter
EQ-012   payload            MicaSense RedEdge-MX          —             multispectral
EQ-045   gnss_receiver      Emlid Reach RS2               —             —
```

## API

```
GET /equipments
```

## Roles

Pilot.
