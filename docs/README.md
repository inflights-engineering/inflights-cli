# Inflights CLI

Command-line tool for managing drone flights, quotes, and data processing on [inflights.com](https://inflights.com).

## Who is this for?

| Role | What they do |
|------|-------------|
| **Customer** | Requests quotes, orders services, downloads deliverables |
| **Pilot** | Manages gear, flies missions, uploads drone imagery |
| **Processor** | Downloads raw data, uploads processed deliverables |

## Quick start

```bash
# Install (TBD — npm, brew, or binary)
inflights login          # authenticate via browser
inflights flights        # list your flights
inflights quotes         # list quotes assigned to you
```

## Documentation map

- [[Command Reference]] — full list of every command
- **Roles:** [[Customer]], [[Pilot]], [[Processor]]
- **Workflows:** [[Requesting a Quote]], [[Flight Lifecycle]], [[Data Pipeline]]
- [[Authentication]] — how login and tokens work
- [[API Mapping]] — which API endpoints each command hits
- [[Glossary]] — terms used throughout the CLI
