package main

import(
    "fmt"
    "github.com/appcrash/media/server"
    "github.com/wernerd/GoRTP/src/net/rtp"
)

func main() {
    s := rtp.Session{}
    fmt.Println(s)
    server.InitServer(4000)
}

