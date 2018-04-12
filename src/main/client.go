package main

import (
    "fmt"
    // "errors"
    "math/rand"
    "net"
    "protocol"
)

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:6010")
    if err != nil {
        fmt.Println(err);
        return
    }
    defer func(){
        fmt.Println("main defer do : conn closed")
        conn.Close()
    }()
    go writeFromServer(conn)

    for {
        var talkContent string
        fmt.Scanln(&talkContent)

        if len(talkContent) > 0 {
            _, err = sendMessageTo(conn, talkContent)
            if err != nil {
                fmt.Println("write to server error")
                return
            }
        }
    }
}

func sendMessageTo(conn net.Conn, message string) (bool, error){
    pack := protocol.Packet([]byte(message))
    _, err := conn.Write(pack)
    if err != nil {
        return false, err
    }
    return true, nil
}

func writeFromServer(conn net.Conn) {
    defer func(){
        fmt.Println("writeFromServer defer do : conn closed")
        conn.Close()
    }()

    for {
        data := make([]byte, 1024)
        c, err := conn.Read(data)
        if err != nil {
            fmt.Println("rand", rand.Intn(10), "have no server write", err)
            return
        }
        fmt.Println(string(data[0:c]) + "\n ")
    }
}