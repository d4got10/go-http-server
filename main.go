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
        var response_data response_data = create_response_data();
        response_data = set_html(response_data, html);

        var response []byte = create_response(response_data);
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
    _, err := conn.Write(data)
    conn.Close()
    if err != nil {
        return err
    }
    return nil
}

func create_html(route string) string {
    var template string = "<!DOCTYPE html><html lang=\"ru\"><head><meta charset=\"UTF-8\"><title>%s</title></head><body>%s</body></html>\n";
    var title string = "Test server";
    var body string = route;
    var html string = fmt.Sprintf(template, title, body);
    return html;
}

type response_data struct {
    status string
    date time.Time
    server string
    last_modified time.Time
    content_language string
    content_type string
    content_length int
    connection string
    html string
}

func set_html(data response_data, html string) response_data {
    data.html = html;
    data.content_length = len(html);
    return data
}

func create_response_data() response_data {
    var data response_data = response_data {
        status: "HTTP/1.1 200 OK\n",
        date: time.Now(),
        server: "Server: Custom\n",
        last_modified: time.Now(),
        content_language: "Content-Language: ru\n",
        content_type: "Content-Type: text/html; charset=utf-8\n",
        content_length: 0,
        connection: "Connection: close\n",
        html: "",
    };
    return data
}

var date_format string = "Mon, 2 Jan 2006 15:04:05 GMT";

func create_response(data response_data) []byte {

    date := fmt.Sprintf("Date: %s\n", data.date.Format(date_format));
    content_length := fmt.Sprintf("Content-Length: %d\n", data.content_length);
    last_modified := fmt.Sprintf("Last-Modified: %s\n", data.last_modified.Format(date_format));
    response := fmt.Sprintf("%s%s%s%s%s%s%s%s\n%s\n",
                  data.status,
                  date,
                  data.server,
                  last_modified,
                  data.content_language,
                  data.content_type,
                  content_length,
                  data.connection,
                  data.html);

    fmt.Println(response);
    return []byte(response)
}
