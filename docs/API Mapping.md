# API Mapping

Every CLI command and the API endpoint it calls. Base URL: `https://api.inflights.com/v1`

All requests include `Authorization: Bearer <token>`.

## Auth

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights login` | POST | `/auth/token-exchange` |
| `inflights whoami` | GET | `/auth/me` |

## Quotes

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights quotes` | GET | `/quotes?assignee=me` |
| `inflights quote show <id>` | GET | `/quotes/<id>` |
| `inflights quote assign <id> <user>` | PATCH | `/quotes/<id>/assign` |
| `inflights quote confirm <id>` | POST | `/quotes/<id>/confirm` |

## Flights

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights flights` | GET | `/flights?geo=…&status=…` |
| `inflights flight show <id>` | GET | `/flights/<id>` |
| `inflights flight add` | POST | `/flights` |
| `inflights flight status <id> <s>` | PATCH | `/flights/<id>` |
| `inflights flight assign-pilot` | PATCH | `/flights/<id>/assign` |
| `inflights flight assign-processor` | PATCH | `/flights/<id>/assign` |

## Data

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights upload images <id>` | POST | `/flights/<id>/images` |
| `inflights upload data <id>` | POST | `/flights/<id>/data` |
| `inflights download <id>` | GET | `/flights/<id>/download` |

## Gear

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights gear list` | GET | `/gear` |
| `inflights gear add` | POST | `/gear` |

## Services & Orders

| Command | Method | Endpoint |
|---------|--------|----------|
| `inflights services` | GET | `/services` |
| `inflights order` | POST | `/orders` |
