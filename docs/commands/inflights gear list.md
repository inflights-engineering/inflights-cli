# inflights gear list

List equipment registered to your profile.

## Usage

```bash
inflights gear list [--format table|json]
```

## Output

```
ID       TYPE     MODEL               SERIAL
GR-01    drone    DJI Matrice 350     SN-12345
GR-02    camera   MicaSense RedEdge   SN-67890
GR-03    drone    DJI Mavic 3E        SN-11223
```

## API

```
GET /v1/gear
```

## Roles

Pilot.
