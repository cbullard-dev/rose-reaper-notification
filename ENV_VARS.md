# Environment Variables

This service requires the following environment variables to be set:

## Required Environment Variables

### Twitch Configuration
- `TWITCH_WEBSOCKET_URL` - Twitch EventSub WebSocket URL (e.g., `wss://eventsub.wss.twitch.tv/ws`)
- `TWITCH_CLIENT_ID` - Your Twitch application client ID
- `TWITCH_CLIENT_SECRET` - Your Twitch application client secret
- `TWITCH_REFRESH_TOKEN` - Twitch OAuth refresh token
- `TWITCH_BEARER_TOKEN` - Twitch OAuth bearer/access token
- `TWITCH_CLIENT_CODE` - Twitch OAuth client code
- `TWITCH_BROADCASTER_USER_ID` - The Twitch user ID of the broadcaster to monitor

### Discord Configuration
- `DISCORD_WEBHOOK_URL` - Discord webhook URL for sending notifications

## Optional Environment Variables

- `LOG_LEVEL` - Logging level (default: `info`). Valid values: `debug`, `info`, `warn`, `error`
- `DATABASE_URL` - PostgreSQL database connection string (currently unused, reserved for future use)

## Example

Create a `.env` file (make sure it's in `.gitignore`):

```bash
TWITCH_WEBSOCKET_URL=wss://eventsub.wss.twitch.tv/ws
TWITCH_CLIENT_ID=your_client_id
TWITCH_CLIENT_SECRET=your_client_secret
TWITCH_REFRESH_TOKEN=your_refresh_token
TWITCH_BEARER_TOKEN=your_bearer_token
TWITCH_CLIENT_CODE=your_client_code
TWITCH_BROADCASTER_USER_ID=your_broadcaster_user_id
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/your_webhook_url
LOG_LEVEL=info
```

## Security Notes

⚠️ **Never commit `.env` files or hardcode these values in your source code.**

These credentials should be:
- Stored securely in your deployment environment
- Passed as environment variables or secrets
- Never exposed in logs or error messages

