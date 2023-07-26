package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
    address := ":5000"
    if len(os.Args) > 1 {
	    address = os.Args[1]
    }

    start_server(address)
}

func start_server(address string) {
    fmt.Printf("Starting server...\n")

    listener, err := net.Listen("tcp", address)
    if err != nil {
        fmt.Printf("Error while launching server: %s\n", err.Error())
        return
    }

    fmt.Printf("Server started!\n")

    for {
        fmt.Printf("Waiting for connection from client...\n")
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Error accepting new connection: %s\n", err.Error())
            return
        }

        fmt.Printf("Client connected!\n")

        err = receive_data(conn)
        if err != nil {
            fmt.Printf("Error receiving data from client: %s", err.Error())
            continue
        }
        err = send_data(conn)
        if err != nil {
            fmt.Printf("Error sending data to client: %s", err.Error())
            continue
        }
    }
}

func receive_data(conn net.Conn) error {
    fmt.Printf("Receiving data from client...\n")
    buffer := make([]byte, 1024)
    read, err := conn.Read(buffer)
    if err != nil {
        return err
    }
    data := string(buffer[:read])
    fmt.Printf("Received data from client: %s", data)
    return nil
}

func send_data(conn net.Conn) error {
    fmt.Printf("Sending data to client...\n")

    html := "<!DOCTYPE html><html lang=\"ru\"><head><meta charset=\"UTF-8\"><title>Test</title></head><body>Ку, Андрюх!</body></html>"

    header := "HTTP/1.1 200 OK\n"
    date := "Date: Wed, 11 Feb 2009 11:20:59 GMT\n"
    server := "Server: Custom\n"
    last_modified := "Last-Modified: Wed, 11 Feb 2009 11:20:59 GMT\n"
    content_language := "Content-Language: ru\n"
    content_type := "Content-Type: text/html; charset=utf-8\n"
    content_length := "Content-Length: " + string(len(html)) + "\n"
    connection := "Connection: close\n"


    text := header + date + server + last_modified + content_language + content_type + content_length + connection + "\n" + html
    _, err := conn.Write([]byte(text))
    conn.Close()
    if err != nil {
        return err
    }
    fmt.Printf("Sent data to client!\n")
    return nil
}
