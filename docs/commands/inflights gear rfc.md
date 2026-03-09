# inflights gear rfc

Request a new piece of equipment to be added to the inflights predefined list.

## Usage

```bash
inflights gear rfc [options]
```

## Options

| Flag              | Description                                                        | Required |
| ----------------- | ------------------------------------------------------------------ | -------- |
| `--name <name>`   | Equipment name/model                                               | yes      |
| `--type <type>`   | Category: `drone`, `payload`, `drone_and_payload`, `gnss_receiver` | yes      |
| `--notes <text>`  | Additional details (e.g. specs, links)                             | no       |

## Example

```bash
inflights gear rfc --name "DJI Matrice 4T" --type drone_and_payload --notes "Released Q1 2026, thermal + visual"
# → Request submitted. You will be notified when reviewed.
```

## API

```
POST /equipment_type_requests
Body: { "name": "...", "category": "...", "notes": "..." }
```

## Roles

Pilot.
