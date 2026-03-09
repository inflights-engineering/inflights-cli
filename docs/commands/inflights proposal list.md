# inflights proposal list

List flight proposals assigned to you as a pilot.

## Usage

```bash
inflights proposal list [options]
```

## Options

| Flag                | Description                                      |
| ------------------- | ------------------------------------------------ |
| `--status <status>` | Filter: `pending`, `accepted`, `rejected`        |
| `--format <fmt>`    | `table` (default), `json`                        |

## Example

```bash
inflights proposal list --status pending
```

```
ID       PRODUCT                  AREA                STATUS    DATE
PR-401   Roof Mapping with CAD    Brussels warehouse  pending   2026-03-15
PR-398   Terrain Mapping          Antwerp site C      pending   2026-03-18
```

## API

```
GET /proposals
```

## Roles

Pilot.
