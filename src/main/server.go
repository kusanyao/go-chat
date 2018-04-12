package main

import (
    "fmt"
    "net"
    "log"
    "strconv"
    "errors"
)



var clientMap = make(map[int]chan string)

func main() {
    // 建立监听链接
    ln, err := net.Listen("tcp", "127.0.0.1:6010")
    fmt.Println(ln, err)

    for {
        fmt.Println("wait connect...")
        conn, err := ln.Accept()
        fmt.Println("CONNECT：", conn, err)
        if err != nil {
            log.Fatal("get client connection error: ", err)
        }
        go handleConnection(conn)
    }
}

func writeMessageToClient(uid int, message string)(bool){
    clientMap[uid] <- message
    return true
}

/**
 * 处理链接
 */
func handleConnection(conn net.Conn){
    defer func() { // 函数返回前关掉当前链接
        fmt.Println("defer do : conn closed")
        conn.Close()
    }()
    var closed = make(chan bool)
    uid, err := checkUserIdentity(conn)
    if  err != nil {
        fmt.Println("check f", uid,err);
        return
    }
    go func(){
        for {
            data   := make([]byte, 1024)
            i, err := conn.Read(data)
            if(err != nil){
                fmt.Println("conn.Read error:", err);
                return;
            }
            message := string(data[0:i])
            fmt.Println("MESSAGE", message)
            go processingBusiness(conn, message)
        }
    }()

    // 从chan 里读出给这个客户端的数据 然后写到该客户端里
    go func(){
        for {
            message := <-clientMap[uid]
            _, err := conn.Write([]byte(message))
            if err != nil {
                closed <- true
            }
        }
    }()

    for {
        if <-closed {
            return
        }
    }
}

/**
 * 认证用户
 */
func checkUserIdentity(conn net.Conn) (int, error){
    data   := make([]byte, 10)
    i, err := conn.Read(data)
    if(err != nil){
        return 0, errors.New("checkUserIdentity ERROR");
    }
    var uidStr string = string(data[0:i])
    uid, err := strconv.Atoi(uidStr)
    if err != nil{
        return 0, errors.New("checkUserIdentity ERROR2");
    }
    fmt.Println("check",uid,err);
    clientMap[uid] = make(chan string)
    fmt.Println("check3333333333333",uid,err);
    return uid, nil
}

/**
 * 处理业务
 */
func processingBusiness(conn net.Conn, message string){
    r, err := conn.Write([]byte(message))
    fmt.Println("processingBusiness", r, err, message)
    if(err != nil){
        return
    }
}