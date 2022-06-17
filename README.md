# chat-go [![Go](https://github.com/cothi/chat-go/actions/workflows/go.yml/badge.svg)](https://github.com/cothi/chat-go/actions/workflows/go.yml)

Terminal based chat made with go

## demo
![demo](/assets/demo_view2.png)


## start

```bash
git clone https://github.com/cothi/chat-go/new/main?readme=1

cd chat-go
go mod tidy

# server
go run server/server.go

# other terminal (client)
go run main.go

# other terminal (client 2)
go run main.go
```

## TODO
1. [ ] cli (server, client, set port)
2. [ ] chennel and lobby
3. [ ] command (set nickname, leave room)
