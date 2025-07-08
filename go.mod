module lesta-battleship/cli

go 1.24.4

require github.com/gorilla/websocket v1.5.3

replace github.com/lesta-battleship/server-core => github.com/lesta-start-battleship/server-core v1.0.0

require github.com/lesta-battleship/matchmaking v0.2.0

require github.com/golang-jwt/jwt/v5 v5.2.2

replace github.com/lesta-battleship/matchmaking => github.com/lesta-start-battleship/matchmaking v0.2.0
