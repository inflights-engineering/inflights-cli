# API Mapping

Every CLI command and the API endpoint it calls. Base URL: `https://inflights.com/api/v1`

All authenticated requests include `Authorization: Bearer <token>`.

Flight endpoints accept either numeric ID or public UID (e.g. `FL-1042`).

Interactive docs: https://developer.inflights.com/v1
OpenAPI spec: https://developer.inflights.com/v1/openapi

## Auth

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights login` | POST | `/auth/login_tokens` → `/auth/token_exchange` |
| `inflights whoami` | GET | `/auth/me` |

## Services

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights services` | GET | `/services` |

## Flights

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights flights` | GET | `/flights?status=...&public_uid=...` |
| `inflights flight <id>` | GET | `/flights/:id` |
| `inflights order <geojson>` | POST | `/flights` |

## Quotes

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights quotes` | GET | `/quotes?status=...` |
| `inflights quote show <id>` | GET | `/quotes/:id` |
| `inflights quote confirm <id>` | POST | `/quotes/:id/accept` or `/quotes/accept_estimate` |
| `inflights quote reject <id>` | POST | `/quotes/:id/reject` or `/quotes/reject_estimate` |

## Proposals

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights proposals` | GET | `/proposals?status=...` |
| `inflights proposal show <id>` | GET | `/proposals/:id` |
| `inflights proposal accept <id>` | POST | `/proposals/:id/accept` |
| `inflights proposal reject <id>` | POST | `/proposals/:id/reject` |

## Equipment

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights gear list` | GET | `/equipment_types?category=...` |
| `inflights gear mine` | GET | `/equipments` |
| `inflights gear add` | POST | `/equipments` |
| `inflights gear remove` | DELETE | `/equipments/:id` |

## Uploads (Data)

| Command | Method | Endpoint | Notes |
|---------|--------|----------|-------|
| `inflights upload data <id> [files...]` | POST | `/flights/:id/uploads/presign` | 1. Presign |
| | POST | S3 presigned URL | 2. Upload to S3 |
| | POST | `/flights/:id/uploads/confirm` | 3. Confirm |

## Uploads (Images)

| Command | Method | Endpoint | Notes |
|---------|--------|----------|-------|
| `inflights upload images <id> [path...]` | POST | `/flights/:id/images/presign` | 1. Presign (per image) |
| | POST | S3 presigned URL | 2. Upload to S3 |
| | POST | `/flights/:id/images/confirm` | 3. Confirm (per image, with EXIF) |
| | POST | `/flights/:id/images/finalize` | 4. Finalize (once, all succeeded) |

## Downloads

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights download <id>` | GET | `/flights/:id/downloads` |
