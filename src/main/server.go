package main

import (
    "fmt"
    "net"
    "log"
    "strconv"
    "errors"
    "protocol"
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
    uid, err := checkUserIdentity(conn)
    if  err != nil {
        fmt.Println("check f", uid,err);
        return
    }
    var closed = make(chan bool)

    //声明一个管道用于接收解包的数据
    readerChannel := make(chan []byte, 16)
    
    tmpBuffer := make([]byte, 0)
    go reader(readerChannel)

    go func(){
        for {
            buffer := make([]byte, 1024)
            i, err := conn.Read(buffer)
            if(err != nil){
                fmt.Println("conn.Read error:", err);
                return;
            }
            tmpBuffer = protocol.Unpack(append(tmpBuffer, buffer[:i]...), readerChannel)
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

    if <-closed {
        return
    }
}

func reader(readerChannel chan []byte) {
    for {
        select {
            case data := <-readerChannel:
                fmt.Println("message" , string(data))
        }
    }
}

/**
 * 认证用户
 */
func checkUserIdentity(conn net.Conn) (int, error){
    for {
        _, err := conn.Write([]byte("请注册用户，输入\"I AM {YOUR NAME}.\""))
        if(err != nil){
            return 0, err;
        }

        headerBuffer := make([]byte, 7)
        _, err = conn.Read(headerBuffer)
        if err != nil  {
            return 0, errors.New("checkUserIdentity ERROR");
        }

        uidBuffer := make([]byte, 7)
        i, err := conn.Read(uidBuffer)
        if err != nil  {
            return 0, errors.New("checkUserIdentity ERROR");
        }

        var uidStr string = string(uidBuffer[0:i])
        uid, err := strconv.Atoi(uidStr)
        if err != nil{
            return 0, err;
        }
        fmt.Println("check",uid,err);
        clientMap[uid] = make(chan string)
        fmt.Println("check3333333333333",uid,err);
        return uid, nil
    }
}

/**
 * 处理业务
 */
func processingBusiness(conn net.Conn, message string){
    r, err := conn.Write([]byte("you said "+message))
    fmt.Println("processingBusiness", r, err, message)
    if(err != nil){
        return
    }
}