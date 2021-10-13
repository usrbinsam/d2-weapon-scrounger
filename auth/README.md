# OAuth2-client Server

This package provides the browser-based OAuth2 flow for Bungie.

## Setup

The following environmental variables are used:

### `CLIENT_ID`

OAuth2 Client ID from Bungie.net

### `CLIENT_SECRET`

OAuth2 Client Secret from Bungie.net

### `AUTH_LISTEN_ADDR`

(passed in from main.go) The address:port to listen on. The default is `":80"`

### `BASE_URL`

The protocol, domain, and port that this server is publically accessible on. For local environment, use ngrok and update this value at Bungie.net/en/Application


## Flow

1. The user visits `/start?discord_id=XYZ` where the Discord ID is their Discord API user ID. The auth server responds with HTTP 302 and a Location header to redirect the user to the Bungie.net authorization site. This also sets a `state` cookie to prevent CSRF.
2. Upon authorization, Bungie.net will redirect the user back to `/auth`. The User is looked up by the State value in their browser
3. Using the auth code, the backend requests an access & refresh token, then stores them in the DB associated to the user's Discord ID.