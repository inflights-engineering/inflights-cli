# inflights upload data

Upload processed data (deliverables) for a given flight.

## Usage

```bash
inflights upload data <flightId> <path>
```

`<path>` can be a directory or zip of processed outputs (orthomosaics, point clouds, reports, etc.).

## Example

```bash
inflights upload data FL-1042 ./processed-output/
# → Uploading processed data (1.8 GB)…
# → ████████████████████ 100%
# → FL-1042 processed data uploaded. Status → completed.
```

## API

```
POST /v1/flights/<flightId>/data
Content-Type: multipart/form-data
```

## Roles

Processor.
