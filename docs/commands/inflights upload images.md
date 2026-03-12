# inflights upload images

Upload drone images to a flight.

## Usage

```bash
inflights upload images <flight-id or public-uid> <path...> [options]
```

Accepts individual files or directories. When a directory is given, it is scanned for image files (`.jpg`, `.jpeg`, `.png`, `.tif`, `.tiff`, `.dng`).

## Options

| Flag | Description |
|------|-------------|
| `-c, --concurrency <n>` | Number of parallel uploads (default: 5) |
| `--json` | Output as JSON |

## Upload flow

Each image goes through a 3-step process:

1. **Presign** — `POST /flights/:id/images/presign` → get S3 upload URL
2. **Upload to S3** — multipart POST to presigned URL
3. **Confirm** — `POST /flights/:id/images/confirm` with file info and EXIF metadata

After all images are uploaded successfully, the dataset is finalized via `POST /flights/:id/images/finalize`.

If any uploads fail, finalize is skipped. Fix the issues and re-run.

## EXIF metadata

Automatically extracted from each image and sent with the confirm request:
- GPS latitude / longitude
- GPS altitude
- Capture date/time
- Camera make / model

## Example

```bash
inflights upload images FL-1042 ./flight-photos/
# → Found 247 images. Uploading with 5 workers...
# → [1/247] IMG_0001.jpg — done
# → [2/247] IMG_0002.jpg — done
# → ...
# → 247/247 images uploaded.
# → Dataset finalized (247 pictures).
```

Upload with higher concurrency:

```bash
inflights upload images FL-1042 ./photos/ -c 10
```

## API

| Step | Method | Endpoint |
|------|--------|----------|
| Presign | POST | `/flights/:id/images/presign` |
| Upload | POST | S3 presigned URL |
| Confirm | POST | `/flights/:id/images/confirm` |
| Finalize | POST | `/flights/:id/images/finalize` |

## Roles

Pilot.
