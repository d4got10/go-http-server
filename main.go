package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type headers = map[string]string;
var status_separator string = " ";

func main() {
    address := ":5000"
    if len(os.Args) > 1 {
	    address = os.Args[1]
    }

    start_server(address)
}

func start_server(address string) {
    fmt.Printf("Starting server...\n");

    listener, err := net.Listen("tcp", address);
    if err != nil {
        fmt.Printf("Error while launching server: %s\n", err.Error());
        return;
    }

    fmt.Printf("Server listening on %s!\n", listener.Addr().String());

    for {
        //fmt.Printf("Waiting for connection from client...\n")
        conn, err := listener.Accept();
        if err != nil {
            fmt.Printf("Error accepting new connection: %s\n", err.Error());
            return;
        }

        //fmt.Printf("Client connected!\n")

        headers, err := receive_request(conn);
        if err != nil {
            fmt.Printf("Error receiving data from client: %s", err.Error());
            continue;
        }
        fmt.Printf("%s\n", headers["Status"]);

        var html string = create_html(headers["Route"]);
        var response []byte = create_response(html);
        err = send_response(response, conn);
        if err != nil {
            fmt.Printf("Error sending data to client: %s", err.Error());
            continue;
        }
    }
}

func receive_request(conn net.Conn) (headers, error) {
    buffer := make([]byte, 1024);
    read, err := conn.Read(buffer);
    if err != nil {
        return nil, err;
    }
    data := string(buffer[:read]);
    lines := strings.Split(data, "\n");
    headers := make(headers);

    var status string = lines[0];
    var splitted_status []string = strings.Split(status, status_separator);

    headers["Method"] = splitted_status[0];
    headers["Route"] = splitted_status[1];
    for _, line := range lines[1:] {
        split_index := strings.Index(line, ":");
        if split_index != -1 {
            headers[line[:split_index]] = line[split_index+1:];
        }
    }
    return headers, nil;
}

func send_response(data []byte, conn net.Conn) error {
    //fmt.Printf("Sending data to client...\n")

    _, err := conn.Write(data)
    conn.Close()
    if err != nil {
        return err
    }
    //fmt.Printf("Sent data to client!\n")
    return nil
}

func create_html(route string) string {
    var template string = "<!DOCTYPE html><html lang=\"ru\"><head><meta charset=\"UTF-8\"><title>%s</title></head><body>%s</body></html>\n";
    var title string = "Test server";
    var body string = route;
    var html string = fmt.Sprintf(template, title, body);
    return html;
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
