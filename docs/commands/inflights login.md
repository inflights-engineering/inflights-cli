# inflights login

Authenticate with inflights.com via browser-based OAuth flow.

## Usage

```bash
inflights login
```

## Behavior

1. CLI generates a random UUID token
2. Opens `https://inflights.com/login/<token>` in default browser
3. User logs in on the web page
4. CLI receives the bearer token via callback / polling
5. Token saved to `~/.inflights/credentials`

## API

```
POST /v1/auth/token-exchange
```

## Roles

All roles. See [[Authentication]].
