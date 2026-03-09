# inflights gear list

List all available equipment on the inflights platform.

## Usage

```bash
inflights gear list [--type <type>] [--format table|json]
```

## Options

| Flag              | Description                                              |
| ----------------- | -------------------------------------------------------- |
| `--type <type>`   | Filter by category: `drone`, `payload`, `drone_and_payload`, `gnss_receiver` |
| `--format <fmt>`  | Output format: `table` (default) or `json`               |

## Example

```bash
inflights gear list --type drone
```

```
ID       CATEGORY   NAME                          DRONE TYPE    SENSOR TYPES
EQ-001   drone      DJI Phantom 4 RTK             rotary_wing   rolling_shutter
EQ-002   drone      DJI M300 + Zenmuse P1         rotary_wing   global_shutter
EQ-003   drone      Wingtra One GEN II            fixed_wing    global_shutter
EQ-004   drone      senseFly eBee X               fixed_wing    global_shutter
EQ-005   drone      Autel EVO II Pro              rotary_wing   rolling_shutter
```

## API

```
GET /equipment_types
```

## Roles

Pilot.
