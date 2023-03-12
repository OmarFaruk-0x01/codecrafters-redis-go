package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// connections := map[string]net.Conn{}
	// var rm sync.RWMutex

	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	nilChecker(err)

	defer l.Close()

	newConnection := make(chan net.Conn)
	closeConnection := make(chan net.Conn)

	go func() {
		for {
			conn, err := l.Accept()
			nilChecker(err)
			newConnection <- conn
		}
	}()

	go func() {
		conn := <-closeConnection
		fmt.Printf("Disconnected: %v", conn.RemoteAddr().String())
	}()

	for {
		select {
		case conn := <-newConnection:
			go handleConnection(conn, &closeConnection)
		}
	}
}

func handleConnection(conn net.Conn, closeConnection *chan net.Conn) {
	defer conn.Close()
	// defer wg.Done()

	fmt.Println(conn.RemoteAddr())

	for {
		buf := make([]byte, 1024)
		data, err := conn.Read(buf)
		if err != nil {
			sendError(conn, "command not found")
			*closeConnection <- conn
			return
		}
		if data == 0 {
			// Connection closed by client.
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}
		if commands, err := parser(string(buf[:data])); err == nil {
			if command, err := commandParser(commands); err == nil {
				handleCommand(conn, command)
			} else {
				sendError(conn, err.Error())
			}
		} else {
			sendError(conn, err.Error())
		}
	}
}

func handleCommand(conn net.Conn, command *Command) {
	switch strings.ToUpper(command.cmd) {
	case "PING":
		sendMessage(conn, "PONG")
		return
	case "ECHO":
		sendMessage(conn, command.arguments[0])
	default:
		sendError(conn, "command not found")
		return
	}

}
