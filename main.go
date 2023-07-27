package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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

    fmt.Printf("Server listening on %s!\n", listener.Addr().String())

    for {
        //fmt.Printf("Waiting for connection from client...\n")
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("Error accepting new connection: %s\n", err.Error())
            return
        }

        //fmt.Printf("Client connected!\n")

        headers, err := receive_data(conn)
        if err != nil {
            fmt.Printf("Error receiving data from client: %s", err.Error())
            continue
        }
        fmt.Printf("%s\n", headers["Status"])

        err = send_data(conn)
        if err != nil {
            fmt.Printf("Error sending data to client: %s", err.Error())
            continue
        }
    }
}

func receive_data(conn net.Conn) (map[string]string, error) {
    //fmt.Printf("Receiving data from client...\n")
    buffer := make([]byte, 1024)
    read, err := conn.Read(buffer)
    if err != nil {
        return nil, err
    }
    data := string(buffer[:read])
    lines := strings.Split(data, "\n")
    headers := make(map[string]string)

    headers["Status"] = lines[0]
    for _, line := range lines[1:] {
        split_index := strings.Index(line, ":")
        if split_index != -1 {
            headers[line[:split_index]] = line[split_index+1:]
        }
    }
    return headers, nil
}

func send_data(conn net.Conn) error {
    //fmt.Printf("Sending data to client...\n")

    //html := "<!DOCTYPE html><html lang=\"ru\"><head><meta charset=\"UTF-8\"><title>Test</title></head><body>Ку, Андрюх!</body></html>"
    html := "Ку, Андрюх!"
    response := create_response(html)

    _, err := conn.Write([]byte(response))
    conn.Close()
    if err != nil {
        return err
    }
    //fmt.Printf("Sent data to client!\n")
    return nil
}

func create_response(html string) []byte {
    status := "HTTP/1.1 200 OK\n"
    date := "Date: " + time.Now().Format("Wed, 11 Feb 2009 11:20:59 GMT") + "\n"
    server := "Server: Custom\n"
    last_modified := "Last-Modified: " + time.Now().Format("Wed, 11 Feb 2009 11:20:59 GMT") + "\n"
    content_language := "Content-Language: ru\n"
    content_type := "Content-Type: text/html; charset=utf-8\n"
    content_length := "Content-Length: " + string(len(html)) + "\n"
    connection := "Connection: close\n"
    headers := status + date + server + last_modified + content_language + content_type + content_length + connection

    response := headers + "\n" + html

    return []byte(response)
}
