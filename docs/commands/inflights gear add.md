# inflights gear add

Register a new piece of equipment (drone, camera, sensor) to your profile.

## Usage

```bash
inflights gear add [options]
```

## Options

| Flag | Description | Required |
|------|-------------|----------|
| `--type <type>` | `drone`, `camera`, `sensor`, `other` | yes |
| `--model <model>` | Model name | yes |
| `--serial <serial>` | Serial number | no |
| `--notes <text>` | Additional details | no |

## Example

```bash
inflights gear add --type drone --model "DJI Matrice 350" --serial SN-12345
# → Gear GR-04 registered.
```

## API

```
POST /v1/gear
```

## Roles

Pilot.
