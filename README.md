# sc2hub-server

Repository of sc2hub.net server side code

Start the app:

    go run ./cmd/sc2hub 
        -secret {JWT_CLIENT_SECRET} 
        -twitchClientId {TWITCH_CLIENT_ID} 
        -twitchClientSecret {TWITCH_CLIENT_SECRET} 
        -youtube_api_key {YOUTUBE_API_KEY} 

Other flags:

    -prod: if true, start HTTPS server instead of HTTP. Default: false
    -addr: HTTP network address. Default: :443
    -dsn: MySQL data source name. Default: root:root@/sc2hub
    -origin: Client app origin URL (for handling preflight requests). Default: http://localhost:4200

Create DB users:

    go run ./cmd/create_user -dsn {DSN}
    Enter username, press enter
    Enter password, press enter
