# Image Upload Flow — CLI Implementation

The API provides a 3-step flow for uploading raw flight images (photos from the drone). This is separate from the document upload flow (processed deliverables).

## Endpoints

All nested under `POST /api/v1/flights/:flight_id/images/...`

1. **presign** — Get S3 presigned URL for one image
2. **confirm** — Register the uploaded image as a Picture record
3. **finalize** — Mark the flight as uploaded, triggers zip generation + notifications

## How the CLI should handle `inflights upload-images <flight_id> <folder>`

### 1. Discover files

- Scan the folder for image files (`.jpg`, `.jpeg`, `.png`, `.tif`, `.tiff`, `.dng`)
- Count total files and total size for progress display
- Optionally read EXIF from each file (GPS coords, altitude, camera model)

### 2. Upload loop (per file, with concurrency)

For each image file, run these two steps:

**a) Presign**
```
POST /flights/:flight_id/images/presign
Body: { filename: "IMG_0001.jpg" }
Response: { file_id, presign_data: { url, fields } }
```

**b) Upload to S3**
```
POST presign_data.url
Content-Type: multipart/form-data
Fields: all presign_data.fields + file=@<path>
```

**c) Confirm**
```
POST /flights/:flight_id/images/confirm
Body: { filename, file_id, size, exif: { latitude, longitude, altitude, ... } }
Response: { id, filename, size }
```

### 3. Finalize (once, after all images confirmed)

```
POST /flights/:flight_id/images/finalize
Response: { dataset_id, dataset_status, picture_count }
```

This triggers server-side zip generation and email notifications.

## Concurrency

- Use a worker pool (e.g., 4-8 goroutines) for parallel uploads
- Each worker runs the full presign → S3 upload → confirm cycle for one image
- Only call finalize after ALL workers complete successfully

## Progress display

- Show a progress bar or counter: `Uploading 12/48 images...`
- Show per-file status: presigning → uploading → confirming → done
- Print summary at end: `48 images uploaded. Dataset status: uploaded`

## Error handling

- **Presign fails**: Retry up to 3 times with backoff
- **S3 upload fails**: Retry up to 3 times (presign URL may expire, re-presign if needed)
- **Confirm fails**: Log error, continue with other files, report failures at end
- **Partial failure**: Don't call finalize if any images failed. Show which files failed and let user retry
- **Resume support** (nice-to-have): Track confirmed file_ids locally so re-running skips already-uploaded files

## EXIF extraction

Use a Go EXIF library (e.g., `rwcarlsen/goexif` or `dsoprea/go-exif`) to extract:
- `GPSLatitude` / `GPSLongitude` — decimal degrees
- `GPSAltitude` — meters
- `DateTimeOriginal` — capture timestamp
- `Make` / `Model` — camera info

Send as the `exif` param in the confirm request. The server stores this on the Picture record.
