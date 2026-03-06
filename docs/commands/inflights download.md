# inflights download

Download data for a given flight (images or processed outputs, depending on your role).

## Usage

```bash
inflights download <flightId> [options]
```

## Options

| Flag | Description |
|------|-------------|
| `--type <type>` | `images`, `data`, or `all` (default: role-appropriate) |
| `--output <dir>` | Destination directory (default: `./<flightId>/`) |

## Example

```bash
inflights download FL-1042 --type data --output ./deliverables/
# → Downloading processed data for FL-1042 (1.8 GB)…
# → Saved to ./deliverables/FL-1042/
```

## API

```
GET /v1/flights/<flightId>/download?type=...
```

Returns a signed URL or streams the archive.

## Roles

- **Processor:** downloads raw images
- **Customer:** downloads processed data (deliverables)
- **Pilot:** downloads their own uploads
