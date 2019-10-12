# sc2hub-server

Server side repository for sc2hub.net

Copy .env.template to .env and fill in values

Start the app:

    go run ./cmd/sc2hub -secret {JWT_CLIENT_SECRET}
    
Create DB users:

    go run ./cmd/create_user
    Enter username, press enter
    Enter password, press enter

Other flags:

    -prod: if true, start HTTPS server instead of HTTP. Default: false
    -addr: HTTP network address. Default: :443
    -dsn: MySQL data source name. Default: root:root@/sc2hub
    -origin: Client app origin URL (for handling preflight requests). Default: http://localhost:4200
