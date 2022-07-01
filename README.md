# chat-go [![Go](https://github.com/cothi/chat-go/actions/workflows/go.yml/badge.svg)](https://github.com/cothi/chat-go/actions/workflows/go.yml)

Terminal based chat made with go

## demo
![demo](/assets/demo_view2.png)


## start

```bash
git clone https://github.com/cothi/chat-go.git

cd chat-go
go mod tidy

go build ./



# start server

./chat-go server --serverPort=8000


# start client 
./chat-go client --serverPort=8000

```

## TODO
1. [x] cli (server, client, set port)
2. [x] chennel
3. [x] lobby
4. [x] command (set nickname,
5. [ ] command (leave room)
