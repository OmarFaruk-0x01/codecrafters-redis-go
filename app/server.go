package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
	"os"
)

func nilChecker(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
		return
	}
}

func parser(data string) (string, error) {
	p := regexp.MustCompile("\\$[0-9]+\r\n([^\r\n]+)\r\n")
	matches := p.FindAllStringSubmatch(data, -1)

	for _, match := range matches {
		return strings.Trim(match[1], "\r\n"), nil
	}
	return "", errors.New("Unable to parse command.")
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("+%v\r\n", message)))
}

func sendError(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("-%v\r\n", message)))
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	nilChecker(err)

	defer l.Close()

	for {
		conn, err := l.Accept()
		nilChecker(err)

		go func(conn net.Conn) {
			buf := make([]byte, 1024)

			for {
				data, err := conn.Read(buf)
				nilChecker(err)
				if cmd, err := parser(string(buf[:data])); err == nil {
					if strings.ToUpper(cmd) == "PING" {
						sendMessage(conn, "PONG")
					} else {
						sendError(conn, "command not found")
					}
				} else {
					sendError(conn, "unable to parse command.")
				}
				//
			}

		}(conn)

	}

}
