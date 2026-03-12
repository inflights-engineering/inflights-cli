# inflights download

Download data for a given flight.

## Usage

```bash
inflights download <flight-id or public-uid> [options]
```

## Options

| Flag | Description |
|------|-------------|
| `-o, --output <dir>` | Destination directory (default: current directory) |
| `--json` | Output as JSON (list of files without downloading) |

## Example

```bash
inflights download FL-1042 -o ./deliverables/
# → Downloading orthomosaic.tif (450 MB)...
# → Downloading report.pdf (2 MB)...
# → Downloaded 2 files to ./deliverables/
```

## What gets downloaded

The command fetches all available downloads for the flight, which may include:
- Processed documents (deliverables)
- Picture set ZIP (drone imagery archive)

## API

```
GET /v1/flights/:id/downloads
```

Returns a list of downloadable files with URLs.

## Roles

- **Processor:** downloads raw images
- **Customer:** downloads processed data (deliverables)
- **Pilot:** downloads their own uploads
