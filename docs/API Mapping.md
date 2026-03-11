# API Mapping

Every CLI command and the API endpoint it calls. Base URL: `https://inflights.com/api/v1`

All authenticated requests include `Authorization: Bearer <token>`.

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
| `inflights order` | POST | `/flights` |

## Equipment

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights gear list` | GET | `/equipment_types?category=...` |
| `inflights gear mine` | GET | `/equipments` |
| `inflights gear add` | POST | `/equipments` |
| `inflights gear remove` | DELETE | `/equipments/:id` |
| `inflights gear rfc` | — | Not yet implemented in API v1 |

## Flights

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights flights` | GET | `/flights?status=...&public_uid=...` |
| `inflights flight show <id>` | GET | `/flights/:id` |
| `inflights flight add` | POST | `/flights` |
| `inflights flight status <id> <s>` | PATCH | `/flights/:id` |
| `inflights flight assign-pilot` | PATCH | `/flights/:id/assign` |
| `inflights flight assign-processor` | PATCH | `/flights/:id/assign` |

## Quotes

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights quotes` | GET | `/quotes?status=...` |
| `inflights quote show <id>` | GET | `/quotes/:id` |
| `inflights quote assign <id> <user>` | PATCH | `/quotes/:id/assign` |
| `inflights quote confirm <id>` | POST | `/quotes/:id/accept` or `/quotes/accept_estimate` |

## Proposals

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights proposal list` | GET | `/proposals?status=...` |
| `inflights proposal show <id>` | GET | `/proposals/:id` |
| `inflights proposal accept <id>` | POST | `/proposals/:id/accept` |
| `inflights proposal reject <id>` | POST | `/proposals/:id/reject` |

## Uploads (Data)

| Command | Method | Endpoint | Notes |
|---------|--------|----------|-------|
| `inflights upload data <id>` | POST | `/flights/:id/uploads/presign` | 1. Presign |
| | POST | S3 presigned URL | 2. Upload to S3 |
| | POST | `/flights/:id/uploads/confirm` | 3. Confirm |

## Uploads (Images)

| Command | Method | Endpoint | Notes |
|---------|--------|----------|-------|
| `inflights upload images <id>` | POST | `/flights/:id/images/presign` | 1. Presign (per image) |
| | POST | S3 presigned URL | 2. Upload to S3 |
| | POST | `/flights/:id/images/confirm` | 3. Confirm (per image) |
| | POST | `/flights/:id/images/finalize` | 4. Finalize (once) |

## Downloads

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights download <id>` | GET | `/flights/:id/downloads` |
