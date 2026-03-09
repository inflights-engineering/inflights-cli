# inflights gear add

Select equipment from the inflights predefined list and add it to your pilot profile.

## Usage

```bash
inflights gear add <equipmentTypeId>
```

## Example

```bash
inflights gear add EQ-001
# → Equipment "DJI Phantom 4 RTK" added to your profile.
```

## API

```
POST /equipments
Body: { "equipment_type_id": 1 }
```

## Roles

Pilot.
