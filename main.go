package main

import (
	"errors"
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
        conn, err := listener.Accept();
        if err != nil {
            fmt.Printf("Error accepting new connection: %s\n", err.Error());
            return;
        }

        headers, err := receive_request(conn);
        if err != nil {
            fmt.Printf("Error receiving data from client: %s", err.Error());
            continue;
        }
        fmt.Printf("%s\n", headers["Status"]);

        var html string = create_html(headers["Route"]);
        var response_data response_data = create_response_data().set_html(html).set_status("200 OK");

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
    var handler handler = get_handler(route);
    var body string = handler.get_body();
    var html string = fmt.Sprintf(template, title, body);
    return html;
}

type handler interface {
    get_body() string
}

var (
    route_index string = "/";
    greeter_index string = "/greet";
);


func get_handler(route string) handler {
    if(route == "/"){
        return index{};
    }
    fmt.Printf("Route: %s\n", route);
    if(string_starts_with(route, "/greet")){
        rest, err := string_get_rest(route, "/greet")
        if err != nil {
            return server_error{error_message: err.Error()}
        }

        rest, err = string_get_rest(rest, "?");
        if err != nil {
            return server_error{error_message: "You did not specify your name!"}
        }

        name, err := string_get_rest(rest, "name=");
        if err != nil {
            return server_error{error_message: "You did not specify your name!"}
        }

        return greeter{name: name}
    }

    return not_found{};
}

func string_starts_with(source string, target string) bool {
    var target_len int = len(target);
    var source_len int = len(source);
    return source_len >= target_len &&  source[0:target_len] == target;
}

func string_get_rest(route string, after string) (string, error) {
    if(!string_starts_with(route, after)) {
        return "", errors.New("Route doesn't start with specified string");
    }
    var target_len int = len(after);
    return route[target_len:], nil
}

type greeter struct {
    name string
}
func (data greeter) get_body() string {
    return fmt.Sprintf("Hello, %s!\n", data.name);
}

type not_found struct {}
func (_ not_found) get_body() string {
    return "Page not found";
}

type server_error struct {
    error_message string
}
func (data server_error) get_body() string {
    return fmt.Sprintf("Server error: %s\n", data.error_message);
}

type index struct {}
func (_ index) get_body() string {
    return "Index page";
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

func (data response_data) set_html(html string) response_data {
    data.html = html;
    data.content_length = len(html);
    return data
}

func (data response_data) set_status(status string) response_data {
    data.status = fmt.Sprintf("HTTP/1.1 %s\n", status);
    return data;
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
