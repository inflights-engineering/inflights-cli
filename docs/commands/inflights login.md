# inflights login

Authenticate with inflights.com via browser-based OAuth flow.

## Usage

```bash
inflights login
```

## Behavior

1. CLI generates a random UUID token
2. Opens `https://inflights.com/login/<token>` in default browser
3. If already signed in: user sees a confirmation page ("Authorize CLI as **email**?") and can authorize or choose a different account
4. If not signed in: user logs in with email and password
5. CLI receives the bearer token via polling
6. Token saved to `~/.inflights/credentials`

## API

```
POST /v1/auth/token-exchange
```

## Roles

All roles. See [[Authentication]].
