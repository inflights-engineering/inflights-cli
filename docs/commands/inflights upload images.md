# inflights upload images

Upload drone images for a given flight.

## Usage

```bash
inflights upload images <flightId> <path>
```

`<path>` can be a directory of images or a zip archive.

## Example

```bash
inflights upload images FL-1042 ./flight-photos/
# → Uploading 247 images (3.2 GB)…
# → ████████████████████ 100%
# → FL-1042 images uploaded. Status → images-uploaded.
```

Automatically updates flight status to `images-uploaded`.

## API

```
POST /v1/flights/<flightId>/images
Content-Type: multipart/form-data
```

Uses chunked / resumable upload for large payloads.

## Roles

Pilot.
