# inflights upload data

Upload processed deliverables to a flight.

## Usage

```bash
inflights upload data <flight-id or public-uid> <files...> [options]
```

## Options

| Flag | Description |
|------|-------------|
| `-c, --concurrency <n>` | Number of parallel uploads (default: 5) |
| `--deliverable <id>` | Deliverable type ID (see `inflights flight <uid>` for available types) |
| `--json` | Output as JSON |

## Upload flow

Each file goes through a 3-step process:

1. **Presign** — `POST /flights/:id/uploads/presign` → get S3 upload URL
2. **Upload to S3** — multipart POST to presigned URL
3. **Confirm** — `POST /flights/:id/uploads/confirm` with file info

## Example

```bash
inflights upload data FL-1042 ./orthomosaic.tif --deliverable 3
# → Uploading 1 files...
# → [1/1] orthomosaic.tif — done
# → 1/1 files uploaded.
```

Upload multiple files:

```bash
inflights upload data FL-1042 report.pdf pointcloud.las -c 10
```

To see available deliverable type IDs for a flight:

```bash
inflights flight FL-1042
# → Shows deliverable types table with IDs
```

## API

| Step | Method | Endpoint |
|------|--------|----------|
| Presign | POST | `/flights/:id/uploads/presign` |
| Upload | POST | S3 presigned URL |
| Confirm | POST | `/flights/:id/uploads/confirm` |

## Roles

Processor.
